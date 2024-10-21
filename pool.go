package gopacket

import (
	"reflect"
	"sync"
)

var pools = make(map[reflect.Type]interface{})

func newPool[T any]() Pool[T] {
	pp := &pool[T]{}
	pp.p = sync.Pool{
		New: func() any {
			r := new(recycler[T])
			r.p = pp
			return r
		},
	}
	return pp
}

type pool[T any] struct {
	p sync.Pool
}

func (p *pool[T]) Get() Recycler[T] {
	r := p.p.Get()
	t, ok := r.(*recycler[T])
	if ok {
		return t
	}

	return nil
}

func (p *pool[T]) Free(r Recycler[T]) {
	//t := r.(*recycler[T])
	//t.t = reflect.Zero(reflect.TypeFor[T]()).Interface().(T)
	p.p.Put(r)
}

type recycler[T any] struct {
	p *pool[T]
	t T
}

func (t *recycler[T]) Handle(HandleFunc func(p *T) error) error {
	defer t.Free()
	return HandleFunc(&t.t)
}
func (t *recycler[T]) Get() *T {
	return &t.t
}

func (t *recycler[T]) Free() {
	t.p.Free(t)
}

func Get[R Recycler[T], T any]() Recycler[T] {
	p, ok := pools[reflect.TypeFor[T]()]
	if ok {
		return p.(*pool[T]).Get()
	}

	p = newPool[T]()
	pools[reflect.TypeFor[T]()] = p
	return p.(*pool[T]).Get()
}
