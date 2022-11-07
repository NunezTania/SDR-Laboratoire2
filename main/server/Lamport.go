package main

type time uint64

type Lamport struct {
	counterTime time
}

// NewLamport Create a new Lamport clock
func (l *Lamport) NewLamport() *Lamport {
	l.counterTime = 0
	return l
}

// GetTime Get the current time of the Lamport clock
func (l *Lamport) GetTime() time {
	return l.counterTime
}

// Increment the Lamport clock
func (l *Lamport) Increment() time {
	l.counterTime++
	return l.counterTime
}

// Update the Lamport clock
func (l *Lamport) Update(otherTime time) time {
	if otherTime > l.counterTime {
		l.counterTime = otherTime
	}
	l.Increment()
	return l.counterTime
}
