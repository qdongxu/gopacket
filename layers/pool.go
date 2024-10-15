package layers

import (
	"github.com/google/gopacket"
	"reflect"
	"sync"
)

func init() {
	gopacket.LayerFreeFunc = Free
	gopacket.LayerGetIPv4 = GetIpv4
	gopacket.LayerGetIPv4 = GetIpv6
	gopacket.LayerGetTCP = GetTCP
	gopacket.LayerGetUDP = GetUDP
}

func New() gopacket.Packet {
	p := &gopacket.EagerPacket{}
	p.Network = &IPv4{}
	p.Transport = &TCP{}
	return p
}

var emptyTCP = TCP{}
var tcpPool = sync.Pool{
	New: func() any {
		return &TCP{}
	},
}

var emptyUDP = UDP{}
var udpPool = sync.Pool{
	New: func() any {
		return &UDP{}
	},
}

var emptyIpv4 = IPv4{}
var ipv4Pool = sync.Pool{
	New: func() any {
		return &IPv4{}
	},
}

var emptyIpv6 = IPv6{}
var ipv6Pool = sync.Pool{
	New: func() any {
		return &IPv6{}
	},
}

func GetIpv4() gopacket.NetworkLayer {
	i := ipv4Pool.Get().(*IPv4)
	return i
}

func GetIpv6() gopacket.NetworkLayer {
	i := ipv4Pool.Get().(*IPv6)
	return i
}

func GetTCP() gopacket.TransportLayer {
	i := tcpPool.Get().(*TCP)
	return i
}

func GetUDP() gopacket.TransportLayer {
	i := tcpPool.Get().(*UDP)
	return i
}

// Free put the packet back into the pool.
func Free(p gopacket.Layer) {
	if p == nil {
		return
	}

	if reflect.ValueOf(p).IsNil() {
		return
	}

	switch p.(type) {
	case *TCP:
		t := p.(*TCP)
		*t = emptyTCP
		tcpPool.Put(t)
	case *UDP:
		u := p.(*UDP)
		*u = emptyUDP
		udpPool.Put(u)
	case *IPv4:
		i := p.(*IPv4)
		*i = emptyIpv4
		ipv4Pool.Put(i)
	case *IPv6:
		i := p.(*IPv6)
		*i = emptyIpv6
		ipv6Pool.Put(i)
	default:
		// ignore any others
	}
}
