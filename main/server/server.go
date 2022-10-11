package main

/*
list of commands the client can ask the server
- Create a new event						: CREATE
- Close the event 							: CLOSE
- Add a benevole to a post					: ADD
- List all the manifestations				: LISTM
- List all the posts of a manifestation		: LISTP
*/

import (
	"fmt"
	"log"
	"net"
	"strconv"
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

type Event struct {
	id      int
	name    string
	ownerId int
	isOpen  bool
}

type Post struct {
	id       int
	name     string
	capacity int
	eventId  int
}

type User struct {
	name     string
	password string
	post     Post
}

// array of events
var events []Event

// array of posts
var posts []Post

// array of users
var users []User

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is listening on port 8080")

	defer listen.Close()

	createUsers()
	conn := net.Conn(nil)

	conn, err = listen.Accept()
	if err != nil {
		log.Fatal(err)
	}

	for {
		go handleRequest(conn)
	}
	conn.Close()
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("La reponse gere la requete")

	conn.Write([]byte(parseBuffer(buf)))

}

func parseBuffer(buf []byte) string {
	str := string(buf)
	slice := strings.Split(str, " ")

	switch slice[0] {
	case "CREATE":
		return createEvent(slice)
	case "CLOSE":
		return closeEvent(slice)
	case "ADD":
		return addBenevole(slice)
	case "LISTM":
		return listEvents()
	case "LISTP":
		return listPosts(slice)
	default:
		return "Command not found"
	}
}

func createUsers() {
	events = append(events, Event{0, "Festival de la bière", 1, true})
	eventCounter++
	posts = append(posts, Post{0, "Post1", 3, 1})
	postCounter++
	posts = append(posts, Post{1, "Post2", 2, 1})
	postCounter++
	users = append(users, User{"Ficelle", "1234", Post{}})
	users = append(users, User{"Taro", "1234", Post{}})
	users = append(users, User{"Catanne", "1234", Post{}})
	users = append(users, User{"Willi", "1234", Post{}})
}

func authentification(username string, password string) bool {
	for _, user := range users {
		if user.name == username && user.password == password {
			return true
		}
	}
	return false
}

func createEvent(slice []string) string {
	fmt.Println("Starting an event")
	if authentification(slice[1], slice[2]) {
		for i := 3; i < len(slice)-1; i++ {
			capacity, _ := strconv.Atoi(slice[i+1])
			posts = append(posts, Post{postCounter, slice[i], capacity, eventCounter})
			postCounter++
		}
		events = append(events, Event{eventCounter, slice[3], eventCounter, true})
		eventCounter++
		return "Event created"
	} else {
		return "Authentication failed"
	}
}

func closeEvent(slice []string) string {
	fmt.Println("Closing an event")
	if authentification(slice[1], slice[2]) {
		for i := 0; i < len(events); i++ {
			id, _ := strconv.Atoi(slice[3])
			if events[i].id == id {
				events[i].isOpen = false
				return "Event closed"
			}
		}
		return "Event not found"
	} else {
		return "Authentication failed"
	}
}

func addBenevole(slice []string) string {
	fmt.Println("Adding a benevole")
	if authentification(slice[1], slice[2]) {
		for i := 0; i < len(posts); i++ {
			idEvent, _ := strconv.Atoi(slice[3])
			idPost, _ := strconv.Atoi(slice[4])
			if posts[i].id == idPost && posts[i].eventId == idEvent && posts[i].capacity > 0 {
				// todo change the idpost of the user
				posts[i].capacity--
				return "Benevole added"
			}
		}
		return "Post not found"
	} else {
		return "Authentication failed"
	}
}

func listEvents() string {
	var str string
	for i := 0; i < len(events); i++ {
		if events[i].isOpen {
			str += "Event's id: " + strconv.Itoa(events[i].id) + " Event's name: " + events[i].name + " OwnerID: " + strconv.Itoa(events[i].ownerId) + " is open:" + strconv.FormatBool(events[i].isOpen)
		}
	}
	return str
}

func listPosts(slice []string) string {
	var str string
	for i := 0; i < len(posts); i++ {
		idEvent, _ := strconv.Atoi(slice[1])
		if posts[i].eventId == idEvent {
			str += "Post's id: " + strconv.Itoa(posts[i].id) + " Post's name: " + posts[i].name + " Capacity: " + strconv.Itoa(posts[i].capacity) + " Event's id: " + strconv.Itoa(posts[i].eventId)
			// todo do a array with the postID and the user name
		}
	}
	return str
}
