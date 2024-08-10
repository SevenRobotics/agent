package publishers

import "sync"

type Publisher[P any] interface {
	Configure() error
	Send(P) error
	SetSerializer(func(msg P) ([]byte, error)) Publisher[P]
	SetTopic(string) Publisher[P]
	SetAddress(string) Publisher[P]
	SetKey(string) Publisher[P]
	Serialize(P) ([]byte, error)
	Run(<-chan P, chan int, chan error, *sync.WaitGroup)
}
