package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// this class is used to manage the communication between all the servers
// it uses the Lamport clock to manage the time
// the critical section is the dataRW package

// the server need to send a request if the want to enter in critical section
// the server need to send a release if the want to leave the critical section
// the server will send a acknowledgement ACK if it receive a request

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

// init the Lamport clock
func initialisation(identifier int) {
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
	var request = strToMessage(string(buf))
	if request.rType == "req" {
		clock = clock.Update(request.time)
		var r = Request{"ack", clock, id}
		SendMessageTo(request.id, r)
		// todo add the new message in the mutex array

	} else if request.rType == "rel" {
		clock = clock.Update(request.time)

	} else if request.rType == "ack" {
		// todo
	}
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

func strToMessage(str string) Request {
	var request Request
	var args = strings.Split(str, " ")
	request.rType = args[0]
	request.time = strToLamport(args[1])
	request.id, _ = strconv.Atoi(args[2])
	return request
}

func MessageToStr(request Request) string {
	return request.rType + " " + strconv.Itoa(request.time.counterTime) + " " + strconv.Itoa(request.id)
}

func SendMessageTo(id int, request Request) {
	msg := MessageToStr(request)
	var currConn net.Conn
	currConn, err = net.Dial("tcp", "localhost:"+strconv.Itoa(8000+id))
	defer currConn.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = currConn.Write([]byte(msg))
}

func SendToAll(request Request) {
	for i := 0; i < nbServ; i++ {
		if i != id {
			SendMessageTo(i, request)
		}
	}
}
