package gopacket

type Recycler[T any] interface {
	Handle(HandleFunc func(t *T) error) error
	Get() *T
	Free()
	SetPool(p Pool[T])
}

type Pool[T any] interface {
	Get() Recycler[T]
	Free(r Recycler[T])
}
