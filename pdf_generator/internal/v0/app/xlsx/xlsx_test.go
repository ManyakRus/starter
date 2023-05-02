package xlsx

import (
	"github.com/ManyakRus/starter/pdf_generator/internal/v0/app/programdir"
	"testing"
)

func TestCreateXLSX1(t *testing.T) {

	ProgramDir := programdir.ProgramDir()
	filename := ProgramDir + "templates/test.xlsx"
	FilenameOut := ProgramDir + "templates/test_ready.xlsx"

	map1 := make(map[string]interface{})
	map1["name"] = "Никитин А.В."

	err := CreateXLSX(filename, FilenameOut, map1)
	if err != nil {
		t.Error("xlsx_test.TestCreateXLSX1() error: ", err)
	}
}
