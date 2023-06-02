package input

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"go-airline-crew-rostering/airline"

	"golang.org/x/exp/slices"
)

func ReadFile(fileName string, scheduleStartDate time.Time) []*airline.Pair {
	// Read a csv file containing pairings and create a list of pairings
	// The pairings in the file must be sorted by Id from smallest to biggest
	// Returns the list of pairings
	pairs := []*airline.Pair{}
	file, err := os.Open(fileName)
	var year, month, day, hour, minute int
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ';'
	for {
		flightLeg, err := reader.Read()
		if err == io.EOF {
			break
		}
		pairId, _ := strconv.Atoi(flightLeg[0]) // Id of pairing
		// legId, _ := strconv.Atoi(flightLeg[1])
		// source := flightLeg[2]
		// destination := flightLeg[3]
		fmt.Sscanf(flightLeg[4], "%d-%d-%d", &year, &month, &day)                      // start date
		fmt.Sscanf(flightLeg[5], "%d:%d", &hour, &minute)                              // start time
		start := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC) // start date and time
		fmt.Sscanf(flightLeg[6], "%d-%d-%d", &year, &month, &day)                      // end date
		fmt.Sscanf(flightLeg[7], "%d:%d", &hour, &minute)                              // end time
		end := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)   // end date and time
		if pairId > len(pairs) {
			// Check if the pairing already exists and create a new one if it does not
			pair := new(airline.Pair)
			pair.Initialization(pairId, scheduleStartDate)
			pair.Add(pairId, start, end, scheduleStartDate)
			pairs = slices.Insert(pairs, pairId-1, pair)
		} else {
			// Otherwise add the flightleg to the appropriate pairing
			pairs[pairId-1].Add(pairId, start, end, scheduleStartDate)
		}
	}
	return pairs
}

func FilterPairs(pairsArray []*airline.Pair, scheduleStart time.Time, scheduleEnd time.Time) []*airline.Pair {
	// Filter list of pairings based on their start and end datetimes, using "scheduleStart"
	// and "scheduleEnd" as cutoffs, and return the filtered list
	filteredPairs := []*airline.Pair{}
	for _, pair := range pairsArray {
		if pair.Start.After(scheduleStart) && pair.End.Before(scheduleEnd) {
			filteredPairs = append(filteredPairs, pair)
		}
	}
	return filteredPairs
}

func SortPairs(pairsArray []*airline.Pair) []*airline.Pair {
	// Sort pairings by their start datetime, from oldest to most recent
	sort.Slice(pairsArray, func(i, j int) bool {
		return pairsArray[i].Start.Before(pairsArray[j].Start)
	})
	return pairsArray
}
