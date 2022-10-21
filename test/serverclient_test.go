package test

import (
	"SDR-Laboratoire1/main/server"
	"bytes"
	"fmt"
	"log"
	"net"
	"testing"
)

const (
	HOST = "localhost"
	PORT = "5555"
	TYPE = "tcp"
)

func Test(t *testing.T) {
	TestLISTM(t)
}

func TestLISTM(t *testing.T) {
	fmt.Println("Test")
	go server.Run()
	conn, err := net.Dial(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new event

	// Display all the events
	bufferClientIn := make([]byte, 1024)
	_, errW := conn.Write([]byte("LISTM "))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClientIn)
	bufferClientIn = bytes.Trim(bufferClientIn, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}
	if string(bufferClientIn) != "Event's id: 0, Event's name: Festival de la musique, Owner: Bob, is open:true\nEvent's id: 1, Event's name: Festival de la bière, Owner: Bob, is open:true" {
		fmt.Println("Test failed")
	} else {
		fmt.Println("Test passed")
	}
	fmt.Println(fmt.Sprintf("%s", bufferClientIn))

}

// Test idea

// 1. CREATE an event and do LISTM, check that the new event is there

// 2. CREATE an event, ADD a staff, LISTU and check that the staff was added to the event

// 3. ADD a staff to an event he already is in, check that the users is on the new post and was erase from the old post

// 4. CREATE an event with a post capacity of 1, ADD a person, answer = the person was added, ADD an other person, the answer = "Could not add user to post because post is full"

// 5. CLOSE an event with the authentification of the owner, answer should be = "Event closed"

// 6. CLOSE an event with the authentification of someone else than the owner = "Event couldn't be closed"

// 7. CREATE an event with the wrong authentification, answer should be = "Authentication failed"
