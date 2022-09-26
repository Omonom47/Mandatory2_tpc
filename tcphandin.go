package main

import (
	"fmt"
	"math/rand"
	"sort"
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
	data        byte
}

func MakePacket(data byte, seqNum uint32, srcPort uint16, desPort uint16, mesHash uint16) packet {
	var p packet
	p.data = data
	p.source = srcPort
	p.destination = desPort
	p.sequenceNum = seqNum
	p.checksum = mesHash
	p.timeStamp = time.Now()
	p.lifeTime = 5

	return p
}

func FragmentMessage(message string) []packet {
	mesLen := len(message)
	toReturn := make([]packet, 0, 4)
	var seq uint32 = 0

	mesHash := PacketHash(message)

	source := uint16(rand.Int31n(1024))
	destination := uint16(rand.Int31n(1024))

	for seq < uint32(mesLen) {
		p := MakePacket(message[seq], seq, source, destination, mesHash)
		p.mesLen = uint32(mesLen)
		toReturn = append(toReturn, p)
		seq++
	}

	return toReturn
}

func PacketHash(message string) uint16 {
	var h uint16
	for i := 0; i < len(message); i++ {
		h += uint16(message[i]) * IntPow(53, i) % 17959
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

var finishvar int

func main() {

	finishvar = 0
	go Server()

	for {
		if finishvar == 5 {
			break
		}
	}

}

func RequestHandle(packetChan chan packet, info int, threewayChan chan [2]int, confChan chan int, client int) {

	clientInfo := info
	randomSeq := rand.Int()

	threewayChan <- [2]int{clientInfo + 1, randomSeq}

	recieved := <-threewayChan
	if recieved[0] == randomSeq+1 {
		//confirmation of recieving packet
		var p packet
		dataRecived := make([]packet, 0, 4)
		p = <-packetChan
		dataRecived = append(dataRecived, p)
		confChan <- 1

		for i := 0; i < int(p.mesLen)-1; i++ {

			p = <-packetChan
			dataRecived = append(dataRecived, p)
			if i != int(p.mesLen)-1 {
				confChan <- 1
			}
		}

		sort.SliceStable(dataRecived, func(i, j int) bool {
			return dataRecived[i].sequenceNum < dataRecived[j].sequenceNum
		})

		var message string
		for i := 0; i < int(p.mesLen); i++ {
			message += string(dataRecived[i].data)
		}
		fmt.Println("Client ", client, message)
		finishvar++
	}

}

func Client(name int, serverChan chan [2]int, threewayChan chan [2]int, packetChan chan packet, confChan chan int) {

	datasize := rand.Intn(100) + 1
	data := CreateRandomData(datasize)

	//fmt.Println(data)
	//fmt.Println(len(data))

	available := <-serverChan

	//fmt.Println("appr is: ", approved)
	if available[1] == 0 {
		check := rand.Int()
		serverChan <- [2]int{name, check}
		time.Sleep(5)
		confirmation := <-threewayChan
		if confirmation[0] == check+1 {

			packets := FragmentMessage(data)
			seqNum := confirmation[1]
			threewayChan <- [2]int{seqNum + 1, check + 1}

			for i := 0; i < len(packets); i++ {

				packetChan <- packets[i]

				time.Sleep(2)

				select {
				case <-confChan:

				default:
					i--
				}

			}
		}
	} else {
		fmt.Println("Client", name, "server not available")
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

func MiddleWare(fromClient chan packet, toServer chan packet, clientName string, threewayChan chan [2]int, confChan chan int) {

	for {
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(101)
		p := <-fromClient

		switch randNum {
		case 100, 99:
		case 49, 51:
			corrupted := p.data
			corrupted = corrupted << 2
			corrupted = ^corrupted
			p.data = corrupted
			toServer <- p
		default:
			toServer <- p
		}

	}
}

func Server() {

	serverChanSlice := make([]chan [2]int, 0, 5)
	channelSlice := make([]chan packet, 0, 5)
	ackSlice := make([]chan [2]int, 0, 5)
	confirmationSlice := make([]chan int, 0, 5)

	for i := 0; i < 5; i++ {
		serverChan := make(chan [2]int)
		serverChanSlice = append(serverChanSlice, serverChan)
		channel := make(chan packet)
		channelSlice = append(channelSlice, channel)
		ack := make(chan [2]int)
		ackSlice = append(ackSlice, ack)
		confirmationChan := make(chan int)
		confirmationSlice = append(confirmationSlice, confirmationChan)
		go Client(i+1, serverChan, ack, channel, confirmationChan)
	}

	for {
		for i := 0; i < 5; i++ {
			serverChanSlice[i] <- [2]int{0, 0}
			recieved := <-serverChanSlice[i]
			if recieved[0] != 0 {
				clientNumber := recieved[0] - 1
				go RequestHandle(channelSlice[clientNumber], recieved[1], ackSlice[clientNumber], confirmationSlice[clientNumber], clientNumber)
			}
		}
	}
}
