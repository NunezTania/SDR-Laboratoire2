package processMutex

import (
	"fmt"
)

type Message struct {
	rType string
	time  Lamport
	id    int
}

var mapMessage map[int][]Message

func StartClock(clock *Lamport) {
	clock.NewLamport()
	mapMessage = make(map[int][]Message)
	for i := 0; i < Config.NServ; i++ {
		mapMessage[i] = make([]Message, Config.NServ)
		for j := 0; j < Config.NServ; j++ {
			mapMessage[i][j].id = -1
		}
	}
}

func AskForSC(id int, clock *Lamport) {
	fmt.Println("I'm id = ", id, " and I ask for sc at time ", clock.counterTime)
	clock.Increment()
	sendRequests(clock, id)
	mapMessage[id][id] = Message{"req", *clock, id}
}

func FreeSC(id int, clock *Lamport, inSC *bool) {
	fmt.Println("Server ", id, " leaves sc at time", clock.counterTime)
	clock.Increment()
	sendReleases(clock, id)
	mapMessage[id][id] = Message{"rel", *clock, id}
	*inSC = false
}

func NoteNewMessage(message Message, index int, id int, inSC *bool, ChannelSc *chan string) {
	mapMessage[id][index] = message
	checkSCAvailable(id, inSC)
	if *inSC {
		*ChannelSc <- "SC"
	}
}

// Check if the SC is available
// the SC is available if all the servers have sent : an ack, a rel, or a req with an bigger clock
// the SC is not available if one server has sent a req with a smaller clock or if a sever haven't answered
func checkSCAvailable(id int, inSC *bool) {
	if mapMessage[id][id].rType != "req" {
		*inSC = false
		return
	}
	if Config.NServ == 1 {
		*inSC = true
		return
	}
	for _, msg := range mapMessage[id] {
		if msg.id != id {
			if msg.id == -1 || (msg.rType == "req" && msg.time.counterTime < mapMessage[id][id].time.counterTime) {
				*inSC = false
				return
			}
		}
	}
	*inSC = true
}
