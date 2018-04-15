package main

import (
	"fmt"
	"time"
)

/* given some recent lap data predict which lap the car is on
Known issues, if all sensors are faulty will count incorrectly
will miss the pit stop sensor perhaps

We could make it a go routine, which would be able to act as a re-entrant
procedure
*/

func predictLap(data []Event) int {
	fmt.Printf("data = %+v\n", data)
	// count how many time we have seen a sensors
	sensors := make(map[string]int, 0)
	max := 1
	for _, v := range data {
		if c, ok := sensors[v.Emitter]; ok {
			sensors[v.Emitter] = c + 1
			if c+1 > max {
				max = c + 1
			}
		} else { //we haven't seen this sensor before
			sensors[v.Emitter] = 1
		}
	}
	return max
}

// a re-entrant version using go routines and a channel
// but how do we ensure that its thread safe? I.e. we need to ask it for data
// use three channels, one for reading events, another for requests and a
// final one for responses

func TrackLaps(stream <-chan Event, request <-chan string, result chan<- string) {
	for {
		select {
		// read from stream
		case msg := <-stream:
			fmt.Println("received from stream")
			fmt.Printf("msg = %+v\n", msg)

		// read from requests
		case msg := <-request:
			fmt.Println("received from request")
			fmt.Printf("msg = %+v\n", msg)
			result <- "foo"
		case <-time.Tick(time.Second):
			// while testing break after one second
			break
		}
	}

}

func testTrackLaps() {
	events := make(chan Event)
	req := make(chan string)
	res := make(chan string)

	go TrackLaps(events, req, res)
	req <- "request"
}
