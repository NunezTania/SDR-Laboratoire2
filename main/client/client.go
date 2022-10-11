package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	helpMenu()

	// connection to the server
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
		if !checkCommands(command) {
			return
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

		fmt.Println("La reponse est : " + string(buf))
	}
	conn.Close()
}

func helpMenu() {
	fmt.Println("Hello and Welcome to the Event Manager")
	fmt.Println("Here is a list of possible commands:")
	fmt.Println("1. CREATE userName userPassword eventName postName1 capacity1 postName2 capacity2 ... ")
	fmt.Println("2. CLOSE eventName")
	fmt.Println("3. ADD eventId postId")
	fmt.Println("4. List all the event")
	fmt.Println("5. List all the posts of a event")
}

func checkCommands(command string) bool {
	cmdStart := strings.SplitAfter(command, " ")[0]

	switch cmdStart {
	case "CREATE ":
		return checkCreate(command)
	case "CLOSE ":
		return checkClose(command)
	case "ADD ":
		return checkAdd(command)
	case "LISTM ":
		return checkListM(command)
	case "LISTP ":
		return checkListP(command)
	default:
		fmt.Println("Command not recognized")
		return false
	}
}

func checkListP(command string) bool {
	return true
}

func checkListM(command string) bool {
	return true
}

func checkAdd(command string) bool {
	return true
}

func checkClose(command string) bool {
	return true
}

func checkCreate(command string) bool {
	// check that there is at least 4 words
	if len(strings.SplitAfter(command, " ")) < 4 {
		fmt.Println("Command not recognized 1")
		return false
	}
	// check that there is a pair number of words
	if (len(strings.SplitAfter(command, " ")))%2 != 0 {
		fmt.Println("Command not recognized 2")
		return false
	}

	commandRest := strings.SplitAfter(command, " ")[4:]
	for i := 0; i < len(commandRest)-1; i += 2 {
		if len(commandRest[i]) < 1 && len(commandRest[i+1]) < 1 {
			fmt.Println("Command not recognized 3")
			return false
		}
		// todo verify that the capacity is a number
		/*
			fmt.Println("DEBUG " + commandRest[i+1] + "result")
			_, err := strconv.Atoi(commandRest[i+1])
			if err != nil {
				fmt.Println("Command not recognized 4")
				return false
			}
		*/
	}
	return true
}
