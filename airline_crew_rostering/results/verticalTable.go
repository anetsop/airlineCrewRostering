package results

import "github.com/xuri/excelize/v2"

func drawVerticalTable(f *excelize.File, sheetName string, startCell string, numberOfRows int, headerFill string, dataFill string) {
	tableStyleTitle := excelize.Style{
		Font:      &excelize.Font{Size: 16, Color: "FFFFFF", Bold: true, Underline: "single"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{headerFill}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    []excelize.Border{{Type: "bottom", Color: "FFFFFF", Style: 5}},
	}

	x, y, _ := excelize.CellNameToCoordinates(startCell)
	end, _ := excelize.CoordinatesToCellName(x+3, y+2)
	f.MergeCell(sheetName, startCell, end)
	styleId, _ := f.NewStyle(&tableStyleTitle)
	f.SetCellStyle(sheetName, startCell, end, styleId)

	y += 3
	tableStyleCellsHeader := excelize.Style{
		Font:      &excelize.Font{Size: 12, Color: "FFFFFF", Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{headerFill}},
		Border:    []excelize.Border{{Type: "top", Color: "FFFFFF", Style: 5}, {Type: "right", Color: "FFFFFF", Style: 1}, {Type: "bottom", Color: "FFFFFF", Style: 1}},
	}
	tableStyleCellsData := excelize.Style{
		Font:      &excelize.Font{Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{dataFill}},
		Border:    []excelize.Border{{Type: "top", Color: "FFFFFF", Style: 5}, {Type: "left", Color: "FFFFFF", Style: 1}, {Type: "bottom", Color: "FFFFFF", Style: 1}},
	}

	start, _ := excelize.CoordinatesToCellName(x, y)
	end, _ = excelize.CoordinatesToCellName(x+1, y)
	f.MergeCell(sheetName, start, end)
	start2, _ := excelize.CoordinatesToCellName(x+2, y)
	end2, _ := excelize.CoordinatesToCellName(x+3, y)
	f.MergeCell(sheetName, start2, end2)
	styleId, _ = f.NewStyle(&tableStyleCellsHeader)
	f.SetCellStyle(sheetName, start, end, styleId)
	styleId, _ = f.NewStyle(&tableStyleCellsData)
	f.SetCellStyle(sheetName, start2, end2, styleId)

	top := &tableStyleCellsHeader.Border[0].Style
	*top = 1
	top = &tableStyleCellsData.Border[0].Style
	*top = 1

	for i := 0; i < numberOfRows; i++ {
		start, _ = excelize.CoordinatesToCellName(x, y+i)
		end, _ = excelize.CoordinatesToCellName(x+1, y+i)
		f.MergeCell(sheetName, start, end)

		start2, _ = excelize.CoordinatesToCellName(x+2, y+i)
		end2, _ = excelize.CoordinatesToCellName(x+3, y+i)
		f.MergeCell(sheetName, start2, end2)
	}
	start, _ = excelize.CoordinatesToCellName(x, y)
	end, _ = excelize.CoordinatesToCellName(x+1, y+numberOfRows-2)
	start2, _ = excelize.CoordinatesToCellName(x+2, y)
	end2, _ = excelize.CoordinatesToCellName(x+3, y+numberOfRows-2)

	styleId, _ = f.NewStyle(&tableStyleCellsHeader)
	f.SetCellStyle(sheetName, start, end, styleId)
	styleId, _ = f.NewStyle(&tableStyleCellsData)
	f.SetCellStyle(sheetName, start2, end2, styleId)

	bottom := &tableStyleCellsHeader.Border[2].Style
	*bottom = 0
	bottom = &tableStyleCellsData.Border[2].Style
	*bottom = 0
	y += numberOfRows - 1
	start, _ = excelize.CoordinatesToCellName(x, y)
	end, _ = excelize.CoordinatesToCellName(x+1, y)
	start2, _ = excelize.CoordinatesToCellName(x+2, y)
	end2, _ = excelize.CoordinatesToCellName(x+3, y)

	styleId, _ = f.NewStyle(&tableStyleCellsHeader)
	f.SetCellStyle(sheetName, start, end, styleId)
	styleId, _ = f.NewStyle(&tableStyleCellsData)
	f.SetCellStyle(sheetName, start2, end2, styleId)
}
