package server

// The server package contains the server of the program.
// It is used to communicate with the client.
// And can be used to create, close, add and list events.
// In order to manage data, it uses the dataRW package.

import (
	"SDR-Laboratoire1/main/dataRW"
	"fmt"
	"log"
	"net"
)

const (
	HOST = "localhost"
	PORT = "5555"
	TYPE = "tcp"
)

var eventCounter = 0
var postCounter = 0

type Event struct {
	id     int
	name   string
	owner  User
	isOpen bool
	posts  []Post
}

type Post struct {
	id       int
	name     string
	capacity int
	eventId  int
	staff    []User
}

type User struct {
	name     string
	password string
}

var events []Event
var posts []Post
var users []User

// Run function of the server
func Run() {
	go dataRW.HandleRWActions()
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is listening on port 8080")
	defer listen.Close()
	createUsersAndEvents()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleRequest(conn)
	}
}

// handleRequest handles the requests from the clients
func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	for commandTreatment := askDataRW(buf); commandTreatment != "q"; commandTreatment = askDataRW(buf) {
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

// askDataRW asks the dataRW to treat the command
func askDataRW(commandParameters []byte) string {
	clientChannel := make(chan []byte)
	dataRW.DataChannel <- clientChannel
	clientChannel <- commandParameters
	response := <-clientChannel
	return string(response)
}

// createUsersAndEvents creates some users and some events
func createUsersAndEvents() {
	// creation of users
	users = append(users, User{"Bob", "1234"})
	users = append(users, User{"Lea", "1234"})
	users = append(users, User{"Leo", "1234"})
	users = append(users, User{"Willi", "1234"})
	users = append(users, User{"Lili", "1234"})
	users = append(users, User{"T", "1234"})
	// creation of posts
	posts = append(posts, Post{postCounter, "Bar à bière", 3, 0, users[0:1]})
	postCounter++
	posts = append(posts, Post{postCounter, "Securité", 2, 0, users[2:4]})
	postCounter++
	posts = append(posts, Post{postCounter, "Vente de ticket", 5, 1, users[0:1]})
	postCounter++
	posts = append(posts, Post{postCounter, "Logistique", 1, 1, users[2:4]})
	postCounter++
	posts = append(posts, Post{postCounter, "Securité", 2, 1, users[4:5]})
	postCounter++
	// creation of events
	events = append(events, Event{eventCounter, "Festival de la musique", users[0], true, posts[0:2]})
	eventCounter++
	events = append(events, Event{eventCounter, "Festival de la bière", users[0], true, posts[2:5]})
	eventCounter++
}
