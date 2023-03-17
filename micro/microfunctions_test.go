package micro

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestAddSeparator(t *testing.T) {
	s := ""
	s2 := AddSeparator(s)
	if s2 != SeparatorFile() {
		t.Error("microfunctions_test.TestAddSeparator() AddSeparator() error !")
	}
}

func TestFileExists(t *testing.T) {
	filename := CurrentFilename()
	ok, err := FileExists(filename)
	if ok == false {
		t.Error("microfunctions_test.TestFileExists() FileExists() !=true !")
	}

	if err != nil {
		t.Error("microfunctions_test.TestFileExists() FileExists() error!=nil!")
	}
}

func TestFindDirUp(t *testing.T) {

	dir := ProgramDir_Common()
	dir2 := FindDirUp(dir)
	if dir2 == "" || dir == dir2 {
		t.Error("microfunctions_test.TestFindDirUp() FindDirUp() error !")
	}
}

func TestIsTestApp(t *testing.T) {
	stage0 := os.Getenv("STAGE")
	err := os.Setenv("STAGE", "local")
	if err != nil {
		t.Error("microfunctions_test.TestIsTestApp() os.Setenv() error: ", err)
	}

	isTestApp := IsTestApp()
	if isTestApp != true {
		t.Error("microfunctions_test.TestIsTestApp() <>true !")
	}
	err = os.Setenv("STAGE", stage0)
	if err != nil {
		t.Error("microfunctions_test.TestIsTestApp() os.Setenv() error: ", err)
	}
}

func TestPause(t *testing.T) {
	Pause(1)
}

//func TestProgramDir(t *testing.T) {
//	dir := ProgramDir()
//	if dir == "" {
//		t.Error("microfunctions_test.TestProgramDir() ProgramDir() empty !")
//	}
//}

func TestSeparatorFile(t *testing.T) {
	s := SeparatorFile()
	if s != "/" && s != "\\" {
		t.Error("microfunctions_test.TestSeparatorFile() SeparatorFile() error !")
	}
}

func TestSleep(t *testing.T) {
	Sleep(1)
}

func TestCurrentFilename(t *testing.T) {
	filename := CurrentFilename()
	if filename == "" {
		t.Error("microfunctions_test.TestCurrentFilename() CurrentFilename() error !")
	}
}

func TestErrorJoin(t *testing.T) {
	err1 := errors.New("1")
	err2 := errors.New("2")
	err := ErrorJoin(err1, err2)
	if err == nil {
		t.Error("microfunctions_test.TestErrorJoin() ErrorJoin()=nil !")
	}
}

func TestSubstringLeft(t *testing.T) {
	otvet := SubstringLeft("123", 1)
	if otvet != "1" {
		t.Error("microfunctions.TestSubstringLeft() error !")
	}
}

func TestSubstringRight(t *testing.T) {
	otvet := SubstringRight("123", 1)
	if otvet != "3" {
		t.Error("microfunctions.TestSubstringLeft() error !")
	}
}

func TestStringBetween(t *testing.T) {
	s := "123"
	otvet := StringBetween(s, "1", "3")
	if otvet != "2" {
		t.Error("microfunctions_test.TestStringBetween() error !")
	}
}

func TestLastWord(t *testing.T) {
	s := "а.б_б"
	s2 := LastWord(s)
	if s2 != "б_б" {
		t.Error("TestLastWord error")
	}
}

func TestFileNameWithoutExtension(t *testing.T) {

	filename := "test.xlsx"
	filename2 := FileNameWithoutExtension(filename)
	if filename2 != "test" {
		t.Error("microfunctions_test.TestFileNameWithoutExtension() error !")
	}

}

func TestBeginningOfMonth(t *testing.T) {

	l := time.Local
	DateTest := time.Date(2022, 1, 10, 0, 0, 0, 0, l)
	DateGood := time.Date(2022, 1, 1, 0, 0, 0, 0, l)

	Date1 := BeginningOfMonth(DateTest)
	if Date1 != DateGood {
		t.Error("microfunctions_test TestBeginningOfMonth() error Date1 != DateGood")
	}

}

func TestEndOfMonth(t *testing.T) {
	l := time.Local
	DateTest := time.Date(2022, 1, 10, 0, 0, 0, 0, l)
	DateGood := time.Date(2022, 1, 31, 0, 0, 0, 0, l)

	Date1 := EndOfMonth(DateTest)
	if Date1 != DateGood {
		t.Error("microfunctions_test TestBeginningOfMonth() error Date1 != DateGood")
	}

}

func TestStringAfter(t *testing.T) {
	s := "123456"
	s2 := "34"
	Otvet := StringAfter(s, s2)
	if Otvet != "56" {
		t.Error("TestStringAfter error")
	}
}

func TestStringFrom(t *testing.T) {
	s := "123456"
	s2 := "34"
	Otvet := StringFrom(s, s2)
	if Otvet != "3456" {
		t.Error("TestStringAfter error")
	}
}

func TestTrim(t *testing.T) {

	s := ` 1234		
`
	Otvet := Trim(s)
	if Otvet != "1234" {
		t.Error("microfunctions_test.TestTrim() error")
	}
}
