package main

import "fmt"

// this class is used to manage the communication between all the servers
// it uses the Lamport clock to manage the time
// the critical section is the dataRW package

// the server need to send a request if the want to enter in critical section
// the server need to send a release if the want to leave the critical section
// the server will send a acknowledgement ACK if it receive a request

var clock = Lamport{}

func main() {
	clock.NewLamport()
	fmt.Println(clock.GetTime())
}

// init the Lamport clock
func init() {
	clock.NewLamport()

}

func sendRequest() {
	// increment the Lamport clock
	clock.Increment()

	// send a request to all the servers

}
