package rmq

import (
	"fmt"
	"go_agent/config"
	"go_agent/publishers"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

type rmqPublisher[P any] struct {
	name       string
	config     config.RMQClientConfig
	serializer func(msg P) ([]byte, error)
	client     RabbitClient
	ready      bool
}

func (r *rmqPublisher[P]) Send(msg P) error {
	wiremsg, err := r.Serialize(msg)
	if err != nil {
		return err
	}

	err = r.client.Send(r.config.Ctx, r.config.Exchange, r.config.RoutingKey, amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Persistent,
		Body:         []byte(wiremsg),
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *rmqPublisher[P]) Run(in <-chan P, done chan int, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		case msg, ok := <-in:
			if !ok {
				continue
			}
			err := r.Send(msg)
			if err != nil {
				errCh <- fmt.Errorf("publisher %s Send failed: %v", r.name, err)
			}
		}
	}
}

func (r *rmqPublisher[P]) SetTopic(topic string) publishers.Publisher[P] {
	r.config.Topic = topic
	return r
}

func (r *rmqPublisher[P]) Serialize(msg P) ([]byte, error) {
	if r.serializer == nil {
		return []byte{}, fmt.Errorf("Serializer not set")
	}
	return r.serializer(msg)
}

func (r *rmqPublisher[P]) SetAddress(addr string) publishers.Publisher[P] {
	r.config.Exchange = addr
	return r
}

func (r *rmqPublisher[P]) SetKey(key string) publishers.Publisher[P] {
	r.config.RoutingKey = key
	return r
}

func (r *rmqPublisher[P]) SetSerializer(s func(msg P) ([]byte, error)) publishers.Publisher[P] {
	r.serializer = s
	return r
}

func (r *rmqPublisher[P]) Configure() error {
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

func NewRMQPublisher[P any](conf config.RMQClientConfig, client RabbitClient) (*rmqPublisher[P], error) {
	p := &rmqPublisher[P]{
		config: conf,
		client: client,
	}

	err := p.Configure()

	if err != nil {
		return nil, err
	}

	return p, nil
}
