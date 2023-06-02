package problem

import (
	"math/rand"

	"go-airline-crew-rostering/airline"
	"go-airline-crew-rostering/graph"
)

// struct representing a pilot that could accept a pair
type candidate struct {
	pilot    *airline.Pilot
	position float64
	index    int // position in the list of the pilot's assigned pairs, where the new pair would be inserted
}

func ConstructSolution(al *airline.Airline, graph *graph.Graph, id int) ([]*airline.Pilot, []int, bool) {
	// Build a solution for agent with id "id"
	// returns the solution in 2 different forms
	// and whether it is valid

	// Create a list of pilots to assign pairs
	pilotsArray := []*airline.Pilot{}
	condensedSolution := []int{}
	for i := 0; i < al.NumberOfPilots; i++ {
		pilot := new(airline.Pilot)
		pilot.Initialization(i, al.ScheduleDuration, al.PairsArray[0])
		pilotsArray = append(pilotsArray, pilot)
	}
	validSolution := true
	// Take each pair and try to assign it to a pilot
	for _, pair := range al.PairsArray[1:] {
		candidates := []*candidate{} // list of pilots that can accept the pair "pair"
		index := 0
		for _, pilot := range pilotsArray {
			// Check if the assignment of "pair" to "pilot" obeys to the rules
			index = al.RestPeriodRule(pilot, pair)
			if index > -1 && al.DaysOffRule(pilot, pair, true) {
				// Find the position of the pair that is just before the "pair"
				// in the pilot's schedule by taking the corresponding edge
				previousPair := pilot.AssignedPairs[index-1]
				edge, EdgeExists := graph.Nodes[previousPair.Id].Edges[pair.Id]
				var position float64
				if EdgeExists {
					position = edge.Position[id]
				} else {
					// if the edge does not exist (new connection) we use a default value
					position = 1
				}

				// Heuristic mechanism to reinfonce more compact pilots' schedules
				if previousPair.Id > 0 {
					restPeriod := pair.Start.Sub(pilot.AssignedPairs[index-1].End).Hours()
					restPeriod = restPeriod - float64(al.RestPeriod)/60 + 1
					position = position / restPeriod
				}
				candidatepilot := &candidate{
					pilot:    pilot,
					position: position,
					index:    index,
				}
				candidates = append(candidates, candidatepilot)
			}
		}
		// select a pilot from the "candidates" list
		selectedPilot := selectPilot(candidates)
		if selectedPilot == nil {
			validSolution = false
		} else {
			// add the pairing "pair" to the selected pilot's schedule
			selectedPilot.pilot.Add(pair, selectedPilot.index)
			condensedSolution = append(condensedSolution, selectedPilot.pilot.Id)
		}
	}
	// optimize the solution
	al.EqualizeWorkload(pilotsArray)
	if validSolution {
		// try to optimize again since it is valid
		al.EqualizeWorkload(pilotsArray)
	}
	return pilotsArray, condensedSolution, validSolution
}

func selectPilot(candidates []*candidate) *candidate {
	// select a pilot from "candidates" list
	if len(candidates) == 0 { // empty list
		return nil
	} else if len(candidates) == 1 { // only 1 to choose
		return candidates[0]
	} else {
		offset := 0.0

		// find candidate with minimum position
		for _, candidate := range candidates {
			if candidate.position < offset {
				offset = candidate.position
			}
		}
		// normalize all positions with respect to 0
		if offset < 0 {
			for _, candidate := range candidates {
				candidate.position -= offset
			}
		}

		// sum all normalized positions
		sum := 0.0
		for _, candidate := range candidates {
			sum += candidate.position
		}
		random := rand.Float64() // generate a random number in range [0,1)
		random = random * sum    // take a percentage of the summed positions
		// choose a pilot randomly
		for _, candidate := range candidates {
			random -= candidate.position
			if random <= 0 {
				return candidate
			}
		}
	}
	return nil
}
