package main

import "fmt"
import "strings"
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
	StartPos     int
}

func main() {
	fmt.Println("vrooom!")

	// basic design is to have a channel that implements sim ticks
	// each go-routine will emit 0 or more events when it hears a sim tick

	events := make(chan Event, 10)

	// start a event reader

	board := make(map[int]int)
	go leaderBoard(events, board)

	// use a wait group to only quit when all closed
	var wg sync.WaitGroup
	c := Car{33, 5, 1}
	wg.Add(1)
	go c.speedcar(events, &wg)
	c1 := Car{5, 6, 2}
	wg.Add(1)
	go c1.speedcar(events, &wg)

	wg.Wait()
	close(events)
	fmt.Println("all routines are done")

}

func eventReader(events <-chan Event) {
	for v := range events {
		fmt.Printf("v = %+v\n", v)
	}
}

func eventPipe(events <-chan Event) {

	//read our events and push them to another service?

	// rabbit mq?
	// pubsub?
	// splunk
}

func leaderBoard(events <-chan Event, board map[int]int) {
	gaps := make(map[int]time.Duration)
	lastTime := time.Now()
	fmt.Println("starting leaderboard")
	for v := range events {
		// filter only sector times
		if strings.Index(v.Emitter,"sector" ) >= 0 {
			if _, ok := board[v.Car]; !ok {
				board[v.Car] = 0
			}
			// append a time gap
			//fmt.Println(lastTime, time.Now())
			gaps[v.Car] = time.Since(lastTime)
			lastTime = time.Now()
			board[v.Car] = board[v.Car] + 1
		}
		fmt.Println("----------")
		printBoard(board, gaps)
		fmt.Println("----------")
	}
}

func printBoard(data map[int]int, gaps map[int]time.Duration) {
	// we need to sort based on the value not key
	//find min max in values
	max := 0
	min := 99
	for _, v := range data {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}

	for i := max; i >= min; i-- {
		for k, v := range data {
			if v == i {
				fmt.Println("car", k, "Lap:", i)
				if i == max {
					fmt.Println("gap to car", "0.000")
				} else {
					fmt.Println("gap to car", gaps[k])
				}
			}
		}
	}
}

func (c *Car) speedcar(events chan<- Event, wg *sync.WaitGroup) {
	// delay amount based on car startpos
	time.Sleep(time.Duration(c.StartPos) * time.Second)
	var sector,cs string
	for i := 0; i < 10; i++ {
	  cs = fmt.Sprintf("current speed %d" , c.CurrentSpeed)
		sector = fmt.Sprintf("sector %d", (i%3))
		d := Event{c.Number, sector, "time", cs}
		events <- d
		sectorTime := 1000 - (c.CurrentSpeed * 100)
		time.Sleep(time.Duration(sectorTime) * time.Millisecond)
	}
	wg.Done()
}
