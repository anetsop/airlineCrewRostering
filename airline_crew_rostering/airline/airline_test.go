package airline_test

import (
	"testing"
	"time"

	"go-airline-crew-rostering/airline"
)

func TestAirline(t *testing.T) {
	// Simple test for airline package
	pilot := new(airline.Pilot)
	al := new(airline.Airline)
	start1 := time.Date(2011, 11, 1, 3, 20, 0, 0, time.UTC)
	start2 := time.Date(2011, 11, 1, 5, 0, 0, 0, time.UTC)
	start3 := time.Date(2011, 11, 1, 12, 0, 0, 0, time.UTC)
	end1 := time.Date(2011, 11, 1, 4, 20, 0, 0, time.UTC)
	end2 := time.Date(2011, 11, 1, 6, 5, 0, 0, time.UTC)
	end3 := time.Date(2011, 11, 1, 13, 5, 0, 0, time.UTC)
	startSchedule := time.Date(2011, 11, 1, 0, 0, 0, 0, time.UTC)
	endSchedule := time.Date(2011, 11, 2, 0, 0, 0, 0, time.UTC)
	al.Initialization(660, 7, 2, startSchedule, endSchedule, 45)
	root := new(airline.Pair)
	root.Initialization(0, startSchedule)
	pair1 := new(airline.Pair)
	pair2 := new(airline.Pair)
	pair1.Initialization(1, startSchedule)
	pair2.Initialization(2, startSchedule)
	pilot.Initialization(0, al.ScheduleDuration, root)
	pair1.Add(1, start3, end3, startSchedule)
	pair1.Add(1, start2, end2, startSchedule)
	pair1.Add(1, start1, end1, startSchedule)

	pair2.Add(2, start1, end1, startSchedule)
	// t.Log(pair1.Duration, "minutes")
	// t.Log("id:", pair1.Id)
	// t.Log("start:", pair1.Start)
	// t.Log("end:", pair1.End)

	pilot.Add(pair1, 1)
	// t.Log("pilot first pair:", pilot.AssignedPairs[0].Id, pilot.AssignedPairs[0].Start)
	// t.Log("pilot second pair:", pilot.AssignedPairs[1].Id)
	pilot.Add(pair2, 1)
	// t.Log("pilot second pair:", pilot.AssignedPairs[1].Id)
	// t.Log("pilot third pair:", pilot.AssignedPairs[2].Id)
	pilot.Remove(pair2)
	// t.Log("pilot second pair:", pilot.AssignedPairs[1].Id)

	pilot2 := new(airline.Pilot)
	pilot2.Initialization(1, al.ScheduleDuration, root)
	t.Log("pilot2 pairs:", pilot2.AssignedLength, pilot.AssignedPairs[0])
	t.Log(al.OverlappingPairs(pair1, pair2))
	start4 := time.Date(2011, 11, 2, 18, 20, 0, 0, time.UTC)
	end4 := time.Date(2011, 11, 2, 19, 20, 0, 0, time.UTC)
	pair3 := new(airline.Pair)
	pair3.Initialization(3, startSchedule)
	pair3.Add(3, start4, end4, startSchedule)
	t.Log(al.OverlappingPairs(pair1, pair3))
	t.Log(al.OverlappingPairs(pair3, pair1))
	i := al.RestPeriodRule(pilot, pair1)
	t.Log("index for pilot and pair1:", i)
	i = al.RestPeriodRule(pilot, pair3)
	t.Log("index for pilot and pair3:", i)
	i = al.RestPeriodRule(pilot2, pair3)
	t.Log("index for pilot2 and pair3:", i)
	pilot2.Add(pair3, i)
	i = al.RestPeriodRule(pilot2, pair1)
	t.Log("index for pilot2 and pair1:", i)
	pilot2.Add(pair1, i)

}
