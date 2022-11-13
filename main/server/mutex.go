package main

import "fmt"

var clock = Lamport{}
var inSC = false

type Message struct {
	rType string
	time  Lamport
	id    int
}

var msgArray []Message

func start() {
	clock.NewLamport()
	fmt.Println(clock.GetTime())
}

func askForSC() {
	clock.Increment()
	sendRequests(clock)
	msgArray[id] = Message{"req", clock, id}
}

func freeSC() {
	clock.Increment()
	sendReleases(clock)
	msgArray[id] = Message{"rel", clock, id}
	inSC = false
}

func NoteNewMessage(message Message, index int) {
	msgArray[index] = message
	checkSCAvailable()
}

func checkSCAvailable() {
	for _, msg := range msgArray {
		if msg.id != id {
			if &msg == nil || (msg.rType == "req" && msg.time.counterTime < clock.counterTime) {
				inSC = false
				return
			}
		}
	}
	inSC = false
}
