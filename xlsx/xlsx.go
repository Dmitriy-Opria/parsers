package main

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
)

var (
	xlsxPath = "./TestResults.xlsx"
	errIndex = errors.New("invalid sheet index")
)

func main() {
	if err := readXlsx(xlsxPath); err != nil {
		fmt.Println(err.Error())
	}
}

func readXlsx(path string) (err error) {

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Printf("Get XLSX parsing err: %s", err.Error())
	}
	if sheetValue, err := getSheetData(*xlFile, 0); err == nil {
		fmt.Println(sheetValue)
	}
	if sheetValue, err := getSheetData(*xlFile, 1); err == nil {
		fmt.Println(sheetValue)
	}

	return
}

func getSheetData(file xlsx.File, sheetIndex int) (value [][]string, err error) {

	if sheetIndex < len(file.Sheets)-1 {

		sheet := file.Sheets[sheetIndex]
		value = make([][]string, 0, len(sheet.Rows))

		for _, row := range sheet.Rows {
			cells := make([]string, 0, len(row.Cells))

			for _, cell := range row.Cells {
				cells = append(cells, cell.Value)
			}
			value = append(value, cells)
		}
	} else {
		err = errIndex
	}
	return
}
