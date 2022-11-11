package main

import "fmt"

var clock = Lamport{}

type Request struct {
	rType string
	time  Lamport
	id    int
}

// tableau des process

func main() {
	clock.NewLamport()
	fmt.Println(clock.GetTime())
}

// check if SC is available
