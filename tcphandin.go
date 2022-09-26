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

	serverChanSlice := make([]chan [2]int, 0, 5)
	channelSlice := make([]chan packet, 0, 5)
	threewaySlice := make([]chan [2]int, 0, 5)
	confirmationSlice := make([]chan int, 0, 5)

	for i := 0; i < 5; i++ {
		fromServerChan := make(chan [2]int)
		fromClientChan := make(chan [2]int)
		serverChanSlice = append(serverChanSlice, fromServerChan)

		fromServer := make(chan packet)
		fromClient := make(chan packet)
		channelSlice = append(channelSlice, fromServer)

		fromServerThreeway := make(chan [2]int)
		fromClientThreeway := make(chan [2]int)
		threewaySlice = append(threewaySlice, fromServerThreeway)

		fromServerConfirmation := make(chan int)
		fromClientConfirmation := make(chan int)
		confirmationSlice = append(confirmationSlice, fromServerConfirmation)
		go Client(i+1, fromClientChan, fromClientThreeway, fromClient, fromClientConfirmation)
		go MiddleWare(fromClient, fromServer, i+1, fromServerThreeway, fromClientThreeway, fromClientConfirmation, fromServerConfirmation)
	}

	go Server(serverChanSlice, channelSlice, threewaySlice, confirmationSlice)

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

	available := <-serverChan

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

				time.Sleep(time.Duration(packets[i].lifeTime))

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

func MiddleWare(fromClientMes chan packet, toServerMes chan packet, clientName int,
	fromServerThreeway chan [2]int, fromClientThreeway chan [2]int,
	toClientConfirmation chan int, fromServerConfirmation chan int) {

	for {
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(101)
		select {
		case ack := <-fromServerThreeway:
			fromClientThreeway <- ack

		case syn := <-fromClientThreeway:
			fromServerThreeway <- syn
		case p := <-fromClientMes:
			switch randNum {
			case 100, 99:
			default:
				toServerMes <- p
			}
		case conf := <-fromServerConfirmation:
			toClientConfirmation <- conf

		default:

		}

	}
}

func Server(serverChanSlice []chan [2]int, messageSlice []chan packet, threewaySlice []chan [2]int, confirmationSlice []chan int) {

	for {
		for i := 0; i < 5; i++ {
			serverChanSlice[i] <- [2]int{0, 0}
			recieved := <-serverChanSlice[i]
			if recieved[0] != 0 {
				clientNumber := recieved[0] - 1
				go RequestHandle(messageSlice[clientNumber], recieved[1], threewaySlice[clientNumber], confirmationSlice[clientNumber], clientNumber+1)
			}
		}
	}
}
