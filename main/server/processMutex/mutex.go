// The mutex file is used to managed the requests and the releases of the critical section
// It's role is to create and read the messages exchanged between the servers and to stock them in a map
// The map is used to check if the SC is available or not
package processMutex

import (
	"fmt"
)

type Message struct {
	rType string
	time  Lamport
	id    int
}

// StartClock is used to start the clock
func StartClock(clock *Lamport) {
	clock.NewLamport()
}

// AskForSC is used to ask for the SC, it starts by sending a request and then wait until the SC is available
func AskForSC(id int, clock *Lamport, messages *[]Message) {
	fmt.Println("I'm id = ", id, " and I ask for sc at time ", clock.counterTime)
	clock.Increment()
	sendRequests(clock, id)
	(*messages)[id] = Message{"req", *clock, id}
}

// FreeSC is used to free the SC, it sends a release to all the servers
func FreeSC(id int, clock *Lamport, inSC *bool, messages *[]Message) {
	fmt.Println("Server ", id, " leaves sc at time", clock.counterTime)
	clock.Increment()
	sendReleases(clock, id)
	(*messages)[id] = Message{"rel", *clock, id}
	*inSC = false
}

// NoteNewMessage is used to update the array of messages. It's called when a new message is received
func NoteNewMessage(message Message, index int, id int, inSC *bool, ChannelSc *chan string, messages *[]Message) {
	(*messages)[index] = message
	checkSCAvailable(id, inSC, messages)
	if *inSC {
		*ChannelSc <- "SC"
	}
}

// checkSCAvailable check if the SC is available
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
