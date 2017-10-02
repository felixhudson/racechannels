package main

import (
	"testing"
	"time"
)

func TestFunction(t *testing.T) {
	pos := make([]positions, 0)
	pos = append(pos, positions{1, 1, time.Now()})
	results := groupy(pos)
	if len(results[1]) != 1 {
		t.Log("Error sorting 1 car!")
		t.Fail()
	}

}

func Test_multisort(t *testing.T) {
	pos := make([]positions, 0)
	pos = append(pos, positions{9, 2, time.Now()})
	pos = append(pos, positions{1, 1, time.Now()})
	pos = append(pos, positions{2, 1, time.Now()})
	results := groupy(pos)
	if len(results[1]) != 2 {
		t.Log("Error sorting 1 car!")
		t.Fail()
	}
	if len(results[2]) != 1 {
		t.Log("Error sorting 1 car!")
		t.Fail()
	}

}
func Test_timesort(t *testing.T) {
	pos := make([]positions, 0)
	pos = append(pos, positions{9, 1, time.Unix(3, 0)})
	pos = append(pos, positions{1, 1, time.Unix(1, 0)})
	pos = append(pos, positions{2, 1, time.Unix(2, 0)})
	results := sorty(pos)
	if len(results) != 3 {
		t.Fail()
	}
	if results[0].Car != 1 {
		t.Log("car 1 should be first!")
		t.Log(results)
		t.Fail()
	}
	if results[2].Car != 9 {
		t.Log("car 9 should be third!")
		t.Log(results)
		t.Fail()
	}

}
func Test_timesortgroups(t *testing.T) {
	pos := make([]positions, 0)
	pos = append(pos, positions{9, 1, time.Unix(3, 0)})
	pos = append(pos, positions{1, 2, time.Unix(1, 0)})
	pos = append(pos, positions{2, 2, time.Unix(2, 0)})
	results := sorty(pos)
	if len(results) != 3 {
		t.Fail()
	}
	if results[0].Car != 1 {
		t.Log("car 1 should be first!")
		t.Log(results)
		t.Fail()
	}
	if results[2].Car != 9 {
		t.Log("car 9 should be third!")
		t.Log(results)
		t.Fail()
	}

}
