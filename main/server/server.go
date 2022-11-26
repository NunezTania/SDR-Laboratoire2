// Package server contains the server of the program.
// It is used to communicate with the client.
// And can be used to create, close, add and list events.
// In order to manage data, it uses the dataRW package.
package main

import (
	"SDR-Laboratoire1/main/dataRW"
	pm "SDR-Laboratoire1/main/server/processMutex"
	"fmt"
	"log"
	"net"
	"strconv"
)

var conf = pm.Config

func main() {
	doneChans := make([]chan bool, conf.NServ)
	listenChans := make([]chan bool, conf.NServ)
	for i := 0; i < conf.NServ; i++ {
		doneChans[i] = make(chan bool)
	}
	for i := 0; i < conf.NServ; i++ {
		listenChans[i] = make(chan bool)
	}
	for i := 0; i < conf.NServ; i++ {
		go Launch(i, doneChans[i], &listenChans)
	}
	for i := 0; i < conf.NServ; i++ {
		<-doneChans[i]
	}
}

func Launch(idServer int, doneChan chan bool, listenChans *[]chan bool) {
	listenConn, err := net.Listen(pm.Config.Type, pm.Config.Host+":"+strconv.Itoa(pm.Config.PortServ+idServer))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for i := 0; i < pm.Config.NServ-1; i++ { // Unblocks other servers when they try to dial this one
			(*listenChans)[idServer] <- true
		}
	}() // In a goroutine to avoid blocking
	pm.WaitForEveryBody(idServer, listenConn, listenChans)
	var clock = pm.Lamport{}
	var inSC = false
	var ChannelSc = make(chan string)
	var DataChannel = make(chan chan []byte)
	doneClient := make(chan bool)
	doneServer := make(chan bool)
	pm.StartClock(&clock)
	go pm.RunBtwServer(idServer, &clock, &inSC, &ChannelSc, &DataChannel, doneServer, listenConn)
	go RunBtwClient(idServer, &ChannelSc, &clock, &inSC, &DataChannel, doneClient)
	<-doneClient
	<-doneServer
	doneChan <- true
}

func RunBtwClient(id int, ChannelSC *chan string, clock *pm.Lamport, inSC *bool, DataChannel *chan chan []byte, done chan bool) {
	var DataModified = false
	go dataRW.HandleRWActions(DataChannel, &DataModified)
	listen, err := net.Listen(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+id))

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
		go HandleRequest(conn, id, ChannelSC, clock, inSC, DataChannel, &DataModified)
	}
	done <- true
}

// HandleRequest handles the requests from the clients
func HandleRequest(conn net.Conn, id int, ChannelSC *chan string, clock *pm.Lamport, inSC *bool, DataChannel *chan chan []byte, DataModified *bool) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	for commandTreatment := AskDataRW(buf, id, ChannelSC, clock, inSC, DataChannel, DataModified); commandTreatment != "q"; commandTreatment = AskDataRW(buf, id, ChannelSC, clock, inSC, DataChannel, DataModified) {
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
func AskDataRW(commandParameters []byte, id int, ChannelSC *chan string, clock *pm.Lamport, inSC *bool, DataChannel *chan chan []byte, DataModified *bool) string {
	waitForSC(id, ChannelSC, clock)
	clientChannel := make(chan []byte)
	*DataChannel <- clientChannel
	clientChannel <- commandParameters
	response := <-clientChannel
	leaveSC(id, clock, inSC)
	if *DataModified {
		pm.SendDataSyncToAll(commandParameters, id)
	}
	return string(response)
}

func waitForSC(id int, ChannelSc *chan string, clock *pm.Lamport) {
	pm.AskForSC(id, clock)
	<-*ChannelSc
}

func leaveSC(id int, clock *pm.Lamport, inSC *bool) {
	pm.FreeSC(id, clock, inSC)
}
