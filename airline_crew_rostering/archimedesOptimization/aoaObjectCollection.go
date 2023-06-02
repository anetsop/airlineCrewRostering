package archimedesOptimization

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"go-airline-crew-rostering/airline"
	"go-airline-crew-rostering/fitness"
	"go-airline-crew-rostering/graph"
	"go-airline-crew-rostering/metrics"
)

// Container for the functions related to a collection of AOA objects
type AOARepo interface {
	Initialization() *AOAObjectCollection
	AOA() []*AOAObject
	solution() []*airline.Pilot
	collectionUpdate()
	sort() []*AOAObject
}

// Collection of AOA objects
type AOAObjectCollection struct {
	Collection     []*AOAObject     // List of the objects of the collection
	population     int              // Number of objects in the collection
	maxGenerations int              // Maximum number of iterations executed by the optimization algorithm
	params         *aoaParameters   // all parameters of the algoritmh that are shared between the objects
	Mtr            *metrics.Metrics // Metrics used to evaluate the algorithm's efficiency
	costList       []float64        // List of solutions' cost found by all objects of the collection in the current generation
}

type aoaParameters struct {
	// Parameters used in the update step that are the same for all objects
	TF         float64
	d          float64
	C1         float64
	C2         float64
	C3         float64
	C4         float64
	F          int
	T          float64
	bestObject *AOAObject // Object with the best fitness found so far
}

func (collection *AOAObjectCollection) Initialization(population int, maxGenerations int,
	C1 float64, C2 float64, C3 float64, C4 float64, al *airline.Airline, pairGraph *graph.Graph,
	Mtr *metrics.Metrics) *AOAObjectCollection {
	// Initialize an instance of a collection of AOA objects
	collection.Collection = []*AOAObject{}
	collection.population = population
	collection.maxGenerations = maxGenerations
	collection.params = &aoaParameters{
		TF:         0,
		d:          0,
		C1:         C1,
		C2:         C2,
		C3:         C3,
		C4:         C4,
		F:          0,
		T:          0,
		bestObject: nil,
	}
	collection.costList = []float64{}
	collection.Mtr = Mtr

	// Create the objects and build the initial solutions for each one
	for agent := 0; agent < population; agent++ {
		object := new(AOAObject)
		object.Initialization(agent)
		collection.Collection = append(collection.Collection, object)
		collection.solution(al, pairGraph, agent)
	}
	return collection
}

func (collection *AOAObjectCollection) solution(al *airline.Airline, pairGraph *graph.Graph, objectId int) {
	// Build a new solution for the object with id "objectID"
	object := collection.Collection[objectId]
	validSolution := object.ConstructSolution(al, pairGraph)
	if validSolution {
		collection.Mtr.ValidSolutions++
	}
	// Calculate the new solution's fitness and cost
	object.NewFitness, object.NewCost = fitness.FitnessFunction(object.ProposedSolution, len(al.PairsArray)-1, al.AverageWorkload)
	object.NewCost = object.NewCost * collection.Mtr.UnitCost

	collection.costList = append(collection.costList, object.NewCost)

	// Adopt the new solution as the object's solution if the new solution is better
	if object.NewFitness > object.Fitness {
		object.Evaluate()
		// Replace the best object if the new solution is the best found overall
		if collection.params.bestObject == nil || collection.params.bestObject.Fitness < object.Fitness {
			collection.params.bestObject = object
		}
	}
}

func (collection *AOAObjectCollection) collectionUpdate(pairGraph *graph.Graph, generation int) {
	// Update all relevant edges of the graph for all objects of the collection

	edgesToUpdate := make(map[int]*graph.Edge) // List of edges that will be updated

	// Calculated shared parameters
	collection.params.TF = math.Exp((float64(generation-collection.maxGenerations) / float64(collection.maxGenerations)))
	collection.params.d = math.Exp((float64(collection.maxGenerations-generation) / float64(collection.maxGenerations))) - (float64(generation) / float64(collection.maxGenerations))
	P := 2*rand.Float64() - collection.params.C4
	if P <= 0.5 {
		collection.params.F = 1
	} else if P > 0.5 {
		collection.params.F = -1
	}
	collection.params.T = collection.params.C3 * collection.params.TF

	maxAccelerationObject := collection.Collection[0]
	minAccelerationObject := collection.Collection[0]
	// Calculate individual parameters for each object
	for _, object := range collection.Collection {
		object.SetUpdateParameters(collection.Collection, collection.params.bestObject.Id, collection.params.TF)
		if object.accelleration > maxAccelerationObject.accelleration {
			maxAccelerationObject = object
		}
		if object.accelleration < minAccelerationObject.accelleration {
			minAccelerationObject = object
		}
	}
	maxAcceleration := maxAccelerationObject.accelleration
	minAcceleration := minAccelerationObject.accelleration

	// normalize acceleration
	for _, object := range collection.Collection {
		object.accelleration = 0.9*(object.accelleration-minAcceleration)/(maxAcceleration-minAcceleration) + 0.1
	}

	// Gather all edges that need to be updated by scanning the
	// current solution of each object
	for _, object := range collection.Collection {
		for _, pilot := range object.Solution {
			for i := 0; i < pilot.AssignedLength; i++ {
				sourcePairId := pilot.AssignedPairs[i].Id
				goalPairId := pilot.AssignedPairs[i+1].Id
				edge, edgeExists := pairGraph.Nodes[sourcePairId].Edges[goalPairId]
				// Create a new edge if a pair of successive pairings has not be
				// found in another solution so far
				if !edgeExists {
					pairGraph.AddEdge(sourcePairId, goalPairId)
					edge = pairGraph.Nodes[sourcePairId].Edges[goalPairId]
					for i := 0; i < collection.population; i++ {
						edge.Position[i] = 0.95 + rand.Float64()*0.05
					}
				}
				edgesToUpdate[edge.Id] = edge
			}
		}
	}

	// Update all relevant edges for all objects
	for _, edge := range edgesToUpdate {
		for _, object := range collection.Collection {
			object.UpdatePosition(edge, collection.params)
		}
	}

}

func (collection *AOAObjectCollection) sort() {
	// sort the list of AOA objects by their fitness, from best to worst
	sort.Slice(collection.Collection, func(i int, j int) bool {
		return collection.Collection[i].Fitness > collection.Collection[j].Fitness
	})
}

func (collection *AOAObjectCollection) AOA(al *airline.Airline, pairGraph *graph.Graph) []*AOAObject {
	// Main body of the optimization algorithm
	// Returns the list of objects of the collection

	// Calculate Metrics for the initialization step
	globalbestobject, _ := collection.Mtr.SetUpIterationMetrics(collection.costList)
	collection.Mtr.Jumps++
	object := collection.Collection[globalbestobject]

	collection.Mtr.GlobalBestSolutionCost = object.Cost
	collection.Mtr.GlobalBestString = collection.Mtr.SolutionEncoding(object.CondensedSolution, object.Solution)
	for _, object := range collection.Collection {
		if object.Id == globalbestobject {
			continue
		}
		normalisedSolution := collection.Mtr.SolutionEncoding(object.CondensedSolution, object.Solution)
		collection.Mtr.AverageSimilarity += collection.Mtr.SolutionSimilarity(collection.Mtr.GlobalBestString, normalisedSolution)
	}

	// Execute the algorithm for iterations equal to "maxGenerations"
	for t := 1; t < collection.maxGenerations; t++ {

		// Update the positions of the collection
		collection.collectionUpdate(pairGraph, t-1)

		collection.costList = []float64{} // empty the cost list from the previous iteration

		// build solutions for all objects
		for _, object := range collection.Collection {
			collection.solution(al, pairGraph, object.Id)
		}

		// Calculate the metrics of the current iteration
		globalbestobject = collection.calculateMetrics(globalbestobject, t)

	}
	// Try to optimize the solutions of each object
	for _, object := range collection.Collection {
		al.EqualizeWorkload(object.Solution)
		object.Fitness, object.Cost = fitness.FitnessFunction(object.Solution, len(al.PairsArray)-1, al.AverageWorkload)
		object.Cost *= collection.Mtr.UnitCost
	}
	collection.Mtr.AverageSimilarity = collection.Mtr.AverageSimilarity / float64(collection.population*collection.maxGenerations)
	collection.sort()
	return collection.Collection
}

func (collection *AOAObjectCollection) calculateMetrics(globalBest int, generation int) int {
	// Calculate metrics of current iteration and update all metrics
	globalbestobject := globalBest
	globalbestFitness := collection.Collection[globalbestobject].Fitness
	bestObject, _ := collection.Mtr.SetUpIterationMetrics(collection.costList)
	bestFitness := collection.Collection[bestObject].Fitness
	bestCost := collection.costList[bestObject]
	bestSolutionString := []string{}
	for _, object := range collection.Collection {
		normalisedSolution := collection.Mtr.SolutionEncoding(object.NewCondensedSolution, object.ProposedSolution)
		collection.Mtr.AverageSimilarity += collection.Mtr.SolutionSimilarity(collection.Mtr.GlobalBestString, normalisedSolution)
		if object.Id == bestObject {
			bestSolutionString = normalisedSolution
		}
	}
	if bestFitness > globalbestFitness {
		globalbestobject = bestObject
		collection.Mtr.GlobalBestSolutionCost = bestCost
		collection.Mtr.GlobalBestString = bestSolutionString
		collection.Mtr.Jumps++
	}
	if generation%200 == 0 {
		fmt.Println("Generation", generation)
	}
	return globalbestobject
}
