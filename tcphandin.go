package main	

import (
	"fmt"
	"math/rand"
)

type packet struct {
	source      uint16
	destination uint16
	sequenceNum uint32
	checksum    uint16
	data        [4]byte
}

func PacketHash(p packet) uint16 {
	var h uint16
	for i := 0; i < 4; i++ {
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

func Host(name string, comChan chan int, packetChan chan packet) {
	datasize := ran
	CreateRandomData()

}

func CreateRandomData(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

}
