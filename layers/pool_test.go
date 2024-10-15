package layers

import (
	"testing"

	"github.com/google/gopacket"
)

func BenchmarkFreePacket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]*gopacket.EagerPacket, 0, 100)
		for j := 10; j < 100; j++ {
			for i := 0; i < 100; i++ {
				s = append(s, gopacket.GetPacket())
			}
			for i := 0; i < 100; i++ {
				gopacket.Free(s[i])
			}
		}
		s = s[:0]
	}
}

func BenchmarkGC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]gopacket.Packet, 0, 100)
		for j := 10; j < 100; j++ {
			for i := 0; i < 100; i++ {
				p := New()
				s = append(s, p)
			}
			for i := 0; i < 100; i++ {
				//Free(s[i])
			}
			s = s[:0]
		}
	}
}
