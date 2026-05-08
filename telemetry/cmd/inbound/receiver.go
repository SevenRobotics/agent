package inbound

import (
	"context"
	"fmt"
	"go_agent/config"
	"go_agent/publishers/rmq"
	rosgeometrymsgs "go_agent/telemetry/gengo/ros/geometry_msgs"
	geometrymsgs "go_agent/telemetry/genproto/ros/geometry_msgs"
	"log"
	"sync"
	"time"

	"github.com/bluenviron/goroslib/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

const (
	twistMessageType     = "geometry_msgs/Twist"
	consumerRetryDelay   = 5 * time.Second
	receiverClientPrefix = "inbound-"
)

type Receiver struct {
	rmqConfig config.RMQConfig
	rosConfig config.RosNodeConfig
	config    config.RMQInboundReceiverConfig
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

	deliveries, err := client.Receive(ctx, r.config.Queue, r.config.Consumer, r.config.AutoAck)
	if err != nil {
		return fmt.Errorf("consume queue %s: %w", r.config.Queue, err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
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

	printTwist(r.config.Queue, delivery, msg)

	rosMsg := protoTwistToROS(msg)
	rosPub.publish(rosMsg)
	log.Printf(
		"Published %s from queue=%q delivery_tag=%d to ROS topic=%q",
		twistMessageType,
		r.config.Queue,
		delivery.DeliveryTag,
		r.config.RosTopic,
	)

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

type twistPublisher struct {
	node *goroslib.Node
	pub  *goroslib.Publisher
}

func (r *Receiver) newROSPublisher() (*twistPublisher, error) {
	node, err := goroslib.NewNode(goroslib.NodeConf{
		Name:          r.config.RosNodeName,
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

func printTwist(queue string, delivery amqp.Delivery, msg *geometrymsgs.Twist) {
	linear := msg.GetLinear()
	angular := msg.GetAngular()

	log.Printf(
		"Received %s queue=%q routing_key=%q delivery_tag=%d content_type=%q linear=(x=%f y=%f z=%f) angular=(x=%f y=%f z=%f)",
		twistMessageType,
		queue,
		delivery.RoutingKey,
		delivery.DeliveryTag,
		delivery.ContentType,
		linear.GetX(),
		linear.GetY(),
		linear.GetZ(),
		angular.GetX(),
		angular.GetY(),
		angular.GetZ(),
	)
}
