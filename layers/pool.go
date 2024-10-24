package layers

import (
	"github.com/google/gopacket"
	"unsafe"
)

func init() {
	gopacket.TCPFreeFunc = FreeTransportLayer
	gopacket.IPFreeFunc = FreeIPLayer
	gopacket.LinkFreeFunc = FreeLinkLayer
}
func FreeTransportLayer(d gopacket.TransportLayer) {
	var u interface{} = d
	var v = unsafe.Pointer(&u)
	var m = (*interface{})(v)
	if w, ok := (*m).(*gopacket.BaseRecycler[TCP]); ok {
		w.Free()
	} else if w, ok := (*m).(*gopacket.BaseRecycler[UDP]); ok {
		w.Free()
	} else if w, ok := (*m).(*gopacket.BaseRecycler[RUDP]); ok {
		w.Free()
	} else if w, ok := (*m).(*gopacket.BaseRecycler[SCTP]); ok {
		w.Free()
	} else if w, ok := (*m).(*gopacket.BaseRecycler[UDPLite]); ok {
		w.Free()
	}
}

func FreeIPLayer(d gopacket.NetworkLayer) {
	var u interface{} = d
	var v = unsafe.Pointer(&u)
	var m = (*interface{})(v)

	if w, ok := (*m).(*gopacket.BaseRecycler[IPv4]); ok {
		w.Free()
	} else if w, ok := (*m).(*gopacket.BaseRecycler[IPv6]); ok {
		w.Free()
	}
}

func FreeLinkLayer(d gopacket.LinkLayer) {
	var u interface{} = d
	var v = unsafe.Pointer(&u)
	var m = (*interface{})(v)

	if w, ok := (*m).(*gopacket.BaseRecycler[Ethernet]); ok {
		w.Free()
	} else if w, ok := (*m).(*gopacket.BaseRecycler[FDDI]); ok {
		w.Free()
	} else if w, ok := (*m).(*gopacket.BaseRecycler[LinuxSLL]); ok {
		w.Free()
	} else if w, ok := (*m).(*gopacket.BaseRecycler[PPP]); ok {
		w.Free()
	}
}
