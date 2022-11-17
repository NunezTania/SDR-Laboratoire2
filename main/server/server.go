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
	PORT = "5557"
	TYPE = "tcp"
)

var Identifier int

func main() {
	numberServer := 2

	for i := 0; i < numberServer; i++ {
		go Launch(i)
	}
	for {
	}
}

func Launch(idServer int) {
	fmt.Println("id from launch : ", idServer)
	Identifier = idServer
	//WaitForEveryBody()
	StartClock()
	go RunBtwServer()
	go RunBtwClient()
}

func RunBtwClient() {

	go dataRW.HandleRWActions()
	port := strconv.Itoa(5557 + Identifier)
	fmt.Println("port : ", port)
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
		go HandleRequest(conn)
	}
}

// HandleRequest handles the requests from the clients
func HandleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	for commandTreatment := AskDataRW(buf); commandTreatment != "q"; commandTreatment = AskDataRW(buf) {
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
func AskDataRW(commandParameters []byte) string {
	waitForSC()
	clientChannel := make(chan []byte)
	dataRW.DataChannel <- clientChannel
	clientChannel <- commandParameters
	response := <-clientChannel
	leaveSC()
	if dataRW.DataModified {
		SendDataSyncToAll(commandParameters)
	}
	return string(response)
}

func waitForSC() {
	AskForSC()
	<-ChannelSc
}

func leaveSC() {
	FreeSC()
}
