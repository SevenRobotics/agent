package subscribers

type Subscriber[T any] interface {
	callback(msg *T)
	Initialise() error
}
