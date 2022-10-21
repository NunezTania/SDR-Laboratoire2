package main

import (
	"SDR-Laboratoire1/main/client"
	"SDR-Laboratoire1/main/server"
)

func main() {
	go server.Run()
	client.Run()
}
