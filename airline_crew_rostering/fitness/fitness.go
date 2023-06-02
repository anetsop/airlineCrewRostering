package fitness

import (
	"math"

	"go-airline-crew-rostering/airline"
)

func FitnessFunction(solution []*airline.Pilot, totalPairs int, averageWorkload float64) (float64, float64) {
	// Fitness function for the airline crew rostering problem
	// returns the solution's fitness and cost (cost is the
	// sum of each pilot's deviation from the average workload)
	pairsCovered := 0 // pairs covered by the solution
	deviation := 0.0  // sum of each pilot's deviation from the average workload
	for _, pilot := range solution {
		pairsCovered += pilot.AssignedLength
		deviation += math.Abs(averageWorkload - pilot.FlightTime)
	}
	cost := deviation
	deviation += 1
	deviation /= 750
	fitness := 1/deviation + 1/float64(totalPairs-pairsCovered+1)*0.75
	return fitness, cost
}
