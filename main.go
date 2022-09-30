package main

import (
	"strings"
)

/*
list of commands the client can ask to the server
- Create a new event						: CREATE
- Close the event 							: CLOSE
- Add a benevole to a post					: ADD
- List all the manifestations				: LISTM
- List all the posts of a manifestation		: LISTP
- List all the benevoles of a manifestation : LISTB
*/

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

type Event struct {
	id      int
	name    string
	ownerId int
}

type Post struct {
	id       int
	name     string
	capacity int
	event    Event
}

type User struct {
	name     string
	password int
	post     Post
}

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

	// parse the buffer
	go parseBuffer(buf)

	// send a response back to person contacting us
	conn.Write([]byte("Message received."))
	// close the connection when you're done with it
	conn.Close()
}

/*
this function parses the buffer and call the appropriate function
*/
func parseBuffer(buf []byte) {
	// transform the buffer into a string
	str := string(buf)

	// split the string into a slice
	slice := strings.Split(str, " ")

	// do a switch on the first element of slice
	switch slice[0] {
	case "CREATE":
		fmt.Println("CREATE")
	case "CLOSE":
		fmt.Println("CLOSE")
	case "ADD":
		fmt.Println("ADD")
	case "LISTM":
		fmt.Println("LISTM")
	case "LISTP":
		fmt.Println("LISTP")
	case "LISTB":
		fmt.Println("LISTB")
	default:
		fmt.Println("Command not found")
	}
}
