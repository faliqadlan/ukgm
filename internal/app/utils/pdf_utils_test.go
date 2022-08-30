package utils_test

import (
	"os"
	"testing"

	"github.com/faliqadlan/ukgm/internal/app/utils"

	loremipsum "gopkg.in/loremipsum.v1"
)

func TestCreatePDF(t *testing.T) {
	loremIpsumGenerator := loremipsum.New()
	loremIpsumGenerator.Paragraph()
	mydir, _ := os.Getwd()
	pdf := utils.NewPDF("Testing")
	pdf.SetLocation(mydir + "/storage/temp/")
	pdf.InitPDF()
	pdf.PrintHeader1(true, "JUDUL COBAIN LOCATION", "R", "B")
	pdf.PrintText(false, "Saya anak baik, suka menabung, tidak jahat dan masih hidup")
	pdf.PrintText(false, loremIpsumGenerator.Paragraph()+loremIpsumGenerator.Paragraph())
	pdf.LineBreak(5)
	pdf.PrintText(false, "Saya anak baik, suka menabung, tidak jahat dan masih hidup")
	pdf.LineBreak()
	pdf.PrintHorizontalLine(1000, true)
	pdf.PrintVerticalLine(1000, true, "R")
	pdf.PrintVerticalLine(1000, true, "L")
	pdf.PrintText(false, "Saya anak baik, suka menabung, tidak jahat dan masih hidup")
	pdf.PrintVerticalLine(1000, true)

	pdf.LineBreak()
	pdf.PrintText(false, loremIpsumGenerator.Paragraph()+loremIpsumGenerator.Paragraph())
	pdf.LineBreak()
	pdf.AddPage()
	pdf.PrintHeader1(true, "CELL", "B")
	pdf.PrintText(false, "Saya anak baik, suka menabung, tidak jahat dan masih hidup ")
	pdf.PrintInCell("loremIpsumGenerator.Paragraph()loremIpsumGenerator.Paragraph()loremIpsumGenerator.Paragraph()loremIpsumGenerator.Paragraph()", 20.0, true)
	pdf.PrintText(false, "Saya anak baik, suka menabung, tidak jahat dan masih hidup")
	pdf.LineBreak()
	pdf.PrintText(false, loremIpsumGenerator.Paragraph()+loremIpsumGenerator.Paragraph())
	pdf.AddNewPage()
	pdf.PrintText(false, loremIpsumGenerator.Paragraph()+loremIpsumGenerator.Paragraph())
	pdf.SavePDF()
}
