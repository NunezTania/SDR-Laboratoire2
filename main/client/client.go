package main

import (
	"bufio"
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
	helpMenu()

	// connection to the server
	conn, err := net.Dial(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})

	// read the user input
	fmt.Println("Please enter a command:")
	var command string
	for {
		Scanner := bufio.NewScanner(os.Stdin)
		Scanner.Scan()
		if err := Scanner.Err(); err != nil {
			log.Fatal(err)
		}
		command = Scanner.Text()
		if command != "" {
			break
		}
	}

	// todo check if the command is valid

	// send the command to the server
	_, err = conn.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}

	conn.Close()
	<-done
}

func helpMenu() {
	fmt.Println("Hello and Welcome to the Event Manager")
	fmt.Println("Here is a list of possible commands:")
	fmt.Println("1. CREATE eventName postName1 capacity1 postName2 capacity2 ... ")
	fmt.Println("2. CLOSE eventName")
	fmt.Println("3. ADD eventId postId")
	fmt.Println("4. List all the event")
	fmt.Println("5. List all the posts of a event")
	fmt.Println("6. List all the benevoles of a event")
}
