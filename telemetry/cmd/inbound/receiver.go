package inbound

import (
	"context"
	"fmt"
	"go_agent/config"
	"go_agent/publishers/rmq"
	geometrymsgs "go_agent/telemetry/genproto/ros/geometry_msgs"
	"log"
	"sync"
	"time"

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
	config    config.RMQInboundReceiverConfig
}

func NewReceiver(rmqConfig config.RMQConfig, receiverConfig config.RMQInboundReceiverConfig) (*Receiver, error) {
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
	if receiverConfig.Consumer == "" {
		receiverConfig.Consumer = receiverClientPrefix + receiverConfig.Name
	}

	return &Receiver{
		rmqConfig: rmqConfig,
		config:    receiverConfig,
	}, nil
}

func StartEnabled(ctx context.Context, rmqConfig config.RMQConfig, inboundConfig config.RMQInboundConfig, wg *sync.WaitGroup) error {
	started := 0
	for _, receiverConfig := range inboundConfig.Receivers {
		if !receiverConfig.Enabled {
			continue
		}

		receiver, err := NewReceiver(rmqConfig, receiverConfig)
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
	conn, err := rmq.NewRabbitMQ(r.rmqConfig)
	if err != nil {
		return fmt.Errorf("connect RabbitMQ: %w", err)
	}

	client, err := conn.NewClient(receiverClientPrefix + r.config.Name)
	if err != nil {
		return fmt.Errorf("create RabbitMQ client: %w", err)
	}

	log.Printf(
		"Starting inbound receiver %s: exchange=%q queue=%q consumer=%q message_type=%q auto_ack=%t",
		r.config.Name,
		r.config.Exchange,
		r.config.Queue,
		r.config.Consumer,
		r.config.MessageType,
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
			r.handleDelivery(delivery)
		}
	}
}

func (r *Receiver) handleDelivery(delivery amqp.Delivery) {
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

	// Future ROS publisher integration:
	// 1. Convert proto geometry_msgs.Twist to telemetry/gengo/ros/geometry_msgs.Twist.
	// 2. Publish it to the configured ROS topic.
	// 3. Keep the Ack below after a successful ROS publish.

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
