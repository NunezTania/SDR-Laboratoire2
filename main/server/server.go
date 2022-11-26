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

var ChanReady = make(chan bool)
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

func sayReadyToAll(id int, conns *[]net.Conn) {
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

	// Création des canaux pour la synchronisation
	chanSC := make(chan bool)
	done := make(chan bool)
	isReady := false

	// Variables pour processus mutex
	msgArray := make([]Message, nbServ)
	var clock = Lamport{}
	StartClock(&clock)

	//Envoie d'un message pour dire que le serveur est prêt
	sayReadyToAll(idServer, &connsWithOtherServers)

	// Lancement des goroutines pour la réception des messages
	go handleCommunicationWithServers(idServer, processListener, &connsWithOtherServers, done, &msgArray, chanSC, &clock, isReady)

	// Attente que tout les serveurs soient prêts
	<-ChanReady

	// Lancement de la boucle d'écoute pour les clients
	RunBtwClient(idServer, clientListener, &connsWithOtherServers, chanSC, &clock, &msgArray)

	<-done // Attends que les goroutines d'écoute des autres processus soient terminées

	fmt.Println("Server is closing")
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

func RunBtwClient(id int, listener net.Listener, conns *[]net.Conn, chanSC chan bool, clock *Lamport, msgArray *[]Message) {

	go dataRW.HandleRWActions()

	fmt.Println("Server is listening")
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleRequest(id, conn, conns, chanSC, clock, msgArray)
	}
}

// HandleRequest handles the requests from the clients
func HandleRequest(id int, clientConn net.Conn, processesConns *[]net.Conn, chanSC chan bool, clock *Lamport, msgArray *[]Message) {
	buf := make([]byte, 1024)
	_, err := clientConn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	for commandTreatment := AskDataRW(id, buf, processesConns, chanSC, clock, msgArray); commandTreatment != "q"; commandTreatment = AskDataRW(id, buf, processesConns, chanSC, clock, msgArray) {
		fmt.Println("Handling request")
		_, err := clientConn.Write([]byte(commandTreatment))
		if err != nil {
			return
		}
		buf = make([]byte, 1024)
		_, err = clientConn.Read(buf)
		if err != nil {
			return
		}
	}
	_, writeErr := clientConn.Write([]byte("Bye"))
	if writeErr != nil {
		log.Fatal(writeErr)
	}
	closeErr := clientConn.Close()
	if closeErr != nil {
		log.Fatal(closeErr)
	}
}

// AskDataRW asks the dataRW to treat the command
func AskDataRW(id int, commandParameters []byte, conns *[]net.Conn, chanSC chan bool, clock *Lamport, msgArray *[]Message) string {
	waitForSC(id, conns, clock, chanSC, msgArray)
	fmt.Println("Entering SC")
	clientChannel := make(chan []byte)
	dataRW.DataChannel <- clientChannel
	clientChannel <- commandParameters
	response := <-clientChannel
	if dataRW.DataModified {
		SendDataSyncToAll(id, conns, commandParameters)
	}
	leaveSC(id, conns, clock, msgArray)
	return string(response)
}

func waitForSC(id int, conns *[]net.Conn, clock *Lamport, chanSC chan bool, msgArray *[]Message) {
	AskForSC(id, conns, clock, msgArray)
	print(fmt.Println("Waiting for SC"))
	<-chanSC // Devrait bloquer jusqu'à ce que les goroutines gérant la reception des messages finissent par envoyer un message sur le channel
}

func leaveSC(id int, conns *[]net.Conn, clock *Lamport, msgArray *[]Message) {
	FreeSC(id, conns, clock, msgArray)
	fmt.Println("Leaving SC")
}

func CountReady() (IncrementReady func()) {
	var nbReady = 0
	IncrementReady = func() {
		nbReady++
	}
	checkServersReady(nbReady)
	return IncrementReady
}

func checkServersReady(nbReady int) {
	conf := ReadConfigFile()
	numberServer := conf.NServ
	if nbReady == numberServer {
		fmt.Println("wow all servers are ready")
		ChanReady <- true
	}
}
