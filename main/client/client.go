package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var servAddr = "127.0.0.1:8080"

func RunClient(id string) {
	conn, err := net.Dial("tcp", servAddr)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan string)
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- "bsus" // signal the main goroutine
	}()
	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Fatal(err)
	}
	conn.Close()
	<-done // wait for background goroutine to finish
}

func main() {
	fmt.Println("Coucou")
}
