package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"go-airline-crew-rostering/airline"
	"go-airline-crew-rostering/archimedesOptimization"
	"go-airline-crew-rostering/graph"
	"go-airline-crew-rostering/input"
	"go-airline-crew-rostering/metrics"
	"go-airline-crew-rostering/multicso"
	"go-airline-crew-rostering/results"

	"golang.org/x/exp/slices"
)

func main() {
	startOfExecution := time.Now()
	args := input.SetUpParser() // parse arguments
	if args == nil {
		return
	}
	// Initialize random number generator
	if *args.Seed != -1 {
		rand.Seed(int64(*args.Seed))
	}
	al := AirlineSetup(args)
	pairGraph := GraphSetup(*args.Agents, al.PairsArray)

	var metric *metrics.Metrics
	if args.Algorithm == "multiCSO" {
		// execute multi-step CSO algorithm
		swarm := SwarmSetup(al, pairGraph, *args.Agents, *args.Generations, *args.FL)
		swarm.MultiCSO(al, pairGraph)
		chicken := swarm.Swarm[0]
		al.PilotsArray = chicken.Solution
		metric = swarm.Mtr

	} else if args.Algorithm == "AOA" {
		// execute AOA algorithm
		collection := CollectionSetup(al, pairGraph, *args.Agents, *args.Generations, args.Constants)
		collection.AOA(al, pairGraph)
		object := collection.Collection[0]
		al.PilotsArray = object.Solution
		metric = collection.Mtr
	}
	difference := 0.0
	rest := 0.0
	pairsCovered := 0
	totalDaysOff := 0
	// calculate rest between two pairs and days off
	for _, pilot := range al.PilotsArray {
		difference = difference + math.Abs(al.AverageWorkload-pilot.FlightTime)
		pairsCovered += pilot.AssignedLength
		restPerPilot := 0.0
		restPerPilot = pilot.TotalRestPeriod(al)
		totalDaysOff += pilot.DaysOff()
		rest += restPerPilot
	}
	metric.AverageRestPeriod = rest / float64((len(al.PairsArray)-1)*60)
	metric.AverageDaysOff = al.AverageDaysOff(totalDaysOff)
	metric.TotalTime = time.Since(startOfExecution)
	metric.TotalAssignedPairs = pairsCovered

	// check again if the solution obeys the rules
	if SolutionChecker(al.PilotsArray) {
		fmt.Println("Valid solution")
	} else {
		fmt.Println("Invalid solution")
	}
	results.PrintResults(metric, args, al) // store the results

}

func AirlineSetup(args *input.ArgumentCollection) *airline.Airline {
	// Create, initialize and set up an airline instance
	// returns pointer to the airline
	al := new(airline.Airline)
	al.Initialization(660, 7, 2, args.StartDate, args.EndDate, *args.Pilots)

	// create pairings
	al.PairsArray = input.ReadFile(*args.Filename, al.ScheduleStart)
	al.PairsArray = input.FilterPairs(al.PairsArray, al.ScheduleStart, al.ScheduleEnd)
	al.PairsArray = input.SortPairs(al.PairsArray)
	root := new(airline.Pair) // create special root pair
	root.Initialization(0, al.ScheduleStart)
	al.PairsArray = slices.Insert(al.PairsArray, 0, root)

	// create array of pilots
	for i := 0; i < al.NumberOfPilots; i++ {
		pilot := new(airline.Pilot)
		pilot.Initialization(i, al.ScheduleDuration, root)
		al.PilotsArray = append(al.PilotsArray, pilot)
	}
	al.CalculateAverageWorkload()
	return al
}

func GraphSetup(agents int, pairsArray []*airline.Pair) *graph.Graph {
	// Create and Initialize a graph with the pairings from the pairsArray
	// returns pointer to the graph
	pairGraph := new(graph.Graph)
	pairGraph.Initialization(agents)
	pairGraph.Populate(pairsArray)
	return pairGraph
}

func SwarmSetup(al *airline.Airline, pairGraph *graph.Graph, agents int, maxGenerations int, fl float64) *multicso.MultiCSO {
	// Create and Initialize a chicken swarm
	Mtr := new(metrics.Metrics)
	Mtr.Initialization(maxGenerations*agents, 32)
	swarm := new(multicso.MultiCSO)
	swarm.Initialization(al, pairGraph, agents, maxGenerations, fl, Mtr)
	return swarm
}

func CollectionSetup(al *airline.Airline, pairGraph *graph.Graph, agents int, maxGenerations int, constants []float64) *archimedesOptimization.AOAObjectCollection {
	//  Create and Initialize an object collection
	Mtr := new(metrics.Metrics)
	Mtr.Initialization(maxGenerations*agents, 32)
	collection := new(archimedesOptimization.AOAObjectCollection)
	collection.Initialization(agents, maxGenerations, constants[0], constants[1], constants[2], constants[3], al, pairGraph, Mtr)
	return collection
}

func SolutionChecker(solution []*airline.Pilot) bool {
	// checks if the given solution obeys the rules
	// returns true on success
	for _, pilot := range solution {
		for i := 2; i <= pilot.AssignedLength; i++ {
			rest := pilot.AssignedPairs[i].Start.Sub(pilot.AssignedPairs[i-1].End).Minutes()
			if rest < 660 {
				fmt.Printf("pilot %d: pair %d(source) and pair %d(goal) are only %e minutes apart\n", pilot.Id, pilot.AssignedPairs[i].Id, pilot.AssignedPairs[i-1].Id, rest)
				return false
			}
		}
		if !pilot.DaysOffRuleChecker() {
			fmt.Printf("pilot %d: not enough days off\n", pilot.Id)
			return false
		}
	}
	return true
}
