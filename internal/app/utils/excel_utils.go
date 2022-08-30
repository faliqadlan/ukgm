package utils

import (
	"strconv"

	excelize "github.com/xuri/excelize/v2"
)

func GenerateExcel() (excel *excelize.File) {
	excel = excelize.NewFile()
	return
}

// AlphabetIndex : Create Index of String based on Alphabet
func AlphabetIndex() []string {
	alphabetIndex := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	return alphabetIndex
}

// IndexToColumn : Change Index to Culumn (Limited only two digits)
func IndexToColumn(columnIndex int) string {
	indexColumn := AlphabetIndex()
	var column string

	if columnIndex == len(indexColumn) {
		column = indexColumn[0] + indexColumn[0]
	} else if columnIndex > len(indexColumn) {
		column = indexColumn[(columnIndex/len(indexColumn))-1] + indexColumn[(columnIndex%len(indexColumn))]
	} else {
		column = indexColumn[columnIndex]
	}

	return column
}

// WriteCell : Write Value to Excel using Index for Column
func WriteCell(row, columnIndex int, value interface{}, sheetName string, xlsx *excelize.File, colStyle, cellstyle *excelize.Style, bold bool) int {
	var (
		column                 string = IndexToColumn(columnIndex)
		styleInt, cellStyleInt int
	)

	cell := column
	cell += strconv.Itoa(row)

	if colStyle != nil {
		styleInt, _ = xlsx.NewStyle(colStyle)
		xlsx.SetColStyle(sheetName, column, styleInt)
	}
	xlsx.SetCellValue(sheetName, cell, value)
	if cellstyle != nil || bold {
		if cellstyle == nil {
			cellstyle = &excelize.Style{}
		}
		if bold {
			if cellstyle.Font == nil {
				cellstyle.Font = &excelize.Font{}
			}
			cellstyle.Font.Bold = true
		}
		cellStyleInt, _ = xlsx.NewStyle(cellstyle)
		xlsx.SetCellStyle(sheetName, cell, cell, cellStyleInt)
	}

	return (columnIndex + 1)
}

func GetTimeStyle() *excelize.Style {
	timeStyle := excelize.Style{
		NumFmt: 15,
	}
	return &timeStyle
}

func GetTextStyle() *excelize.Style {
	textStyle := excelize.Style{
		NumFmt: 49,
	}
	return &textStyle
}
