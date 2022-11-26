package processMutex

import (
	"fmt"
)

type Message struct {
	rType string
	time  Lamport
	id    int
}

var msgArray []Message

func StartClock(clock *Lamport) {
	clock.NewLamport()
	msgArray = make([]Message, Config.NServ)
	for id := range msgArray {
		msgArray[id].id = -1
	}
}

func AskForSC(id int, clock *Lamport) {
	fmt.Println("I'm id = ", id, " and I ask for sc at time ", clock.counterTime)
	clock.Increment()
	sendRequests(clock, id)
	msgArray[id] = Message{"req", *clock, id}
}

func FreeSC(id int, clock *Lamport, inSC *bool) {
	fmt.Println("I'm id = ", id, " and I free sc at the time", clock.counterTime)
	clock.Increment()
	sendReleases(clock, id)
	msgArray[id] = Message{"rel", *clock, id}
	*inSC = false
}

func NoteNewMessage(message Message, index int, id int, inSC *bool, ChannelSc *chan string, clock *Lamport) {
	msgArray[index] = message
	checkSCAvailable(id, inSC, clock)
	fmt.Println("I'm id = ", id, " and the checkSCAvailable is ", *inSC)
	if *inSC {
		*ChannelSc <- "SC"
	}
	fmt.Println("The msgArray is ", msgArray)
}

// Check if the SC is available
// the SC is available if all the servers have sent : an ack, a rel, or a req with an bigger clock
// the SC is not available if one server has sent a req with a smaller clock or if a sever haven't answered
func checkSCAvailable(id int, inSC *bool, clock *Lamport) {
	fmt.Println("I'm id = ", id, " and I check if the SC is available at the time", clock.counterTime)
	fmt.Println("The msgArray is ", msgArray)
	if msgArray[id].rType != "req" {
		*inSC = false
		return
	}
	for _, msg := range msgArray {
		if msg.id != id {
			if msg.id == -1 || (msg.rType == "req" && msg.time.counterTime < msgArray[id].time.counterTime) {
				*inSC = false
				return
			}
		}
	}
	*inSC = true
}
