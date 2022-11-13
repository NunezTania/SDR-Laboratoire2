package main

import (
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

var listenConn net.Listener
var err error
var nbServ int
var id int

func Initialisation(identifier int) {
	id = identifier
	config := ReadConfigFile()
	nbServ = config.NServ
	listenConn, err = net.Listen(config.Type, config.Host+":"+strconv.Itoa(config.Port+id))
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
		go handleMessage(buf)
	}
}

func handleMessage(buf []byte) {
	var msg = strToMessage(string(buf))
	if msg.rType == "req" {
		clock = clock.Update(msg.time)
		NoteNewMessage(msg, msg.id)
		var r = Message{"ack", clock, id}
		SendMessageTo(msg.id, r)

	} else if msg.rType == "rel" {
		clock = clock.Update(msg.time)
		NoteNewMessage(msg, msg.id)

	} else if msg.rType == "ack" {
		clock = clock.Update(msg.time)
		NoteNewMessage(msg, msg.id)
	}
}

func sendRequests(clock Lamport) {
	var request = Message{"req", clock, id}
	SendToAll(request)
}

func sendReleases(clock Lamport) {
	var request = Message{"rel", clock, id}
	SendToAll(request)
}

func ReadConfigFile() conf {
	yamlFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	var c conf
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	if err != nil {
		log.Fatal(err)
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
	currConn, err = net.Dial("tcp", "localhost:"+strconv.Itoa(8000+id))
	defer currConn.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = currConn.Write([]byte(msg))
}

func SendToAll(request Message) {
	for i := 0; i < nbServ; i++ {
		if i != id {
			SendMessageTo(i, request)
		}
	}
}
