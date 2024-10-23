package gopacket

import (
	"reflect"
	"sync"
)

var pools = make(map[reflect.Type]interface{})

var TcpFreeFunc func(layer Layer)

func newPool[R Recycler[T], T any]() Pool[T] {
	pp := &pool[T]{}
	pp.p = sync.Pool{
		New: func() any {
			r := reflect.New(reflect.TypeFor[R]().Elem())
			r.Interface().(R).SetPool(pp)
			return r.Interface().(R)
		},
	}
	return pp
}

type pool[T any] struct {
	p sync.Pool
}

func (p *pool[T]) Get() Recycler[T] {
	r := p.p.Get()
	t, ok := r.(*BaseRecycler[T])
	if ok {
		return t
	}

	return nil
}

func (p *pool[T]) Free(r Recycler[T]) {
	t := r.Get()
	*t = reflect.Zero(reflect.TypeFor[T]()).Interface().(T)
	p.p.Put(r)
}

type BaseRecycler[T any] struct {
	t T // t must be the first elem to align the memory bound, ensuring the addr of &t is the addr of &BaseRecycler
	p Pool[T]
}

func (t *BaseRecycler[T]) Handle(HandleFunc func(p *T) error) error {
	defer t.Free()
	return HandleFunc(&t.t)
}
func (t *BaseRecycler[T]) Get() *T {
	return &t.t
}

func (t *BaseRecycler[T]) Free() {
	t.p.Free(t)
}

func (t *BaseRecycler[T]) SetPool(p Pool[T]) {
	t.p = p
}

func Get[R Recycler[T], T any]() Recycler[T] {
	p := GetPool[R, T]()

	return p.(*pool[T]).Get()
}

func GetPoolByType(t reflect.Type) interface{} {
	return pools[t]
}

func GetPool[R Recycler[T], T any]() Pool[T] {
	p, ok := pools[reflect.TypeFor[T]()]
	if ok {
		return p.(*pool[T])
	}

	p = newPool[R, T]()
	pools[reflect.TypeFor[T]()] = p

	return p.(Pool[T])
}
