package results

import (
	"time"

	"github.com/xuri/excelize/v2"
)

// struct used for the visualization of the schedule
type dateInfo struct {
	year       int
	month      int
	dateString string
	startCell  string
}

func daysInMonth(year int, month string) int {
	// returns the number of days a month has
	daysInMonth := make(map[string]int)
	daysInMonth["January"] = 31
	daysInMonth["February"] = 28
	daysInMonth["March"] = 31
	daysInMonth["April"] = 30
	daysInMonth["May"] = 31
	daysInMonth["June"] = 30
	daysInMonth["July"] = 31
	daysInMonth["August"] = 31
	daysInMonth["September"] = 30
	daysInMonth["October"] = 31
	daysInMonth["November"] = 30
	daysInMonth["December"] = 31

	days := daysInMonth[month]
	if month == "February" && checkLeapYear(year) {
		days = 29
	}
	return days
}

func checkLeapYear(year int) bool {
	// checks if a year is a leap year
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}

func drawScheduleHeader(f *excelize.File, sheetName string, startCell string) {
	// function that visualizes the headers of an airline crew rostering schedule
	col, row, _ := excelize.CellNameToCoordinates(startCell)
	endCell, _ := excelize.CoordinatesToCellName(col+1, row+1)
	f.MergeCell(sheetName, startCell, endCell)
	f.SetCellValue(sheetName, startCell, "         Dates\n  Pilots")

	style := &excelize.Style{Font: &excelize.Font{Size: 16, Color: "FFFFFF", Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"376091"}},
		Border: []excelize.Border{{Type: "top", Color: "FFFFFF", Style: 5},
			{Type: "right", Color: "FFFFFF", Style: 5},
			{Type: "bottom", Color: "FFFFFF", Style: 5},
			{Type: "left", Color: "FFFFFF", Style: 5},
			{Type: "diagonalDown", Color: "FFFFFF", Style: 1},
		}}
	styleId, _ := f.NewStyle(style)
	f.SetCellStyle(sheetName, startCell, endCell, styleId)
	f.SetRowHeight(sheetName, row, 30.0)
	f.SetRowHeight(sheetName, row+1, 30.0)
}

func drawPilotColumn(f *excelize.File, sheetName string, startCell string, pilots int) {
	// function that visualizes the column containing the pilots' ids of an airline crew rostering schedule
	col, row, _ := excelize.CellNameToCoordinates(startCell)
	endCell, _ := excelize.CoordinatesToCellName(col+1, row+pilots-1)
	style := &excelize.Style{Font: &excelize.Font{Size: 12, Color: "FFFFFF", Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4F81BD"}},
		Border: []excelize.Border{{Type: "right", Color: "FFFFFF", Style: 5},
			{Type: "bottom", Color: "FFFFFF", Style: 2},
			{Type: "left", Color: "FFFFFF", Style: 5},
		}}
	styleId, _ := f.NewStyle(style)
	f.SetCellStyle(sheetName, startCell, endCell, styleId)
	for pilot := 0; pilot < pilots; pilot++ {
		cell, _ := excelize.CoordinatesToCellName(col, row+pilot)
		endCell, _ := excelize.CoordinatesToCellName(col+1, row+pilot)
		f.SetRowHeight(sheetName, row+pilot, 15.0)
		f.MergeCell(sheetName, cell, endCell)
		f.SetCellValue(sheetName, cell, pilot)
	}
}

func drawMonths(f *excelize.File, sheetName string, startCell string, months []*dateInfo) {
	// function that visualizes the months and days of an airline crew rostering schedule
	col, row, _ := excelize.CellNameToCoordinates(startCell)
	monthStyleId, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 16, Color: "FFFFFF", Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4F81BD"}},
		Border: []excelize.Border{{Type: "top", Color: "FFFFFF", Style: 5},
			{Type: "right", Color: "FFFFFF", Style: 5},
			{Type: "bottom", Color: "FFFFFF", Style: 5},
			{Type: "left", Color: "FFFFFF", Style: 5},
		}})
	dayStyleId, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 14, Color: "FFFFFF", Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4F81BD"}},
		Border: []excelize.Border{{Type: "right", Color: "FFFFFF", Style: 2},
			{Type: "bottom", Color: "FFFFFF", Style: 5},
			{Type: "left", Color: "FFFFFF", Style: 2},
		}})
	lastDayStyleId, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 14, Color: "FFFFFF", Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"4F81BD"}},
		Border: []excelize.Border{{Type: "right", Color: "FFFFFF", Style: 5},
			{Type: "bottom", Color: "FFFFFF", Style: 5},
			{Type: "left", Color: "FFFFFF", Style: 2},
		}})
	for _, month := range months {
		days := daysInMonth(month.year, time.Month(month.month).String())
		monthStartCell, _ := excelize.CoordinatesToCellName(col, row)
		monthEndCell, _ := excelize.CoordinatesToCellName(col+2*days-1, row)
		month.startCell = monthStartCell
		f.MergeCell(sheetName, monthStartCell, monthEndCell)
		f.SetCellStyle(sheetName, monthStartCell, monthEndCell, monthStyleId)
		f.SetCellValue(sheetName, monthStartCell, month.dateString)
		colStartStr, _ := excelize.ColumnNumberToName(col)
		colEndStr, _ := excelize.ColumnNumberToName(col + 2*days - 1)
		f.SetColWidth(sheetName, colStartStr, colEndStr, 2.15)

		daysStartCell, _ := excelize.CoordinatesToCellName(col, row+1)
		daysEndCell, _ := excelize.CoordinatesToCellName(col+2*days-3, row+1)
		f.SetCellStyle(sheetName, daysStartCell, daysEndCell, dayStyleId)

		lastDayCell1, _ := excelize.CoordinatesToCellName(col+2*days-2, row+1)
		lastDayCell2, _ := excelize.CoordinatesToCellName(col+2*days-1, row+1)
		f.SetCellStyle(sheetName, lastDayCell1, lastDayCell2, lastDayStyleId)

		for i, d := 0, 1; i < 2*days; i, d = i+2, d+1 {
			startCell, _ := excelize.CoordinatesToCellName(col+i, row+1)
			endCell, _ := excelize.CoordinatesToCellName(col+i+1, row+1)
			f.MergeCell(sheetName, startCell, endCell)
			f.SetCellValue(sheetName, startCell, d)
		}
		col += 2 * days
	}
}

func drawDataArea(f *excelize.File, sheetName string, months []*dateInfo, pilots int) {
	// function that visualizes the area containing the pair assignments of an airline crew rostering schedule
	halfDayEvenStyleId, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 11, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"B8CCE4"}},
		Border:    []excelize.Border{{Type: "right", Color: "FFFFFF", Style: 1}, {Type: "bottom", Color: "FFFFFF", Style: 2}}})
	halfDayOddStyleId, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 11, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"DEEBF6"}},
		Border:    []excelize.Border{{Type: "right", Color: "FFFFFF", Style: 1}, {Type: "bottom", Color: "FFFFFF", Style: 2}}})
	dayEvenStyleId, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 11, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"B8CCE4"}},
		Border:    []excelize.Border{{Type: "right", Color: "FFFFFF", Style: 2}, {Type: "bottom", Color: "FFFFFF", Style: 2}}})
	dayOddStyleId, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 11, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"DEEBF6"}},
		Border:    []excelize.Border{{Type: "right", Color: "FFFFFF", Style: 2}, {Type: "bottom", Color: "FFFFFF", Style: 2}}})
	lastDayEvenStyleId, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 11, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"B8CCE4"}},
		Border:    []excelize.Border{{Type: "right", Color: "FFFFFF", Style: 5}, {Type: "bottom", Color: "FFFFFF", Style: 2}}})
	lastDayOddStyleId, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 11, Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"DEEBF6"}},
		Border:    []excelize.Border{{Type: "right", Color: "FFFFFF", Style: 5}, {Type: "bottom", Color: "FFFFFF", Style: 2}}})

	for _, month := range months {
		col, row, _ := excelize.CellNameToCoordinates(month.startCell)
		row += 2
		days := daysInMonth(month.year, time.Month(month.month).String())
		for pilot := 0; pilot < pilots; pilot++ {
			for d := 0; d < 2*days; d += 2 {
				halfCell, _ := excelize.CoordinatesToCellName(col+d, row+pilot)
				cell, _ := excelize.CoordinatesToCellName(col+d+1, row+pilot)
				if pilot%2 == 0 {
					f.SetCellStyle(sheetName, halfCell, halfCell, halfDayEvenStyleId)
					if d == 2*days-2 {
						f.SetCellStyle(sheetName, cell, cell, lastDayEvenStyleId)
					} else {
						f.SetCellStyle(sheetName, cell, cell, dayEvenStyleId)
					}
				} else {
					f.SetCellStyle(sheetName, halfCell, halfCell, halfDayOddStyleId)
					if d == 2*days-2 {
						f.SetCellStyle(sheetName, cell, cell, lastDayOddStyleId)
					} else {
						f.SetCellStyle(sheetName, cell, cell, dayOddStyleId)
					}
				}
			}
		}
	}
}

func findCell(startCell string, date time.Time, pilotId int) string {
	// find a cell in a visualized airline crew rostering schedule based on the date and the pilot
	col, row, _ := excelize.CellNameToCoordinates(startCell)
	row += 2 + pilotId
	col += 2*date.Day() - 2
	if date.Hour() > 12 {
		col++
	}
	cell, _ := excelize.CoordinatesToCellName(col, row)
	return cell
}

func drawDashes(f *excelize.File, sheetName string, startCell string, endCell string) {
	// draw dashes in cells to show that a pilot is busy during the time represented by
	// these cells
	col, row, _ := excelize.CellNameToCoordinates(startCell)
	endCol, _, _ := excelize.CellNameToCoordinates(endCell)
	for j := col + 1; j <= endCol; j++ {
		cell, _ := excelize.CoordinatesToCellName(j, row)
		value, _ := f.GetCellValue(sheetName, cell)
		if value != "F" {
			f.SetCellValue(sheetName, cell, "-")
		}
	}
}
