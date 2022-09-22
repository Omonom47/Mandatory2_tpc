package main

import (
	"math/rand"
	"time"
)

type packet struct {
	source      uint16
	destination uint16
	sequenceNum uint32
	checksum    uint16
	timeStamp   time.Time
	lifeTime    uint8
	data        [8]byte
}

func MakePacket(data [8]byte, seqNum uint32, srcPort uint16, desPort uint16) packet {
	var p packet
	p.data = data
	p.source = srcPort
	p.destination = desPort
	p.sequenceNum = seqNum
	p.checksum = PacketHash(p)
	p.timeStamp = time.Now()
	p.lifeTime = 5

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

func Client(name string, comChan chan int, comChan2 chan int, packetChan chan packet, confChan chan int) {

	//createpacketfunc
	//3wayhandshakefunc
	for true {
		senddata := rand.Int31n(2)
		if senddata == 1 {
			datasize := rand.Int()
			data := CreateRandomData(datasize)
			comChan <- 1
			time.Sleep(5)
			if <-comChan2 == 2 {
				comChan <- 3
				FragmentMessage(data)
				//packetChan <- packet
				//for loop{
				if <-confChan == 1 {
					//packetChan <- packet
				}
			}
		}
	}

	/*if 3wayhandshake=accepted {
		packpacketChan <- //packet
	}*/

}

func CreateRandomData(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
