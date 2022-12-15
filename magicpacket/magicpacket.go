package magicpacket

import (
	"net"
)

const (
	headerSize     = 6
	macRepetitions = 16
)

func New(mac net.HardwareAddr) []byte {
	mpSize := headerSize + (len(mac) * macRepetitions)
	var mp []byte = make([]byte, mpSize)

	for i := 0; i < headerSize; i++ {
		mp[i] = 0xFF
	}

	for i := headerSize; i < mpSize; i = i + len(mac) {
		for j := range mac {
			mp[i+j] = mac[j]
		}
	}

	return mp
}
