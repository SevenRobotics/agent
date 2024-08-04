package publishers

type Publisher interface {
	Send(interface{}) error
	SetSerializer(func(msg interface{}) ([]byte, error)) *Publisher
	SetTopic(string) *Publisher
	SetAddress(string) *Publisher
	SetKey(string) *Publisher
	Serialize(msg interface{}) ([]byte, error)
}
