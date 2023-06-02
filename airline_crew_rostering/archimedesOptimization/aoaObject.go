package archimedesOptimization

import (
	"math/rand"

	"go-airline-crew-rostering/airline"
	"go-airline-crew-rostering/graph"
	"go-airline-crew-rostering/problem"
)

// Container for the functions related to AOAObject
type AOAObjectRepo interface {
	Initialization() *AOAObject
	SetUpdateParameters()
	UpdatePosition()
	ConstructSolution() bool
	Evaluate()
}

// struct representing an archimedes optimization object
type AOAObject struct {
	Id                   int
	Fitness              float64          // Fitness of current solution
	NewFitness           float64          // Fitness of the new proposed solution
	Cost                 float64          // Cost of the current solution
	NewCost              float64          // Cost of the new proposed solution
	Solution             []*airline.Pilot // Array of pilots representing the current solution found by the object
	ProposedSolution     []*airline.Pilot // Array of pilots representing a new solution
	CondensedSolution    []int            // Current solution in another form (used for easier calculation of metrics)
	NewCondensedSolution []int            // New solution in another form (used for easier calculation of metrics)
	density              float64
	volume               float64
	accelleration        float64
	randomObject         *AOAObject // object whose position is used in the update step
}

func (obj *AOAObject) Initialization(id int) *AOAObject {
	// Initialization of an instance of an AOA object
	obj.Id = id
	obj.Fitness, obj.NewFitness, obj.Cost, obj.NewCost = 0, 0, 0, 0
	obj.Solution, obj.ProposedSolution = []*airline.Pilot{}, []*airline.Pilot{}
	obj.CondensedSolution, obj.NewCondensedSolution = []int{}, []int{}
	obj.density = rand.Float64()
	obj.volume = rand.Float64()
	obj.accelleration = rand.Float64()
	return obj
}

func (obj *AOAObject) SetUpdateParameters(objects []*AOAObject, bestObjectIndex int, TF float64) {
	// Calculate the update parameters for the next position update of the object
	bestObject := objects[bestObjectIndex]

	// Select random object used for the position update
	randomObjectList := make([]*AOAObject, len(objects))
	copy(randomObjectList, objects)
	index := obj.Id
	randomObjectList = append(randomObjectList[:index], randomObjectList[index+1:]...)
	k := rand.Intn(len(randomObjectList))
	obj.randomObject = randomObjectList[k]

	// Calculate density and volume
	obj.density = obj.density + rand.Float64()*(bestObject.density-obj.density)
	obj.volume = obj.volume + rand.Float64()*(bestObject.volume-obj.volume)

	updateObject := bestObject // Select object with best fitness if we are in exploitation phase
	if TF <= 0.5 {
		// Select random object if we are in exploration phase
		updateObject = obj.randomObject
	}
	// Calculate non-normalized acceleration
	obj.accelleration = (updateObject.density + updateObject.volume*updateObject.accelleration) / (obj.density * obj.volume)

}

func (obj *AOAObject) UpdatePosition(edge *graph.Edge, updateParams *aoaParameters) {
	// Update the position of a graph edge
	position := edge.Position[obj.Id]
	if updateParams.TF <= 0.5 {
		// Exploration phase
		randomObjectPosition := edge.Position[obj.randomObject.Id]
		position = position +
			updateParams.C1*rand.Float64()*obj.accelleration*updateParams.d*(randomObjectPosition-position)
	} else {
		// Exploitation phase
		bestPosition := edge.Position[updateParams.bestObject.Id]
		position = bestPosition +
			float64(updateParams.F)*updateParams.C2*rand.Float64()*obj.accelleration*updateParams.d*(updateParams.T*bestPosition-position)
	}
	// Store new position to the graph
	edge.Position[obj.Id] = position
}

func (obj *AOAObject) ConstructSolution(al *airline.Airline, graph *graph.Graph) bool {
	// Build a new solution for the object
	// Returns true if the solution covers all given pairings, false otherwise
	pilotsArray, condensedSolution, validSolution := problem.ConstructSolution(al, graph, obj.Id)
	obj.ProposedSolution = append([]*airline.Pilot(nil), pilotsArray...)
	obj.NewCondensedSolution = append([]int{}, condensedSolution...)
	return validSolution
}

func (obj *AOAObject) Evaluate() {
	// Adopt the solution that is proposed as the best solution found by the object
	obj.Solution = obj.ProposedSolution
	obj.CondensedSolution = obj.NewCondensedSolution
	obj.Fitness = obj.NewFitness
	obj.Cost = obj.NewCost
}
