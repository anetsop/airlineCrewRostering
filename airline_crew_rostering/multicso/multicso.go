package multicso

import (
	"fmt"
	"sort"

	"go-airline-crew-rostering/airline"
	"go-airline-crew-rostering/fitness"
	"go-airline-crew-rostering/graph"
	"go-airline-crew-rostering/metrics"
)

// Container for the functions related to a chicken swarm
type MultiCSORepo interface {
	Initialization() *MultiCSO
	MultiCSO() []*Chicken
	solution() []*airline.Pilot
	swarmUpdate()
	sort() []*Chicken
}

// Swarm of chickens
type MultiCSO struct {
	Swarm          []*Chicken       // List of the chickens of the swarm
	population     int              // Number of objects in the swarm
	maxGenerations int              // Maximum number of iterations executed by the optimization algorithm
	FL             float64          // algorithm parameter shared by all chickens
	Mtr            *metrics.Metrics // Metrics used to evaluate the algorithm's efficiency
	costList       []float64        // List of solutions' cost found by all chickens of the swarm in the current generation
}

func (swarm *MultiCSO) Initialization(al *airline.Airline, pairGraph *graph.Graph,
	population int, maxGenerations int, FL float64, Mtr *metrics.Metrics) *MultiCSO {
	// Initialize an instance of a chicken swarm
	swarm.Swarm = []*Chicken{}
	swarm.population = population
	swarm.maxGenerations = maxGenerations
	swarm.FL = FL
	swarm.costList = []float64{}
	swarm.Mtr = Mtr

	// Create the chickens and build the initial solutions for each one
	for agent := 0; agent < population; agent++ {
		chicken := new(Chicken)
		chicken.Initialization(agent)
		swarm.Swarm = append(swarm.Swarm, chicken)
		swarm.solution(al, pairGraph, agent)
	}

	return swarm
}

func (swarm *MultiCSO) solution(al *airline.Airline, pairGraph *graph.Graph, chickenId int /*wg *sync.WaitGroup*/) {
	// Build a new solution for the chicken with id "chickenId"
	chicken := swarm.Swarm[chickenId]
	validSolution := chicken.ConstructSolution(al, pairGraph)
	if validSolution {
		swarm.Mtr.ValidSolutions++
	}
	// Calculate the new solution's fitness and cost
	chicken.NewFitness, chicken.NewCost = fitness.FitnessFunction(chicken.ProposedSolution, len(al.PairsArray)-1, al.AverageWorkload)
	chicken.NewCost = chicken.NewCost * swarm.Mtr.UnitCost
	swarm.costList = append(swarm.costList, chicken.NewCost)

	// Adopt the new solution as the chicken's solution if the new solution is better
	if chicken.NewFitness > chicken.Fitness {
		chicken.Evaluate()
	}
}

func (swarm *MultiCSO) swarmUpdate(pairGraph *graph.Graph) {
	// Update all relevant edges of the graph for all chickens of the swarm

	edgesToUpdate := make(map[int]*graph.Edge) // List of edges that will be updated

	// Calculate individual parameters for each chicken
	for _, chicken := range swarm.Swarm {
		chicken.SetUpdateParameters(swarm.Swarm)
	}

	// Gather all edges that need to be updated by scanning the
	// current solution of each chicken
	for _, chicken := range swarm.Swarm {
		for _, pilot := range chicken.Solution {
			for i := 0; i < pilot.AssignedLength; i++ {
				sourcePairId := pilot.AssignedPairs[i].Id
				goalPairId := pilot.AssignedPairs[i+1].Id
				edge, edgeExists := pairGraph.Nodes[sourcePairId].Edges[goalPairId]
				// Create a new edge if a pair of successive pairings has not be
				// found in another solution so far
				if !edgeExists {
					pairGraph.AddEdge(sourcePairId, goalPairId)
					edge = pairGraph.Nodes[sourcePairId].Edges[goalPairId]
				}
				edgesToUpdate[edge.Id] = edge
			}
		}
	}

	// Update all relevant edges for all chickens
	for _, edge := range edgesToUpdate {
		for _, chicken := range swarm.Swarm {
			chicken.UpdatePosition(edge, swarm.FL)
		}
	}

}

func (swarm *MultiCSO) sort() {
	// sort the list of chickens by their fitness, from best to worst
	sort.Slice(swarm.Swarm, func(i int, j int) bool {
		return swarm.Swarm[i].Fitness > swarm.Swarm[j].Fitness
	})
}

func (swarm *MultiCSO) MultiCSO(al *airline.Airline, pairGraph *graph.Graph) []*Chicken {
	// Main body of the optimization algorithm
	// Returns the list of chickens in the swarm

	// Calculate Metrics for the initialization step
	globalbestchicken, _ := swarm.Mtr.SetUpIterationMetrics(swarm.costList)
	swarm.Mtr.Jumps++
	chicken := swarm.Swarm[globalbestchicken]
	swarm.Mtr.GlobalBestSolutionCost = chicken.Cost
	swarm.Mtr.GlobalBestString = swarm.Mtr.SolutionEncoding(chicken.CondensedSolution, chicken.Solution)
	for _, chicken := range swarm.Swarm {
		if chicken.Id == globalbestchicken {
			continue
		}
		normalisedSolution := swarm.Mtr.SolutionEncoding(chicken.CondensedSolution, chicken.Solution)
		swarm.Mtr.AverageSimilarity += swarm.Mtr.SolutionSimilarity(swarm.Mtr.GlobalBestString, normalisedSolution)
	}

	// Execute the algorithm for iterations equal to "maxGenerations"
	for t := 1; t < swarm.maxGenerations; t++ {

		// Update the positions of the swarm
		swarm.swarmUpdate(pairGraph)

		swarm.costList = []float64{} // empty the cost list from the previous iteration

		// build solutions for all chickens
		for _, chicken := range swarm.Swarm {
			swarm.solution(al, pairGraph, chicken.Id)
		}

		// Calculate the metrics of the current iteration
		globalbestchicken = swarm.calculateMetrics(globalbestchicken, t)
	}

	// Try to optimize the solutions of each object
	for _, chicken := range swarm.Swarm {
		al.EqualizeWorkload(chicken.Solution)
		chicken.Fitness, chicken.Cost = fitness.FitnessFunction(chicken.Solution, len(al.PairsArray)-1, al.AverageWorkload)
		chicken.Cost *= swarm.Mtr.UnitCost
	}
	swarm.Mtr.AverageSimilarity = swarm.Mtr.AverageSimilarity / float64(swarm.population*swarm.maxGenerations)
	swarm.sort()
	return swarm.Swarm
}

func (swarm *MultiCSO) calculateMetrics(globalBest int, generation int) int {
	// Calculate metrics of current iteration and update all metrics
	globalbestchicken := globalBest
	globalbestFitness := swarm.Swarm[globalbestchicken].Fitness
	bestchicken, _ := swarm.Mtr.SetUpIterationMetrics(swarm.costList)
	bestFitness := swarm.Swarm[bestchicken].Fitness
	bestCost := swarm.costList[bestchicken]
	bestSolutionString := []string{}
	for _, chicken := range swarm.Swarm {
		normalisedSolution := swarm.Mtr.SolutionEncoding(chicken.NewCondensedSolution, chicken.ProposedSolution)
		swarm.Mtr.AverageSimilarity += swarm.Mtr.SolutionSimilarity(swarm.Mtr.GlobalBestString, normalisedSolution)
		if chicken.Id == bestchicken {
			bestSolutionString = normalisedSolution
		}
	}
	if bestFitness > globalbestFitness {
		globalbestchicken = bestchicken
		swarm.Mtr.GlobalBestSolutionCost = bestCost
		swarm.Mtr.GlobalBestString = bestSolutionString
		swarm.Mtr.Jumps++
	}
	if generation%200 == 0 {
		fmt.Println("Generation", generation)
	}
	return globalbestchicken
}
