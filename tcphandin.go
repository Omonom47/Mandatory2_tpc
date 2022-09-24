package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type packet struct {
	source      uint16
	destination uint16
	sequenceNum uint32
	checksum    uint16
	timeStamp   time.Time
	lifeTime    uint8
	mesLen      uint32
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

	for i := 0; i < len(toReturn); i++ {
		toReturn[i].mesLen = uint32(len(toReturn))
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

var finish sync.WaitGroup

func main() {

	channel := make(chan packet)
	ack := make(chan [2]int)
	confirmationChan := make(chan int)

	finish.Add(2)

	go Client(CreateRandomData(rand.Intn(10-1)+1), ack, channel, confirmationChan)
	go Server(channel, ack, confirmationChan)

}

func Server(packetChan chan packet, threewayChan chan [2]int, confChan chan int) {

	defer finish.Done()
	fmt.Print("HER")
	// initSeq := rand.Intn(2-1) + 1 //number between 1-2
	recieved := <-threewayChan
	randomSeq := rand.Int()
	var message string
	if recieved[0] == 1 {
		threewayChan <- [2]int{recieved[1] + 1, randomSeq}
	} else {
		time.Sleep(10)
	}
	// fmt.Print(recieved[0], " ----- ", randomSeq)
	if recieved[0] == randomSeq+1 {
		fmt.Print("her")
		//confirmation of recieving packet
		p := <-packetChan
		dataRecived := make([]packet, p.mesLen) //queue with packet that sort into right order

		for i := 0; i < int(p.mesLen)-1; i++ {
			p = <-packetChan
			dataRecived = append(dataRecived, p)
			confChan <- 1
		}
		confChan <- 0

		sort.SliceStable(dataRecived, func(i, j int) bool {
			return dataRecived[i].sequenceNum < dataRecived[j].sequenceNum
		})

		for i := 0; i < int(p.mesLen); i++ {
			for j := 0; j < 8; j++ {
				message += string(dataRecived[i].data[j])
			}
		}
		fmt.Println(message)
		//end of message
		//need to have the real fin?
		time.Sleep(10)
		threewayChan <- [2]int{0, 0}
	}
	time.Sleep(10)

}

func Client(name string, threewayChan chan [2]int, packetChan chan packet, confChan chan int) {

	defer finish.Done()
	//3wayhandshakefunc
	for true {
		senddata := rand.Int31n(2)
		if senddata == 1 {
			datasize := rand.Intn(100-1) + 1
			data := CreateRandomData(datasize)
			check := rand.Int()
			threewayChan <- [2]int{1, check}
			time.Sleep(5)
			confirmation := <-threewayChan
			if confirmation[0] == check+1 {
				packets := FragmentMessage(data)
				seqNum := confirmation[1]
				threewayChan <- [2]int{seqNum + 1, check + 1}
				for i := 0; i < len(packets); i++ {
					packetChan <- packets[i]
					time.Sleep(2)
					if <-confChan == 1 {
					} else {
						break
					}
				}
			}
		}
	}

}

func CreateRandomData(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
