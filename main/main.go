package main

import (
	"SDR-Laboratoire1/main/client"
	"SDR-Laboratoire1/main/server"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println(os.Args[0])
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Missing program to execute argument.\n Try adding client or server to go run command. Ex: go run main.go client")
	}
	if strings.Compare(args[0], "client") == 0 {
		client.Run()
	} else if strings.Compare(args[0], "server") == 0 {
		server.Run()
	}
}
