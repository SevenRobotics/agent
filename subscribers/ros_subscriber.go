package subscribers

import (
	"github.com/bluenviron/goroslib/v2"
)

type rosSubscriber[T any] struct {
	sub         *goroslib.Subscriber
	conf        goroslib.SubscriberConf
	outChannel  chan T
	initiliased bool
}

func (s *rosSubscriber[T]) callback(msg *T) {
	if !s.initiliased {
		return
	}

	if msg != nil {
		s.outChannel <- *msg
	}
}

func NewRosSubscriber[T any](topicName string, node *goroslib.Node, out chan T) *rosSubscriber[T] {
	sub := &rosSubscriber[T]{
		outChannel: out,
		conf: goroslib.SubscriberConf{
			Node:  node,
			Topic: topicName,
		},
	}

	return sub
}

func (s *rosSubscriber[T]) Initialise() error {

	s.conf.Callback = s.callback

	sub, err := goroslib.NewSubscriber(s.conf)

	if err != nil {
		return err
	}

	s.sub = sub
	s.initiliased = true

	return nil
}
