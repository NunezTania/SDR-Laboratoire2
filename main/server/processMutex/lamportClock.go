// The Lamport Clock is a logical clock that is used to order events in a distributed system.
package processMutex

import "strconv"

type Lamport struct {
	counterTime int
}

// NewLamport Create a new Lamport clock and initialize it to 0
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

// strToLamport Convert a string to a Lamport clock
func strToLamport(str string) Lamport {
	var clock Lamport
	clock.counterTime, _ = strconv.Atoi(str)
	return clock
}
