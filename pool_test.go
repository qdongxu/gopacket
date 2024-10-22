package gopacket

import (
	"fmt"
	"testing"
	"unsafe"
)

type Object interface {
	Display()
}

type Int struct {
	I int
}

func (i *Int) Display() {
}

type Float struct {
	f float64
}

func (f *Float) Display() {
}

func TestPool(t *testing.T) {

	i := Get[*recycler[Int], Int]()
	i.Get().I = 13
	err := i.Handle(func(o *Int) error {
		fmt.Println(o.I)
		return nil
	})

	i.Free()

	i = Get[*recycler[Int], Int]()
	err = i.Handle(func(o *Int) error {
		if i.Get().I != 0 {
			t.Errorf("instance was not reset to zero value")
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

}

func TestPoolMap(t *testing.T) {
	i := Get[*recycler[Int], Int]()
	i.Get().I = 13
	i.Handle(func(t *Int) error {
		fmt.Println(t.I)
		return nil
	})
	i.Free()

}

func TestPointer(t *testing.T) {
	r := &recycler[Int]{}
	r.Get().I = 10
	i := &r.t

	r2 := (*recycler[Int])(unsafe.Pointer(i))

	i2 := &r2.t

	if i.I != i2.I {
		t.Errorf("addr changed, the value does not equal: i=%d, i2=%d", i.I, i2.I)
	}

}
