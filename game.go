package main

import (
	"fmt"
	"math"
	"os"
	"sync/atomic"
	"time"
)

// ok we need to use milliseconds
const milli = 1000
const updatesPerSecond = 20

// tick rate in updates per second
func game() {

	defer os.Exit(0)

	tick := 0
	previous := 0
	sleepMs := 1 / float64(updatesPerSecond) * milli
	// for checking that the time is working properly
	counter := float64(0)

	// TODO decouple from game. realising i need to improve my async/threading/routines knowledge

	for atomic.LoadInt64(&kill) != 1 {
		now := int(time.Now().UTC().UnixMilli())

		// ms since last update
		deltatime := now - previous
		previous = now

		// delta value
		delta := sleepMs / float64(deltatime)
		counter += delta / updatesPerSecond

		tick++
		sleepBy := time.Millisecond * time.Duration(sleepMs)

		// takes nanos as args
		time.Sleep(sleepBy)

		// log every second
		if tick%updatesPerSecond == 0 {
			// TODO broadcast to connected users for testing purposes
			writeMessageString(fmt.Sprintf("tick: %d, counter: %f (%f), delta ms %d, frame time %f, tickrate: %d/s", tick, math.Ceil(counter), counter, deltatime, delta, updatesPerSecond))
		}
	}

	writeSystemString("Kill detected. Stopping the game\n\n")
	time.Sleep(time.Second)
}
