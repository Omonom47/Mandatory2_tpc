package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
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
	mesLen := len(message)
	toReturn := make([]packet, 0, 4)
	var seq uint32 = 0

	source := uint16(rand.Int31n(1024))
	destination := uint16(rand.Int31n(1024))

	var fragment [8]byte
	for i := 0; i < mesLen; i++ {
		fragment[i%8] = message[i]
		if (i+1)%8 == 0 {
			toReturn = append(toReturn, MakePacket(fragment, seq, source, destination))
			seq++
		}
		checkNum := mesLen - i
		if checkNum < 8 {
			for x := 0; x < checkNum; x++ {
				fragment[x%8] = message[x+i]
			}
			toReturn = append(toReturn, MakePacket(fragment, seq, source, destination))
			break
		}
	}

	fmt.Println(mesLen, ":mesLen, ", len(toReturn), ":toReturn")

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
var finishvar int

func main() {

	channel := make(chan packet)
	ack := make(chan [2]int)
	confirmationChan := make(chan int)
	finishvar = 0

	finish.Add(2)
	go Client("client", ack, channel, confirmationChan)
	go Server(channel, ack, confirmationChan)

	for {
		if finishvar == 1 {
			break
		}
	}

}

func Server(packetChan chan packet, threewayChan chan [2]int, confChan chan int) {

	defer finish.Done()

	recieved := <-threewayChan
	randomSeq := rand.Int()
	var message string
	if recieved[0] == 1 {
		threewayChan <- [2]int{recieved[1] + 1, randomSeq}
	} else {
		time.Sleep(10)
	}
	recieved = <-threewayChan
	if recieved[0] == randomSeq+1 {
		//confirmation of recieving packet
		var p packet
		dataRecived := make([]packet, 4) //queue with packet that sort into right order
		p = <-packetChan
		confChan <- 1

		for i := 0; i < int(p.mesLen)-1; i++ {

			dataRecived = append(dataRecived, p)
			p = <-packetChan
			if i != int(p.mesLen)-1 {
				confChan <- 1
			}
		}

		sort.SliceStable(dataRecived, func(i, j int) bool {
			return dataRecived[i].sequenceNum < dataRecived[j].sequenceNum
		})

		for i := 0; i < int(p.mesLen); i++ {
			message += string(dataRecived[i].data[:])
			// fmt.Println(i, " -- ", dataRecived[i])
		}
		fmt.Println(message)
		fmt.Println(len(strings.Trim(message, " ")))
		time.Sleep(2)
		finishvar = 1
	}

}

func Client(name string, threewayChan chan [2]int, packetChan chan packet, confChan chan int) {

	defer finish.Done()

	senddata := 1 //rand.Int31n(2)
	if senddata == 1 {
		datasize := rand.Intn(100) + 1
		data := CreateRandomData(datasize)

		fmt.Println(data)
		fmt.Println(len(data))

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
				conf := <-confChan

				if conf != 1 {
					break
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
