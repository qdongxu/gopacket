package gopacket

import (
	"fmt"
	"testing"
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
		fmt.Println(o.I)
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
