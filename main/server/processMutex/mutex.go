package processMutex

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

func AskForSC(id int, clock *Lamport, messages *[]Message) {
	fmt.Println("I'm id = ", id, " and I ask for sc at time ", clock.counterTime)
	clock.Increment()
	sendRequests(clock, id)
	(*messages)[id] = Message{"req", *clock, id}
}

func FreeSC(id int, clock *Lamport, inSC *bool, messages *[]Message) {
	fmt.Println("Server ", id, " leaves sc at time", clock.counterTime)
	clock.Increment()
	sendReleases(clock, id)
	(*messages)[id] = Message{"rel", *clock, id}
	*inSC = false
}

func NoteNewMessage(message Message, index int, id int, inSC *bool, ChannelSc *chan string, messages *[]Message) {
	(*messages)[index] = message
	checkSCAvailable(id, inSC, messages)
	if *inSC {
		*ChannelSc <- "SC"
	}
}

// Check if the SC is available
// the SC is available if all the servers have sent : an ack, a rel, or a req with an bigger clock
// the SC is not available if one server has sent a req with a smaller clock or if a sever haven't answered
func checkSCAvailable(id int, inSC *bool, messages *[]Message) {
	if (*messages)[id].rType != "req" {
		*inSC = false
		return
	}
	if Config.NServ == 1 {
		*inSC = true
		return
	}
	for _, msg := range *messages {
		if msg.id != id {
			if msg.id == -1 || (msg.rType == "req" && msg.time.counterTime < (*messages)[id].time.counterTime) {
				*inSC = false
				return
			}
		}
	}
	*inSC = true
}
