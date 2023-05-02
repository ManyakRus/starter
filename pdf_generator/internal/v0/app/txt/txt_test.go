package txt

import (
	"github.com/ManyakRus/starter/pdf_generator/internal/v0/app/programdir"
	"testing"
)

func TestCreateTxt(t *testing.T) {
	ProgramDir := programdir.ProgramDir()
	filename := ProgramDir + "templates/test.txt"

	FilenameOut := ProgramDir + "templates/test_ready.txt"

	map1 := make(map[string]string)
	map1["{{name}}"] = "Никитин А.В."

	err := CreateTxt(filename, FilenameOut, map1)
	if err != nil {
		t.Error("fodt_test.TestCreateTxt() error: ", err)
	}

}

func TestCreateFodt(t *testing.T) {
	ProgramDir := programdir.ProgramDir()
	filename := ProgramDir + "templates/test.fodt"

	FilenameOut := ProgramDir + "templates/test_ready.fodt"

	map1 := make(map[string]string)
	map1["{{name}}"] = "Никитин А.В."

	err := CreateTxt(filename, FilenameOut, map1)
	if err != nil {
		t.Error("fodt_test.TestCreateFodt() error: ", err)
	}

}

func TestCreateFods(t *testing.T) {
	ProgramDir := programdir.ProgramDir()
	filename := ProgramDir + "templates/test.fods"

	FilenameOut := ProgramDir + "templates/test_ready.fods"

	map1 := make(map[string]string)
	map1["{{name}}"] = "Никитин А.В."

	err := CreateTxt(filename, FilenameOut, map1)
	if err != nil {
		t.Error("fodt_test.TestCreateFods() error: ", err)
	}

}
