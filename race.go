package main

import "fmt"
import "sync"
import "time"

type Event struct {
	Name    string
	Emitter string
	Time    string
	Data    string
}

func main() {
	fmt.Println("vrooom!")

	// basic design is to have a channel that implements sim ticks
	// each go-routine will emit 0 or more events when it hears a sim tick

	SimTicks := make(chan bool)

	events := make(chan Event, 5)
	// use a wait group to only quit when all closed
	var wg sync.WaitGroup
	wg.Add(1)
	go car(SimTicks, events, &wg)
	wg.Add(1)
	done := make(chan bool)
	go forevercar(done, SimTicks, events, &wg)
	// send a tick to start sim
	SimTicks <- true
	SimTicks <- true
	SimTicks <- true
	SimTicks <- true
	SimTicks <- true
	done <- true

	wg.Wait()
	fmt.Println("all routines are done")

}

func car(tick <-chan bool, events chan<- Event, wg *sync.WaitGroup) {
	// we can only listen to ticks and send events
	// wait a little time then send an event and quit
	fmt.Println("starting car, wait for tick")
	t := <-tick
	if t {
		d := Event{"test", "car", "time", ""}
		events <- d
	}
	fmt.Println("car finished")
	wg.Done()
}

func forevercar(done <-chan bool, tick <-chan bool, events chan<- Event, wg *sync.WaitGroup) {
	// get a tick but loop until done
	ever := true
	for ever {
		time.Sleep(1 * time.Second)
		select {
		case <-tick:
			{
				fmt.Println("got a tick!")
				d := Event{"test", "car", "time", ""}
				events <- d
			}
		case <-done:
			{
				ever = false
			}
		}
	}
	wg.Done()

}
