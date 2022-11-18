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
	msgArray = make([]Message, ReadConfigFile().NServ)
}

func AskForSC(id int) {
	clock.Increment()
	sendRequests(clock, id)
	msgArray[id] = Message{"req", clock, id}

}

func FreeSC(id int) {
	clock.Increment()
	sendReleases(clock, id)
	msgArray[id] = Message{"rel", clock, id}
	inSC = false
}

func NoteNewMessage(message Message, index int, id int) {
	msgArray[index] = message
	checkSCAvailable(id)
	if inSC {
		ChannelSc <- "SC"
	}
}

// todo verify is function is correct
func checkSCAvailable(id int) {
	for _, msg := range msgArray {
		if msg.id != id {
			if &msg == nil || (msg.rType == "req" && msg.time.counterTime < msgArray[id].time.counterTime && msgArray[id].rType == "req") {
				inSC = false
				return
			}
		}
	}
	inSC = true
}
