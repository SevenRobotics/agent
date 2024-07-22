package internal

type Topic[T any] struct {
	Name      string
	TopicType T
}
