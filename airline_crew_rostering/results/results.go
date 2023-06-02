package results

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"go-airline-crew-rostering/airline"
	"go-airline-crew-rostering/input"
	"go-airline-crew-rostering/metrics"

	"github.com/xuri/excelize/v2"
)

func PrintResults(m *metrics.Metrics, args *input.ArgumentCollection, al *airline.Airline) {
	// Creates an excel file to store an airline crew rostering schedule, along with various
	// statistics
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	docProperties, _ := f.GetDocProps()
	docProperties.Language = "en-UK"
	f.SetDocProps(docProperties)

	f.SetDefaultFont("Arial")

	generalSheetName := "General Information"
	solutionSheetName := "Solution Statistics"
	algorithmSheetName := "Optimization Algorithm"
	scheduleSheetName := "Schedule"
	pairingsSheetName := "Pairings"

	f.SetSheetName("Sheet1", generalSheetName)

	if _, err := f.NewSheet(solutionSheetName); err != nil {
		fmt.Println(err)
		return
	}
	if _, err := f.NewSheet(algorithmSheetName); err != nil {
		fmt.Println(err)
		return
	}
	if _, err := f.NewSheet(scheduleSheetName); err != nil {
		fmt.Println(err)
		return
	}
	if _, err := f.NewSheet(pairingsSheetName); err != nil {
		fmt.Println(err)
		return
	}

	drawGeneralSheet(f, m, args, al)
	drawSolutionStatisticsSheet(f, m, args, al)
	drawOptimizationAlgorithmSheet(f, m, args, al)
	drawScheduleSheet(f, al)
	drawPairingsSheet(f, m, args, al)

	index, _ := f.GetSheetIndex(scheduleSheetName)
	f.SetActiveSheet(index)

	if err := f.SaveAs(*args.ResultsFile); err != nil {
		fmt.Println(err)
	}

}

func setView(f *excelize.File, sheetName string, zoom float64) {
	// set view options for an excel sheet
	options, err := f.GetSheetView(sheetName, 0)
	if err != nil {
		fmt.Println(err)
	}
	*options.ShowGridLines = false
	*options.ShowRowColHeaders = false
	*options.ZoomScale = zoom
	f.SetSheetView(sheetName, 0, &options)
}

func drawGeneralSheet(f *excelize.File, m *metrics.Metrics, args *input.ArgumentCollection, al *airline.Airline) {
	// create an excel sheet containing general information of the application
	sheetName := "General Information"
	algorithmName := ""
	if args.Algorithm == "multiCSO" {
		algorithmName = "Multi-step CSO"
	} else if args.Algorithm == "AOA" {
		algorithmName = "Archimedes Optimization"
	}

	setView(f, sheetName, 100.0)

	f.SetColWidth(sheetName, "G", "J", 11.62)
	f.SetColWidth(sheetName, "L", "O", 11.62)
	for i := 5; i <= 11; i++ {
		f.SetRowHeight(sheetName, i, 30)
	}

	f.SetCellValue(sheetName, "B2", "General Information")
	f.SetCellValue(sheetName, "B5", "Input File")
	f.SetCellValue(sheetName, "B6", "Algorithm")
	f.SetCellValue(sheetName, "B7", "Execution Time")
	f.SetCellValue(sheetName, "D5", *args.Filename)
	f.SetCellValue(sheetName, "D6", algorithmName)

	executionstring := m.TotalTime.String()

	i := strings.Index(executionstring, "m")
	if i > -1 && i < len(executionstring) && executionstring[i+1:i+2] != "s" {
		executionstring = executionstring[:i] + " min. " + executionstring[i+1:]
	}
	i = strings.Index(executionstring, "s")
	if i > -1 && executionstring[i-1:i] != "m" {
		executionstring = executionstring[:i] + " sec." + executionstring[i+1:]
	}
	f.SetCellValue("General Information", "D7", executionstring)

	drawVerticalTable(f, sheetName, "B2", 3, "4F81BD", "B8CCE4")

	f.SetCellValue(sheetName, "G2", "Airline Crew Rostering Information")
	f.SetCellValue(sheetName, "G5", "Start Date")
	f.SetCellValue(sheetName, "G6", "End Date")
	f.SetCellValue(sheetName, "G7", "Pilots")
	f.SetCellValue(sheetName, "G8", "Pairs")
	f.SetCellValue(sheetName, "G9", "Optimal Workload")

	startdatestring := strconv.Itoa(args.StartDate.Day()) + " " + args.StartDate.Month().String() + " " + strconv.Itoa(args.StartDate.Year())
	enddatestring := strconv.Itoa(args.EndDate.Day()) + " " + args.EndDate.Month().String() + " " + strconv.Itoa(args.EndDate.Year())

	f.SetCellValue(sheetName, "I5", startdatestring)
	f.SetCellValue(sheetName, "I6", enddatestring)
	f.SetCellValue(sheetName, "I7", *args.Pilots)
	f.SetCellValue(sheetName, "I8", len(al.PairsArray)-1)
	f.SetCellValue(sheetName, "I9", math.Round(al.AverageWorkload))

	drawVerticalTable(f, sheetName, "G2", 5, "C0504D", "E6B9B8")

	f.SetCellValue(sheetName, "L2", "Optimization Algorithm Information")
	f.SetCellValue(sheetName, "L5", "Agents")
	f.SetCellValue(sheetName, "L6", "Generations")
	f.SetCellValue(sheetName, "L7", "Seed")
	f.SetCellValue(sheetName, "N5", *args.Agents)
	f.SetCellValue(sheetName, "N6", *args.Generations)

	var rows int
	if args.Algorithm == "multiCSO" {
		f.SetCellValue(sheetName, "L8", "FL")
		f.SetCellValue(sheetName, "N8", *args.FL)
		rows = 4
	} else if args.Algorithm == "AOA" {
		f.SetCellValue(sheetName, "L8", "C1")
		f.SetCellValue(sheetName, "L9", "C2")
		f.SetCellValue(sheetName, "L10", "C3")
		f.SetCellValue(sheetName, "L11", "C4")
		f.SetCellValue(sheetName, "N8", args.Constants[0])
		f.SetCellValue(sheetName, "N9", args.Constants[1])
		f.SetCellValue(sheetName, "N10", args.Constants[2])
		f.SetCellValue(sheetName, "N11", args.Constants[3])
		rows = 7
	}

	drawVerticalTable(f, sheetName, "L2", rows, "4CB7C8", "B6E9CE")
	if *args.Seed > -1 {
		styleId, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Size: 12},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
			Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"B6E9CE"}},
			Border:    []excelize.Border{{Type: "top", Color: "FFFFFF", Style: 1}, {Type: "left", Color: "FFFFFF", Style: 1}, {Type: "bottom", Color: "FFFFFF", Style: 1}},
			NumFmt:    49,
		})
		f.SetCellStyle(sheetName, "N7", "O7", styleId)
		f.SetCellValue(sheetName, "N7", strconv.Itoa(*args.Seed))
	} else {
		f.SetCellValue(sheetName, "N7", "None")
	}
}

func drawSolutionStatisticsSheet(f *excelize.File, m *metrics.Metrics, args *input.ArgumentCollection, al *airline.Airline) {
	// create an excel sheet containing statistics related to the solution
	sheetName := "Solution Statistics"

	setView(f, sheetName, 100.0)

	tableStyleTitle := excelize.Style{
		Font:      &excelize.Font{Size: 16, Color: "FFFFFF", Bold: true, Underline: "single"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4CB3C8"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    []excelize.Border{{Type: "bottom", Color: "FFFFFF", Style: 5}},
	}
	tableStyleCellsHeader := excelize.Style{
		Font:      &excelize.Font{Size: 12, Color: "FFFFFF", Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4CB3C8"}},
		Border: []excelize.Border{{Type: "top", Color: "FFFFFF", Style: 5}, {Type: "right", Color: "FFFFFF", Style: 1},
			{Type: "bottom", Color: "FFFFFF", Style: 5}, {Type: "left", Color: "FFFFFF", Style: 1}},
	}
	tableStyleCellsData := excelize.Style{
		Font:      &excelize.Font{Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"B6DDE8"}},
		Border: []excelize.Border{{Type: "top", Color: "FFFFFF", Style: 1}, {Type: "left", Color: "FFFFFF", Style: 1},
			{Type: "bottom", Color: "FFFFFF", Style: 1}, {Type: "left", Color: "FFFFFF", Style: 1}},
	}

	f.SetColWidth(sheetName, "B", "E", 11.62)
	f.SetColWidth(sheetName, "G", "J", 22.62)

	f.SetCellValue(sheetName, "B2", "Solution Information")
	f.SetCellValue(sheetName, "B5", "Assigned Pairs")
	f.SetCellValue(sheetName, "B6", "Unassigned Pairs")
	f.SetCellValue(sheetName, "B7", "Average Rest Period")
	f.SetCellValue(sheetName, "B8", "Average Days Off")
	f.SetCellValue(sheetName, "B9", "Best Cost")
	f.SetCellValue(sheetName, "B10", "Worst Cost")
	f.SetCellValue(sheetName, "B11", "Total Deviation")

	f.SetCellValue(sheetName, "D5", m.TotalAssignedPairs)
	f.SetCellValue(sheetName, "D6", len(al.PairsArray)-1-m.TotalAssignedPairs)
	f.SetCellValue(sheetName, "D7", math.Round(m.AverageRestPeriod*100)/100)
	f.SetCellValue(sheetName, "D8", math.Round(m.AverageDaysOff*100)/100)
	f.SetCellValue(sheetName, "D9", math.Round(m.GlobalBestSolutionCost))

	worstCost := m.IterWorstCost[0]
	for i := range m.IterWorstCost {
		if worstCost < m.IterWorstCost[i] {
			worstCost = m.IterWorstCost[i]
		}
	}
	f.SetCellValue(sheetName, "D10", math.Round(worstCost))
	f.SetCellValue(sheetName, "D11", math.Round(m.GlobalBestSolutionCost/m.UnitCost))

	drawVerticalTable(f, sheetName, "B2", 7, "7666A4", "CCC0DA")

	f.SetCellValue(sheetName, "G2", "Pilot Statistics")
	f.SetCellValue(sheetName, "G5", "Pilot")
	f.SetCellValue(sheetName, "H5", "Flight time\n (in minutes)")
	f.SetCellValue(sheetName, "I5", "Deviation from\n optimal workload")
	f.SetCellValue(sheetName, "J5", "Assigned Pairs")

	styleId, _ := f.NewStyle(&tableStyleTitle)
	f.SetCellStyle(sheetName, "G2", "J4", styleId)
	styleId, _ = f.NewStyle(&tableStyleCellsHeader)
	f.SetCellStyle(sheetName, "G5", "J5", styleId)
	f.MergeCell(sheetName, "G2", "J4")

	styleId, _ = f.NewStyle(&tableStyleCellsData)
	fillColor := &tableStyleCellsData.Fill.Color[0]
	*fillColor = "DBEEF3"
	styleId2, _ := f.NewStyle(&tableStyleCellsData)

	f.SetRowHeight(sheetName, 5, 40.0)
	for i, pilot := range al.PilotsArray {
		f.SetRowHeight(sheetName, 6+i, 40.0)
		cell, _ := excelize.CoordinatesToCellName(7, 6+i)
		f.SetCellValue(sheetName, cell, pilot.Id)
		cell, _ = excelize.CoordinatesToCellName(8, 6+i)
		f.SetCellValue(sheetName, cell, math.Round(pilot.FlightTime))
		cell, _ = excelize.CoordinatesToCellName(9, 6+i)
		f.SetCellValue(sheetName, cell, math.Round(math.Abs(al.AverageWorkload-pilot.FlightTime)))
		cell, _ = excelize.CoordinatesToCellName(10, 6+i)
		f.SetCellValue(sheetName, cell, pilot.AssignedLength)
		start, _ := excelize.CoordinatesToCellName(7, 6+i)
		end, _ := excelize.CoordinatesToCellName(10, 6+i)
		if i%2 == 0 {
			f.SetCellStyle(sheetName, start, end, styleId)
		} else {
			f.SetCellStyle(sheetName, start, end, styleId2)
		}
	}

	start, _ := excelize.CoordinatesToCellName(7, 5)
	end, _ := excelize.CoordinatesToCellName(10, 5+(*args.Pilots))
	enable := true
	f.AddTable(sheetName, &excelize.Table{
		Range:             start + ":" + end,
		Name:              "PilotStatistics",
		StyleName:         "TableStyleMedium20",
		ShowFirstColumn:   false,
		ShowLastColumn:    false,
		ShowRowStripes:    &enable,
		ShowColumnStripes: false,
	})

}

func drawOptimizationAlgorithmSheet(f *excelize.File, m *metrics.Metrics, args *input.ArgumentCollection, al *airline.Airline) {
	// create an excel sheet containing statistics related to the optimization algorithm
	sheetName := "Optimization Algorithm"

	setView(f, sheetName, 100.0)

	f.SetColWidth(sheetName, "B", "E", 11.62)
	for i := 5; i <= 10; i++ {
		f.SetRowHeight(sheetName, i, 30.0)
	}

	f.SetCellValue(sheetName, "B2", "Optimization Algorithm Statistics")
	f.SetCellValue(sheetName, "B5", "Valid Solutions")
	f.SetCellValue(sheetName, "B6", "Invalid Solutions")
	f.SetCellValue(sheetName, "B7", "Total Solutions")
	f.SetCellValue(sheetName, "B8", "Unique Solutions")
	f.SetCellValue(sheetName, "B9", "Jumps")
	f.SetCellValue(sheetName, "B10", "Similarity")
	f.SetCellValue(sheetName, "D5", m.ValidSolutions)

	f.SetCellValue(sheetName, "D6", m.TotalSolutions-m.ValidSolutions)
	f.SetCellValue(sheetName, "D7", m.TotalSolutions)
	f.SetCellValue(sheetName, "D8", m.UniqueCount)
	f.SetCellValue(sheetName, "D9", m.Jumps)
	f.SetCellValue(sheetName, "D10", math.Round(m.AverageSimilarity*100)/10000)

	drawVerticalTable(f, sheetName, "B2", 6, "4F81BD", "B8CCE4")

	styleId, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"B8CCE4"}},
		Border:    []excelize.Border{{Type: "top", Color: "FFFFFF", Style: 1}, {Type: "left", Color: "FFFFFF", Style: 1}},
		NumFmt:    10,
	})
	f.SetCellStyle(sheetName, "D10", "E10", styleId)
	outputFile := strings.Split(*args.ResultsFile, ".")
	plotFilename := outputFile[0] + "_iterations.png"
	drawIterationMetricsPlot(plotFilename, m, args)
	if err := f.AddPicture(sheetName, "H2", plotFilename, &excelize.GraphicOptions{AltText: "Iterations' Comparison"}); err != nil {
		fmt.Println(err)
	}
	os.Remove(plotFilename)
}

func drawScheduleSheet(f *excelize.File, al *airline.Airline) {
	// create an excel sheet containing the airline crew rostering schedule
	sheetName := "Schedule"

	f.SetRowHeight(sheetName, 1, 7)
	f.SetColWidth(sheetName, "A", "A", 1)
	f.SetColWidth(sheetName, "B", "C", 7.65)

	setView(f, sheetName, 82.0)

	scheduleStartDate := al.ScheduleStart
	scheduleEndDate := al.ScheduleEnd
	if len(al.PairsArray) > 1 {
		firstPair := al.PairsArray[1].Start
		scheduleStartDate = time.Date(firstPair.Year(), firstPair.Month(), 1, 0, 0, 0, 0, time.UTC)
		lastPair := al.PairsArray[len(al.PairsArray)-1].End
		for _, pair := range al.PairsArray {
			if pair.End.After(lastPair) {
				lastPair = pair.End
			}
		}
		days := daysInMonth(lastPair.Year(), lastPair.Month().String())
		scheduleEndDate = time.Date(lastPair.Year(), lastPair.Month(), days, 0, 0, 0, 0, time.UTC)
	}

	monthsIndex := make(map[string]int)
	monthsInSchedule := []*dateInfo{}
	for year := scheduleStartDate.Year(); year <= scheduleEndDate.Year(); year++ {
		month := 1
		if year == scheduleStartDate.Year() {
			month = int(scheduleStartDate.Month())
		}
		for ; month <= 12; month++ {
			dateString := time.Month(month).String() + " " + strconv.Itoa(year)
			monthInfo := &dateInfo{year: year, month: month, dateString: dateString, startCell: ""}
			monthsInSchedule = append(monthsInSchedule, monthInfo)
			monthsIndex[monthInfo.dateString] = len(monthsInSchedule) - 1
			if year == scheduleEndDate.Year() && month == int(scheduleEndDate.Month()) {
				break
			}
		}
	}
	startCell := "B2"
	for i := 0; i < len(monthsInSchedule); i += 2 {
		j := i + 2
		if i == len(monthsInSchedule)-1 {
			j = i + 1
		}
		col, row, _ := excelize.CellNameToCoordinates(startCell)
		drawScheduleHeader(f, sheetName, startCell)
		cell, _ := excelize.CoordinatesToCellName(col, row+2)
		drawPilotColumn(f, sheetName, cell, al.NumberOfPilots)
		cell, _ = excelize.CoordinatesToCellName(col+2, row)
		drawMonths(f, sheetName, cell, monthsInSchedule[i:j])
		drawDataArea(f, sheetName, monthsInSchedule[i:j], al.NumberOfPilots)
		startCell, _ = excelize.CoordinatesToCellName(col, row+al.NumberOfPilots+3)
	}

	for _, pilot := range al.PilotsArray {
		for i := 1; i <= pilot.AssignedLength; i++ {
			pair := pilot.AssignedPairs[i]
			dateString := pair.Start.Month().String() + " " + strconv.Itoa(pair.Start.Year())
			monthInfo := monthsInSchedule[monthsIndex[dateString]]

			pairStartCell := findCell(monthInfo.startCell, pair.Start, pilot.Id)

			f.SetCellValue(sheetName, pairStartCell, "F")

			if pair.Start.Month() == pair.End.Month() {
				pairEndCell := findCell(monthInfo.startCell, pair.End, pilot.Id)
				drawDashes(f, sheetName, pairStartCell, pairEndCell)
			} else if pair.Start.Month() != pair.End.Month() {
				days := daysInMonth(pair.Start.Year(), pair.Start.Month().String())
				endOfStartMonth := time.Date(pair.Start.Year(), pair.Start.Month(), days, 0, 0, 0, 0, time.UTC)
				cell := findCell(monthInfo.startCell, endOfStartMonth, pilot.Id)
				drawDashes(f, sheetName, pairStartCell, cell)

				startOfEndMonth := time.Date(pair.End.Year(), pair.End.Month(), 1, 0, 0, 0, 0, time.UTC)
				dateString = pair.End.Month().String() + " " + strconv.Itoa(pair.End.Year())
				monthInfo = monthsInSchedule[monthsIndex[dateString]]
				pairEndCell := findCell(monthInfo.startCell, pair.End, pilot.Id)
				cell = findCell(monthInfo.startCell, startOfEndMonth, pilot.Id)
				col, row, _ := excelize.CellNameToCoordinates(cell)
				cell, _ = excelize.CoordinatesToCellName(col-1, row)
				drawDashes(f, sheetName, cell, pairEndCell)

			}

			dateString = fmt.Sprintf("Departure: %02d/%02d/%d %02d:%02d\nArrival: %02d/%02d/%d %02d:%02d",
				pair.Start.Day(), pair.Start.Month(), pair.Start.Year(), pair.Start.Hour(), pair.Start.Minute(),
				pair.End.Day(), pair.End.Month(), pair.End.Year(), pair.End.Hour(), pair.End.Minute())
			titleString := fmt.Sprintf("Pair %04d", pair.Id)
			dataValidation := excelize.NewDataValidation(true)
			dataValidation.Sqref = pairStartCell + ":" + pairStartCell
			dataValidation.SetInput(titleString, dateString)
			f.AddDataValidation(sheetName, dataValidation)
		}
	}

}

func drawPairingsSheet(f *excelize.File, m *metrics.Metrics, args *input.ArgumentCollection, al *airline.Airline) {
	// create an excel sheet containing all the pairings along with their start and end datetimes
	sheetName := "Pairings"
	setView(f, sheetName, 100)

	f.SetRowHeight(sheetName, 2, 30)
	f.SetRowHeight(sheetName, 3, 30)
	f.SetRowHeight(sheetName, 4, 36)
	f.SetColWidth(sheetName, "G", "I", 25.65)

	rows := len(al.PairsArray) - 1
	f.MergeCell(sheetName, "G2", "I3")
	f.SetCellValue(sheetName, "G2", "Pairings")
	styleId, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 16, Color: "FFFFFF", Bold: true, Underline: "single"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"7266A4"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{{Type: "bottom", Color: "FFFFFF", Style: 5},
			{Type: "left", Color: "FFFFFF", Style: 2},
			{Type: "right", Color: "FFFFFF", Style: 2}},
	})
	f.SetCellStyle(sheetName, "G2", "I3", styleId)

	f.SetCellValue(sheetName, "G4", "Pairings")
	f.SetCellValue(sheetName, "H4", "Departure")
	f.SetCellValue(sheetName, "I4", "Arrival")
	styleId, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12, Color: "FFFFFF", Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"7266A4"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{{Type: "bottom", Color: "FFFFFF", Style: 5},
			{Type: "left", Color: "FFFFFF", Style: 2},
			{Type: "right", Color: "FFFFFF", Style: 2}},
	})
	f.SetCellStyle(sheetName, "G4", "I4", styleId)

	dateTimefmt := "dd/mm/yyyy hh:mm"
	pairingsfmt := "000#"

	cellStyleId1, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CCC0DA"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{{Type: "bottom", Color: "FFFFFF", Style: 1},
			{Type: "left", Color: "FFFFFF", Style: 2},
			{Type: "right", Color: "FFFFFF", Style: 2}},
		CustomNumFmt: &pairingsfmt,
	})

	cellStyleId2, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"E5E0EC"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{{Type: "bottom", Color: "FFFFFF", Style: 1},
			{Type: "left", Color: "FFFFFF", Style: 2},
			{Type: "right", Color: "FFFFFF", Style: 2}},
		CustomNumFmt: &pairingsfmt,
	})

	cellStyleId3, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CCC0DA"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{{Type: "bottom", Color: "FFFFFF", Style: 1},
			{Type: "left", Color: "FFFFFF", Style: 2},
			{Type: "right", Color: "FFFFFF", Style: 2}},
		CustomNumFmt: &dateTimefmt,
	})

	cellStyleId4, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"E5E0EC"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{{Type: "bottom", Color: "FFFFFF", Style: 1},
			{Type: "left", Color: "FFFFFF", Style: 2},
			{Type: "right", Color: "FFFFFF", Style: 2}},
		CustomNumFmt: &dateTimefmt,
	})

	for i := 1; i < len(al.PairsArray); i++ {
		departure := al.PairsArray[i].Start
		arrival := al.PairsArray[i].End
		departureString := fmt.Sprintf("%d/%d/%d %d:%02d",
			departure.Day(), departure.Month(), departure.Year(), departure.Hour(), departure.Minute())
		arrivalString := fmt.Sprintf("%d/%d/%d %d:%02d",
			arrival.Day(), arrival.Month(), arrival.Year(), arrival.Hour(), arrival.Minute())

		f.SetRowHeight(sheetName, 4+al.PairsArray[i].Id, 20.0)

		var pairStyle int
		var dateStyle int
		if al.PairsArray[i].Id%2 != 0 {
			pairStyle = cellStyleId1
			dateStyle = cellStyleId3
		} else {
			pairStyle = cellStyleId2
			dateStyle = cellStyleId4
		}

		cell, _ := excelize.CoordinatesToCellName(7, 4+al.PairsArray[i].Id)
		f.SetCellValue(sheetName, cell, al.PairsArray[i].Id)
		f.SetCellStyle(sheetName, cell, cell, pairStyle)

		cell, _ = excelize.CoordinatesToCellName(8, 4+al.PairsArray[i].Id)
		f.SetCellValue(sheetName, cell, departureString)
		f.SetCellStyle(sheetName, cell, cell, dateStyle)

		cell, _ = excelize.CoordinatesToCellName(9, 4+al.PairsArray[i].Id)
		f.SetCellValue(sheetName, cell, arrivalString)
		f.SetCellStyle(sheetName, cell, cell, dateStyle)
	}

	start, _ := excelize.CoordinatesToCellName(7, 4)
	end, _ := excelize.CoordinatesToCellName(9, 4+rows)
	enable := true
	f.AddTable(sheetName, &excelize.Table{
		Range:             start + ":" + end,
		Name:              "Pairings",
		StyleName:         "TableStyleMedium12",
		ShowFirstColumn:   false,
		ShowLastColumn:    false,
		ShowRowStripes:    &enable,
		ShowColumnStripes: false,
	})

}
