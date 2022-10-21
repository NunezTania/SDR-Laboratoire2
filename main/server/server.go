// Package server contains the server of the program.
// It is used to communicate with the client.
// And can be used to create, close, add and list events.
// In order to manage data, it uses the dataRW package.
package server

import (
	"SDR-Laboratoire1/main/dataRW"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "5557"
	TYPE = "tcp"
)

var eventCounter = 0
var postCounter = 0

type conf struct {
	NServ int    `yaml:"nServ"`
	Port  int    `yaml:"port"`
	Host  string `yaml:"host"`
}

type Event struct {
	id     int
	name   string
	owner  User
	isOpen bool
	posts  []Post
}

type Post struct {
	id       int
	name     string
	capacity int
	eventId  int
	staff    []User
}

type User struct {
	name     string
	password string
}

var events []Event
var posts []Post
var users []User

// Run the main function of the server
func Run() {
	go dataRW.HandleRWActions()
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is listening on port 8080")
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

func ReadConfigFile() conf {
	yamlFile, err := os.ReadFile("./main/server/config.yaml")
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
