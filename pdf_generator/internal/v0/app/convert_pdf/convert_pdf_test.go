package convert_pdf

import (
	"github.com/ManyakRus/starter/pdf_generator/internal/v0/app/programdir"
	"testing"
)

func TestConvertXlsxPdf(t *testing.T) {

	ProgramDir := programdir.ProgramDir()
	FilenameIn := ProgramDir + "templates/test.xlsx"
	//FilenamePDF := ProgramDir + "templates/test_xlsx.pdf"

	err := ConvertToPdf(FilenameIn)
	if err != nil {
		t.Error("convert_pdf_test.TestConvertXlsxPdf() error: ", err)
	} else {
		t.Log("TestConvertXlsxPdf() ok. FilenameIn: ", FilenameIn)
	}

}

func TestConvertDocxPdf(t *testing.T) {

	ProgramDir := programdir.ProgramDir()
	FilenameIn := ProgramDir + "templates/test.docx"
	//FilenamePDF := ProgramDir + "templates/test_docx.pdf"

	err := ConvertToPdf(FilenameIn)
	if err != nil {
		t.Error("convert_pdf_test.TestConvertDocxPdf() error: ", err)
	} else {
		t.Log("TestConvertDocxPdf() ok. FilenameIn: ", FilenameIn)
	}

}

func TestConvertTxtPdf(t *testing.T) {

	ProgramDir := programdir.ProgramDir()
	FilenameIn := ProgramDir + "templates/test.txt"
	//FilenamePDF := ProgramDir + "templates/test_docx.pdf"

	err := ConvertToPdf(FilenameIn)
	if err != nil {
		t.Error("convert_pdf_test.TestConvertTxtPdf() error: ", err)
	} else {
		t.Log("TestConvertTxtPdf() ok. FilenameIn: ", FilenameIn)
	}

}

func TestConvertFodtPdf(t *testing.T) {

	ProgramDir := programdir.ProgramDir()
	FilenameIn := ProgramDir + "templates/test.fodt"
	//FilenamePDF := ProgramDir + "templates/test_docx.pdf"

	err := ConvertToPdf(FilenameIn)
	if err != nil {
		t.Error("convert_pdf_test.TestConvertFodtPdf() error: ", err)
	} else {
		t.Log("TestConvertFodtPdf() ok. FilenameIn: ", FilenameIn)
	}

}

func TestConvertFodsPdf(t *testing.T) {

	ProgramDir := programdir.ProgramDir()
	FilenameIn := ProgramDir + "templates/test.fods"
	//FilenamePDF := ProgramDir + "templates/test_docx.pdf"

	err := ConvertToPdf(FilenameIn)
	if err != nil {
		t.Error("convert_pdf_test.TestConvertFodsPdf() error: ", err)
	} else {
		t.Log("TestConvertFodsPdf() ok. FilenameIn: ", FilenameIn)
	}

}

func TestConvertClaim(t *testing.T) {

	ProgramDir := programdir.ProgramDir()
	FilenameIn := ProgramDir + "templates/Претензия.docx"
	//FilenamePDF := ProgramDir + "templates/test_docx.pdf"

	err := ConvertToPdf(FilenameIn)
	if err != nil {
		t.Error("convert_pdf_test.TestConvertClaim() error: ", err)
	} else {
		t.Log("TestConvertClaim() ok. FilenameIn: ", FilenameIn)
	}

}
