package main

import "math/rand"

type packet struct {
	source      uint16
	destination uint16
	sequenceNum uint32
	checksum    uint16
	data        [4]byte
}

type Pair[T, U any] struct {
	First  T
	Second U
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

func Host(name string, syn chan Pair[int, int], ack chan int, packetchan chan packet) {

}

func threeWayHandshakeClient(syncChan [2]chan int, ackChan [2]chan int) {
	initSeq := rand.Intn(2-1) + 1 //number between 1-2
	syncChan[0] <- 1
	syncChan[1] <- initSeq

	if 1 == 1 {

	} else {
		//try again
	}
}

func threeWayHandshakeServer(syncChan [2]chan int, ackChan [2]chan int) {
	initSeq := rand.Intn(2-1) + 1 //number between 1-2
	seqRecived := <-syncChan[1]
	if initSeq == seqRecived {
		ackChan[0] <- 1 //can establish contact
		ackChan[1] <- initSeq + 1
	} else {
		ackChan[0] <- 0 //cannot establish contact
	}
}

func Main() {
	channel := make(chan packet)
	tup := tuple.New2(5, "hi!")
}

//server kører herinde
