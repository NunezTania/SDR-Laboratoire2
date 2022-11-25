package main

import (
	"SDR-Laboratoire1/main/dataRW"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func handleCommunicationWithServers(id int, listenConn net.Listener, conns *[]net.Conn, done chan bool, requests *[]Message, chanSC chan bool, clock *Lamport) {
	for {
		fmt.Println("Process ", id, " listening on port ", listenConn.Addr())
		conn, err := listenConn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 1024)
		_, rErr := conn.Read(buf)
		if rErr != nil {
			log.Fatal(err)
		}
		fmt.Println("handleCommunicationWith : Message received from process ", strToMessage(string(buf)).id)
		handleMessage(id, buf, clock, requests, chanSC, (*conns)[strToMessage(string(buf)).id])
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
