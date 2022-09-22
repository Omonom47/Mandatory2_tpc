package main

import (
	"math/rand"
)

type packet struct {
	source      uint16
	destination uint16
	sequenceNum uint32
	checksum    uint16
	data        [8]byte
}

func MakePacket(data [8]byte, seqNum uint32, srcPort uint16, desPort uint16) packet {
	var p packet
	p.data = data
	p.source = srcPort
	p.destination = desPort
	p.sequenceNum = seqNum
	p.checksum = PacketHash(p)

	return p
}

func FragmentMessage(message string) []packet {
	asBytes := []byte(message)
	mesLen := len(asBytes)
	toReturn := make([]packet, mesLen/8)
	var seq uint32 = 0

	source := uint16(rand.Int31n(1024))
	destination := uint16(rand.Int31n(1024))
	var fragment [8]byte
	for i := 0; i < mesLen; i++ {
		fragment[i%8] = asBytes[i]
		if i+1%8 == 0 {
			toReturn = append(toReturn, MakePacket(fragment, seq, source, destination))
			seq++
		}
	}

	return toReturn
}

func PacketHash(p packet) uint16 {
	var h uint16
	for i := 0; i < len(p.data); i++ {
		h += uint16(p.data[i]) * IntPow(53, i) % 17959
	}
	return h
}

func IntPow(base uint16, exp int) uint16 {
	if exp == 0 {
		return 1
	}
	result := base
	for i := 2; i <= exp; i++ {
		result *= base
	}

	return result
}

func Main() {

	channel := make(chan packet)
}
