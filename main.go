package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {

	// Creating a simple TCP server

	// listen to incoming connections
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println("Server is listening on port 8080")

	// close the listener when the application closes
	defer listen.Close()

	// accept connection on port
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// call the handleRequest function
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {

	// make a buffer to hold incoming data
	buf := make([]byte, 1024)

	// read the incoming connection into the buffer
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// display the buffer
	log.Println("La reponse du server est :", string(buf))

	// send a response back to person contacting us
	conn.Write([]byte("Message received."))
	// close the connection when you're done with it
	conn.Close()
}
