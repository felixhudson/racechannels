package main

import (
	"testing"
)

func expect(a interface{}, b interface{}, t *testing.T) {
	if a != b {
		t.Fatal("expect ", a, " got ", b)
	}
}

func Test_predict(t *testing.T) {
	event := Event{11, "s1", "0", "foo"}
	event2 := Event{11, "s1", "1", "foo"}
	events2 := Event{11, "s2", "1", "foo"}
	foo := predictLap([]Event{event})
	if foo != 1 {
		expect(1, foo, t)
	}

	foo = predictLap([]Event{event, event2})
	if foo != 2 {
		expect(2, foo, t)
	}
	foo = predictLap([]Event{event, events2, event2})
	if foo != 2 {
		expect(2, foo, t)
	}
}

func Test_TrackLaps(t *testing.T) {
	testTrackLaps()
}
