package main

import (
	"SDR-Laboratoire1/main/dataRW"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func NetworkProcess(id int, conns *[]net.Conn, chanMutexNetwork chan Message, chanNetworkMutex chan Message, done chan bool) {
	/*conf := ReadConfigFile()
	// Init connexions with other network processes
	connections := make([]net.Conn, nbServ)
	connections[id] = nil
	for i := 0; i < nbServ; i++ {
		if i != id {
			connections[i], err = net.Dial("tcp", "localhost:"+strconv.Itoa(conf.Port + i))
			if err != nil {
				log.Fatal(err)
			}
		}
	}*/
	doneProcesses := make(chan bool)
	doneMutex := make(chan bool)
	//go handleCommunicationWithOtherProcesses(id, conns, chanNetworkMutex, doneProcesses, nil, nil, nil)
	go handleCommunicationWithMutexProcess(id, conns, chanMutexNetwork, doneMutex)

	<-doneProcesses
	<-doneMutex

	done <- true

}

func handleCommunicationWithMutexProcess(id int, conns *[]net.Conn, chanMutexNetwork chan Message, done chan bool) {
	for {
		msg := <-chanMutexNetwork
		fmt.Println("handleCommunicationWithMutexProcess : Message received in process ", id)
		// Si on reçoit un ACK depuis le processus mutex, on doit changer son id puisque celui passé est celui du processus
		// auquel il faut envoyer le message
		idTo := msg.id
		if msg.rType == "ack" {
			msg.id = id
			SendMessageTo(msg, (*conns)[idTo])
			return
		}
		// Si on reçoit un REQ ou un REL depuis le processus mutex, on doit envoyer le message à tous les autres processus
		SendToAll(id, msg, conns)
	}
	done <- true
}

func handleCommunicationWithOtherProcesses(id int, conns *[]net.Conn, done chan bool, requests *[]Message, chanSC chan bool, clock *Lamport) {
	fmt.Println("handleCommunicationWithOtherProcesses : Starting process ", id)
	for i, conn := range *conns {
		if i != id {
			fmt.Println("Listening to process ", i)
			go handleCommunicationWith(id, conn, done, requests, chanSC, clock)
		}
	}
}

func handleCommunicationWith(id int, conn net.Conn, done chan bool, requests *[]Message, chanSC chan bool, clock *Lamport) {
	for {
		fmt.Println("Process ", id, " listening to process ", conn.RemoteAddr())
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("handleCommunicationWith : Message received from process ", strToMessage(string(buf)).id)
		handleMessage(id, buf, clock, requests, chanSC, conn)
	}
	done <- true
}

func handleMessage(id int, buf []byte, clock *Lamport, requests *[]Message, chanSC chan bool, conn net.Conn) {

	var res = strings.Split(string(buf), " ")

	if res[0] == "data" { // message is a data sync
		dataRW.ProcessCommand(res[1:])

	} else { // message is a SC message

		var msg = strToMessage(string(buf))
		if msg.rType == "req" {
			*clock = clock.Update(msg.time)
			NoteNewMessage(id, msg, requests, chanSC)
			var r = Message{"ack", *clock, id}
			SendMessageTo(r, conn)

		} else if msg.rType == "rel" {
			*clock = clock.Update(msg.time)
			NoteNewMessage(id, msg, requests, chanSC)

		} else if msg.rType == "ack" {
			*clock = clock.Update(msg.time)
			NoteNewMessage(id, msg, requests, chanSC)
		}
	}
}

func strToMessage(str string) Message {
	var request Message
	var args = strings.Split(str, " ")
	request.rType = args[0]
	request.time = strToLamport(args[1])
	request.id, _ = strconv.Atoi(args[2])
	return request
}

func MessageToStr(request Message) string {
	return request.rType + " " + strconv.Itoa(request.time.counterTime) + " " + strconv.Itoa(request.id)
}

func SendMessageTo(request Message, conn net.Conn) {
	fmt.Println("SendMessageTo : Sending message to", conn.RemoteAddr())
	msg := MessageToStr(request)
	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
	}
}

func SendDataSyncTo(conn net.Conn, data []byte) {
	_, err := conn.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func SendToAll(id int, request Message, conns *[]net.Conn) {
	for i, conn := range *conns {
		if i != id {
			SendMessageTo(request, conn)
		}
	}
}

func SendDataSyncToAll(id int, conns *[]net.Conn, command []byte) {
	msg := append([]byte("data "), command...)
	for i, conn := range *conns {
		if i != id {
			SendDataSyncTo(conn, msg)
		}
	}
}

func AskForSC(id int, conns *[]net.Conn, clock *Lamport, messages *[]Message) {
	clock.Increment()
	msg := Message{"req", *clock, id}
	(*messages)[id] = msg
	SendToAll(id, msg, conns)
}

func FreeSC(id int, conns *[]net.Conn, clock *Lamport, messages *[]Message) {
	clock.Increment()
	msg := Message{"rel", *clock, id}
	(*messages)[id] = msg
	SendToAll(id, msg, conns)
}

/*
func RunBtwServer() {
	config := ReadConfigFile()
	nbServ = config.NServ
	listenConn, err := net.Listen(config.Type, config.Host+":"+strconv.Itoa(config.Port+Identifier))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("id in run btw server", Identifier)

	for {
		conn, err := listenConn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		go handleMessage(0, buf, nil, nil, nil)
	}
}

func sendRequests(clock Lamport) {
	var request = Message{"req", clock, Identifier}
	SendToAll(request)
}

func sendReleases(clock Lamport) {
	var request = Message{"rel", clock, Identifier}
	SendToAll(request)
}

*/
