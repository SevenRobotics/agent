package inbound

import (
	"context"
	"fmt"
	"go_agent/config"
	"go_agent/publishers/rmq"
	rosgeometrymsgs "go_agent/telemetry/gengo/ros/geometry_msgs"
	geometrymsgs "go_agent/telemetry/genproto/ros/geometry_msgs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bluenviron/goroslib/v2"
	"github.com/bluenviron/goroslib/v2/pkg/apimaster"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

const (
	twistMessageType     = "geometry_msgs/Twist"
	consumerRetryDelay   = 5 * time.Second
	rosStatePollInterval = 1 * time.Second
	receiverClientPrefix = "inbound-"
)

type Receiver struct {
	rmqConfig config.RMQConfig
	rosConfig config.RosNodeConfig
	config    config.RMQInboundReceiverConfig

	rosSubscriberReady func(context.Context) (bool, error)
	rosPollInterval    time.Duration
}

func NewReceiver(rmqConfig config.RMQConfig, rosConfig config.RosNodeConfig, receiverConfig config.RMQInboundReceiverConfig) (*Receiver, error) {
	if receiverConfig.Name == "" {
		return nil, fmt.Errorf("inbound receiver name is required")
	}
	if receiverConfig.Queue == "" {
		return nil, fmt.Errorf("inbound receiver %s queue is required", receiverConfig.Name)
	}
	if receiverConfig.MessageType == "" {
		return nil, fmt.Errorf("inbound receiver %s message_type is required", receiverConfig.Name)
	}
	if receiverConfig.MessageType != twistMessageType {
		return nil, fmt.Errorf("inbound receiver %s unsupported message_type %q", receiverConfig.Name, receiverConfig.MessageType)
	}
	if receiverConfig.RosTopic == "" {
		return nil, fmt.Errorf("inbound receiver %s ros_topic is required", receiverConfig.Name)
	}
	if rosConfig.Address == "" {
		return nil, fmt.Errorf("inbound receiver %s ROS master address is required", receiverConfig.Name)
	}
	if receiverConfig.Consumer == "" {
		receiverConfig.Consumer = receiverClientPrefix + receiverConfig.Name
	}
	if receiverConfig.RosNodeName == "" {
		receiverConfig.RosNodeName = receiverConfig.Name + "_cmd_vel_publisher"
	}
	if receiverConfig.WaitForRosSubscriber == "" {
		receiverConfig.WaitForRosSubscriber = expectedTelemetrySubscriberNode(receiverConfig.RosTopic)
	}

	return &Receiver{
		rmqConfig: rmqConfig,
		rosConfig: rosConfig,
		config:    receiverConfig,
	}, nil
}

func StartEnabled(ctx context.Context, rmqConfig config.RMQConfig, rosConfig config.RosNodeConfig, inboundConfig config.RMQInboundConfig, wg *sync.WaitGroup) error {
	started := 0
	for _, receiverConfig := range inboundConfig.Receivers {
		if !receiverConfig.Enabled {
			continue
		}

		receiver, err := NewReceiver(rmqConfig, rosConfig, receiverConfig)
		if err != nil {
			return err
		}

		wg.Add(1)
		go receiver.Run(ctx, wg)
		started++
	}

	if started == 0 {
		log.Printf("No enabled inbound RabbitMQ receivers configured")
	}

	return nil
}

func (r *Receiver) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if err := r.consume(ctx); err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Inbound receiver %s stopped: %v. Retrying in %s", r.config.Name, err, consumerRetryDelay)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(consumerRetryDelay):
		}
	}
}

func (r *Receiver) consume(ctx context.Context) error {
	if err := r.waitForROSSubscriber(ctx); err != nil {
		return err
	}

	rosPub, err := r.newROSPublisher()
	if err != nil {
		return fmt.Errorf("create ROS publisher: %w", err)
	}
	defer rosPub.close()

	conn, err := rmq.NewRabbitMQ(r.rmqConfig)
	if err != nil {
		return fmt.Errorf("connect RabbitMQ: %w", err)
	}

	client, err := conn.NewClient(receiverClientPrefix + r.config.Name)
	if err != nil {
		return fmt.Errorf("create RabbitMQ client: %w", err)
	}

	log.Printf(
		"Starting inbound receiver %s: exchange=%q queue=%q consumer=%q message_type=%q ros_topic=%q auto_ack=%t",
		r.config.Name,
		r.config.Exchange,
		r.config.Queue,
		r.config.Consumer,
		r.config.MessageType,
		r.config.RosTopic,
		r.config.AutoAck,
	)

	consumeCtx, cancelConsume := context.WithCancel(ctx)
	defer cancelConsume()

	deliveries, err := client.Receive(consumeCtx, r.config.Queue, r.config.Consumer, r.config.AutoAck)
	if err != nil {
		return fmt.Errorf("consume queue %s: %w", r.config.Queue, err)
	}

	rosMonitor := r.monitorROSSubscriber(consumeCtx)

	for {
		select {
		case <-ctx.Done():
			return nil
		case err, ok := <-rosMonitor:
			if ok && err != nil {
				return err
			}
			rosMonitor = nil
		case delivery, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("delivery stream closed for queue %s", r.config.Queue)
			}
			r.handleDelivery(delivery, rosPub)
		}
	}
}

func (r *Receiver) handleDelivery(delivery amqp.Delivery, rosPub *twistPublisher) {
	msg := &geometrymsgs.Twist{}
	if err := proto.Unmarshal(delivery.Body, msg); err != nil {
		log.Printf(
			"Inbound receiver %s failed to deserialize %s from queue=%q delivery_tag=%d: %v",
			r.config.Name,
			r.config.MessageType,
			r.config.Queue,
			delivery.DeliveryTag,
			err,
		)
		r.rejectMalformed(delivery)
		return
	}

	rosMsg := protoTwistToROS(msg)
	rosPub.publish(rosMsg)

	if !r.config.AutoAck {
		if err := delivery.Ack(false); err != nil {
			log.Printf(
				"Inbound receiver %s failed to ack queue=%q delivery_tag=%d: %v",
				r.config.Name,
				r.config.Queue,
				delivery.DeliveryTag,
				err,
			)
		}
	}
}

func (r *Receiver) waitForROSSubscriber(ctx context.Context) error {
	if r.config.WaitForRosSubscriber == "" {
		return nil
	}

	ticker := time.NewTicker(r.getROSPollInterval())
	defer ticker.Stop()

	for {
		ready, err := r.isROSSubscriberReady(ctx)
		if err == nil && ready {
			return nil
		}
		r.logROSWait(err)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (r *Receiver) monitorROSSubscriber(ctx context.Context) <-chan error {
	monitorErr := make(chan error, 1)
	if r.config.WaitForRosSubscriber == "" {
		close(monitorErr)
		return monitorErr
	}

	go func() {
		defer close(monitorErr)

		ticker := time.NewTicker(r.getROSPollInterval())
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			ready, err := r.isROSSubscriberReady(ctx)
			if ctx.Err() != nil {
				return
			}
			if err != nil {
				monitorErr <- fmt.Errorf(
					"check ROS subscriber %s for topic %s: %w",
					r.config.WaitForRosSubscriber,
					r.config.RosTopic,
					err,
				)
				return
			}
			if !ready {
				monitorErr <- fmt.Errorf(
					"ROS subscriber %q not available for topic %q",
					r.config.WaitForRosSubscriber,
					r.config.RosTopic,
				)
				return
			}
		}
	}()

	return monitorErr
}

func (r *Receiver) isROSSubscriberReady(ctx context.Context) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if r.rosSubscriberReady != nil {
		return r.rosSubscriberReady(ctx)
	}

	client := apimaster.NewClient(r.rosConfig.Address, r.config.RosNodeName, &http.Client{})
	return hasTopicSubscriber(client, r.config.RosTopic, r.config.WaitForRosSubscriber)
}

func (r *Receiver) logROSWait(err error) {
	if err != nil {
		log.Printf(
			"Inbound receiver %s waiting for ROS master/subscriber %q on topic %q before publishing: %v",
			r.config.Name,
			r.config.WaitForRosSubscriber,
			r.config.RosTopic,
			err,
		)
		return
	}

	log.Printf(
		"Inbound receiver %s waiting for ROS subscriber %q on topic %q before publishing",
		r.config.Name,
		r.config.WaitForRosSubscriber,
		r.config.RosTopic,
	)
}

func (r *Receiver) getROSPollInterval() time.Duration {
	if r.rosPollInterval > 0 {
		return r.rosPollInterval
	}
	return rosStatePollInterval
}

func hasTopicSubscriber(client *apimaster.Client, topic string, subscriberName string) (bool, error) {
	state, err := client.GetSystemState()
	if err != nil {
		return false, err
	}

	for _, sub := range state.SubscribedTopics {
		if sub.Name != topic {
			continue
		}
		for _, node := range sub.Nodes {
			if sameROSName(node, subscriberName) {
				return true, nil
			}
		}
	}

	return false, nil
}

func expectedTelemetrySubscriberNode(topic string) string {
	nodeName := strings.Trim(strings.TrimSpace(topic), "/")
	nodeName = strings.ReplaceAll(nodeName, "/", "_")
	nodeName = strings.ReplaceAll(nodeName, ".", "_")
	if nodeName == "" {
		return ""
	}
	return "/" + nodeName + "_node"
}

func sameROSName(left string, right string) bool {
	return normalizeROSName(left) == normalizeROSName(right)
}

func normalizeROSName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" || strings.HasPrefix(name, "/") {
		return name
	}
	return "/" + name
}

type twistPublisher struct {
	node *goroslib.Node
	pub  *goroslib.Publisher
}

func (r *Receiver) newROSPublisher() (*twistPublisher, error) {
	nodeName := uniqueROSNodeName(r.config.RosNodeName)
	node, err := goroslib.NewNode(goroslib.NodeConf{
		Name:          nodeName,
		MasterAddress: r.rosConfig.Address,
	})
	if err != nil {
		return nil, err
	}

	pub, err := goroslib.NewPublisher(goroslib.PublisherConf{
		Node:  node,
		Topic: r.config.RosTopic,
		Msg:   &rosgeometrymsgs.Twist{},
	})
	if err != nil {
		node.Close()
		return nil, err
	}

	return &twistPublisher{
		node: node,
		pub:  pub,
	}, nil
}

func uniqueROSNodeName(base string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		base = "inbound_cmd_vel_publisher"
	}

	host, err := os.Hostname()
	if err != nil || strings.TrimSpace(host) == "" {
		host = "host"
	}

	suffix := sanitizeROSNameToken(host) + "_" + strconv.Itoa(os.Getpid())
	return strings.TrimRight(base, "_") + "_" + suffix
}

func sanitizeROSNameToken(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "node"
	}

	var builder strings.Builder
	lastWasUnderscore := false
	for _, char := range value {
		valid := (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9')

		if valid {
			builder.WriteRune(char)
			lastWasUnderscore = false
			continue
		}

		if !lastWasUnderscore {
			builder.WriteByte('_')
			lastWasUnderscore = true
		}
	}

	token := strings.Trim(builder.String(), "_")
	if token == "" {
		return "node"
	}
	return token
}

func (p *twistPublisher) publish(msg *rosgeometrymsgs.Twist) {
	p.pub.Write(msg)
}

func (p *twistPublisher) close() {
	if p.pub != nil {
		p.pub.Close()
	}
	if p.node != nil {
		p.node.Close()
	}
}

func protoTwistToROS(msg *geometrymsgs.Twist) *rosgeometrymsgs.Twist {
	linear := msg.GetLinear()
	angular := msg.GetAngular()

	return &rosgeometrymsgs.Twist{
		Linear: rosgeometrymsgs.Vector3{
			X: linear.GetX(),
			Y: linear.GetY(),
			Z: linear.GetZ(),
		},
		Angular: rosgeometrymsgs.Vector3{
			X: angular.GetX(),
			Y: angular.GetY(),
			Z: angular.GetZ(),
		},
	}
}

func (r *Receiver) rejectMalformed(delivery amqp.Delivery) {
	if r.config.AutoAck {
		return
	}

	if err := delivery.Nack(false, false); err != nil {
		log.Printf(
			"Inbound receiver %s failed to reject malformed queue=%q delivery_tag=%d: %v",
			r.config.Name,
			r.config.Queue,
			delivery.DeliveryTag,
			err,
		)
	}
}
