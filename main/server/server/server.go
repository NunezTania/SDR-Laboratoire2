// Package server contains the server of the program.
// It is used to communicate with the client.
// And can be used to create, close, add and list events.
// In order to manage data, it uses the dataRW package.
package server

import (
	"SDR-Laboratoire1/main/dataRW"
	pm "SDR-Laboratoire1/main/server/processMutex"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

var conf = pm.Config

// LaunchNServ launches n servers
func LaunchNServ(nServ int) {
	doneChans := make([]chan bool, nServ)
	for i := 0; i < nServ; i++ {
		doneChans[i] = make(chan bool)
	}
	for i := 0; i < nServ; i++ {
		go Launch(i, doneChans[i])
	}
	for i := 0; i < nServ; i++ {
		<-doneChans[i]
	}
}

// createUsersAndEventsFromConf creates the users and events from the configuration file
func createUsersAndEventsFromConf(users *[]dataRW.User, events *[]dataRW.Event, eventCounter *int) {
	for _, user := range conf.Users {
		*users = append(*users, dataRW.User{user.Name, user.Password})
	}
	// creation of events
	for _, event := range conf.Events {
		*eventCounter++
		*events = append(*events,
			dataRW.Event{event.Id, event.Name, event.Owner, true, event.Posts})
	}
}

// Launch launches a specific server
func Launch(idServer int, doneChan chan bool) {
	listenConn, err := net.Listen(pm.Config.Type, pm.Config.Host+":"+strconv.Itoa(pm.Config.PortServ+idServer))
	if err != nil {
		log.Fatal(err)
	}
	pm.WaitForEveryBody(idServer, listenConn)
	var clock = pm.Lamport{}
	var inSC = false
	var ChannelSc = make(chan string)
	var DataChannel = make(chan chan []byte)
	doneClient := make(chan bool)
	doneServer := make(chan bool)
	messages := make([]pm.Message, conf.NServ)
	pm.StartClock(&clock)
	go pm.RunBtwServer(idServer, &clock, &inSC, &ChannelSc, &DataChannel, doneServer, listenConn, &messages)
	go RunBtwClient(idServer, &ChannelSc, &clock, &inSC, &DataChannel, doneClient, &messages)
	<-doneClient
	<-doneServer
	doneChan <- true
}

// RunBtwClient server listen to the clients
func RunBtwClient(id int, ChannelSC *chan string, clock *pm.Lamport, inSC *bool, DataChannel *chan chan []byte, done chan bool, messages *[]pm.Message) {
	var DataModified = false
	var eventCounter = 0
	var postCounter = 0

	var events = make([]dataRW.Event, 0)
	var users = make([]dataRW.User, 0)
	createUsersAndEventsFromConf(&users, &events, &eventCounter)
	go dataRW.HandleRWActions(DataChannel, &DataModified, &users, &events, &postCounter, &eventCounter)
	listen, err := net.Listen(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+id))

	for err != nil {
		listen, err = net.Listen(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+id))
	}

	fmt.Println("Server ", id, " is listening")
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listen)
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go HandleRequest(conn, id, ChannelSC, clock, inSC, DataChannel, &DataModified, messages)
	}
	done <- true
}

// HandleRequest handles the requests from the clients
func HandleRequest(conn net.Conn, id int, ChannelSC *chan string, clock *pm.Lamport, inSC *bool, DataChannel *chan chan []byte, DataModified *bool, messages *[]pm.Message) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	for commandTreatment := AskDataRW(buf, id, ChannelSC, clock, inSC, DataChannel, DataModified, messages); commandTreatment != "q"; commandTreatment = AskDataRW(buf, id, ChannelSC, clock, inSC, DataChannel, DataModified, messages) {
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
func AskDataRW(commandParameters []byte, id int, ChannelSC *chan string, clock *pm.Lamport, inSC *bool, DataChannel *chan chan []byte, DataModified *bool, messages *[]pm.Message) string {
	waitForSC(id, ChannelSC, clock, messages)
	if pm.Config.Debug == 1 {
		fmt.Println("Sleeping for debug...")
		time.Sleep(10 * time.Second)
	}
	fmt.Println("Server ", id, " is in SC")
	clientChannel := make(chan []byte)
	*DataChannel <- clientChannel
	clientChannel <- commandParameters
	response := <-clientChannel
	leaveSC(id, clock, inSC, messages)
	if *DataModified {
		pm.SendDataSyncToAll(commandParameters, id)
	}
	return string(response)
}

// waitForSC waits for the SC
func waitForSC(id int, ChannelSc *chan string, clock *pm.Lamport, messages *[]pm.Message) {
	if pm.Config.NServ == 1 {
		return
	}
	pm.AskForSC(id, clock, messages)
	<-*ChannelSc
}

// leaveSC leaves the SC
func leaveSC(id int, clock *pm.Lamport, inSC *bool, messages *[]pm.Message) {
	pm.FreeSC(id, clock, inSC, messages)
}
