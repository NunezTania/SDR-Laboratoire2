// Package client contains the client of the program.
// It is used to communicate with the server.
// And can be used to create, close, add and list events.
package main

import (
	pm "SDR-Laboratoire1/main/server/processMutex"
	"bufio"
	"fmt"
	"golang.org/x/term"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// main function of the program client
func main() {
	RunRandomServ()
}

// RunRandomServ connects the client to a random server
func RunRandomServ() {
	Run(-1)
}

// Run connects the client to a server with the id given in parameter
// if the id is -1, it connects to a random server
func Run(idServ int) {
	HelpMenu()
	conf := pm.Config
	if idServ == -1 {
		rand.Seed(time.Now().UnixNano())
		idServ = rand.Intn(conf.NServ)
	}

	for {
		fmt.Println("Please enter a command:")
		var command string
		Scanner := bufio.NewScanner(os.Stdin)
		Scanner.Scan()
		errRead := Scanner.Err()
		for errRead != nil {
			fmt.Println("Error while reading the command, trying again...")
			Scanner = bufio.NewScanner(os.Stdin)
			Scanner.Scan()
			errRead = Scanner.Err()
		}

		command = Scanner.Text()

		if !ProcessCommands(&command) {
			continue
		}

		if command == "QUIT " {
			break
		}

		conn, err := net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+idServ))
		for err != nil {
			fmt.Println("Error while connecting to the server, trying again...")
			conn, errRead = net.Dial(conf.Type, conf.Host+":"+strconv.Itoa(conf.PortClient+idServ))
		}

		_, errRead = conn.Write([]byte(command))
		for errRead != nil {
			fmt.Println("Error while sending the command, trying again...")
			_, errRead = conn.Write([]byte(command))
		}

		buf := make([]byte, 1024)
		_, errRead = conn.Read(buf)
		for errRead != nil {
			fmt.Println("Error while reading the response, trying again...")
			_, errRead = conn.Read(buf)
		}

		fmt.Println("La reponse est : \n" + string(buf))

		errClosedConn := conn.Close()
		if errClosedConn != nil {
			fmt.Println("Error while closing the connection")
		}
	}
}

// HelpMenu displays all the commands that can be used by the user
func HelpMenu() {
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

// GetUserAnswer returns the answer of the user
func GetUserAnswer() string {
	var answer string
	Scanner := bufio.NewScanner(os.Stdin)
	Scanner.Scan()
	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
	answer = Scanner.Text()
	return answer
}

// AuthentificationProcess returns the login and the password of the user
func AuthentificationProcess() string {
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

// ProcessCommands returns true if the command is valid
func ProcessCommands(command *string) bool {
	switch *command {
	case "CREATE":
		return ProcessCreate(command)
	case "CLOSE":
		return ProcessClose(command)
	case "ADD":
		return ProcessAdd(command)
	case "LISTM":
		*command += " "
		return true
	case "LISTP":
		return ProcessList(command)
	case "LISTU":
		return ProcessList(command)
	case "QUIT":
		*command += " "
		return true
	default:
		fmt.Println("Command not recognized")
		return false
	}
}

// ProcessCreate process the creation of a new event returns true if the command is valid
func ProcessCreate(command *string) bool {
	login := AuthentificationProcess()
	fmt.Println("Enter the name of the event:")
	eventName := GetUserAnswer()
	fmt.Println("Enter the posts:")
	posts := GetUserAnswer()
	if !CheckPosts(posts) || eventName == "" {
		return false
	}
	*command += " " + login + " " + eventName + " " + posts
	return true
}

// CheckPosts verify that the posts created are valid and returns true if they are
func CheckPosts(command string) bool {
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

// ProcessClose process the closure of an event returns true if the command is valid
func ProcessClose(command *string) bool {
	login := AuthentificationProcess()
	fmt.Println("Enter the event id")
	eventId := GetUserAnswer()
	if !isNumber(eventId) {
		return false
	}
	*command += " " + login + " " + eventId
	return true
}

// isNumber returns true if the string is a number
func isNumber(x string) bool {
	_, err := strconv.Atoi(x)
	return err == nil
}

// ProcessAdd process the addition of a staff in an event returns true if the command is valid
func ProcessAdd(command *string) bool {
	login := AuthentificationProcess()
	fmt.Println("Enter the event id:")
	eventId := GetUserAnswer()
	fmt.Println("Enter the post id:")
	postId := GetUserAnswer()
	if !isNumber(eventId) || !isNumber(postId) {
		return false
	}
	*command += " " + login + " " + eventId + " " + postId
	return true
}

// ProcessList process the listing command returns true if the command is valid
func ProcessList(command *string) bool {
	fmt.Println("Enter the event id:")
	eventId := GetUserAnswer()
	if !isNumber(eventId) {
		return false
	}
	*command += " " + eventId
	return true
}
