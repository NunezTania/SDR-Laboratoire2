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
	StartClock()
	go RunBtwServer(idServer)
	go RunBtwClient(idServer)
}

func RunBtwClient(id int) {

	go dataRW.HandleRWActions()
	port := strconv.Itoa(PORT + id)
	listen, err := net.Listen(TYPE, HOST+":"+port)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is listening")
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleRequest(conn, id)
	}
}

// HandleRequest handles the requests from the clients
func HandleRequest(conn net.Conn, id int) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	for commandTreatment := AskDataRW(buf, id); commandTreatment != "q"; commandTreatment = AskDataRW(buf, id) {
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
func AskDataRW(commandParameters []byte, id int) string {
	waitForSC(id)
	clientChannel := make(chan []byte)
	dataRW.DataChannel <- clientChannel
	clientChannel <- commandParameters
	response := <-clientChannel
	leaveSC(id)
	if dataRW.DataModified {
		SendDataSyncToAll(commandParameters, id)
	}
	return string(response)
}

func waitForSC(id int) {
	AskForSC(id)
	<-ChannelSc
}

func leaveSC(id int) {
	FreeSC(id)
}
