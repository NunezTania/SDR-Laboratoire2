package main

/*
list of commands the client can ask the server
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
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

// counters
var eventCounter int = 0
var postCounter int = 0
var userCounter int = 0

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

// array of events
var events []Event

// array of posts
var posts []Post

// array of users
var users []User

var isAuthentified = false

func main() {

	// Creating a simple TCP server

	// listen to incoming connections
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is listening on port 8080")

	// close the listener when the application closes
	defer listen.Close()

	createUsers()

	// accept connection on port
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
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
	}

	// display the buffer
	log.Println("La reponse du server est :", string(buf))

	// parse the buffer
	parseBuffer(buf)

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
		createEvent(slice)
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

// Create a few users
func createUsers() {
	events = append(events, Event{1, "Festival de la bière", 1})
	posts = append(posts, Post{1, "Post1", 3, events[0]})
	posts = append(posts, Post{2, "Post2", 2, events[0]})
	users = append(users, User{"Ficelle", 1234, Post{}})
	users = append(users, User{"Taro", 1234, Post{}})
	users = append(users, User{"Catanne", 1234, posts[0]})
	users = append(users, User{"Willi", 1234, posts[1]})
}

func authentification() {
	// ask for the username
	fmt.Println("Please enter your username : ")
	var username string
	fmt.Scanln(&username)

	// check if the username exists
	for _, user := range users {
		if user.name == username {

			// ask for the password
			fmt.Println("Please enter your password : ")
			var password int
			fmt.Scanln(&password)

			// check if the password is correct
			if user.password == password {
				isAuthentified = true
				fmt.Println("Authentification successful")
				return
			} else {
				fmt.Println("Wrong password")
			}
		}
	}
	fmt.Println("Wrong username")
	authentification()
}

func createEvent(slice []string) {
	fmt.Println("Starting an event")
	for !isAuthentified {
		authentification()
	}
	isAuthentified = false
	//events = append(events, Event{1, "Festival de la bière", 1})
}
