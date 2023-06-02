package multicso

import (
	"math"
	"math/rand"

	"go-airline-crew-rostering/airline"
	"go-airline-crew-rostering/graph"
	"go-airline-crew-rostering/problem"
)

var e float64 = 1e-11 // constant to avoid division by zero

// // Container for the functions related to Chicken
type ChickenRepo interface {
	Initialization() *Chicken
	SetUpdateParameters()
	UpdatePosition()
	ConstructSolution() bool
	Evaluate()
}

type Chicken struct {
	Id                   int
	Fitness              float64          // Fitness of current solution
	NewFitness           float64          // Fitness of the new proposed solution
	Cost                 float64          // Cost of the current solution
	NewCost              float64          // Cost of the new proposed solution
	Solution             []*airline.Pilot // Array of pilots representing the current solution found by the object
	ProposedSolution     []*airline.Pilot // Array of pilots representing a new solution
	CondensedSolution    []int            // Current solution in another form (used for easier calculation of metrics)
	NewCondensedSolution []int            // New solution in another form (used for easier calculation of metrics)
	S1                   float64          // CSO parameter
	S2                   float64          // CSO parameter
	random               float64          // uniform random number use in the update step
	randN                float64          // gaussian random number use in the update step
	randomChickens       []*Chicken       // list of random chickens of the swarm used in the update step
}

func (chicken *Chicken) Initialization(id int) *Chicken {
	// Initialization of an instance of a chicken
	chicken.Id = id
	chicken.Fitness, chicken.NewFitness, chicken.Cost, chicken.NewCost = 0, 0, 0, 0
	chicken.Solution, chicken.ProposedSolution = []*airline.Pilot{}, []*airline.Pilot{}
	chicken.CondensedSolution, chicken.NewCondensedSolution = []int{}, []int{}
	chicken.S1, chicken.S2, chicken.random, chicken.randN = 0, 0, 0, 0
	chicken.randomChickens = []*Chicken{}
	return chicken
}

func (chicken *Chicken) SetUpdateParameters(swarm []*Chicken) {
	// Calculate the update parameters for the next position update of the chicken

	// Find 4 random chickens in the swarm
	randomChickenList := make([]*Chicken, len(swarm))
	copy(randomChickenList, swarm)
	index := chicken.Id
	randomChickenList = append(randomChickenList[:index], randomChickenList[index+1:]...)
	k := rand.Intn(len(randomChickenList))
	randomChicken1 := randomChickenList[k]
	k = rand.Intn(len(randomChickenList))
	randomChicken3 := randomChickenList[k]
	k = rand.Intn(len(randomChickenList))
	randomChicken4 := randomChickenList[k]
	for index = 0; index <= len(randomChickenList); index++ {
		if randomChickenList[index].Id == randomChicken1.Id {
			break
		}
	}
	randomChickenList = append(randomChickenList[:index], randomChickenList[index+1:]...)
	k = rand.Intn(len(randomChickenList))
	randomChicken2 := randomChickenList[k]

	// use 1 random chicken for the calculation of S1
	chicken.S1 = math.Exp((chicken.Fitness - randomChicken1.Fitness) / (math.Abs(chicken.Fitness) + e))
	// use 1 random chicken for the calculation of S2
	chicken.S2 = math.Exp(randomChicken2.Fitness - chicken.Fitness)

	// generate random uniform number
	chicken.random = rand.Float64()

	// use 1 random chicken for the calculation of sigma
	var sigma float64
	if randomChicken3.Fitness >= chicken.Fitness {
		sigma = 1
	} else {
		sigma = math.Exp((randomChicken3.Fitness - chicken.Fitness) / (math.Abs(chicken.Fitness) + e))
	}

	// generate random gaussian number
	chicken.randN = rand.NormFloat64() * sigma

	// store the random chickens in a list
	chicken.randomChickens = []*Chicken{randomChicken1, randomChicken2, randomChicken3, randomChicken4}
}

func (chicken *Chicken) UpdatePosition(edge *graph.Edge, FL float64) {
	// Update the position of a graph edge

	position := edge.Position[chicken.Id]
	randomChickenPosition1 := edge.Position[chicken.randomChickens[0].Id]
	randomChickenPosition2 := edge.Position[chicken.randomChickens[1].Id]
	randomChickenPosition4 := edge.Position[chicken.randomChickens[3].Id]

	// use the CSO formulas to calculate the new position
	position = position +
		chicken.S1*chicken.random*(randomChickenPosition1-position) +
		chicken.S2*chicken.random*(randomChickenPosition2-position)
	position = position * (1 + chicken.randN)
	position = position + FL*(randomChickenPosition4-position)
	// Store new position to the graph
	edge.Position[chicken.Id] = position
}

func (chicken *Chicken) ConstructSolution(al *airline.Airline, graph *graph.Graph) bool {
	// Build a new solution for the chicken
	// Returns true if the solution covers all given pairings, false otherwise
	pilotsArray, condensedSolution, validSolution := problem.ConstructSolution(al, graph, chicken.Id)
	chicken.ProposedSolution = append([]*airline.Pilot(nil), pilotsArray...)
	chicken.NewCondensedSolution = append([]int{}, condensedSolution...)
	return validSolution
}

func (chicken *Chicken) Evaluate() {
	// Adopt the solution that is proposed as the best solution found by the object
	chicken.Solution = chicken.ProposedSolution
	chicken.CondensedSolution = chicken.NewCondensedSolution
	chicken.Fitness = chicken.NewFitness
	chicken.Cost = chicken.NewCost
}
