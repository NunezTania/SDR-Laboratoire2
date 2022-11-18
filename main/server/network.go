package main

import (
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
	go handleCommunicationWithOtherProcesses(id, conns, chanNetworkMutex, done)
	go handleCommunicationWithMutexProcess(id, conns, chanMutexNetwork, done)

}

func handleCommunicationWithMutexProcess(id int, conns *[]net.Conn, chanMutexNetwork chan Message, done chan bool) {
	for {
		msg := <-chanMutexNetwork
		fmt.Println("Message received from mutex process ", id)
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

func handleCommunicationWithOtherProcesses(id int, conns *[]net.Conn, chanMutexNetwork chan Message, done chan bool) {
	for i, conn := range *conns {
		if i != id {
			go handleCommunicationWith(conn, chanMutexNetwork, done)
		}
	}
}

func handleCommunicationWith(conn net.Conn, chanMutexNetwork chan Message, done chan bool) {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		chanMutexNetwork <- strToMessage(string(buf))
		fmt.Println("Message received from", conn.RemoteAddr())
	}
	done <- true
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
		go handleMessage(buf)
	}
}

func handleMessage(buf []byte) {

	var res = strings.Split(string(buf), " ")

	if res[0] == "data" { // message is a data sync
		dataRW.ProcessCommand(res[1:])

	} else { // message is a SC message

		var msg = strToMessage(string(buf))
		if msg.rType == "req" {
			clock = clock.Update(msg.time)
			NoteNewMessage(msg, msg.id)
			var r = Message{"ack", clock, Identifier}
			SendMessageTo(msg.id, r)

		} else if msg.rType == "rel" {
			clock = clock.Update(msg.time)
			NoteNewMessage(msg, msg.id)

		} else if msg.rType == "ack" {
			clock = clock.Update(msg.time)
			NoteNewMessage(msg, msg.id)
		}
	}
}

func sendRequests(clock Lamport) {
	var request = Message{"req", clock, Identifier}
	SendToAll(request)
}

func sendReleases(clock Lamport) {
	var request = Message{"rel", clock, Identifier}
	SendToAll(request)
}*/

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
