package main

var ChannelSc = make(chan string)

type Message struct {
	rType string
	time  Lamport
	id    int
}

/*
channelMutexNetwork is used to send messages to the network
channelNetworkMutex is used to receive messages from the network and process them
*/
func MutexProcess(idServ int, chanClientMutex chan string, chanMutexNetwork chan Message, chanNetworkMutex chan Message, chanSC chan bool, done chan bool) {
	var msgArray []Message
	var clock = Lamport{}
	StartClock(&clock)
	for {
		select {
		case clientMsg := <-chanClientMutex:
			handleClientMsg(idServ, clientMsg, &msgArray, &clock, chanMutexNetwork)
		case networkMsg := <-chanNetworkMutex:
			handleNetworkMsg(idServ, networkMsg, chanMutexNetwork, chanSC, &msgArray, &clock)
		}
	}
	done <- true
}

func handleNetworkMsg(id int, msg Message, chanMutexNetwork chan Message, chanSC chan bool, messages *[]Message, clock *Lamport) {
	clock.Update(msg.time)
	NoteNewMessage(id, msg, messages, chanSC)
	if msg.rType == "req" {
		// Attention l'identifiant ici est celui à qui on envoie le message
		chanMutexNetwork <- Message{"ack", *clock, msg.id}
	}
}

func handleClientMsg(id int, msg string, messages *[]Message, clock *Lamport, chanMutexNetwork chan Message) {
	switch msg {
	case "ask":
		AskForSC(id, messages, clock, chanMutexNetwork)
	case "free":
		FreeSC(id, messages, clock, chanMutexNetwork)
	}
}

func StartClock(clock *Lamport) {
	clock.NewLamport()
}

func AskForSC(id int, messages *[]Message, clock *Lamport, chanMutexNetwork chan Message) {
	clock.Increment()
	var msg = Message{"req", *clock, id}
	(*messages)[id] = msg
	chanMutexNetwork <- msg
}

func FreeSC(id int, messages *[]Message, clock *Lamport, chanMutexNetwork chan Message) {
	clock.Increment()
	msg := Message{"rel", *clock, id}
	(*messages)[id] = msg
	chanMutexNetwork <- msg
}

func NoteNewMessage(id int, msg Message, messages *[]Message, chanSC chan bool) {
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
