package main

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
}

func AskForSC() {
	clock.Increment()
	sendRequests(clock)
	msgArray[Identifier] = Message{"req", clock, Identifier}

}

func FreeSC() {
	clock.Increment()
	sendReleases(clock)
	msgArray[Identifier] = Message{"rel", clock, Identifier}
	inSC = false
}

func NoteNewMessage(message Message, index int) {
	msgArray[index] = message
	checkSCAvailable()
	if inSC {
		ChannelSc <- "SC"
	}
}

// todo verify is function is correct
func checkSCAvailable() {
	for _, msg := range msgArray {
		if msg.id != Identifier {
			if &msg == nil || (msg.rType == "req" && msg.time.counterTime < msgArray[Identifier].time.counterTime && msgArray[Identifier].rType == "req") {
				inSC = false
				return
			}
		}
	}
	inSC = true
}
