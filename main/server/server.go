// Package server contains the server of the program.
// It is used to communicate with the client.
// And can be used to create, close, add and list events.
// In order to manage data, it uses the dataRW package.
package main

import (
	"SDR-Laboratoire1/main/dataRW"
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	HOST = "localhost"
	PORT = 5557
	TYPE = "tcp"
)

func main() {
	numberServer := 3

	for i := 0; i < numberServer; i++ {
		Launch(i)
	}
	for {
	}
}

func Launch(idServer int) {
	WaitForEveryBody(idServer)
	var clock = Lamport{}
	var inSC = false
	var ChannelSc = make(chan string)
	StartClock(&clock)
	go RunBtwServer(idServer, &clock, &inSC, &ChannelSc)
	go RunBtwClient(idServer, &ChannelSc, &clock, &inSC)
}

func RunBtwClient(id int, ChannelSC *chan string, clock *Lamport, inSC *bool) {

	go dataRW.HandleRWActions()
	port := strconv.Itoa(PORT + id)
	listen, err := net.Listen(TYPE, HOST+":"+port)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is listening")
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listen)
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleRequest(conn, id, ChannelSC, clock, inSC)
	}
}

// HandleRequest handles the requests from the clients
func HandleRequest(conn net.Conn, id int, ChannelSC *chan string, clock *Lamport, inSC *bool) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	for commandTreatment := AskDataRW(buf, id, ChannelSC, clock, inSC); commandTreatment != "q"; commandTreatment = AskDataRW(buf, id, ChannelSC, clock, inSC) {
		fmt.Println("Handling request")
		_, err := conn.Write([]byte(commandTreatment))
		if err != nil {
			return
		}
		buf = make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			return
		}
	}
	_, writeErr := conn.Write([]byte("Bye"))
	if writeErr != nil {
		log.Fatal(writeErr)
	}
	closeErr := conn.Close()
	if closeErr != nil {
		log.Fatal(closeErr)
	}
}

// AskDataRW asks the dataRW to treat the command
func AskDataRW(commandParameters []byte, id int, ChannelSC *chan string, clock *Lamport, inSC *bool) string {
	waitForSC(id, ChannelSC, clock)
	clientChannel := make(chan []byte)
	dataRW.DataChannel <- clientChannel
	clientChannel <- commandParameters
	response := <-clientChannel
	leaveSC(id, clock, inSC)
	if dataRW.DataModified {
		SendDataSyncToAll(commandParameters, id)
	}
	return string(response)
}

func waitForSC(id int, ChannelSc *chan string, clock *Lamport) {
	AskForSC(id, clock)
	<-*ChannelSc
}

func leaveSC(id int, clock *Lamport, inSC *bool) {
	FreeSC(id, clock, inSC)
}
