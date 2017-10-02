package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

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

	// all goroutines will recieve data in an event queue
	events := make(chan Event, 10)

	// use a wait group to only quit when all closed
	var wg sync.WaitGroup
	c := Car{33, 5, 1}
	wg.Add(1)
	c1 := Car{5, 6, 2}
	wg.Add(1)
	c2 := Car{111, 5, 3}
	//cars := []Car{c,c1,c2}
	wg.Add(1)
	// start a event reader
	//board := make(map[int]int)
	//go leaderBoard(events, board)
	//go matrixBoard(events,cars)
	//go orderBoard(events)
	go sortBoard(events)

	go c.speedcar(events, &wg)
	go c1.speedcar(events, &wg)
	go c2.speedcar(events, &wg)

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

func sortBoard(events <-chan Event) {
	// keep a list of when we last saw the car and sort by that order!
	// keep a count of the lap also!
	// sort by lapnumber Desc then last seen Asc
	data := make(map[int]positions,0)
	tosort := make([]positions,0)
	var lap,size int
	for event := range events {
		//update the data
		size = len(data)
		if _,ok := data[event.Car]; ok {
			lap = data[event.Car].Lap + 1
			data[event.Car] = positions{event.Car,lap,time.Now()}
			if size != len(data) {
				log.Fatal("didnt match")
			}
		}else {
			data[event.Car] = positions{event.Car,1,time.Now()}
			if size + 1 != len(data) {
				log.Fatal("didnt grow")
			}
		}
		// make positions
		tosort = nil
		size = len(data)
		for _, v := range data {
			tosort = append(tosort,v)
		}
		if len(tosort) != size {
			log.Fatal("too much to sort")
		}
		fmt.Println("=====")
		for k, v := range sorty(tosort) {
			fmt.Printf("%d %d %d\n" , k+1,v.Car,v.Lap)
		}

	}

}

func leaderBoard(events <-chan Event, board map[int]int) {
	gaps := make(map[int]time.Duration)
	lastTime := time.Now()
	fmt.Println("starting leaderboard")
	for v := range events {
		// filter only sector times
		if strings.Index(v.Emitter, "sector") >= 0 {
			if _, ok := board[v.Car]; !ok {
				board[v.Car] = 0
			}
			// append a time gap
			gaps[v.Car] = time.Since(lastTime)
			lastTime = time.Now()
			board[v.Car] = board[v.Car] + 1
		}
		fmt.Println("----------")
		printBoard(board, gaps)
		fmt.Println("----------")
	}
}

type CarTimingsMap struct {
	car      *Car
	timing   map[int]time.Duration
	lastseen time.Time
}

func orderBoard(events <-chan Event) {
	//capture the cars in order
	// everytime we see the first car again reset order
	order := make([]int, 0)
	for event := range events {
		if (len(order) > 0) && (order[0] == event.Car) {
			//reset the order
			order = nil
		}
		order = append(order, event.Car)
		fmt.Println("=========")
		for k, v := range order {
			fmt.Println(k+1, "::", v)
		}
	}
}

func matrixBoard(events <-chan Event, cars []Car) {
	// have a matrix of cars. Work out timing differences between them
	var matrix = make(map[int]CarTimingsMap, 0)
	for _, v := range cars {
		matrix[v.Number] = CarTimingsMap{&v, make(map[int]time.Duration), time.Unix(0, 0)}
	}
	nowtime := time.Now()
	var other time.Time

	for event := range events {
		fmt.Printf("event = %+v\n", event)
		// if car doesnt exist in timings add it!
		nowtime = time.Now()
		fmt.Printf("nowtime = %+v\n", nowtime)
		for _, v := range matrix {
			if _, ok := v.timing[event.Car]; !ok {
				//matrix[v.Car] = 0
				fmt.Println("adding a car")
				v.timing[event.Car] = 0
				fmt.Printf("matrix = %+v\n", matrix)
			}
			// calculate gap by
			// find time of other car
			other = matrix[event.Car].lastseen
			v.timing[event.Car] = time.Since(other)
		}
	}
}
func sortGroup(){
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
					max = -1
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
	var sector, cs string
	for i := 0; i < 10; i++ {
		cs = fmt.Sprintf("current speed %d", c.CurrentSpeed)
		sector = fmt.Sprintf("sector %d", (i % 3))
		d := Event{c.Number, sector, "time", cs}
		events <- d
		sectorTime := 1000 - (c.CurrentSpeed * 100)
		time.Sleep(time.Duration(sectorTime) * time.Millisecond)
	}
	wg.Done()
}
