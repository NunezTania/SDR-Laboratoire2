package main

// A faire :
// todo trouver pourquoi le server se deconnecte parfois
// todo ecrire les commentaires
// todo ecriture des tests
// todo terminer la fonction addBenevole()

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
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

func main() {
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

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	for commandTreatment := parseBuffer(buf); commandTreatment != "q"; commandTreatment = parseBuffer(buf) {
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
	_, writeErr := conn.Write([]byte("Bye bye, Xoxo"))
	if writeErr != nil {
		log.Fatal(writeErr)
	}
	closeErr := conn.Close()
	if closeErr != nil {
		log.Fatal(closeErr)
	}
}

func parseBuffer(buf []byte) string {
	slice := strings.Split(string(buf), " ")
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
	case "LISTU":
		return listUsers(slice)
	case "QUIT":
		return "q"
	default:
		return "Command not found"
	}
}

func createUsersAndEvents() {
	// creation of users
	users = append(users, User{"Bob", "1234"})
	users = append(users, User{"Lea", "1234"})
	users = append(users, User{"Leo", "1234"})
	users = append(users, User{"Willi", "1234"})
	users = append(users, User{"Lili", "1234"})
	// creation of posts
	posts = append(posts, Post{postCounter, "Bar à bière", 3, eventCounter, users[0:1]})
	postCounter++
	posts = append(posts, Post{postCounter, "Securité", 2, eventCounter, users[2:4]})
	postCounter++
	posts = append(posts, Post{postCounter, "Vente de ticket", 5, eventCounter, users[0:1]})
	postCounter++
	posts = append(posts, Post{postCounter, "Logistique", 1, eventCounter, users[2:4]})
	postCounter++
	posts = append(posts, Post{postCounter, "Securité", 2, eventCounter, users[2:4]})
	postCounter++
	// creation of events
	events = append(events, Event{eventCounter, "Festival de la musique", users[0], true, posts[0:2]})
	eventCounter++
	events = append(events, Event{eventCounter, "Festival de la bière", users[0], true, posts[2:4]})
	eventCounter++
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
		owner := User{slice[1], slice[2]}
		var newPost []Post
		for i := 3; i < len(slice)-1; i++ {
			capacity, _ := strconv.Atoi(slice[i+1])
			newPost = append(newPost, Post{postCounter, slice[i], capacity, eventCounter, nil})
			postCounter++
		}
		posts = append(posts, newPost...)
		events = append(events, Event{eventCounter, slice[3], owner, true, newPost})
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
			if events[i].id == id && events[i].owner.name == slice[1] {
				events[i].isOpen = false
				return "Event closed"
			}
		}
		return "Event couldn't be closed"
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

			if posts[i].id == idPost && posts[i].eventId == idEvent && posts[i].capacity > 0 && getEventById(slice[3]).isOpen {
				// todo check if user is already in a post of this event
				// if it's the case, erase the user from the old post
				posts[i].staff = append(posts[i].staff, User{slice[1], slice[2]})
				posts[i].capacity--
				// on cherche si le benevole ne fait pas deja partie de l'event
				// on itere sur tout les poste du festival
				for _, post := range getEventById(slice[3]).posts {
					/*
						if (post.staff.contains(User{slice[1], slice[2]})) {
							post.staff.remove(User{slice[1], slice[2]})
						}
					*/
					fmt.Println(post)

				}
			}

			return "Benevole added"
		}
	} else {
		return "Authentication failed"
	}
	return "Benevole couldn't be added"
}

func listEvents() string {
	var str string
	for i := 0; i < len(events); i++ {
		if events[i].isOpen {
			str += "Event's id: " + strconv.Itoa(events[i].id) + ", Event's name: " + events[i].name + ", Owner: " + events[i].owner.name + ", is open:" + strconv.FormatBool(events[i].isOpen) + "\n"
		}
	}
	return str
}

func listPosts(slice []string) string {
	var str string
	for i := 0; i < len(posts); i++ {
		idEvent, _ := strconv.Atoi(slice[1])
		if posts[i].eventId == idEvent {
			str += "Post's id: " + strconv.Itoa(posts[i].id) + ", Post's name: " + posts[i].name + ", Capacity: " + strconv.Itoa(posts[i].capacity) + "\n"
		}
	}
	return str
}

/*
	Event   | postId 1| postId 2| postId 3|

nbInscrit |    1    |    2    |    3    |
user 1    |    x    |         |         |
user 2    |		    |         |    x    |
*/
func listUsers(slice []string) string {
	var event = getEventById(slice[1])

	var tabCell string = "%-15v"
	var tabCellCross string = "%-8v"
	var firstColumn string = "%-25v" // Line label
	var header string = fmt.Sprintf(firstColumn, event.name) + "|"
	var nbInscrit string = fmt.Sprintf(firstColumn, "nbInscrit") + "|"
	var tab string = ""

	for i, post := range event.posts {
		header += fmt.Sprintf(tabCell, post.name+" "+strconv.Itoa(post.id))
		header += "|"
		nbInscrit += fmt.Sprintf(tabCell, len(post.staff))
		nbInscrit += "|"
		for _, user := range post.staff {
			tab += fmt.Sprintf(firstColumn, user.name) + "|"
			for j := 0; j < len(event.posts); j++ {
				if j == i {
					tab += fmt.Sprintf(tabCellCross, "x") + "       "
				} else {
					tab += fmt.Sprintf(tabCellCross, "") + "       "
				}
				tab += "|"
			}
			tab += "\n"
		}
	}
	header += "\n"
	nbInscrit += "\n"

	return header + nbInscrit + tab
}

func getEventById(id string) Event {
	for i := 0; i < len(events); i++ {
		idEvent, _ := strconv.Atoi(id)
		if events[i].id == idEvent {
			return events[i]
		}
	}
	return Event{}
}
