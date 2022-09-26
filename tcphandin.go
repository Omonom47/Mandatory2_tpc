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

var finish sync.WaitGroup
var finishvar int

func main() {

	channel := make(chan packet)
	ack := make(chan [2]int)
	//connectPosChan := make(chan int)
	confirmationChan := make(chan int)
	finishvar = 0

	finish.Add(2)
	go Client("client" /*connectPosChan,*/, ack, channel, confirmationChan)
	//go Client("client 2", connectPosChan, ack, channel, confirmationChan)
	//go Client("client 3", connectPosChan, ack, channel, confirmationChan)
	go Server(channel /*connectPosChan,*/, ack, confirmationChan)

	for {
		if finishvar == 1 {
			break
		}
	}

}

func Server(packetChan chan packet /*conApprChan chan int,*/, threewayChan chan [2]int, confChan chan int) {

	defer finish.Done()
	//available := 1
	recieved := <-threewayChan
	randomSeq := rand.Int()
	//conApprChan <- available
	if recieved[0] == 1 {
		threewayChan <- [2]int{recieved[1] + 1, randomSeq}
		//conApprChan <- available
		//available = 0
	} else {
		time.Sleep(10)
	}
	recieved = <-threewayChan
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
			fmt.Println(p.sequenceNum)
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
		fmt.Println(message)
		//available = 1
		time.Sleep(2)
		finishvar = 1
	}

}

func Client(name string, threewayChan chan [2]int, packetChan chan packet, confChan chan int) {

	defer finish.Done()
	time.Sleep(1)

	senddata := 1 //rand.Int31n(4)
	if senddata == 1 {
		datasize := rand.Intn(100) + 1
		data := CreateRandomData(datasize)

		fmt.Println(data)
		fmt.Println(len(data))

		/*approved := <-conApprChan

		fmt.Println("appr is: ", approved)
		if approved == 1 {*/
		check := rand.Int()
		threewayChan <- [2]int{1, check}
		time.Sleep(5)
		confirmation := <-threewayChan
		if confirmation[0] == check+1 {

			packets := FragmentMessage(data)
			seqNum := confirmation[1]
			threewayChan <- [2]int{seqNum + 1, check + 1}

			rand.Seed(time.Now().UnixNano())
			randInterval := rand.Perm(len(packets)) //random interval til at sende packets
			for i := 0; i < len(packets); i++ {

				packetChan <- packets[randInterval[i]] //packets bliver sendt i random order

				time.Sleep(2)
				conf := <-confChan

				if conf != 1 {
					break
				}
			}
		}
		/*} else {
			fmt.Println("server not available")
		}*/
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
