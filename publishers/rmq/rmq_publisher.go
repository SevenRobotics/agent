package rmq

import (
	"context"
	"fmt"
	"go_agent/publishers/rmq/internal"

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type ClientConfig struct {
	exchange   string
	topic      string
	routingKey string
	durable    bool
	autodelete bool
	ctx        context.Context
}

type rmqPublisher struct {
	name       string
	config     ClientConfig
	serializer func(msg proto.Message) ([]byte, error)
	client     internal.RabbitClient
	ready      bool
}

func (r *rmqPublisher) Send(msg proto.Message) error {
	wiremsg, err := r.Serialize(msg)
	if err != nil {
		return err
	}

	err = r.client.Send(r.config.ctx, r.config.exchange, r.config.routingKey, amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Persistent,
		Body:         []byte(wiremsg),
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *rmqPublisher) SetTopic(topic string) *rmqPublisher {
	r.config.topic = topic
	return r
}

func (r *rmqPublisher) Serialize(msg proto.Message) ([]byte, error) {
	if r.serializer == nil {
		return []byte{}, fmt.Errorf("Serializer not set")
	}
	return r.serializer(msg)
}

func (r *rmqPublisher) SetAddress(addr string) *rmqPublisher {
	r.config.exchange = addr
	return r
}

func (r *rmqPublisher) SetKey(key string) *rmqPublisher {
	r.config.routingKey = key
	return r
}

func (r *rmqPublisher) SetSerializer(s func(msg proto.Message) ([]byte, error)) *rmqPublisher {
	r.serializer = s
	return r
}

func (r *rmqPublisher) Configure() error {
	if r.client == nil {
		return fmt.Errorf("Client not set")
	}

	if err := r.client.NewExchangeDeclare(r.config.exchange, "topic", true, false); err != nil {
		return err
	}

	if err := r.client.NewQueueDeclare(r.config.topic, true, false); err != nil {
		return err
	}

	if err := r.client.CreateBinding(r.config.topic, r.config.routingKey, r.config.exchange); err != nil {
		return err
	}

	r.ready = true
	return nil
}

func NewRMQPublisher(ctx context.Context, exchange, routingKey, topic string, client internal.RabbitClient) (*rmqPublisher, error) {
	p := &rmqPublisher{
		config: ClientConfig{
			ctx:        ctx,
			exchange:   exchange,
			routingKey: routingKey,
			topic:      topic,
		},
		client: client,
	}

	err := p.Configure()

	if err != nil {
		return nil, err
	}

	return p, nil
}
