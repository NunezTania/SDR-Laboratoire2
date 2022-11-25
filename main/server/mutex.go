package main

import (
	"fmt"
)

type Message struct {
	rType string
	time  Lamport
	id    int
}

func StartClock(clock *Lamport) {
	clock.NewLamport()
}

func NoteNewMessage(id int, msg Message, messages *[]Message, chanSC chan bool) {
	fmt.Println("Message received in process", id, " ", msg)
	(*messages)[msg.id] = msg
	if checkSCAvailable(id, messages) {
		chanSC <- true
	}
}

// todo verify is function is correct
func checkSCAvailable(id int, messages *[]Message) bool {
	for _, msg := range *messages {
		if msg.id != id {
			if &msg == nil || (msg.rType == "req" && msg.time.counterTime < (*messages)[id].time.counterTime && (*messages)[id].rType == "req") {
				return false
			}
		}
	}
	return true
}
