package rmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go_agent/config"
)

type rmqSubscriber[S any] struct {
	name   string
	config config.RMQClientConfig
	client RabbitClient
	ready  bool
}

func (r *rmqSubscriber[S]) Configure() error {
	if r.client == nil {
		return fmt.Errorf("Client not set")
	}

	if err := r.client.NewExchangeDeclare(r.config.Exchange, "topic", r.config.Durable, r.config.Autodelete); err != nil {
		return err
	}

	if err := r.client.NewQueueDeclare(r.config.Topic, r.config.Durable, r.config.Autodelete); err != nil {
		return err
	}

	if err := r.client.CreateBinding(r.config.Topic, r.config.RoutingKey, r.config.Exchange); err != nil {
		return err
	}

	r.ready = true
	return nil
}

func (r *rmqSubscriber[S]) Receive() (<-chan amqp.Delivery, error) {
	return r.client.Receive(r.config.Ctx, r.config.Exchange, "", true)
}

func NewRMQSubscriber[S any](conf config.RMQClientConfig, client RabbitClient) (*rmqSubscriber[S], error) {
	s := &rmqSubscriber[S]{
		config: conf,
		client: client,
	}

	err := s.Configure()

	if err != nil {
		return nil, err
	}

	return s, nil
}
