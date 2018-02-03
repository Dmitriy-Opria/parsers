package main

import (
	"fmt"
	"github.com/ledongthuc/pdf"
	"strings"
)

var (
	nutrientKeys = []string{
		"NUTRIENT",
		"Result",
		"Low",
		"Marginal",
		"Sufficient",
		"High",
		"Exces",
		"Sufficiency Range",
	}

	nutrientFields = []string{
		"Nitrogen Total",
		"Phosphorus",
		"Potassium",
		"Sulfur",
		"Calcium",
		"Magnesium",
		"Sodium",
		"Paddock",
		"Chloride",
		"Copper",
		"Zinc",
		"Manganese",
		"Iron",
		"Boron",
		"Nitrogen/Sulphur",
	}

	analysisKeys = "ANALYSIS RESULTS"

	analysisFields = []string{
		"Paddock Name",
		"Sample Depth (cm)",
		"Nitrogen Total (Dumas) %",
		"Nitrate Nitrogen mg/kg",
		"Phosphorus %",
		"Potassium %",
		"Sulfur %",
		"Calcium %",
		"Magnesium %",
		"Sodium %",
		"Chloride %",
		"Copper mg/kg",
		"Zinc mg/kg",
		"Manganese mg/kg",
		"Iron mg/kg",
		"Boron mg/kg",
		"Nitrogen/Phosphorus Ratio",
		"Nitrogen/Potassium Ratio",
		"Nitrogen/Sulphur Ratio",
	}

	filesList = []string{
		"./SM128800 - Auscott Narrabri F6-1 S1.pdf",
		"./SM128801 - Auscott Narrabri F6-1 S2.pdf",
		"./SM128802 - Auscott Narrabri F6-1 S3.pdf",
		"./SM128803 - Auscott Narrabri F6-1 S4.pdf",
		"./SM128804 - Auscott Narrabri F23 S1.pdf",
		"./SM128805 - Auscott Narrabri F23 S2.pdf",
		"./SM128806 - Auscott Narrabri F23 S3.pdf",
		"./SM128807 - Auscott Narrabri F23 S4.pdf",
		"./SM128808 - Auscott Narrabri F48 S1.pdf",
		"./SM128809 - Auscott Narrabri F48 S2.pdf",
		"./SM128810 - Auscott Narrabri F48 S3.pdf",
	}

	date = "DATE"

	dateOfSampling = "Date of sampling"
)

type (
	Record struct {
		Page       int
		StartIndex int
		EndIndex   int
		Content    string
	}
	TableHeader struct {
		Content string
		Start   float64
		End     float64
		Len     float64
	}
	NutrientContent struct {
		Name        string
		Result      string
		Sufficiency string
	}
	AnalysisContent struct {
		Name  string
		Value string
	}

	DocumentContent struct {
		Date            string
		DateOfSampling  string
		AnalysisContent []AnalysisContent
		NutrientContent []NutrientContent
	}
)

func main() {

	for _, v := range filesList {
		content, err := readPdf(v) // Read local pdf file
		if err != nil {
			panic(err)
		}
		fmt.Printf("Table fields\nDate:               %s\nDate of sampling:   %s\n", content.Date, content.DateOfSampling)

		for _, v := range content.NutrientContent {
			fmt.Printf("Name: %15s  |  Result: %10s  | Sufficiency Range: %10s\n", v.Name, v.Result, v.Sufficiency)
		}
		for _, v := range content.AnalysisContent {
			fmt.Printf("Name: %30s  |  Value: %10s\n", v.Name, v.Value)
		}

	}

}

func readPdf(path string) (content DocumentContent, err error) {
	r, err := pdf.Open(path)
	if err != nil {
		return
	}

	totalPage := r.NumPage()

	var text string
	var stringsList = make([]Record, 0, 256)
	var line float64
	var width float64
	var startIndex int

	//GET ALL STRINGS OF PAGES

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		page := r.Page(pageIndex)

		for i, tx := range page.Content().Text {

			if tx.Y != line {
				record := Record{}
				record.Page = pageIndex
				record.EndIndex = i - 1
				record.Content = text
				record.StartIndex = startIndex

				stringsList = append(stringsList, record)

				startIndex = i
				text = tx.S
				line = tx.Y
				width = 0
			} else {
				if tx.X > (width*1.2 + tx.W) {
					if width != 0 {
						text += " "
					}
				}
				width = tx.X
				text += tx.S
			}
		}

	}

	//GET NUTRIENT TABLE FIELDS AND BORDERS

	nutrientHeader := make([]TableHeader, 0, len(nutrientKeys))

	if tableIndex := getStartNutrientIndex(stringsList); tableIndex != -1 {

		tableHead := stringsList[tableIndex]

		page := r.Page(tableHead.Page)
		startindex := tableHead.StartIndex
		endindex := tableHead.EndIndex

		for _, tableField := range nutrientKeys {
			content := page.Content().Text[startindex:endindex]

			headerField := TableHeader{}

			for conInx, v := range content {

				if i := strings.Index(tableField, v.S); i != -1 {

					if headerField.Content == "" {
						headerField.Start = v.X
					}
					headerField.End = v.X
					headerField.Content += v.S
					headerField.Len = headerField.End - headerField.Start

				} else {
					headerField.Content = ""
				}
				if len(headerField.Content) == len(tableField) {
					startindex += conInx
					break
				}
			}
			nutrientHeader = append(nutrientHeader, headerField)
		}

		setMiddleValue(nutrientHeader)

		//GET NUTRIENT CONTENT FROM LOOKING TABLE FIELDS

		content.NutrientContent = getNutrientData(stringsList, r, nutrientHeader, tableIndex)
	}

	if analysisTableIndex := getStartAnalisysIndex(stringsList); analysisTableIndex != -1 {

		header := []TableHeader{
			{
				Content: "ANALYSIS RESULTS",
				Start:   15,
				End:     128,
			},
			{
				Content: "F",
				Start:   128,
				End:     800,
			},
		}

		content.AnalysisContent = getAnalysisData(stringsList, r, header, analysisTableIndex)
	}
	content.Date = getDate(stringsList)
	content.DateOfSampling = getDateOfSampling(stringsList)

	return
}

func getStartNutrientIndex(stringsList []Record) int {

	for i, str := range stringsList {

		currentCounter := 0
		for _, v := range nutrientKeys {
			if strings.Contains(str.Content, v) {
				currentCounter++
			}
		}

		if currentCounter == len(nutrientKeys) {
			return i
		}
	}
	return -1
}

func getStartAnalisysIndex(stringsList []Record) int {

	for i, str := range stringsList {

		if strings.Contains(str.Content, analysisKeys) {
			return i
		}
	}
	return -1
}

func setMiddleValue(header []TableHeader) {

	for index, head := range header {
		if index == 0 {
			continue
		}
		if index == len(header)-1 {
			header[index].Start = header[index-1].End
			header[index].End = header[index-1].End + 80
			header[index].Len = 80
			continue
		}
		header[index-1].End = (head.Start + header[index-1].End) / 2
		header[index].Start = (head.Start + header[index-1].End) / 2
		header[index-1].Len = header[index-1].End - header[index-1].Start
	}
}

func getNutrientData(stringsList []Record, r *pdf.Reader, header []TableHeader, tableIndex int) (tableConten []NutrientContent) {

	tableContent := getContentFields(stringsList, tableIndex)

	for _, content := range tableContent {
		page := r.Page(content.Page)
		cont := page.Content().Text[content.StartIndex:content.EndIndex]

		var nameVal, ratioVal, sufRange string
		for _, v := range cont {

			if len(header) < 8 {
				return
			}

			name := header[0]
			if v.X >= name.Start && v.X < name.End {
				nameVal += v.S
			}

			ratio := header[1]
			if v.X >= ratio.Start && v.X < ratio.End {
				if areValue(v.S) {
					ratioVal += v.S
				}
			}

			suf := header[7]
			if v.X >= suf.Start && v.X < suf.End {
				if areValue(v.S) {
					sufRange += v.S
				}
			}

		}
		tableConten = append(tableConten, NutrientContent{Name: nameVal, Result: ratioVal, Sufficiency: sufRange})
	}

	return
}

func getAnalysisData(stringsList []Record, r *pdf.Reader, header []TableHeader, tableIndex int) (tableConten []AnalysisContent) {

	tableContent := getAnalysContentFields(stringsList, tableIndex)

	for i, content := range tableContent {
		page := r.Page(content.Page)
		if content.StartIndex < 0 {
			content.StartIndex = 0
		}
		if content.EndIndex < 0 {
			content.EndIndex = len(page.Content().Text)
		}
		cont := page.Content().Text[content.StartIndex : content.EndIndex+10]

		var nameVal, value string
		for _, v := range cont {
			if len(header) < 2 {
				return
			}

			nameVal = analysisFields[i]

			ratio := header[1]
			if v.X >= ratio.Start && v.X < ratio.End {
				if i == 0 {
					value += v.S
				} else if areValue(v.S) {
					value += v.S
				}
			}

		}
		value = strings.Trim(value, " - ")
		tableConten = append(tableConten, AnalysisContent{Name: nameVal, Value: value})
	}

	return
}

func getContentFields(stringsList []Record, tableIndex int) (tableContent []Record) {

	tableContent = make([]Record, 0, 15)

	elseStrings := stringsList[tableIndex:]

	for _, field := range nutrientFields {
		for _, line := range elseStrings {
			if strings.Contains(line.Content, field) {
				tableContent = append(tableContent, line)
			}
		}
	}

	return
}

func getAnalysContentFields(stringsList []Record, tableIndex int) (tableContent []Record) {

	tableContent = make([]Record, 0, 15)

	elseStrings := stringsList[tableIndex:]

	for _, field := range analysisFields {
		for _, line := range elseStrings {
			if strings.Contains(line.Content, field) {
				tableContent = append(tableContent, line)
				break
			}
		}
	}

	return
}

func getDate(stringsList []Record) (dateRes string) {
	for _, str := range stringsList {
		if strings.Contains(str.Content, date) {
			dateRes = strings.TrimPrefix(str.Content, "DATE:")
			return
		}
	}
	return
}

func getDateOfSampling(stringsList []Record) (dateOfSamplingRes string) {
	for _, str := range stringsList {
		if strings.Contains(str.Content, dateOfSampling) {
			dateOfSamplingRes = strings.TrimPrefix(str.Content, "Date of sampling ")
			return
		}
	}
	return
}

func areValue(v string) bool {
	switch v {
	case ".":
		fallthrough
	case ",":
		fallthrough
	case "-":
		fallthrough
	case " ":
		fallthrough
	case "0":
		fallthrough
	case "1":
		fallthrough
	case "2":
		fallthrough
	case "3":
		fallthrough
	case "4":
		fallthrough
	case "5":
		fallthrough
	case "6":
		fallthrough
	case "7":
		fallthrough
	case "8":
		fallthrough
	case "9":
		fallthrough
	case "10":
		return true
	}
	return false
}
