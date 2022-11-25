package main

import (
	"fmt"
	"net"
)

type Message struct {
	rType string
	time  Lamport
	id    int
}

/*
channelMutexNetwork is used to send messages to the network
channelNetworkMutex is used to receive messages from the network and process them
*/
func MutexProcess(idServ int, nbServ int, chanSC chan bool, conns *[]net.Conn, done chan bool) {
	msgArray := make([]Message, nbServ)
	var clock = Lamport{}
	StartClock(&clock)
	fmt.Println("Mutex process started with id", idServ, "and clock", clock.counterTime)
	handleCommunicationWithOtherProcesses(idServ, conns, done, &msgArray, chanSC, &clock)
}

func handleNetworkMsg(id int, msg Message, chanMutexNetwork chan Message, chanSC chan bool, messages *[]Message, clock *Lamport) {
	clock.Update(msg.time)
	NoteNewMessage(id, msg, messages, chanSC)
	if msg.rType == "req" {
		// Attention l'identifiant ici est celui à qui on envoie le message
		// TODO appeler directement les méthodes ne network plutôt que de communiquer par channel
		chanMutexNetwork <- Message{"ack", *clock, msg.id}
	}
}

/*
func handleClientMsg(id int, msg string, messages *[]Message, clock *Lamport, chanMutexNetwork chan Message) {
	switch msg {
	case "ask":
		AskForSC(id, messages, clock, chanMutexNetwork)
	case "free":
		FreeSC(id, messages, clock, chanMutexNetwork)
	}
}
*/

func StartClock(clock *Lamport) {
	clock.NewLamport()
}

/*
func AskForSC(id int, messages *[]Message, clock *Lamport, chanMutexNetwork chan Message) {
	clock.Increment()
	var msg = Message{"req", *clock, id}
	fmt.Println("Message sent from ", id, msg)
	(*messages)[id] = msg
	chanMutexNetwork <- msg
}

func FreeSC(id int, messages *[]Message, clock *Lamport, chanMutexNetwork chan Message) {
	clock.Increment()
	msg := Message{"rel", *clock, id}
	(*messages)[id] = msg
	go func() {
		chanMutexNetwork <- msg
	}()
}
*/

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
