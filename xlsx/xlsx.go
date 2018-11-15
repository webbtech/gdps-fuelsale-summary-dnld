package xlsx

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/pulpfree/gdps-fs-sum-dwnld/model"

	log "github.com/sirupsen/logrus"
)

// XLSX struct
type XLSX struct {
	file *excelize.File
}

// Defaults
const (
	abc             = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	floatFrmt       = "#,#0"
	timeShortForm   = "20060102"
	timeMonthForm   = "200601"
	dateDayFormat   = "Jan 2, 2006"
	dateMonthFormat = "January 2006"
)

// NewFile function
func NewFile() (x *XLSX, err error) {

	x = new(XLSX)
	x.file = excelize.NewFile()
	if err != nil {
		log.Errorf("xlsx err %s: ", err)
	}
	return x, err
}

// NLSales method
func (x *XLSX) NLSales(fs *model.Report) (err error) {

	var cell string
	var style int

	xlsx := x.file
	sheetNm := "Sheet1"
	statNmColWidth := 12.0

	sales := fs.FuelSales
	xlsx.SetSheetName(sheetNm, "NoLead Summary")

	title := fmt.Sprintf("NoLead Sales period: %s - %s", sales.DateFrom.Format(dateDayFormat), sales.DateTo.Format(dateDayFormat))
	style, _ = xlsx.NewStyle(`{"font":{"bold":true,"size":12}}`)
	xlsx.SetCellValue(sheetNm, "A1", title)
	xlsx.SetCellStyle(sheetNm, "A1", "A1", style)

	// Fill in data
	col := 1
	row := 2
	style, _ = xlsx.NewStyle(`{"number_format": 3}`)

	xlsx.SetColWidth(sheetNm, "A", "B", statNmColWidth)

	for _, s := range *sales.Sales {
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, s.StationName)

		col++
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, s.Fuels.NL)
		xlsx.SetCellStyle(sheetNm, cell, cell, style)

		col = 1
		row++
	}

	// Set summary
	col = 2
	startRow := strconv.Itoa(2)
	endRow := strconv.Itoa(row - 1)
	cell = toChar(col) + strconv.Itoa(row)
	startCell := toChar(col) + startRow
	endCell := toChar(col) + endRow
	rangeStr := fmt.Sprintf("SUM(%s:%s)", startCell, endCell)
	style, _ = xlsx.NewStyle(`{"font":{"bold": true}}`)
	xlsx.SetCellFormula(sheetNm, cell, rangeStr)
	xlsx.SetCellStyle(sheetNm, cell, cell, style)

	return err
}

// DSLSales method
func (x *XLSX) DSLSales(fs *model.Report) (err error) {
	var cell string
	var style int

	xlsx := x.file
	sheetNm := "Sheet2"
	statNmColWidth := 12.0

	_ = xlsx.NewSheet(sheetNm)
	xlsx.SetSheetName(sheetNm, "Diesel Summary")

	sales := fs.FuelSales
	title := fmt.Sprintf("Diesel Sales period: %s - %s", sales.DateFrom.Format(dateDayFormat), sales.DateTo.Format(dateDayFormat))
	style, _ = xlsx.NewStyle(`{"font":{"bold":true,"size":12}}`)
	xlsx.SetCellValue(sheetNm, "A1", title)
	xlsx.SetCellStyle(sheetNm, "A1", "A1", style)

	// Fill in data
	col := 1
	row := 2
	style, _ = xlsx.NewStyle(`{"number_format": 3}`)

	xlsx.SetColWidth(sheetNm, "A", "B", statNmColWidth)

	// Create new slice for diesel results
	DSLSales := []model.SalesSummary{}
	for _, s := range *sales.Sales {
		if s.HasDSL == true {
			DSLSales = append(DSLSales, s)
		}
	}

	for _, s := range DSLSales {
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, s.StationName)

		col++
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, s.Fuels.DSL)
		xlsx.SetCellStyle(sheetNm, cell, cell, style)

		col = 1
		row++
	}

	// Set summary
	col = 2
	startRow := strconv.Itoa(2)
	endRow := strconv.Itoa(row - 1)
	cell = toChar(col) + strconv.Itoa(row)
	startCell := toChar(col) + startRow
	endCell := toChar(col) + endRow
	rangeStr := fmt.Sprintf("SUM(%s:%s)", startCell, endCell)
	style, _ = xlsx.NewStyle(`{"font":{"bold": true}}`)
	xlsx.SetCellFormula(sheetNm, cell, rangeStr)
	xlsx.SetCellStyle(sheetNm, cell, cell, style)

	return err
}

// OutputFile method
func (x *XLSX) OutputFile() (buf bytes.Buffer, err error) {
	err = x.file.Write(&buf)
	if err != nil {
		log.Errorf("xlsx err: %s", err)
	}
	return buf, err
}

// OutputToDisk method
func (x *XLSX) OutputToDisk(path string) (fp string, err error) {
	err = x.file.SaveAs(path)
	return path, err
}

//
// ======================== Helper Methods ================================= //
//

// see: https://stackoverflow.com/questions/36803999/golang-alphabetic-representation-of-a-number
// for a way to map int to letters
func toChar(i int) string {
	return abc[i-1 : i]
}
