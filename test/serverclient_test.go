// Package test contains the tests for the server
// All the commands are tested
package test

import (
	"SDR-Laboratoire1/main/server"
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
)

const (
	HOST = "localhost"
	PORT = "5556"
	TYPE = "tcp"
)

var setup = false

func setUp() {
	if !setup {
		go server.Run()
		setup = true
		conn, err = net.Dial(TYPE, HOST+":"+PORT)
	}
}

var conn net.Conn
var err error

func TestLISTM(t *testing.T) {
	fmt.Println("CREATE an event and do LISTM, check that the new event is there")
	setUp()
	// Create a new event
	bufferClient := make([]byte, 1024)
	cmdCreateEvent := "CREATE Lili 1234 FestiNeuch PostOne 2 PostTwo 4"
	_, errW := conn.Write([]byte(cmdCreateEvent))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if err != nil {
		log.Fatal(errR)
	}

	// Display all the events
	bufferClientIn := make([]byte, 1024)
	_, errW = conn.Write([]byte("LISTM "))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR = conn.Read(bufferClientIn)
	bufferClientIn = bytes.Trim(bufferClientIn, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}

	expected := "Event's id: 0, Event's name: Festival de la musique, Owner: Bob, is open:true\n" +
		"Event's id: 1, Event's name: Festival de la bière, Owner: Bob, is open:true\n" +
		"Event's id: 2, Event's name: FestiNeuch, Owner: Lili, is open:true\n"
	if strings.Compare(string(bufferClientIn),
		expected) != 0 {
		fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClientIn))
	}

}

func TestFalseCmd(t *testing.T) {
	fmt.Println("Test")
	setUp()
	// write the command LISTX
	bufferClientIn := make([]byte, 1024)
	_, errW := conn.Write([]byte("LISTX "))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClientIn)
	bufferClientIn = bytes.Trim(bufferClientIn, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}
	expected := "Command not found"
	if strings.Compare(string(bufferClientIn),
		expected) != 0 {
		_ = fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClientIn))
	}
}

func TestADDOne(t *testing.T) {
	fmt.Println("CREATE an event, ADD a staff, LISTU and check that the staff was added to the event")
	setUp()
	// Create a new event
	bufferClient := make([]byte, 1024)
	cmdCreateEvent := "CREATE Lili 1234 FestiNeuch PostOne 2 PostTwo 4"
	_, errW := conn.Write([]byte(cmdCreateEvent))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if err != nil {
		log.Fatal(errR)
	}

	// Add a user Bob to the event in the second post
	bufferClientIn := make([]byte, 1024)
	_, errW = conn.Write([]byte("ADD Bob 1234 2 1"))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR = conn.Read(bufferClientIn)
	bufferClientIn = bytes.Trim(bufferClientIn, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}

	bufferClientIn = make([]byte, 1024)
	_, errW = conn.Write([]byte("LISTU 2"))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR = conn.Read(bufferClientIn)
	bufferClientIn = bytes.Trim(bufferClientIn, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}

	if !strings.Contains(string(bufferClientIn), "Bob") {
		_ = fmt.Errorf("Expected string to contain Bob as he was added to event 2, got: %s", string(bufferClientIn))
	}
}

func TestADDTwo(t *testing.T) {
	fmt.Println("ADD a staff to an event he already is in, check that the users is on the new post and was erase from the old post")
	setUp()
	// List user in even 1
	bufferClient := make([]byte, 1024)
	cmdAddUser := "LISTU 1"
	_, errW := conn.Write([]byte(cmdAddUser))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}

	expected := "Festival de la bière     |Vente de ticket 0   |Logistique 1        |Securité 2          |\nnbInscrit                |1                   |2                   |1                   |\nBob                      |x                   |                    |                    |\nLeo                      |                    |x                   |                    |\nWilli                    |                    |x                   |                    |\nLili                     |                    |                    |x                   |"

	if strings.Compare(string(bufferClient), expected) != 0 {
		_ = fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClient))
	}
	setUp()
	// Add user Lili to an event where she already is
	bufferClient = make([]byte, 1024)
	cmdAddUser = "ADD Lili 1234 1 1"
	_, errW = conn.Write([]byte(cmdAddUser))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR = conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}
	bufferClient = make([]byte, 1024)
	cmdAddUser = "LISTU 1"
	_, errW = conn.Write([]byte(cmdAddUser))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR = conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}
	expected = "Festival de la bière     |Vente de ticket 0   |Logistique 1        |Securité 2          |\nnbInscrit                |1                   |3                   |0                   |\nBob                      |x                   |                    |                    |\nLeo                      |                    |x                   |                    |\nWilli                    |                    |x                   |                    |\nLili                     |                    |x                   |                    |"
	if strings.Compare(string(bufferClient), expected) != 0 {
		_ = fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClient))
	}

}

func TestCapacity(t *testing.T) {
	fmt.Println("CREATE an event with a post capacity of 1, ADD a person, answer = the person was added, ADD an other person, the answer = \"Could not add user to post because post is full\"")
	setUp()
	// Create a new event
	bufferClient := make([]byte, 1024)
	cmdCreateEvent := "CREATE Willi 1234 PartySound Vente 1"
	_, errW := conn.Write([]byte(cmdCreateEvent))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if err != nil {
		log.Fatal(errR)
	}

	// Add a first staff
	bufferClient = make([]byte, 1024)
	cmdAddUser := "ADD Lili 1234 2 0"
	_, errW = conn.Write([]byte(cmdAddUser))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR = conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}
	expected := "User successfully added to post"
	if strings.Compare(string(bufferClient), expected) != 0 {
		_ = fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClient))
	}

	// add a second staff
	bufferClient = make([]byte, 1024)
	cmdAddUser = "ADD Leo 1234 2 0"
	_, errW = conn.Write([]byte(cmdAddUser))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR = conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}
	expected = "Could not add user to post because post is full"
	if strings.Compare(string(bufferClient), expected) != 0 {
		_ = fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClient))
	}

}

func TestCloseWithSuccessfulAuth(t *testing.T) {
	fmt.Println("CLOSE an event with a successful authentification")
	setUp()
	bufferClient := make([]byte, 1024)
	cmdAddUser := "CLOSE Bob 1234 1"
	_, errW := conn.Write([]byte(cmdAddUser))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}
	expected := "Event closed"
	if strings.Compare(string(bufferClient), expected) != 0 {
		_ = fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClient))
	}
}

func TestCloseWithUnsuccessfulAuth(t *testing.T) {
	fmt.Println("CLOSE an event with an unsuccessful authentification")
	setUp()
	bufferClient := make([]byte, 1024)
	cmdAddUser := "CLOSE T 1234 1"
	_, errW := conn.Write([]byte(cmdAddUser))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}
	expected := "Event couldn't be closed"
	if strings.Compare(string(bufferClient), expected) != 0 {
		_ = fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClient))
	}
}

func TestCreateEventWithUnsuccessfulAuth(t *testing.T) {
	fmt.Println("CREATE an event with an unsuccessful authentification")
	setUp()

	bufferClient := make([]byte, 1024)
	cmdAddUser := "CREATE Anon pwd PartySound Vente 1"
	_, errW := conn.Write([]byte(cmdAddUser))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}
	expected := "Authentication failed"
	if strings.Compare(string(bufferClient), expected) != 0 {
		_ = fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClient))
	}
}
