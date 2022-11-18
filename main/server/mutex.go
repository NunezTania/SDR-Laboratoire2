package main

import (
	"fmt"
)

var clock = Lamport{}
var inSC = false
var ChannelSc = make(chan string)

type Message struct {
	rType string
	time  Lamport
	id    int
}

var msgArray []Message

func StartClock() {
	clock.NewLamport()
	msgArray = make([]Message, ReadConfigFile().NServ)
}

func AskForSC(id int) {
	fmt.Println("I'm id = ", id, " and I ask for sc at time ", clock.counterTime)
	clock.Increment()
	sendRequests(clock, id)
	msgArray[id] = Message{"req", clock, id}
}

func FreeSC(id int) {
	fmt.Println("I'm id = ", id, " and I free sc at the time", clock.counterTime)
	clock.Increment()
	sendReleases(clock, id)
	msgArray[id] = Message{"rel", clock, id}
	inSC = false
}

func NoteNewMessage(message Message, index int, id int) {
	msgArray[index] = message
	checkSCAvailable(id)
	fmt.Println("I'm id = ", id, " and the checkSCAvailable is ", inSC)
	if inSC {
		ChannelSc <- "SC"
	}
	fmt.Println("The msgArray is ", msgArray)
}

// Check if the SC is available
// the SC is available if all the servers have sent : an ack, a rel, or a req with an bigger clock
// the SC is not available if one server has sent a req with a smaller clock or if a sever haven't answered
func checkSCAvailable(id int) {
	fmt.Println("I'm id = ", id, " and I check if the SC is available at the time", clock.counterTime)
	for i := 0; i < ReadConfigFile().NServ; i++ {
		if i != id {
			if msgArray[i].rType == "req" {
				if msgArray[i].time.counterTime < clock.counterTime {
					inSC = true
				} else {
					inSC = false
					return
				}
			} else if msgArray[i].rType == "ack" || msgArray[i].rType == "rel" {
				inSC = true
			} else {
				inSC = false
				return
			}
		}
	}
}
