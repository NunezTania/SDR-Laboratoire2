package main

import "strconv"

type Lamport struct {
	counterTime int
}

// NewLamport Create a new Lamport clock
func (l *Lamport) NewLamport() *Lamport {
	l.counterTime = 0
	return l
}

// GetTime Get the current time of the Lamport clock
func (l *Lamport) GetTime() int {
	return l.counterTime
}

// Increment the Lamport clock
func (l *Lamport) Increment() int {
	l.counterTime++
	return l.counterTime
}

// Update the Lamport clock
func (l *Lamport) Update(otherTime Lamport) Lamport {
	if otherTime.counterTime > l.counterTime {
		l.counterTime = otherTime.counterTime
	}
	l.Increment()
	return *l
}

func strToLamport(str string) Lamport {
	var clock Lamport
	clock.counterTime, _ = strconv.Atoi(str)
	return clock
}
