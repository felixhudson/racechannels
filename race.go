package main

import "fmt"
import "sync"
import "time"

type Event struct {
	Car     int
	Emitter string
	Time    string
	Debug   string
}

type Car struct {
	Number       int
	CurrentSpeed int
}

func main() {
	fmt.Println("vrooom!")

	// basic design is to have a channel that implements sim ticks
	// each go-routine will emit 0 or more events when it hears a sim tick

	events := make(chan Event, 10)
	//cars := 2
	sim := make([]chan bool, 0)
	s1 := make(chan bool)
	s2 := make(chan bool)
	s3 := make(chan bool)
	sim = append(sim, s1)
	sim = append(sim, s2)
	sim = append(sim, s3)

	// use a wait group to only quit when all closed
	var wg sync.WaitGroup
	wg.Add(1)
	go car(sim[0], events, &wg)
	wg.Add(1)
	go forevercar(sim[1], events, &wg)
	c := Car{33, 5}
	wg.Add(1)
	go c.vroom(sim[2], events, &wg)
	// send a tick to start sim
	fmt.Println("starting ticks")
	for _, v := range sim {
		v <- true
	}
	for _, v := range sim {
		close(v)
	}

	wg.Wait()
	fmt.Println("all routines are done")

}

func car(tick <-chan bool, events chan<- Event, wg *sync.WaitGroup) {
	// we can only listen to ticks and send events
	// wait a little time then send an event and quit
	fmt.Println("starting car, wait for tick")
	t := <-tick
	if t {
		d := Event{1, "car", "time", ""}
		events <- d
	}
	fmt.Println("car finished")
	wg.Done()
}

func (*Car) vroom(tick <-chan bool, events chan<- Event, wg *sync.WaitGroup) {
	// get a tick but loop until done
	ever := true
	for ever {
		time.Sleep(100 * time.Microsecond)
		_, more := <-tick
		if more {
			fmt.Println("vroom a tick!")
			d := Event{2, "car", "time", ""}
			events <- d
		} else {
			ever = false
		}

	}
	fmt.Println("vroom done")
	wg.Done()

}

func forevercar(tick <-chan bool, events chan<- Event, wg *sync.WaitGroup) {
	// get a tick but loop until done
	ever := true
	for ever {
		time.Sleep(100 * time.Microsecond)
		_, more := <-tick
		if more {
			fmt.Println("got a tick!")
			d := Event{2, "car", "time", ""}
			events <- d
		} else {
			ever = false
		}

	}
	fmt.Println("forever done")
	wg.Done()

}
