package input_test

import (
	"testing"
	"time"

	"go-airline-crew-rostering/input"
)

func TestInput(t *testing.T) {
	// simple test for the input package
	filename := "../Pairings.csv"
	startSchedule := time.Date(2011, 11, 5, 0, 0, 0, 0, time.UTC)
	endSchedule := time.Date(2011, 11, 6, 23, 59, 0, 0, time.UTC)
	pairsArray := input.ReadFile(filename, startSchedule)

	t.Log("pair ", pairsArray[0].Id, ":", pairsArray[0].Start, pairsArray[0].End)
	pairsArray = input.FilterPairs(pairsArray, startSchedule, endSchedule)
	t.Log("number of filtered pairs:", len(pairsArray))

	for _, pair := range pairsArray {
		t.Log(pair.Id, pair.Start, pair.End)
	}
}
