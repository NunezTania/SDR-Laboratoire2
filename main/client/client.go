package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	helpMenu()
	conn, err := net.Dial(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Println("Please enter a command:")
		var command string
		Scanner := bufio.NewScanner(os.Stdin)
		Scanner.Scan()
		if err := Scanner.Err(); err != nil {
			log.Fatal(err)
		}
		command = Scanner.Text()
		if !processCommands(&command) {
			if command == "QUIT " {
				break
			}
			continue
		}

		_, err = conn.Write([]byte(command))
		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("La reponse est : \n" + string(buf))
	}
	errClosedConn := conn.Close()
	if err != nil {
		log.Fatal(errClosedConn)
	}
}

func helpMenu() {
	fmt.Println("Hello and Welcome to the Event Manager")
	fmt.Println("Here is a list of possible commands:")
	fmt.Println("1. CREATE")
	fmt.Println("2. CLOSE")
	fmt.Println("3. ADD")
	fmt.Println("4. LISTM")
	fmt.Println("5. LISTP")
	fmt.Println("6. LISTU")
	fmt.Println("7. QUIT")
}

func getUserAnswer() string {
	var answer string
	Scanner := bufio.NewScanner(os.Stdin)
	Scanner.Scan()
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	answer = Scanner.Text()
	return answer
}

func authentificationProcess() string {
	fmt.Println("Please enter your username:")
	var username string
	Scanner := bufio.NewScanner(os.Stdin)
	Scanner.Scan()
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	username = Scanner.Text()

	fmt.Println("Please enter your password:")
	var password string
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	password = string(bytePassword)
	return username + " " + password
}

func processCommands(command *string) bool {
	switch *command {
	case "CREATE":
		return processCreate(command)
	case "CLOSE":
		return processClose(command)
	case "ADD":
		return processAdd(command)
	case "LISTM":
		*command += " "
		return true
	case "LISTP":
		return processList(command)
	case "LISTU":
		return processList(command)
	case "QUIT":
		*command += " "
		return true
	default:
		fmt.Println("Command not recognized")
		return false
	}
}

func processCreate(command *string) bool {
	login := authentificationProcess()
	fmt.Println("Enter the name of the event:")
	eventName := getUserAnswer()
	fmt.Println("Enter the posts:")
	posts := getUserAnswer()
	if !checkPosts(posts) || eventName == "" {
		return false
	}
	*command += " " + login + " " + eventName + " " + posts
	//fmt.Println("Voila la commande formé " + *command)
	return true
}

// todo is it possible to have no post in an event?
func checkPosts(command string) bool {
	posts := strings.Split(command, " ")

	// check that there is at least 2 words
	if len(posts) < 2 {
		fmt.Println("Command not recognized 1")
		return false
	}
	// check that there is a pair number of words
	if len(posts)%2 != 0 {
		fmt.Println("Command not recognized 2")
		return false
	}
	// check that the capacity is a number
	for i := 1; i < len(posts)-1; i += 2 {
		if !isNumber(posts[i]) {
			fmt.Println("Command not recognized 3")
			return false
		}
	}
	return true
}

func processClose(command *string) bool {
	login := authentificationProcess()
	fmt.Println("Enter the event id")
	eventId := getUserAnswer()
	if !isNumber(eventId) {
		return false
	}
	*command += " " + login + " " + eventId
	return true
}

func isNumber(x string) bool {
	_, err := strconv.Atoi(x)
	return err == nil
}

func processAdd(command *string) bool {
	login := authentificationProcess()
	fmt.Println("Enter the event id:")
	eventId := getUserAnswer()
	fmt.Println("Enter the post id:")
	postId := getUserAnswer()
	if !isNumber(eventId) || !isNumber(postId) {
		return false
	}
	*command += " " + login + " " + eventId + " " + postId
	return true
}

func processList(command *string) bool {
	fmt.Println("Enter the event id:")
	eventId := getUserAnswer()
	if !isNumber(eventId) {
		return false
	}
	*command += " " + eventId
	return true
}
