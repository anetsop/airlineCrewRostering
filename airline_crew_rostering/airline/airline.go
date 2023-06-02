package airline

import (
	"math"
	"time"
)

// struct representing an airline
type Airline struct {
	PairsArray       []*Pair   // list of pairings to be assigned
	PilotsArray      []*Pilot  // final solution to the problem
	NumberOfPilots   int       // number of available pilots
	AverageWorkload  float64   // average workload per pilot
	RestPeriod       int       // minimum rest period between two consecutive pairings (in minutes)
	timespan         int       // time period (in days) that must contain a number of days off equal to "minimumDaysOff"
	minimumDaysOff   int       // minimum number of days without duty in a time period equal to "timespan"
	ScheduleStart    time.Time // Start of schedule
	ScheduleEnd      time.Time // End of schedule
	ScheduleDuration int       // Duration of schedule in days
}

// container for functions related to airline struct
type AirlineRepo interface {
	Initialization() interface{}
	Add() bool
	Remove() bool
	TotalRestPeriod() float64
	DaysOff() int
	CalculateAverageWorkload() float64
	OverlappingPairs() int
	RestPeriodRule() int
	DaysOffRule() bool
	DaysOffRuleChronological() bool
	EqualizeWorkload() []*Pilot
}

func (airline *Airline) Initialization(restPeriod int, timespan int, minimumDaysOff int,
	scheduleStart time.Time, scheduleEnd time.Time, numberOfPilots int) interface{} {
	// Initialization of an airline instance
	if airline != nil {
		airline.PairsArray = []*Pair{}
		airline.PilotsArray = []*Pilot{}
		airline.NumberOfPilots = numberOfPilots
		airline.AverageWorkload = 0
		airline.RestPeriod = restPeriod
		airline.timespan = timespan
		airline.minimumDaysOff = minimumDaysOff
		airline.ScheduleStart = scheduleStart
		airline.ScheduleEnd = scheduleEnd
		airline.ScheduleDuration = int(time.Duration.Hours(scheduleEnd.Sub(scheduleStart)) / 24)
	}
	return airline
}

func (airline *Airline) CalculateAverageWorkload() float64 {
	// Calculate and return the average workload per pilot,
	// based on the given pairings
	totalWorkload := 0.0
	if len(airline.PilotsArray) > 0 {
		for _, pair := range airline.PairsArray {
			totalWorkload += pair.Duration
		}
		airline.AverageWorkload = totalWorkload / float64(len(airline.PilotsArray))
		return airline.AverageWorkload
	}
	return -1
}

func (airline *Airline) OverlappingPairs(pair1 *Pair, pair2 *Pair) int {
	// Check if the time between pair1 and pair2 is more than "RestPeriod"
	// returns 0 if the time is less than "RestPeriod", 1 or 2 otherwise.
	// Also, return value of 1 means that pair2 is chronologically after pair1,
	// while return value of 2 means that pair2 is before pair1.
	if pair2.Start.Sub(pair1.End).Minutes() >= float64(airline.RestPeriod) {
		return 1
	}
	if pair1.Start.Sub(pair2.End).Minutes() >= float64(airline.RestPeriod) {
		return 2
	}
	return 0
}

func (airline *Airline) RestPeriodRule(pilot *Pilot, pair *Pair) int {
	// Check if the addition of "pair" to "pilot"' schedule would result
	// in two consecutive pairings having a time difference of less than
	// RestPeriod, in which case the function returns -1. Otherwise, it
	// returns the position in pilot's list of assigned pairs, where the
	// pair should be inserted to preserve the list's chronological order
	currentOverlap := -1
	previousOverlap := -1
	// check if the "pair" is after the pilot's last assigned pair
	if pilot.AssignedLength > 0 {
		currentOverlap = airline.OverlappingPairs(pilot.AssignedPairs[pilot.AssignedLength], pair)
	}
	// if it is after the last, or the pilot's list is empty return the last position of the list
	if pilot.AssignedLength == 0 || currentOverlap == 1 {
		return pilot.AssignedLength + 1
	} else if currentOverlap == 0 {
		// overlap with the last pairing
		return -1
	}
	// check each pairing against "pair", excluding the special "root"
	// pairing at position 0 of the list
	for index := 1; index <= pilot.AssignedLength; index++ {
		pilotPair := pilot.AssignedPairs[index]
		currentOverlap = airline.OverlappingPairs(pilotPair, pair)
		if currentOverlap == 2 {
			if index == 1 || previousOverlap == 1 {
				return index
			}
		} else if currentOverlap == 0 {
			break
		}
		previousOverlap = currentOverlap
	}
	return -1
}

func (airline *Airline) DaysOffRule(pilot *Pilot, pair *Pair, chronologicalOptional ...bool) bool {
	// Check if the addition of "pair" to "pilot"' schedule would result
	// in the pilot having less than "minimumDaysOff" in a "timespan" period
	// chronologicalOptional is a flag set to true if the order
	// we examine the pairings is chronological by their start datetimes
	// if we omit the flag, the default value is false
	daysOff := airline.timespan                    // days in a timespan without duty
	startPoint := pair.StartDay - airline.timespan // first day of first timespan to search
	middlePoint := pair.StartDay                   // last day of first timespan to search
	endPoint := pair.StartDay + airline.timespan   // last day of last timespan to search
	ruleConfirmed := true
	chronological := false

	if len(chronologicalOptional) > 0 {
		chronological = chronologicalOptional[0]
	}
	for i := pair.StartDay; i <= pair.EndDay; i++ {
		pilot.workdays[i]++
	}

	// adjust the startPoint of the search if timespan > pair.startDay
	if startPoint < 0 {
		startPoint = 0
		middlePoint = airline.timespan
	}
	if endPoint > len(pilot.workdays) {
		endPoint = len(pilot.workdays)
	} else if chronological {
		// fewer calculations are needed if we have
		// chronological order of examination
		endPoint = pair.EndDay + 1
	}

	// check the first timespan
	for i := startPoint; i < middlePoint; i++ {
		if pilot.workdays[i] > 0 {
			daysOff--
		}
	}
	if daysOff < airline.minimumDaysOff {
		ruleConfirmed = false
	}

	// check the remaining timespans
	for i := middlePoint; i < endPoint; i++ {
		if pilot.workdays[i-airline.timespan] > 0 {
			daysOff++
		}
		if pilot.workdays[i] > 0 {
			daysOff--
		}
		if daysOff < airline.minimumDaysOff {
			ruleConfirmed = false
			break
		}
	}
	for i := pair.StartDay; i <= pair.EndDay; i++ {
		pilot.workdays[i]--
	}
	return ruleConfirmed
}

func (airline *Airline) EqualizeWorkload(pilots []*Pilot) []*Pilot {
	// Reassign pairings from pilots with heavier workloads
	// to pilots with lighter ones
	// returns the list of pilots with the improved schedules
	for _, pilot := range pilots {
		if pilot.FlightTime > airline.AverageWorkload {
			for i := 1; i <= pilot.AssignedLength; i++ {
				pair := pilot.AssignedPairs[i]
				difference := pilot.FlightTime - airline.AverageWorkload
				if math.Abs(difference-pair.Duration) < difference {
					for _, pilot2 := range pilots {
						if pilot2.Id == pilot.Id {
							continue
						}
						difference2 := airline.AverageWorkload - pilot2.FlightTime
						if math.Abs(difference2-pair.Duration) < math.Abs(difference2) {
							index := airline.RestPeriodRule(pilot2, pair)
							if index > -1 && airline.DaysOffRule(pilot2, pair) {
								pilot.Remove(pair)
								pilot2.Add(pair, index)
								i--
								break
							}
						}
					}
				}
				if pilot.FlightTime < airline.AverageWorkload {
					break
				}
			}
		}
	}
	return pilots
}

func (al *Airline) AverageDaysOff(totalDaysOff int) float64 {
	// Calculate the average days off per pilot per timespan period
	numberOfTimespansInSchedule := float64((al.ScheduleDuration - 1) / al.timespan)
	average := float64(totalDaysOff) / numberOfTimespansInSchedule
	average = average / float64(al.NumberOfPilots)
	return average
}
