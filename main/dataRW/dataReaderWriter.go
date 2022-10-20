package dataRW

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var DataChannel = make(chan chan []byte)

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

var eventCounter = 0
var postCounter = 0

var events []Event
var posts []Post
var users []User

func createUsersAndEvents() {
	// creation of users
	users = append(users, User{"Bob", "1234"})
	users = append(users, User{"Lea", "1234"})
	users = append(users, User{"Leo", "1234"})
	users = append(users, User{"Willi", "1234"})
	users = append(users, User{"Lili", "12345"})
	users = append(users, User{"T", "1234"})
	// creation of posts
	posts = append(posts, Post{postCounter, "Bar à bière", 3, 0, users[0:1]})
	postCounter++
	posts = append(posts, Post{postCounter, "Securité", 2, 0, users[2:4]})
	postCounter = 0
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

func authentification(username string, password string) bool {
	for _, user := range users {
		if user.name == username && user.password == password {
			return true
		}
	}
	return false
}

func removeUserPost(username string, password string, idEvent string) {
	event := getEventById(idEvent)
	for _, post := range event.posts {
		for i, staff := range post.staff {
			if staff.name == username && staff.password == password {
				fmt.Println("Removing user from post")
				post.staff = append(post.staff[:i], post.staff[i+1:]...)
				post.capacity++
			}
		}
	}
}

func getEventById(id string) Event {
	for i := 0; i < len(events); i++ {
		idEvent, _ := strconv.Atoi(string(bytes.Trim([]byte(id), "\x00")))
		if events[i].id == idEvent {
			return events[i]
		}
	}
	return Event{}
}

func contains(users []User, person User) bool {
	for _, a := range users {
		if a == person {
			return true
		}
	}
	return false
}

func createEvent(parameters []string) string {
	uname := parameters[0]
	pwd := parameters[1]
	eventName := parameters[2]
	posts := parameters[3:]
	if len(posts)%2 != 0 {
		return "Invalid Number of arguments"
	}
	if authentification(uname, pwd) {
		postCounter = 0
		owner := User{uname, pwd}
		var newPost []Post
		for i := 0; i < len(posts)-1; i += 2 {
			capacity, _ := strconv.Atoi(string(bytes.Trim([]byte(posts[i+1]), "\x00")))
			newPost = append(newPost, Post{postCounter, posts[i], capacity, eventCounter, nil})
			postCounter++
		}
		events = append(events, Event{eventCounter, eventName, owner, true, newPost})
		eventCounter++
		return "Event Created"
	}
	return "Authentification failed"
}

func closeEvent(commandParameters []string) string {
	fmt.Println("Closing an event")
	if authentification(commandParameters[0], commandParameters[1]) {
		for i := 0; i < len(events); i++ {
			id, _ := strconv.Atoi(string(bytes.Trim([]byte(commandParameters[2]), "\x00")))
			if events[i].id == id && events[i].owner.name == commandParameters[0] {
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
	uname := slice[0]
	pwd := slice[1]
	idEvent := slice[2]
	idPost := slice[3]
	if authentification(slice[0], slice[1]) {
		idPost, _ := strconv.Atoi(string(bytes.Trim([]byte(idPost), "\x00")))
		//removeUserPost(uname, pwd, idEvent)
		event := getEventById(idEvent)
		post := getEventById(idEvent).posts[idPost]
		if post.capacity < 1 {
			return "Could not add user to post because post is full"
		}
		var staff []User
		copy(post.staff, staff)
		if !contains(post.staff, User{uname, pwd}) {
			for _, user := range post.staff {
				staff = append(staff, user)
			}
			event.posts[idPost].staff = append(staff, User{uname, pwd})
		}
		return "User successfully added to post"
	}
	return "Authentication failed"
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
		idEvent, _ := strconv.Atoi(string(bytes.Trim([]byte(slice[0]), "\x00")))
		if posts[i].eventId == idEvent {
			str += "Post's id: " + strconv.Itoa(posts[i].id) + ", Post's name: " + posts[i].name + ", Capacity: " + strconv.Itoa(posts[i].capacity) + "\n"
		}
	}
	return str
}

func listUsers(slice []string) string {
	var event = getEventById(slice[0])

	tabCell := "%-20v"
	tabCellCross := "%-8v"
	firstColumn := "%-25v" // Line label
	header := fmt.Sprintf(firstColumn, event.name) + "|"
	nbInscrit := fmt.Sprintf(firstColumn, "nbInscrit") + "|"
	tab := ""

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

func processCommand(commandParameters []string) string {
	fmt.Println("Processing in dataRW")
	switch commandParameters[0] {
	case "CREATE":
		return createEvent(commandParameters[1:])
	case "CLOSE":
		return closeEvent(commandParameters[1:])
	case "ADD":
		return addBenevole(commandParameters[1:])
	case "LISTM":
		return listEvents()
	case "LISTP":
		return listPosts(commandParameters[1:])
	case "LISTU":
		return listUsers(commandParameters[1:])
	case "QUIT":
		return "q"
	default:
		return "Command not found"
	}
}

func HandleRWActions() {
	createUsersAndEvents()
	for {
		clientChan := <-DataChannel
		// Process request
		command := <-clientChan
		clientChan <- []byte(processCommand(strings.Split(string(command), " ")))
	}
}
