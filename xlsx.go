package main

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"strings"
	"io"
	//"encoding/json"
	//"os"
)

var (
	xlsxPath = "TestResults.xlsx"
	errIndex = errors.New("invalid sheet index")
)

func init() {

	return

	records, err := ReadXLSX(xlsxPath)
	if err != nil {
		fmt.Println(err.Error())
	}

	_ = records
}

func main(){

	records, err := ReadXLSX(xlsxPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(records)
}
//

type (
	Header map[string]int
	Record map[string]string
)

var headerFields = map[string]string{
	"farm":          "farm",
	"field number":  "field",
	"sample number": "number",
	"sample date":   "date",
	"lat":           "lat",
	"long":          "lon",
	"total n":       "total_n",
	"total p":       "total_p",
	"total k":       "total_k",
	"nitrate-ppm":   "nitrate_ppm",
}

func (h Header) Index(name string) int {
	return h[name]
}

func parseHeader(line []string) (header Header) {
	header = make(Header, len(line))
	for index, name := range line {
		name = strings.ToLower(name)
		if value, ok := headerFields[name]; ok {
			header[value] = index
		}
	}
	return
}

func parseRecord(line []string, header Header) (record Record) {
	record = make(Record, len(header))
	for name, index := range header {
		if index < len(line) {
			record[name] = line[index]
		}
	}
	return
}

func Read(r io.ReaderAt, size int64) (records []Record, err error) {

	xlsxFile, err := xlsx.OpenReaderAt(r, size)
	if err != nil {
		return nil, err
	}

	sheetValue, err := getSheetData(xlsxFile, 0)
	if err != nil {
		return nil, err
	}

	if len(sheetValue) > 1 {
		header := parseHeader(sheetValue[0])
		for _, line := range sheetValue[1:] {
			record := parseRecord(line, header)
			records = append(records, record)
			//json.NewEncoder(os.Stdout).Encode(record)
		}
	}
	return
}

func ReadXLSX(path string) (records []Record, err error) {

	xlsxFile, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}

	//for index, sheet := range xlsxFile.Sheets {
	//	fmt.Println("sheet:", index, sheet.Name)
	//}

	sheetValue, err := getSheetData(xlsxFile, 0)
	if err != nil {
		return nil, err
	}

	if len(sheetValue) > 1 {
		header := parseHeader(sheetValue[0])
		for _, line := range sheetValue[1:] {
			record := parseRecord(line, header)
			records = append(records, record)
			//json.NewEncoder(os.Stdout).Encode(record)
		}
	}
	return
}

func getSheetData(file *xlsx.File, sheetIndex int) (value [][]string, err error) {

	if sheetIndex < len(file.Sheets)-1 {

		sheet := file.Sheets[sheetIndex]
		value = make([][]string, 0, len(sheet.Rows))

		for _, row := range sheet.Rows {

			cells := make([]string, 0, len(row.Cells))
			emptyCell := 0

			for _, cell := range row.Cells {

				value := cell.String()

				cells = append(cells, value)
				if value == "" || value == "#VALUE!" {
					emptyCell++
				}
			}

			if emptyCell != len(cells) {
				value = append(value, cells)
			}
		}

	} else {

		err = errIndex
	}
	return
}
