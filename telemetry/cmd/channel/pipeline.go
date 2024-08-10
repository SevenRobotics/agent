package channel

import (
	"fmt"
	"go_agent/config"
	"go_agent/publishers"
	"go_agent/publishers/rmq"
	"go_agent/subscribers"
	"sync"
)

type pipeline[S any, P any] struct {
	subscriber   subscribers.Subscriber[S]
	publisher    publishers.Publisher[P]
	bridge       *Bridge[S, P]
	in           chan S
	out          chan P
	errorChannel chan error
	done         chan int
	name         string
	conn         *rmq.RabbitMQMaster
}

func NewPipeline[S any, P any](conf config.RRPipelineConfig, msgConverter func(input S) (P, error), serializer func(msg P) ([]byte, error)) (*pipeline[S, P], error) {
	in := make(chan S)
	out := make(chan P)
	errCh := make(chan error)
	done := make(chan int)

	sub := subscribers.NewRosSubscriber[S](conf.SubConfig, nil)
	conn, err := rmq.NewRabbitMQ(conf.RMQConnConfig)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to RMQ Server %s:%v", conf.Name, err)
	}

	client, err := conn.NewClient(conf.Name)

	if err != nil {
		return nil, fmt.Errorf("Failed to create RMQ Client %s: %v", conf.Name, err)
	}

	pub, err := rmq.NewRMQPublisher[P](conf.RMQPubConfig, client)
	if err != nil {
		return nil, fmt.Errorf("Failed to create RMQ Publisher %s: %v", conf.Name, err)
	}

	pub.SetSerializer(serializer)

	b := NewBridge(conf.Name, in, out, done, errCh)
	b.SetConverter(msgConverter)
	return &pipeline[S, P]{
		subscriber:   sub,
		publisher:    pub,
		bridge:       b,
		in:           in,
		out:          out,
		errorChannel: errCh,
		done:         done,
		name:         conf.Name,
	}, nil
}

func (p *pipeline[S, P]) Shutdown() {}

func (p *pipeline[S, P]) Start(wg *sync.WaitGroup) {

	defer wg.Done()
	fmt.Printf("Starting pipeline %s\n", p.name)

	wg.Add(1)
	go p.bridge.Run(wg)
	p.subscriber.Initialise(p.in)
	wg.Add(1)
	go p.publisher.Run(p.out, p.done, p.errorChannel, wg)
}

func (p *pipeline[S, P]) GetErrorStream() chan error {
	return p.errorChannel
}

func (p *pipeline[S, P]) Name() string {
	return p.name
}
