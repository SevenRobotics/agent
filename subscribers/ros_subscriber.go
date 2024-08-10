package subscribers

import (
	"fmt"
	"github.com/bluenviron/goroslib/v2"
	"go_agent/config"
)

type rosSubscriber[T any] struct {
	sub         *goroslib.Subscriber
	conf        config.RosSubscriberConfig
	outChannel  chan<- T
	node        *goroslib.Node
	initiliased bool
}

func (s *rosSubscriber[T]) callback(msg *T) {

	if !s.initiliased {
		return
	}

	//how to handle this channel being closed?
	if msg != nil {
		s.outChannel <- *msg
	}
}

func NewRosSubscriber[T any](conf config.RosSubscriberConfig, node *goroslib.Node) *rosSubscriber[T] {
	sub := &rosSubscriber[T]{
		conf: conf,
		node: node,
	}

	return sub
}

func (s *rosSubscriber[T]) Initialise(out chan<- T) error {

	if s.node == nil {
		if s.conf.Node.Name == "" || s.conf.Node.Address == "" {
			return fmt.Errorf("Cannot initialise a subscriber without a node or a node config")
		}

		n, err := goroslib.NewNode(goroslib.NodeConf{
			Name:          s.conf.Node.Name,
			MasterAddress: s.conf.Node.Address,
		})

		if err != nil {
			return fmt.Errorf("Failed initiliasing node for subscriber %s: %v", s.conf.Name, err)
		}

		s.node = n
	}

	if s.outChannel == nil {
		s.outChannel = out
	}

	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     s.node,
		Topic:    s.conf.Topic,
		Callback: s.callback,
	})

	if err != nil {
		return fmt.Errorf("failed to initialise subscriber %s: %v", s.conf.Name, err)
	}

	s.sub = sub
	s.initiliased = true

	return nil
}
