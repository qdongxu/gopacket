package gopacket

import (
	"reflect"
	"sync"
)

var emptyEagerPacket = EagerPacket{}
var eagerPacketPool = sync.Pool{
	New: func() any {
		return &EagerPacket{}
	},
}

var emptyLazyPacket = LazyPacket{}
var lazyPacketPool = sync.Pool{
	New: func() any {
		return &LazyPacket{}
	},
}

func GetPacket() *EagerPacket {
	p := eagerPacketPool.Get().(*EagerPacket)
	p.Transport = LayerGetTCP()
	p.Network = LayerGetIPv4()
	return p
}

var LayerFreeFunc func(p Layer)

var LayerGetIPv4 func() NetworkLayer
var LayerGetIPv6 func() NetworkLayer
var LayerGetTCP func() TransportLayer
var LayerGetUDP func() TransportLayer

// Free put the packet back into the pool.
func Free(p any) {
	if reflect.ValueOf(p).IsNil() {
		return
	}

	switch p.(type) {
	case *EagerPacket:
		ep := p.(*EagerPacket)

		LayerFreeFunc(ep.Network)
		LayerFreeFunc(ep.Transport)

		*ep = emptyEagerPacket
		eagerPacketPool.Put(ep)
	case *LazyPacket:
		lp := p.(*LazyPacket)

		LayerFreeFunc(lp.Network)
		LayerFreeFunc(lp.Transport)

		//*lp = emptyLazyPacket
		lazyPacketPool.Put(lp)
	default:
		// ignore any others
	}
}
