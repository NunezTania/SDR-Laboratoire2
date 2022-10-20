package dataRW

import (
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

func authentification(username string, password string) bool {
	for _, user := range users {
		if user.name == username && user.password == password {
			return true
		}
	}
	return false
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
		owner := User{uname, pwd}
		var newPost []Post
		for i := 0; i < len(posts)-1; i += 2 {
			//var capa = slice[i+1]
			capacity, _ := strconv.Atoi(posts[i+1])
			newPost = append(newPost, Post{postCounter, posts[i], capacity, eventCounter, nil})
			postCounter++
		}
		events = append(events, Event{eventCounter, eventName, owner, true, newPost})
		eventCounter++
	}
	return "Authentification failed"
}

func closeEvent(commandParameters []string) string {
	fmt.Println("Closing an event")
	if authentification(commandParameters[1], commandParameters[2]) {
		for i := 0; i < len(events); i++ {
			var idE = commandParameters[3]
			id, _ := strconv.Atoi(idE)
			if events[i].id == id && events[i].owner.name == commandParameters[1] {
				events[i].isOpen = false
				return "Event closed"
			}
		}
		return "Event couldn't be closed"
	}
	return "Authentification failed"
}

func addBenevole(slice []string) string {
	fmt.Println("Adding a benevole")
	if authentification(slice[1], slice[2]) {
		for i := 0; i < len(posts); i++ {
			idEvent, _ := strconv.Atoi(slice[3])
			idPost, _ := strconv.Atoi(slice[4])

			if posts[i].id == idPost && posts[i].eventId == idEvent && posts[i].capacity > 0 && getEventById(slice[3]).isOpen {
				for i, post := range getEventById(slice[3]).posts {
					if (contains(post.staff, User{slice[1], slice[2]})) {
						post.staff[i] = post.staff[len(post.staff)-1]
						post.staff = post.staff[:len(post.staff)-1]
					}
				}
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
		return "Couldn't add benevole to this post"
	} else {
		return "Authentication failed"
	}
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

func processCommand(commandParameters []string) string {
	switch commandParameters[0] {
	case "CREATE":
		return createEvent(commandParameters[1:])
	case "CLOSE":
		return closeEvent(commandParameters)
	case "ADD":
		return addBenevole(commandParameters)
	case "LISTM":
		return listEvents()
	case "LISTP":
		return listPosts(commandParameters)
	case "LISTU":
		return listUsers(commandParameters)
	case "QUIT":
		return "q"
	default:
		return "Command not found"
	}
}

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

func HandleRWActions() {
	createUsersAndEvents()
	for {
		clientChan := <-DataChannel
		// Process request
		s := <-clientChan
		command := strings.Split(string(s), " ")
		for _, keks := range command {
			fmt.Println(keks)
		}
		clientChan <- []byte("Bonjour vous aussi :)")
	}
}
