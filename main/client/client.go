package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// counters
var eventCounter int = 0
var postCounter int = 0
var userCounter int = 0

func main() {
	helpMenu()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Fatal(err)
	}
	conn.Close()
	<-done // wait for background goroutine to finish
}

func helpMenu() {
	fmt.Println("Hello and Welcome to the Event Manager")
	fmt.Println("Here is a list of possible commands:")
	fmt.Println("1. CREATE eventName postName1 capacity1 postName2 capacity2 ... ")
	fmt.Println("2. CLOSE eventName")
	fmt.Println("3. Add a benevole to a post")
	fmt.Println("4. List all the manifestations")
	fmt.Println("5. List all the posts of a manifestation")
	fmt.Println("6. List all the benevoles of a manifestation")
}
