package airline

import "golang.org/x/exp/slices"

// struct representing an airline pilot
type Pilot struct {
	Id             int
	AssignedPairs  []*Pair // list of assigned pairings
	AssignedLength int     // length of assigned pairs list minus 1 (the root pair)
	FlightTime     float64 // pilot flight time
	workdays       []int   // list representing the days of schedule showing how many
	// pairings the pilot has each day
}

func (pilot *Pilot) Initialization(id int, scheduleDuration int, root *Pair) interface{} {
	// Initialization of a pilot instance
	// returns a pointer to the pilot
	if pilot != nil {
		pilot.Id = id
		pilot.AssignedPairs = []*Pair{root} // the assigned pairs list is initialized with the special root pair
		pilot.AssignedLength = 0
		pilot.FlightTime = 0
		pilot.workdays = make([]int, scheduleDuration)
	}
	return pilot
}

func (pilot *Pilot) Add(pair *Pair, index int) bool {
	// Add a pair to the pilot's schedule in a specific position
	// given by "index"
	// returns true on success
	if index >= pilot.AssignedLength+1 {
		pilot.AssignedPairs = append(pilot.AssignedPairs, pair)
	} else if index > 0 {
		pilot.AssignedPairs = slices.Insert(pilot.AssignedPairs, index, pair)
	} else {
		return false
	}
	pilot.FlightTime += pair.Duration
	pilot.AssignedLength++
	for i := pair.StartDay; i <= pair.EndDay; i++ {
		pilot.workdays[i]++
	}
	return true
}

func (pilot *Pilot) Remove(pair *Pair) bool {
	// removes pair from pilot's schedule
	// returns true on success
	index := 1
	for ; index <= pilot.AssignedLength; index++ {
		if pilot.AssignedPairs[index] == pair {
			break
		}
	}
	if index > pilot.AssignedLength {
		return false
	}
	pilot.AssignedPairs = append(pilot.AssignedPairs[:index], pilot.AssignedPairs[index+1:]...)
	pilot.FlightTime -= pair.Duration
	pilot.AssignedLength--
	for i := pair.StartDay; i <= pair.EndDay; i++ {
		pilot.workdays[i]--
	}
	return true
}

func (pilot *Pilot) TotalRestPeriod(al *Airline) float64 {
	// Calculate the total excess rest period of the pilot
	// (excluding minimum days off and mandatory rests)
	// and returns the duration in minutes
	minimumDaysOff := float64(al.minimumDaysOff)
	minimumDaysOff *= float64(al.ScheduleDuration / al.timespan)
	timeOff := float64(24 * 60 * minimumDaysOff)
	restTime := 0.0
	if pilot.AssignedLength == 0 {
		return float64(al.ScheduleEnd.Sub(al.ScheduleStart).Minutes() - timeOff)
	}
	for i := 1; i <= pilot.AssignedLength; i++ {
		restTime += pilot.AssignedPairs[i].Start.Sub(pilot.AssignedPairs[i-1].End).Minutes()
		if minimumDaysOff > 0 {
			daysInBetween := float64(pilot.AssignedPairs[i].StartDay - pilot.AssignedPairs[i-1].EndDay - 1)
			if daysInBetween <= 0 {
				restTime -= float64(al.RestPeriod)
			} else if daysInBetween > minimumDaysOff {
				daysInBetween -= minimumDaysOff
				restTime -= 1440 * minimumDaysOff
				minimumDaysOff = 0
			} else {
				minimumDaysOff -= daysInBetween
				restTime -= 1440 * daysInBetween
			}
		} else {
			restTime -= float64(al.RestPeriod)
		}
	}
	daysInBetween := float64(al.ScheduleDuration - pilot.AssignedPairs[pilot.AssignedLength].EndDay - 1)
	if minimumDaysOff > 0 {
		daysInBetween -= minimumDaysOff
		restTime += 1440 * daysInBetween
	} else {
		restTime += 1440 * daysInBetween
	}
	return restTime
}

func (pilot *Pilot) DaysOff() int {
	// Calculate the total days off of a pilot
	count := 0
	for _, workday := range pilot.workdays {
		if workday == 0 {
			count++
		}
	}
	return count - 1
}

func (pilot *Pilot) DaysOffRuleChecker() bool {
	// returns true if the pilot's schedule
	// obeys the rule implemented by the
	// DaysOffRule function
	daysOff := 7
	for i := 0; i < 7; i++ {
		if pilot.workdays[i] > 0 {
			daysOff--
		}
	}
	if daysOff < 2 {
		return false
	}
	for i := 7; i < 124; i++ {
		if pilot.workdays[i-7] > 0 {
			daysOff++
		}
		if pilot.workdays[i] > 0 {
			daysOff--
		}
		if daysOff < 2 {
			return false
		}
	}
	return true
}
