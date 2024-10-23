package gopacket

import (
	"fmt"
	"reflect"
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

	i := Get[*BaseRecycler[Int], Int]()
	i.Get().I = 13
	err := i.Handle(func(o *Int) error {
		fmt.Println(o.I)
		return nil
	})

	i.Free()

	i = Get[*BaseRecycler[Int], Int]()
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
	i := Get[*BaseRecycler[Int], Int]()
	i.Get().I = 13
	i.Handle(func(t *Int) error {
		fmt.Println(t.I)
		return nil
	})
	i.Free()

}

func TestPointer(t *testing.T) {
	// r := &BaseRecycler[Int]{}
	r := Get[*BaseRecycler[Int], Int]()
	r.Get().I = 10
	i := r.Get()
	fmt.Println(r, i)

	r2 := (*BaseRecycler[Int])(unsafe.Pointer(i))

	i2 := &r2.t

	if i.I != i2.I {
		t.Errorf("addr changed, the value does not equal: i=%d, i2=%d", i.I, i2.I)
	}

	func(i unsafe.Pointer) {
		v := (*Int)(i)
		fmt.Println(v.I)
		x := (*BaseRecycler[Int])(i)
		fmt.Println(x.Get().I)
	}(unsafe.Pointer(r.Get()))

	// https://stackoverflow.com/questions/56554574/interface-convert-to-unsafe-pointer-problem
	func(i Object) {
		var l interface{} = i
		k := unsafe.Pointer(&l)
		m := (*interface{})(k)
		v := (*m).(*Int)
		fmt.Println(v.I)
		w := unsafe.Pointer(v)
		x := (*BaseRecycler[Int])(w)
		fmt.Println(x.Get().I)
	}(r.Get())

	func() {
		pool := GetPoolByType(reflect.TypeFor[Int]())
		var u = unsafe.Pointer(r.Get())
		var y = (*BaseRecycler[Int])(u)
		pool.(Pool[Int]).Free(y)
	}()

}
