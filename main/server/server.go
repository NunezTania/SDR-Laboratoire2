// Package server contains the server of the program.
// It is used to communicate with the client.
// And can be used to create, close, add and list events.
// In order to manage data, it uses the dataRW package.
package main

import (
	"SDR-Laboratoire1/main/dataRW"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	HOST = "localhost"
	PORT = "5557"
	TYPE = "tcp"
)

type conf struct {
	NServ int    `yaml:"nServ"`
	Port  int    `yaml:"port"`
	Host  string `yaml:"host"`
	Type  string `yaml:"type"`
}

var Identifier int

func main() {
	conf := ReadConfigFile()
	numberServer := conf.NServ

	for i := 0; i < numberServer; i++ {
		go Launch(i, conf)
	}
	for {
		// Waiting for every body to be ready
	}
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

func WaitForEveryBody(id int, nbServ int, conns *[]net.Conn, listener net.Listener) {
	fmt.Println("Waiting for every body to be ready")
	msg := "ready"

	for i, conn := range *conns {
		if i != id {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for i := 0; i < nbServ-1; i++ {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
	}
	fmt.Println("Everybody is ready")
}

func Launch(idServer int, conf conf) {
	fmt.Println("id from launch : ", idServer)

	nbServ := conf.NServ
	port := conf.Port
	host := conf.Host
	typeConn := conf.Type

	// Création du listener pour les clients
	clientListener, errClient := net.Listen(typeConn, host+":"+strconv.Itoa(port+nbServ+idServer))
	if errClient != nil {
		log.Fatal(errClient)
	}

	// Création du listener pour les autres processus
	processListener, errProcess := net.Listen(typeConn, host+":"+strconv.Itoa(port+idServer))
	if errProcess != nil {
		log.Fatal(errProcess)
	}

	// Création des connexions entre les processus
	connsWithOtherServers := make([]net.Conn, nbServ)
	connsWithOtherServers[idServer] = nil
	for i := 0; i < nbServ; i++ {
		if i != idServer {
			var err error
			connsWithOtherServers[i], err = net.Dial(typeConn, host+":"+strconv.Itoa(port+i))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	WaitForEveryBody(idServer, nbServ, &connsWithOtherServers, processListener)

	/*
		chanClientMutex := make(chan string)
		chanMutexNetwork := make(chan Message)
		chanNetworkMutex := make(chan Message)
		chanSC := make(chan bool)
		doneMutex := make(chan bool)
		doneNetwork := make(chan bool)

		go RunBtwClient(idServer, clientListener, &connsWithOtherServers, chanClientMutex, chanSC)
		go MutexProcess(idServer, chanClientMutex, chanMutexNetwork, chanNetworkMutex, chanSC, doneMutex)
		go NetworkProcess(idServer, &connsWithOtherServers, chanMutexNetwork, chanNetworkMutex, doneNetwork)

		<-doneMutex
		<-doneNetwork
	*/
	// Fermeture des connexions avec les autres processus servers
	for i, conn := range connsWithOtherServers {
		if i != idServer {
			err := conn.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Closing client listener
	errClient = clientListener.Close()
	if errClient != nil {
		log.Fatal(errClient)
	}

	// Closing process listener
	errProcess = processListener.Close()
	if errProcess != nil {
		log.Fatal(errProcess)
	}

}

func RunBtwClient(id int, listener net.Listener, conns *[]net.Conn, chanClientMutex chan string, chanSC chan bool) {

	go dataRW.HandleRWActions()

	fmt.Println("Server is listening")
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleRequest(id, conn, conns, chanClientMutex, chanSC)
	}
}

// HandleRequest handles the requests from the clients
func HandleRequest(id int, conn net.Conn, conns *[]net.Conn, chanClientMutex chan string, chanSC chan bool) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	for commandTreatment := AskDataRW(id, buf, conns, chanClientMutex, chanSC); commandTreatment != "q"; commandTreatment = AskDataRW(id, buf, conns, chanClientMutex, chanSC) {
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
func AskDataRW(id int, commandParameters []byte, conns *[]net.Conn, chanClientMutex chan string, chanSC chan bool) string {
	waitForSC(chanClientMutex, chanSC)
	clientChannel := make(chan []byte)
	dataRW.DataChannel <- clientChannel
	clientChannel <- commandParameters
	response := <-clientChannel
	if dataRW.DataModified {
		SendDataSyncToAll(id, conns, commandParameters)
	}
	leaveSC(chanClientMutex)
	return string(response)
}

func waitForSC(chanClientMutex chan string, chanSC chan bool) {
	chanClientMutex <- "ask"
	<-chanSC
}

func leaveSC(chanClientMutex chan string) {
	chanClientMutex <- "free"
}
