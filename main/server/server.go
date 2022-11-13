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
	"os"
	"strconv"
)

var srvNumber int
var countSrv = 0

const (
	HOST = "localhost"
	PORT = "5557"
	TYPE = "tcp"
)

func main() {
	// check if an argument is given
	if len(os.Args) < 2 {
		srvNumber = countSrv
		countSrv++
	} else {
		srvNumber, _ = strconv.Atoi(os.Args[1])
	}
	Launch()
}

// Launch the main function of the server
func Launch() {

	Initialisation(srvNumber)
	// todo doit attendre que les autres serveurs soient prêts

	go dataRW.HandleRWActions()

	listen, err := net.Listen(TYPE, HOST+":"+PORT)

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
	clientChannel := make(chan []byte)
	dataRW.DataChannel <- clientChannel
	clientChannel <- commandParameters
	response := <-clientChannel
	return string(response)
}
