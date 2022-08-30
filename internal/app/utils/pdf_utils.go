package utils

import (
	"fmt"
	"strings"

	"github.com/iman231/go-help"

	gofpdf "github.com/jung-kurt/gofpdf/v2"
	"go.uber.org/zap"
)

/*
	Constants description

	Only use some of the align used in this package

	LeftAlign = L

	CenterAlign = C

	RightAlign = R
*/
const (
	LeftAlign   = "L"
	CenterAlign = "C"
	RightAlign  = "R"
)

/*
	Variables description

	AllAligns is array consists every align used in this package
*/
var (
	AllAligns = []string{
		LeftAlign, CenterAlign, RightAlign,
	}
)

/*
	StyleStringPipeline

	Pipeline to get matching style within the rule, every unmatch style will be ignored

	B = Bold

	I = Italic

	S = Strikethrough

	U = Underline
*/
func StyleStringPipeline(StyleString []string) (result []string) {
	EnabledStyle := []string{
		"B", "I", "S", "U",
	}
	for _, style := range StyleString {
		if len(style) > 1 {
			for _, styleChar := range style {
				if help.InArray(string(styleChar), EnabledStyle) {
					result = append(result, string(styleChar))
				}
			}
		} else {
			if help.InArray(style, EnabledStyle) {
				result = append(result, style)
			}
		}
	}
	return
}

/*
	AlignPipeline

	Pipeline to get matching align within the rule, every unmatch align will be ignored
*/
func AlignPipeline(StyleString []string) (result string) {
	for _, style := range StyleString {
		if help.InArray(style, AllAligns) {
			return style
		}
	}
	return
}

type (
	/*
		PDFStruct

		Struct to help process and customize PDF file using jung-kurt gofpdf
		this struct contain information for construct the PDF

		orientation : Potrait or Landscape (depends on jung-kurt package use) Default is Potrait

		unitStr : Specifies the unit of length used in size parameters (depends on jung-kurt package use) Default is mm

		pageSize : Specifies the page size used for the pdf (depends on jung-kurt package use) Default is A4

		filename : Specifies the file name (There's no default)

		location : Specifies the file location

		topMargin : Specifies top margin of the document, default is 15

		leftMargin : Specifies left margin of the document, default is 15

		rightMargin : Specifies right margin of the document, default is 15

		textAlign : Specifies top margin of the document, default is 15

		font : Specifies font used in document, default is Arial

		h1Size : Specifies font size for Header 1, default size is 16 in mm

		textSize : Specifies basic text size, default size is 11 in mm. This variable is used as basic text but it's not directly used for text print, this variable is replace tempTextSize if the font size is reset

		tempTextSize : Variable used for current text size, default size is follow textSize. This variable is the one used for text print

		lineBreakSize : Specified line break size, default size is 5 in mm
	*/
	PDFStruct struct {
		*gofpdf.Fpdf
		orientation string
		unitStr     string
		pageSize    string
		fileName    string
		location    string

		topMargin   float64
		leftMargin  float64
		rightMargin float64

		textAlign     string
		font          string
		h1Size        float64
		textSize      float64
		tempTextSize  float64
		lineBreakSize float64
	}
)

// NewPDF Create default PDFStruct object and every variables value is default. String parameter is needed for the file name
func NewPDF(filename string) (result *PDFStruct) {
	if filename == "" {
		filename = "NewFile"
	}
	result = &PDFStruct{
		orientation:   "P",
		unitStr:       "mm",
		pageSize:      "A4",
		fileName:      filename,
		h1Size:        16,
		textAlign:     "L",
		font:          "Arial",
		textSize:      11,
		tempTextSize:  11,
		leftMargin:    15,
		topMargin:     15,
		rightMargin:   15,
		lineBreakSize: 5,
	}
	return
}

// AddNewPage Add new page and reset margin, font and reset current position to first position in new page
func (pdf *PDFStruct) AddNewPage() {
	pdf.AddPage()
	pdf.MarginReset()
	pdf.FontSizeReset()
	pdf.ResetPosition()
}

// ResetPosition Reset position following page margin
func (pdf *PDFStruct) ResetPosition() {
	pdf.SetX(pdf.leftMargin)
	pdf.SetY(pdf.topMargin)
}

// SetOrientation Set page orientation, this function can be used only before the pdf initialized
func (pdf *PDFStruct) SetOrientation(orientation string) (err error) {
	if pdf.Fpdf != nil {
		err = fmt.Errorf("pdf has been created, format cannot be changed")
		zap.S().Error(err)
		return
	}
	if !help.InArray(orientation, []string{"P", "L"}) {
		err = fmt.Errorf("orientation is not listed")
		zap.S().Error(err)
		return
	}
	pdf.orientation = orientation
	return
}

// SetUnitStr Set unit measurement for document, this function can be used only before the pdf initialized
func (pdf *PDFStruct) SetUnitStr(unitStr string) (err error) {
	if pdf.Fpdf != nil {
		err = fmt.Errorf("pdf has been created, format cannot be changed")
		zap.S().Error(err)
		return
	}
	if !help.InArray(unitStr, []string{"mm", "cm", "pt"}) {
		err = fmt.Errorf("unitStr is not listed")
		zap.S().Error(err)
		return
	}
	pdf.unitStr = unitStr
	return
}

// SetLineBreak Set document line break size, this function can be used only before the pdf initialized
func (pdf *PDFStruct) SetLineBreak(lineBreak float64) (err error) {
	if pdf.Fpdf != nil {
		err = fmt.Errorf("pdf has been created, format cannot be changed")
		zap.S().Error(err)
		return
	}
	if lineBreak == 0 {
		return
	}
	pdf.lineBreakSize = lineBreak
	return
}

// SetPageSize Set document page size, this function can be used only before the pdf initialized
func (pdf *PDFStruct) SetPageSize(pageSize string) (err error) {
	if pdf.Fpdf != nil {
		err = fmt.Errorf("pdf has been created, format cannot be changed")
		zap.S().Error(err)
		return
	}
	if !help.InArray(pageSize, []string{"a3", "a4", "a5", "a6", "a2", "a1", "letter", "legal", "tabloid"}) {
		err = fmt.Errorf("PageSize is not listed")
		zap.S().Error(err)
		return
	}
	pdf.pageSize = pageSize
	return
}

// SetH1Size Set document Header size, this function can be used only before the pdf initialized
func (pdf *PDFStruct) SetH1Size(h1Size float64) (err error) {
	if pdf.Fpdf != nil {
		err = fmt.Errorf("pdf has been created, format cannot be changed")
		zap.S().Error(err)
		return
	}
	if h1Size <= 0 {
		return
	}
	pdf.h1Size = h1Size
	return
}

// SetLocation Set document save directory
func (pdf *PDFStruct) SetLocation(location string) {
	pdf.location = location
}

// ChangeTextSize Change document text size including current text size, every texts already printed won't be changed
func (pdf *PDFStruct) ChangeTextSize(textSize float64) {
	pdf.textSize = textSize
	pdf.tempTextSize = textSize
}

// LeftMarginReset Change current left margin to left margin that has initialized first
func (pdf *PDFStruct) LeftMarginReset() {
	pdf.SetLeftMargin(pdf.leftMargin)
}

// RightMarginReset Change current right margin to right margin that has initialized first
func (pdf *PDFStruct) RightMarginReset() {
	pdf.SetRightMargin(pdf.rightMargin)
}

// Change current top margin to top margin that has initialized first
func (pdf *PDFStruct) TopMarginReset() {
	pdf.SetTopMargin(pdf.topMargin)
}

// MarginReset Reset every margin use margin sizes that has initialized first
func (pdf *PDFStruct) MarginReset() {
	pdf.TopMarginReset()
	pdf.LeftMarginReset()
	pdf.RightMarginReset()
}

// FontSizeReset Set current text font size to text font size in variable textSize that has initialized
func (pdf *PDFStruct) FontSizeReset() {
	pdf.tempTextSize = pdf.textSize
}

// ChangeCurrentTextSize Set current text font size to desired size
func (pdf *PDFStruct) ChangeCurrentTextSize(textSize float64) {
	pdf.tempTextSize = textSize
}

// InitPDF Initialized pdf file and ready to be filled
func (pdf *PDFStruct) InitPDF() {
	// INIT PDF
	pdf.Fpdf = gofpdf.New(pdf.orientation, pdf.unitStr, pdf.pageSize, "")
	pdf.FontSizeReset()
	pdf.MarginReset()
	pdf.AddPage()
	return
}

/*
	SavePDF.

	Save pdf file to storage use the directory initialized before, return the filename and error if there's problem when save the document

*/
func (pdf *PDFStruct) SavePDF() (fileName string, err error) {
	fileName = pdf.fileName + ".pdf"

	if pdf.location == "" {
		pdf.location = "./internal/storage/temp/"
	}
	if err := pdf.Fpdf.OutputFileAndClose(pdf.location + fileName); err != nil {
		zap.S().Error(err)
		return "", err
	}

	return
}

/*
	PrintHeader1

	Print header using the initialized font and header size and automatically line break once.

	Parameters needed

	Capitalized type boolean. If the value is true then the inputted text will be capitalized

	Text type string. The desired text used for title, if the text is "" then the function will set it to default text as "title"

	Style type string this parameters is non-mandatory parameters, used for text style like Bold, Italic, Striketrough and Underline
	for the style parameter can be combined ex:"BI","BU", etc. also can be used for title alignment position,
	for alignment position used only one character for alignment and separate parameter from text style ex: PrintHeader1(true, "title", "BI", "C")
*/
func (pdf *PDFStruct) PrintHeader1(Capitalized bool, Text string, Style ...string) {
	var StyleString string
	if Text == "" {
		Text = "Title"
	}
	if Capitalized {
		Text = strings.ToUpper(Text)
	}
	var align string = pdf.textAlign
	if tempAlign := AlignPipeline(Style); tempAlign != "" {
		align = tempAlign
	}
	StyleString = strings.Join(StyleStringPipeline(Style), "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont(pdf.font, StyleString, pdf.h1Size)
	pdf.WriteAligned(0, 0, Text, align)
	pdf.Ln(pdf.lineBreakSize)

	return
}

/*
	PrintText

	Print text using the initialized font and text size

	Parameters needed

	Indentation type boolean. If the value is true then the inputted text will have indentationconst

	Text type string. The desired text printed in current position and margin set in pdf

	AlignStyle type string this parameters is non-mandatory parameters, used for text style like Bold, Italic, Striketrough and Underline
	for the style parameter can be combined ex:"BI","BU", etc. also can be used for title alignment position,
	for alignment position used only one character for alignment and separate parameter from text style ex: PrintText(true, "title", "BI", "C")
*/
func (pdf *PDFStruct) PrintText(Indentation bool, Text string, AlignStyle ...string) {
	var Style string
	Style = strings.Join(StyleStringPipeline(AlignStyle), "")
	if Indentation {
		pdf.WriteAligned(0, pdf.tempTextSize/3, "\t\t", pdf.textAlign)
	}
	var align string = pdf.textAlign
	if tempAlign := AlignPipeline(AlignStyle); tempAlign != "" {
		align = tempAlign
	}
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont(pdf.font, Style, pdf.tempTextSize)
	pdf.WriteAligned(0, pdf.tempTextSize/3, Text, align)
	return
}

/*
	PrintHorizontalLine.

	Create horizontal line in document

	Parameters needed

	length type float64. Desired length in document measurement, if the line length exceed the margin width then the line will automatically adjusted across the document

	Under type boolean. Set the line is on the upper side or bottom side on current row

	Align type strings. Set the position on page (use "L","C","R"), if there's no align inputted or invalid align then automatically set to default text align
*/
func (pdf *PDFStruct) PrintHorizontalLine(length float64, Under bool, Align ...string) {
	var align string = pdf.textAlign
	if tempAlign := AlignPipeline(Align); tempAlign != "" {
		align = tempAlign
	}
	Y := pdf.GetY()
	X1 := pdf.GetX()

	if Under {
		Y += (pdf.textSize / 3)
	}
	width, _ := pdf.GetPageSize()

	if align == CenterAlign {
		centerPoint := width / 2
		X1 = centerPoint - (length / 2)
		if X1 < pdf.leftMargin {
			X1 = pdf.leftMargin
		}
		if X1+length >= width {
			length = width - pdf.leftMargin - pdf.rightMargin
		}
	} else if align == LeftAlign {
		if X1+length >= width {
			length = width - pdf.leftMargin - pdf.rightMargin
		}
	} else if align == RightAlign {
		X1 = width - pdf.leftMargin
		if X1-length <= 0 {
			length = -(width - pdf.leftMargin)
		} else {
			length *= -1
		}
	}

	X2 := X1 + length
	pdf.Line(X1, Y, X2, Y)
	return
}

/*
	PrintVerticalLine.

	Create vertical line in document

	Parameters needed

	length type float64. Desired length in document measurement, if the line length exceed the margin width then the line will automatically adjusted across the document

	Under type boolean. Set the line is on the upper side or bottom side on current row

	Align type strings. Set the position on page (use "L","C","R"), if there's no align inputted or invalid align then the line will appear from current X position
*/
func (pdf *PDFStruct) PrintVerticalLine(length float64, Under bool, Align ...string) {
	var X float64
	var align string
	if tempAlign := AlignPipeline(Align); tempAlign != "" {
		align = tempAlign
	}
	width, height := pdf.GetPageSize()
	if align == CenterAlign {
		X = width / 2
	} else if align == LeftAlign {
		X = pdf.leftMargin
	} else if align == RightAlign {
		X = width - pdf.rightMargin
	} else if align == "" {
		X = pdf.GetX()
	}
	Y1 := pdf.GetY()
	if Under {
		Y1 += (pdf.textSize / 3)
	}
	Y2 := Y1 + length
	if Y2 >= height {
		Y2 = height - pdf.topMargin
	}
	pdf.Line(X, Y1, X, Y2)
	return
}

/*
	LineBreak.

	Create line break in document

	Parameter needed

	loop type integers. loop the line break (loop) times. If the parameter is not inputted then line break once
*/
func (pdf *PDFStruct) LineBreak(loop ...int) {
	if len(loop) >= 1 {
		for i := 0; i < loop[0]; i++ {
			pdf.Ln(pdf.lineBreakSize)
		}
	} else {
		pdf.Ln(pdf.lineBreakSize)
	}
	return
}

/*
	PrintInCell.

	Print text inside cell with desired width

	Parameters needed

	text type string. Inputted text will be printed inside the cell

	width type float64. Desired width cell

	backToTop type boolean. If the value is true then the next text will be placed at the top after the cell

*/
func (pdf *PDFStruct) PrintInCell(text string, width float64, backToTop bool) {
	Y := pdf.GetY()
	X := pdf.GetX()
	currentPage := pdf.PageNo()
	PageWidth, _ := pdf.GetPageSize()
	if pdf.GetX()+width >= PageWidth-pdf.rightMargin {
		width = (PageWidth - pdf.GetX()) - pdf.rightMargin
		backToTop = false
	}
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont(pdf.font, "", pdf.tempTextSize)
	pdf.MultiCell(width, pdf.tempTextSize/3, text, "", "", false)
	if currentPage != pdf.PageNo() {
		backToTop = false
	}
	if backToTop {
		pdf.SetLeftMargin(X + width)
		pdf.SetY(Y)
		// pdf.SetX(X + width)
	}
}
