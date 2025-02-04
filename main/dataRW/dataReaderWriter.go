// Package dataRW contains the functions to read and write data
// It allows to create events, close events, list events and add users to Posts
// All of these actions are concurrency safe
package dataRW

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Event struct {
	Id     int
	Name   string
	Owner  User
	IsOpen bool
	Posts  []Post
}

type Post struct {
	Id       int
	Name     string
	Capacity int
	EventId  int
	Staff    []User
}

type User struct {
	Name     string
	Password string
}

// Authentification checks if the user is in the list of users and Password is correct
func Authentification(username string, password string, users *[]User) bool {
	for _, user := range *users {
		if user.Name == username && user.Password == password {
			return true
		}
	}
	return false
}

// RemoveUserPost removes a user from a post
func RemoveUserPost(username string, password string, idEvent string, events *[]Event) {
	event := GetEventById(idEvent, events)
	for _, post := range event.Posts {
		for i, staff := range post.Staff {
			if staff.Name == username && staff.Password == password {
				(*events)[event.Id].Posts[post.Id].Staff = append((*events)[event.Id].Posts[post.Id].Staff[:i], (*events)[event.Id].Posts[post.Id].Staff[i+1:]...)
				post.Capacity++
			}
		}
	}
}

// GetEventById returns the event with the given Id
func GetEventById(id string, events *[]Event) Event {
	for i := 0; i < len(*events); i++ {
		idEvent, _ := strconv.Atoi(string(bytes.Trim([]byte(id), "\x00")))
		if (*events)[i].Id == idEvent {
			return (*events)[i]
		}
	}
	var event Event
	event.Id = -1
	return event
}

// contains check if a person is contains inside users
func contains(users []User, person User) bool {
	for _, a := range users {
		if a == person {
			return true
		}
	}
	return false
}

// CreateEvent creates an event
func CreateEvent(parameters []string, users *[]User, events *[]Event, postCounter *int, eventCounter *int) string {
	uname := parameters[0]
	pwd := parameters[1]
	eventName := parameters[2]
	posts := parameters[3:]
	if len(posts)%2 != 0 {
		return "Invalid Number of arguments"
	}
	if Authentification(uname, pwd, users) {
		*postCounter = 0
		owner := User{uname, pwd}
		var newPost []Post
		for i := 0; i < len(posts)-1; i += 2 {
			capacity, _ := strconv.Atoi(string(bytes.Trim([]byte(posts[i+1]), "\x00")))
			newPost = append(newPost, Post{*postCounter, posts[i], capacity, *eventCounter, nil})
			*postCounter++
		}
		*events = append(*events, Event{*eventCounter, eventName, owner, true, newPost})
		*eventCounter++
		return "Event Created"
	}
	return "Authentification failed"
}

// CloseEvent closes an event
func CloseEvent(commandParameters []string, users *[]User, events *[]Event) string {
	if Authentification(commandParameters[0], commandParameters[1], users) {
		for i := 0; i < len(*events); i++ {
			id, _ := strconv.Atoi(string(bytes.Trim([]byte(commandParameters[2]), "\x00")))
			if (*events)[i].Id == id && (*events)[i].Owner.Name == commandParameters[0] {
				(*events)[i].IsOpen = false
				return "Event closed"
			}
		}
		return "Event couldn't be closed"
	} else {
		return "Authentication failed"
	}
}

// AddBenevole adds a benevole to a post
func AddBenevole(slice []string, users *[]User, events *[]Event) string {
	uname := slice[0]
	pwd := slice[1]
	idEvent := slice[2]
	idPost := slice[3]
	if Authentification(slice[0], slice[1], users) {
		idPost, err := strconv.Atoi(string(bytes.Trim([]byte(idPost), "\x00")))
		if err != nil {
			return "Invalid idPost"
		}
		RemoveUserPost(uname, pwd, idEvent, events)
		evId, err2 := strconv.Atoi(string(bytes.Trim([]byte(idEvent), "\x00")))
		if err2 != nil {
			return "Invalid idEvent"
		}
		event := GetEventById(idEvent, events)
		if !event.IsOpen {
			return "Event is closed"
		}
		if len(*events) < evId {
			return "Event not found"
		}
		if len(event.Posts) < idPost {
			return "Post not found"
		}
		post := GetEventById(idEvent, events).Posts[idPost]
		if post.Capacity < len(post.Staff)+1 {
			return "Could not add user to post because post is full"
		}
		var staff []User
		copy(post.Staff, staff)
		if !contains(post.Staff, User{uname, pwd}) {
			for _, user := range post.Staff {
				staff = append(staff, user)
			}
			event.Posts[idPost].Staff = append(staff, User{uname, pwd})
		}
		return "User successfully added to post"
	}
	return "Authentication failed"
}

// ListEvents lists all the events
func ListEvents(events *[]Event) string {
	var str string
	for i := 0; i < len(*events); i++ {
		str += "Event's Id: " + strconv.Itoa((*events)[i].Id) + ", Event's Name: " + (*events)[i].Name + ", Owner: " + (*events)[i].Owner.Name + ", is open:" + strconv.FormatBool((*events)[i].IsOpen) + "\n"
	}
	return str
}

// ListPosts lists all the Posts of an event
func ListPosts(slice []string, events *[]Event) string {
	event := GetEventById(slice[0], events)
	var str string
	for i := 0; i < len(event.Posts); i++ {
		str += "Post's Id: " + strconv.Itoa(event.Posts[i].Id) + ", Post's Name: " + event.Posts[i].Name + ", Capacity: " + strconv.Itoa(event.Posts[i].Capacity) + "\n"
	}
	return str
}

// ListUsers lists all the users of all the Posts of an event
func ListUsers(slice []string, events *[]Event) string {
	var event = GetEventById(slice[0], events)

	tabCell := "%-20v"
	tabCellCross := "%-10v"
	firstColumn := "%-25v" // Line label
	header := fmt.Sprintf(firstColumn, event.Name) + "|"
	nbInscrit := fmt.Sprintf(firstColumn, "nbInscrit") + "|"
	tab := ""

	for i, post := range event.Posts {
		header += fmt.Sprintf(tabCell, post.Name+" "+strconv.Itoa(post.Id))
		header += "|"
		nbInscrit += fmt.Sprintf(tabCell, len(post.Staff))
		nbInscrit += "|"
		for _, user := range post.Staff {
			tab += fmt.Sprintf(firstColumn, user.Name) + "|"
			for j := 0; j < len(event.Posts); j++ {
				if j == i {
					tab += fmt.Sprintf(tabCellCross, "x") + "          "
				} else {
					tab += fmt.Sprintf(tabCellCross, "") + "          "
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

// ProcessCommand find the first word (the command) of the command and call the corresponding function
func ProcessCommand(commandParameters []string, users *[]User, events *[]Event, DataModified *bool, postCounter *int, eventCounter *int) string {
	switch commandParameters[0] {
	case "CREATE":
		*DataModified = true
		return CreateEvent(commandParameters[1:], users, events, postCounter, eventCounter)
	case "CLOSE":
		*DataModified = true
		return CloseEvent(commandParameters[1:], users, events)
	case "ADD":
		*DataModified = true
		return AddBenevole(commandParameters[1:], users, events)
	case "LISTM":
		*DataModified = false
		return ListEvents(events)
	case "LISTP":
		*DataModified = false
		return ListPosts(commandParameters[1:], events)
	case "LISTU":
		*DataModified = false
		return ListUsers(commandParameters[1:], events)
	case "QUIT":
		*DataModified = false
		return "q"
	default:
		return "Command not found"
	}
}

// HandleRWActions handles the read/write actions
func HandleRWActions(DataChannel *chan chan []byte, DataModified *bool, users *[]User, events *[]Event, postCounter *int, eventCounter *int) {

	//CreateUsersAndEvents(users, events, postCounter, eventCounter)
	for {
		// Blocking eventual other requests for concurrent data access
		clientChan := <-*DataChannel
		// Process request
		command := <-clientChan
		clientChan <- []byte(ProcessCommand(strings.Split(string(command), " "), users, events, DataModified, postCounter, eventCounter))
	}
}
