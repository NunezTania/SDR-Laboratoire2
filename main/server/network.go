package main

import (
	"SDR-Laboratoire1/main/dataRW"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type conf struct {
	NServ int    `yaml:"nServ"`
	Port  int    `yaml:"port"`
	Host  string `yaml:"host"`
	Type  string `yaml:"type"`
}

var err error
var nbServ int

func RunBtwServer(id int) {
	config := ReadConfigFile()
	nbServ = config.NServ
	listenConn, err := net.Listen(config.Type, config.Host+":"+strconv.Itoa(config.Port+id))
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listenConn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		go handleMessage(buf, id)
	}
}

func handleMessage(buf []byte, id int) {

	var res = strings.Split(string(buf), " ")

	if res[0] == "data" { // message is a data sync
		dataRW.ProcessCommand(res[1:])

	} else { // message is a SC message

		var msg = strToMessage(string(buf))
		if msg.rType == "req" {
			clock = clock.Update(msg.time)
			NoteNewMessage(msg, msg.id, id)
			var r = Message{"ack", clock, id}
			SendMessageTo(msg.id, r)
			fmt.Println("I'm id = ", id, " and I sent an ack to ", msg.id, " with clock ", clock.counterTime)

		} else if msg.rType == "rel" {
			clock = clock.Update(msg.time)
			NoteNewMessage(msg, msg.id, id)

		} else if msg.rType == "ack" {
			clock = clock.Update(msg.time)
			NoteNewMessage(msg, msg.id, id)
		}
	}
}

func sendRequests(clock Lamport, id int) {
	var request = Message{"req", clock, id}
	SendToAll(request, id)
}

func sendReleases(clock Lamport, id int) {
	var request = Message{"rel", clock, id}
	SendToAll(request, id)
}

func ReadConfigFile() conf {
	yamlFile, error := os.ReadFile("./main/server/config.yaml")
	if error != nil {
		log.Printf("yamlFile.Get err   #%v ", error)
	}
	var c conf
	error = yaml.Unmarshal(yamlFile, &c)
	if error != nil {
		log.Fatalf("Unmarshal: %v", error)
	}
	if error != nil {
		log.Fatal(error)
	}
	return c
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

func SendMessageTo(id int, request Message) {
	msg := MessageToStr(request)
	var currConn net.Conn
	currConn, err = net.Dial("tcp", "localhost:"+strconv.Itoa(2500+id))
	//defer currConn.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = currConn.Write([]byte(msg))
}

func SendDataSyncTo(id int, data []byte) {
	var currConn net.Conn
	currConn, err = net.Dial("tcp", "localhost:"+strconv.Itoa(2500+id))
	defer currConn.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = currConn.Write(data)
}

func SendToAll(request Message, id int) {
	for i := 0; i < nbServ; i++ {
		if i != id {
			SendMessageTo(i, request)
		}
	}
}

func SendDataSyncToAll(command []byte, id int) {
	msg := append([]byte("data "), command...)
	for i := 0; i < nbServ; i++ {
		if i != id {
			SendDataSyncTo(i, msg)
		}
	}
}

func WaitForEveryBody(id int) {
	fmt.Println("I'm id = ", id, " and im Waiting for every body to be ready")
	msg := "ready"
	var listenConn net.Listener

	for i := 0; i < nbServ; i++ {
		if i != id {
			listenConn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(8000+i))
			//defer listenConn.Close()
			if err != nil {
				log.Fatal(err)
			}
			_, err = listenConn.Write([]byte(msg))
		}
	}

	for i := 0; i < nbServ-1; i++ {
		conn, err := listenConn.Accept()
		defer listenConn.Close()
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
	}
}
