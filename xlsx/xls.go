package main

import (
	"fmt"
	"github.com/extrame/xls"
	"github.com/iizotop/baseweb/utils"
)

type RecordField struct {
	LabNumber     string
	Name          string
	Code          string
	Barcode       string
	Latitude      float64
	Longitude     float64
	Boron         float64
	Calcium       float64
	Chloride      float64
	Copper        float64
	Iron          float64
	Magnesium     float64
	Manganese     float64
	Nitrate       float64
	Phosphorus    float64
	PlantWeight   float64
	Potassium     float64
	Sodium        float64
	Sulfur        float64
	TotalNitrogen float64
	Zinc          float64
}

var (
	charset = "utf-8"
	//testFile = "CSBP_171215_Bristol_002.xls"
	testFile = "CSBP_180102_Airlie_001.xls"
	//testFile = "CSBP_180102_Bella_001.xls"
)

func main() {
	records := readXLS(testFile, charset)

	for _, record := range records {
		fmt.Printf("%#v\n", record)
	}
}

func readXLS(fileName, charset string) (recordList []RecordField) {

	workBook, err := xls.Open(fileName, charset)

	if err != nil {
		fmt.Println("Can`t open file", err.Error())
	}

	sheet := workBook.GetSheet(0)

	rowCount := int(sheet.MaxRow)

	recordList = make([]RecordField, 0, rowCount-4)

	for rowIndex := 5; rowIndex <= rowCount; rowIndex++ {

		currentRow := sheet.Row(rowIndex)

		colMaxCount := currentRow.LastCol()

		record := RecordField{}

		if colMaxCount >= 23 {

			record.LabNumber = currentRow.Col(0)
			record.Name = currentRow.Col(1)
			record.Code = currentRow.Col(3)
			record.Barcode = currentRow.Col(4)
			record.Latitude = utils.ToFloat64(currentRow.Col(6))
			record.Longitude = utils.ToFloat64(currentRow.Col(7))
			record.Boron = utils.ToFloat64(currentRow.Col(8))
			record.Calcium = utils.ToFloat64(currentRow.Col(9))
			record.Chloride = utils.ToFloat64(currentRow.Col(10))
			record.Copper = utils.ToFloat64(currentRow.Col(11))
			record.Iron = utils.ToFloat64(currentRow.Col(12))
			record.Magnesium = utils.ToFloat64(currentRow.Col(13))
			record.Manganese = utils.ToFloat64(currentRow.Col(14))
			record.Nitrate = utils.ToFloat64(currentRow.Col(15))
			record.Phosphorus = utils.ToFloat64(currentRow.Col(16))
			record.PlantWeight = utils.ToFloat64(currentRow.Col(17))
			record.Potassium = utils.ToFloat64(currentRow.Col(18))
			record.Sodium = utils.ToFloat64(currentRow.Col(19))
			record.Sulfur = utils.ToFloat64(currentRow.Col(20))
			record.TotalNitrogen = utils.ToFloat64(currentRow.Col(21))
			record.Zinc = utils.ToFloat64(currentRow.Col(22))

		}

		if record.LabNumber != "" {
			recordList = append(recordList, record)
		}
	}

	return
}
