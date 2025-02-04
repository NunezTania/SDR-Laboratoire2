package test

import (
	pm "SDR-Laboratoire1/main/server/processMutex"
	"SDR-Laboratoire1/main/server/server"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

var setup = false
var servId int
var conf = pm.Config

func setUp() {
	NServ := conf.NServ
	if !setup {

		go server.LaunchNServ(NServ)
		setup = true
	}
	rand.Seed(time.Now().UnixNano())
	servId = rand.Intn(NServ)
}

func TestLISTM(t *testing.T) {
	fmt.Println("CREATE an event and do LISTM, check that the new event is there")
	setUp()
	// Create a new event
	bufferClient := make([]byte, 1024)
	cmdCreateEvent := "CREATE Lili 1234 FestiNeuch PostOne 2 PostTwo 4"
	conn, err := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+servId))
	if err != nil {
		log.Fatal(err)
	}
	_, errW := conn.Write([]byte(cmdCreateEvent))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
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

	expected := fmt.Sprintln("Event's Id: 0, Event's Name: Festival de la musique, Owner: Bob, is open:true\n" +
		"Event's Id: 1, Event's Name: Festival de la bière, Owner: Lea, is open:true\n" +
		"Event's Id: 2, Event's Name: FestiNeuch, Owner: Lili, is open:true")
	if strings.Compare(string(bufferClientIn),
		expected) != 0 {
		log.Fatal(fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClientIn)))
	}
	errClose := conn.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}

}

func TestFalseCmd(t *testing.T) {
	fmt.Println("Test")
	setUp()
	// write the command LISTX
	bufferClientIn := make([]byte, 1024)
	conn, err := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+servId))
	if err != nil {
		log.Fatal(err)
	}
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
		log.Fatal(fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClientIn)))
	}
	errClose := conn.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}
}

func TestADDOne(t *testing.T) {
	fmt.Println("CREATE an event, ADD a staff, LISTU and check that the staff was added to the event")
	setUp()
	// Create a new event
	bufferClient := make([]byte, 1024)
	cmdCreateEvent := "CREATE Lili 1234 FestiNeuch PostOne 2 PostTwo 4"
	conn, err := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+servId))
	if err != nil {
		log.Fatal(err)
	}
	_, errW := conn.Write([]byte(cmdCreateEvent))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
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
		log.Fatal(fmt.Errorf("Expected string to contain Bob as he was added to event 2, got: %s", string(bufferClientIn)))
	}

	errClose := conn.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}
}

func TestADDTwo(t *testing.T) {
	fmt.Println("ADD a staff to an event he already is in, check that the users is on the new post and was erase from the old post")
	setUp()
	// List user in even 1
	bufferClient := make([]byte, 1024)
	cmdAddUser := "LISTU 1"
	conn, err := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+servId))
	if err != nil {
		log.Fatal(err)
	}
	_, errW := conn.Write([]byte(cmdAddUser))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}

	expected := fmt.Sprintln(
		"Festival de la bière     |Vente de ticket 0   |Logistique 1        |Accueil 2           |\n" +
			"nbInscrit                |1                   |1                   |1                   |\n" +
			"Lili                     |x                   |                    |                    |\n" +
			"Lea                      |                    |x                   |                    |\n" +
			"Toto                     |                    |                    |x                   |")

	if strings.Compare(string(bufferClient), expected) != 0 {
		log.Fatal(fmt.Errorf("Expected:\n%s\ngot:\n%s", expected, string(bufferClient)))
	}
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
	expected = fmt.Sprintln("Festival de la bière     |Vente de ticket 0   |Logistique 1        |Accueil 2           |\n" +
		"nbInscrit                |0                   |1                   |1                   |\n" +
		"Lea                      |                    |x                   |                    |\n" +
		"Toto                     |                    |                    |x                   |")
	if strings.Compare(string(bufferClient), expected) != 0 {
		log.Fatal(fmt.Errorf("Expected: %s\n got:\n %s", expected, string(bufferClient)))
	}
	errClose := conn.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}

}

func TestCapacity(t *testing.T) {
	fmt.Println("CREATE an event with a post capacity of 1, ADD a person, answer = the person was added, ADD an other person, the answer = \"Could not add user to post because post is full\"")
	setUp()
	// Create a new event
	bufferClient := make([]byte, 1024)
	cmdCreateEvent := "CREATE Willi 1234 PartySound Vente 1"
	conn, err := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+servId))
	if err != nil {
		log.Fatal(err)
	}
	_, errW := conn.Write([]byte(cmdCreateEvent))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
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
	errClose := conn.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}

}

func TestCloseWithSuccessfulAuth(t *testing.T) {
	fmt.Println("CLOSE an event with a successful authentification")
	setUp()
	bufferClient := make([]byte, 1024)
	cmdCloseEvent := "CLOSE Bob 1234 0"
	conn, err := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+servId))
	if err != nil {
		log.Fatal(err)
	}
	_, errW := conn.Write([]byte(cmdCloseEvent))
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
		log.Fatal(fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClient)))
	}
	errClose := conn.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}
}

func TestCloseWithUnsuccessfulAuth(t *testing.T) {
	fmt.Println("CLOSE an event with an unsuccessful authentification")
	setUp()
	bufferClient := make([]byte, 1024)
	cmdAddUser := "CLOSE T 1234 1"
	conn, err := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+servId))
	if err != nil {
		log.Fatal(err)
	}
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
		log.Fatal(fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClient)))
	}
	errClose := conn.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}
}

func TestCreateEventWithUnsuccessfulAuth(t *testing.T) {
	fmt.Println("CREATE an event with an unsuccessful authentification")
	setUp()

	bufferClient := make([]byte, 1024)
	cmdAddUser := "CREATE Anon pwd PartySound Vente 1"
	conn, err := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+servId))
	if err != nil {
		log.Fatal(err)
	}
	_, errW := conn.Write([]byte(cmdAddUser))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := conn.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}
	expected := "Authentification failed"
	if strings.Compare(string(bufferClient), expected) != 0 {
		log.Fatal(fmt.Errorf("Expected: %s, got: %s", expected, string(bufferClient)))
	}
	errClose := conn.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}
}

func TestCreateOnOneServShouldBeOnAnother(t *testing.T) {
	fmt.Println("CREATE an event on one server should be on another when listing events")
	setUp()
	// Create a new event
	bufferClient := make([]byte, 1024)
	cmdCreateEvent := "CREATE Lili 1234 FestiNeuch PostOne 2 PostTwo 4"
	connServ1, errServ1 := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+servId))
	if errServ1 != nil {
		log.Fatal(errServ1)
	}
	_, errW := connServ1.Write([]byte(cmdCreateEvent))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := connServ1.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}

	// Display all the events
	connServ2, errServ2 := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+((servId+1)%conf.NServ)))
	if errServ2 != nil {
		log.Fatal(errServ2)
	}
	bufferClientIn := make([]byte, 1024)
	_, errW = connServ2.Write([]byte("LISTM "))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR = connServ2.Read(bufferClientIn)
	bufferClientIn = bytes.Trim(bufferClientIn, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}

	substr := fmt.Sprintln("Event's Id: 2, Event's Name: FestiNeuch, Owner: Lili, is open:true")
	if !strings.Contains(string(bufferClientIn), substr) {
		log.Fatal(fmt.Errorf("Expected: %s to contain : %s", string(bufferClientIn), substr))
	}
	errClose := connServ1.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}
	errClose = connServ2.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}
}

func TestAddUserToEventShouldAppearOnOtherServ(t *testing.T) {
	fmt.Println("ADD a user to an event should appear on another server when listing users")
	setUp()
	// Create a new event
	bufferClient := make([]byte, 1024)
	cmdCreateEvent := "ADD Bob 1234 1 0"
	connServ1, errServ1 := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+servId))
	if errServ1 != nil {
		log.Fatal(errServ1)
	}
	_, errW := connServ1.Write([]byte(cmdCreateEvent))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR := connServ1.Read(bufferClient)
	bufferClient = bytes.Trim(bufferClient, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}

	// Display all the events
	connServ2, errServ2 := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+((servId+1)%conf.NServ)))
	if errServ2 != nil {
		log.Fatal(errServ2)
	}
	bufferClientIn := make([]byte, 1024)
	_, errW = connServ2.Write([]byte("LISTU 1"))
	if errW != nil {
		log.Fatal(errW)
	}
	_, errR = connServ2.Read(bufferClientIn)
	bufferClientIn = bytes.Trim(bufferClientIn, "\x00")
	if errR != nil {
		log.Fatal(errR)
	}

	if !strings.Contains(string(bufferClientIn), "Bob") {
		log.Fatal(fmt.Errorf("Expected: %s to contain %s", string(bufferClientIn), "Bob"))
	}
	errClose := connServ1.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}
	errClose = connServ2.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}

}
