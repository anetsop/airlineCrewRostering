package metrics

import (
	"strconv"
	"strings"
	"time"

	"go-airline-crew-rostering/airline"
)

// Container for functions related to metrics
type MetricsRepo interface {
	Initialization()
	SolutionEncoding() []int
	SolutionSimilarity() float64
}

type Metrics struct {
	TotalTime              time.Duration   // Application's execution time
	ValidSolutions         int             // total valid solutions found
	TotalSolutions         int             // total solutions (valid + invalid)
	TotalAssignedPairs     int             // pairs covered by the solution of the app
	AverageRestPeriod      float64         // average rest period per pair of pairings
	AverageDaysOff         float64         // average days off per pilot per timespan
	UnitCost               float64         // cost of 1 unit of the solution's cost
	GlobalBestSolutionCost float64         // Cost of best solution
	GlobalBestString       []string        // solution encoded as a string
	IterBestCost           []float64       // list of the best cost of each iteration
	IterWorstCost          []float64       // list of the worst cost of each iteration
	IterAverageCost        []float64       // list of the average cost of each iteration
	Jumps                  int             // number of times we found a new global best
	uniqueSolutions        map[string]bool // list with all the different solutions found
	UniqueCount            int             // number of the different solutions found
	AverageSimilarity      float64         // average similarity between each solution and the global best
}

func (m *Metrics) Initialization(totalSolutions int, unitCost float64) {
	// Initialize a new instance of metrics
	m.TotalTime = 0
	m.ValidSolutions = 0
	m.TotalSolutions = totalSolutions
	m.TotalAssignedPairs = 0
	m.AverageRestPeriod = 0
	m.UnitCost = unitCost
	m.GlobalBestSolutionCost = 0
	m.GlobalBestString = []string{}
	m.IterBestCost = []float64{}
	m.IterWorstCost = []float64{}
	m.IterAverageCost = []float64{}
	m.Jumps = 0
	m.uniqueSolutions = make(map[string]bool)
	m.UniqueCount = 0
	m.AverageSimilarity = 0
}

func (m *Metrics) SetUpIterationMetrics(costList []float64) (int, int) {
	// Calculate the best, worst and average values of the "costList"
	// (best is smallest, worst is the biggest value)
	bestIndex := 0
	bestCost := costList[0]
	worstIndex := 0
	worstCost := costList[0]
	averageCost := 0.0
	for i, cost := range costList {
		averageCost += cost
		if cost < bestCost {
			bestIndex = i
			bestCost = cost
		}
		if cost > worstCost {
			worstIndex = i
			worstCost = cost
		}
	}
	m.IterAverageCost = append(m.IterAverageCost, averageCost/float64(len(costList)))
	m.IterBestCost = append(m.IterBestCost, bestCost)
	m.IterWorstCost = append(m.IterWorstCost, worstCost)
	return bestIndex, worstIndex
}

func (m *Metrics) SolutionEncoding(condensedSolution []int, completeSolution []*airline.Pilot) []string {
	// encode the solution described by "condensedSolution" and "completeSolution"
	// (they represent the same solution) as a string and check if this is the
	// first time we found this solution
	// returns a list of strings that has the same shape as "condensedSolution"
	normalisedSolution := []string{}
	pilotAliases := make(map[int]string)

	// Create aliases for pilots to account for symmetrical solutions
	for _, pilot := range completeSolution {
		if pilot.AssignedLength > 0 {
			pilotAliases[pilot.Id] = strconv.Itoa(pilot.AssignedPairs[1].Id)
		}
	}

	// create the encoded string based on the aliases
	for _, pilotId := range condensedSolution {
		normalisedSolution = append(normalisedSolution, pilotAliases[pilotId])
	}
	solutionString := strings.Join(normalisedSolution, ",")

	// use the map of encoded solution strings to find out
	// if we have built this solution before
	_, keyExists := m.uniqueSolutions[solutionString]
	if !keyExists {
		// if not update the map and the unique solutions counter
		m.UniqueCount++
		m.uniqueSolutions[solutionString] = true
	}
	return normalisedSolution
}

func (m *Metrics) SolutionSimilarity(solutionA []string, solutionB []string) float64 {
	// Calculate the similarity between two normalized solutions
	// obtained from "SolutionEncoding" function
	similarity := 0.0
	biggestLength := 0
	smallestLength := 0
	if len(solutionA) >= len(solutionB) {
		biggestLength = len(solutionA)
		smallestLength = len(solutionB)
	} else {
		biggestLength = len(solutionB)
		smallestLength = len(solutionA)
	}
	difference := biggestLength - smallestLength
	for i := 0; i < smallestLength; i++ {
		if solutionA[i] != solutionB[i] {
			difference++
		}
	}
	similarity = (1.0 - float64(difference)/float64(biggestLength)) * 100.0
	return similarity
}
