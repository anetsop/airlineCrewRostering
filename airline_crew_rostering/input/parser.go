package input

import (
	"fmt"
	"os"
	"time"

	"github.com/akamensky/argparse"
)

// container for all possible arguments used by the application
type ArgumentCollection struct {
	Filename    *string   // name of the input file (relative or absolute path)
	ResultsFile *string   // name of file to write the results of the application (only the name)
	StartDate   time.Time // start date of the schedule
	EndDate     time.Time // end date of the schedule
	Pilots      *int      // number of available pilots
	Seed        *int      // seed for random number generator
	Generations *int      // maximum iterations of the optimization algorithm
	Agents      *int      // number of agents of the optimization algorithm
	FL          *float64  // FL parameter used by multi-step CSO
	Constants   []float64 // list of parameters (C1, C2, C3, C4) used by AOA
	Algorithm   string    // name of optimization algorithm to be used (options are "multiCSO" or "AOA")
}

func SetUpParser() *ArgumentCollection {
	// Process all arguments given by the command line

	args := new(ArgumentCollection)

	// Set up all shared arguments
	parser := argparse.NewParser("main", "Solve the airline crew rostering problem!")
	args.Filename = parser.String("f", "filename", &argparse.Options{Help: "Name of the file that contains the pairs", Required: true})
	args.ResultsFile = parser.String("", "results", &argparse.Options{Help: "Name of the file to write the results", Required: false, Default: "Output.xlsx"})
	startDateArg := parser.String("", "startDate", &argparse.Options{Help: "Start date of the airline schedule given as YYYY-MM-DD", Required: false, Default: "1999-12-31"})
	endDateArg := parser.String("", "endDate", &argparse.Options{Help: "End date of the airline schedule given as YYYY-MM-DD", Required: false, Default: "2021-1-1"})
	args.Pilots = parser.Int("p", "pilots", &argparse.Options{Help: "Number of pilots available", Required: false, Default: 45})
	args.Seed = parser.Int("", "seed", &argparse.Options{Help: "Seed for random number generator", Required: false, Default: -1})

	args.Generations = parser.Int("g", "generations", &argparse.Options{Help: "Maximum iterations", Required: false, Default: 150})

	// Set up multi-step CSO specific arguments
	multiCSOParser := parser.NewCommand("multiCSO", "Use chicken swarm optimization to solve the problem")
	chickens := multiCSOParser.Int("", "chickens", &argparse.Options{Help: "Number of chickens in swarm", Required: false, Default: 20})
	args.FL = multiCSOParser.Float("", "FL", &argparse.Options{Help: "Parameter for CSO algorithm", Required: false, Default: 0.5})

	// Set up AOA specific arguments
	AOAParser := parser.NewCommand("AOA", "Use Archimedes optimization algorithm to solve the problem")
	objects := AOAParser.Int("", "objects", &argparse.Options{Help: "Number of objects in the object collection", Required: false, Default: 20})
	C1 := AOAParser.Float("", "C1", &argparse.Options{Help: "C1 Parameter for AOA algorithm", Required: false, Default: 2.0})
	C2 := AOAParser.Float("", "C2", &argparse.Options{Help: "C2 Parameter for AOA algorithm", Required: false, Default: 6.0})
	C3 := AOAParser.Float("", "C3", &argparse.Options{Help: "C3 Parameter for AOA algorithm", Required: false, Default: 1.0})
	C4 := AOAParser.Float("", "C4", &argparse.Options{Help: "C4 Parameter for AOA algorithm", Required: false, Default: 0.5})

	err := parser.Parse(os.Args) // parse arguments
	if err != nil {
		fmt.Println(parser.Usage(err))
		return nil
	}

	// form the start and end of the schedule
	var year, month, day int
	fmt.Sscanf(*startDateArg, "%d-%d-%d", &year, &month, &day)
	args.StartDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	fmt.Sscanf(*endDateArg, "%d-%d-%d", &year, &month, &day)
	args.EndDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Adjust the non shared arguments based on the optimization algorithm selection
	if multiCSOParser.Happened() {
		args.Algorithm = "multiCSO"
		args.Agents = chickens
	} else if AOAParser.Happened() {
		args.Algorithm = "AOA"
		args.Agents = objects
		args.Constants = []float64{*C1, *C2, *C3, *C4}
	}

	// Store the path to the output file (it will be saved in the output subfolder)
	*args.ResultsFile = "./output/" + *args.ResultsFile
	os.Mkdir("output", 0777) // create the subfolder, if it does not exist
	return args
}
