package airline

import (
	"time"
)

// struct representing a pairing
type Pair struct {
	Id         int
	Duration   float64   // duration (in minutes) of pairing
	Start      time.Time // Start date and time
	End        time.Time // end date and time
	StartDay   int       // number of days from start of schedule
	EndDay     int       // number of days from end of schedule
	FlightLegs int       // number of flightLegs
}

func (pair *Pair) Initialization(id int, scheduleStart time.Time) interface{} {
	// Initialization of a pairing
	if pair != nil {
		pair.Id = id
		pair.Duration = 0
		pair.Start = scheduleStart
		pair.End = scheduleStart
		pair.StartDay = 0
		pair.EndDay = 0
	}
	return pair
}

func (pair *Pair) Add(id int, start time.Time, end time.Time, scheduleStart time.Time) bool {
	// Adds a flightleg to a pairing
	// returns true on success
	if pair.Id != id {
		return false
	}
	if pair.Duration == 0 || start.Before(pair.Start) {
		pair.Start = start
		pair.StartDay = int(time.Duration.Hours(start.Sub(scheduleStart)) / 24)
	}
	if pair.Duration == 0 || end.After(pair.End) {
		pair.End = end
		pair.EndDay = int(time.Duration.Hours(end.Sub(scheduleStart)) / 24)
	}
	legDuration := time.Duration.Minutes(end.Sub(start))
	pair.Duration += legDuration
	pair.FlightLegs++
	return true
}
