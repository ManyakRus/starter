package micro

import (
	"context"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/contextmain"
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

func TestProgramDir(t *testing.T) {
	dir := ProgramDir()
	if dir == "" {
		t.Error("microfunctions_test.TestProgramDir() ProgramDir() empty !")
	}
}

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

func TestMax(t *testing.T) {
	Otvet := Max(1, 2)
	if Otvet != 2 {
		t.Error("microfunctions_test.TestMax() error: Otvet != 2")
	}
}

func TestMin(t *testing.T) {
	Otvet := Min(1, 2)
	if Otvet != 1 {
		t.Error("microfunctions_test.TestMin() error: Otvet != 1")
	}
}

func TestMaxInt60(t *testing.T) {
	Otvet := MaxInt64(1, 2)
	if Otvet != 2 {
		t.Error("microfunctions_test.TestMax() error: Otvet != 2")
	}
}

func TestMinInt64(t *testing.T) {
	Otvet := MinInt64(1, 2)
	if Otvet != 1 {
		t.Error("microfunctions_test.TestMin() error: Otvet != 1")
	}
}

func TestGoGo(t *testing.T) {
	fn := func() error {
		Pause(2000)
		err := fmt.Errorf("test error")
		return err
	}

	ctxMain := contextmain.GetContext()
	ctx, cancel := context.WithTimeout(ctxMain, 1*time.Second)
	defer cancel()

	err := GoGo(ctx, fn)
	t.Log("Err:", err)
}

func TestMaxDate(t *testing.T) {
	now := time.Now()
	Otvet := MaxDate(now, time.Date(1, 1, 1, 1, 1, 1, 1, time.Local))
	if Otvet != now {
		t.Error("microfunctions_test.TestMaxDate() error: Otvet != ", now)
	}
}

func TestMinDate(t *testing.T) {
	now := time.Now()
	Otvet := MinDate(now, time.Date(9999, 1, 1, 1, 1, 1, 1, time.Local))
	if Otvet != now {
		t.Error("microfunctions_test.TestMinDate() error: Otvet != ", now)
	}
}

func TestCheckINNControlSum(t *testing.T) {
	Inn := ""
	err := CheckINNControlSum(Inn)
	if err == nil {
		t.Error("TestCheckINNControlSum() error")
	}

}

func TestCheckINNControlSum10(t *testing.T) {

	Inn := "5111002549"
	err := CheckINNControlSum10(Inn)
	if err != nil {
		t.Error("TestCheckINNControlSum10() error: ", err)
	}

}

func TestCheckINNControlSum12(t *testing.T) {

	Inn := "510800222725"
	err := CheckINNControlSum12(Inn)
	if err != nil {
		t.Error("TestCheckINNControlSum12() error: ", err)
	}

}

func TestStringFromInt64(t *testing.T) {
	Otvet := StringFromInt64(0)
	if Otvet != "0" {
		t.Error("TestStringFromInt64() error: != '0'")
	}
}

func TestStringDate(t *testing.T) {
	Otvet := StringDate(time.Now())
	if Otvet == "" {
		t.Error("TestStringDate() error: =''")
	}
}

func TestProgramDir_bin(t *testing.T) {
	Otvet := ProgramDir_bin()
	if Otvet == "" {
		t.Error("TestProgramDir_bin() error: =''")
	}
}

func TestSaveTempFile(t *testing.T) {

	bytes := []byte("123")
	Otvet := SaveTempFile(bytes)
	if Otvet == "" {
		t.Error("TestSaveTempFile() error: Otvet =''")
	}

}

func TestHash(t *testing.T) {
	Otvet := Hash("123")
	if Otvet == 0 {
		t.Error("TestHash() error: =0")
	}
}

func TestTextError(t *testing.T) {
	err := errors.New("1")
	s := TextError(err)
	if s != "1" {
		t.Error("TestTextError() error")
	}
}

func TestGetType(t *testing.T) {
	Otvet := GetType(1)
	if Otvet != "int" {
		t.Error("TestGetType() error: Otvet: ", Otvet)
	}
}

func TestFindFileNameShort(t *testing.T) {
	dir := ProgramDir()
	Otvet := FindFileNameShort(dir)
	if Otvet == "" {
		t.Error("TestFindFileNameShort() error: Otvet =''")
	}
}

func TestCurrentDirectory(t *testing.T) {

	Otvet := CurrentDirectory()
	if Otvet == "" {
		t.Error("TestCurrentDirectory() error: Otvet = ''")
	}
}

func TestBoolFromInt64(t *testing.T) {
	Otvet := BoolFromInt64(111)
	if Otvet != true {
		t.Error("TestBoolFromInt64() error: Otvet != true")
	}
}

func TestBoolFromInt(t *testing.T) {
	Otvet := BoolFromInt(111)
	if Otvet != true {
		t.Error("TestBoolFromInt64() error: Otvet != true")
	}
}

func TestDeleteFileSeperator(t *testing.T) {

	dir := "home" + SeparatorFile()
	dir = DeleteFileSeperator(dir)
	if dir != "home" {
		t.Error("TestDeleteFileSeperator() error")
	}
}

func TestCreateFolder(t *testing.T) {
	dir := ProgramDir()
	Filename := dir + "TestCreateFolder"
	err := CreateFolder(Filename, 0)
	if err != nil {
		t.Error("TestCreateFolder() error: ", err)
	}

	err = DeleteFolder(Filename)
	if err != nil {
		t.Error("TestCreateFolder() error: ", err)
	}

}

func TestDeleteFolder(t *testing.T) {
	dir := ProgramDir()
	err := DeleteFolder(dir + "TestCreateFolder")
	if err != nil {
		t.Error("TestDeleteFolder() error: ", err)
	}
}

func TestBoolFromString(t *testing.T) {
	Otvet := BoolFromString(" TrUe ")
	if Otvet != true {
		t.Error("TestBoolFromString() error: Otvet != true")
	}

	Otvet = BoolFromString("TrUe0")
	if Otvet != false {
		t.Error("TestBoolFromString() error: Otvet != true")
	}

}

func TestContextDone(t *testing.T) {
	// Testing when context is done
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if !ContextDone(ctx) {
		t.Error("Expected ContextDone to return true when context is done")
	}

	// Testing when context is not done
	ctx = context.Background()
	if ContextDone(ctx) {
		t.Error("Expected ContextDone to return false when context is not done")
	}
}

func TestStringFromUpperCase(t *testing.T) {
	// Testing empty string
	if result := StringFromUpperCase(""); result != "" {
		t.Errorf("Expected '', but got %s", result)
	}

	// Testing lowercase input
	if result := StringFromUpperCase("hello"); result != "Hello" {
		t.Errorf("Expected 'Hello', but got %s", result)
	}

	// Testing uppercase input
	if result := StringFromUpperCase("WORLD"); result != "WORLD" {
		t.Errorf("Expected 'WORLD', but got %s", result)
	}

	// Testing mixed case input
	if result := StringFromUpperCase("gOoD mOrNiNg"); result != "GOoD mOrNiNg" {
		t.Errorf("Expected 'GOoD mOrNiNg', but got %s", result)
	}
}

func TestStringFromLowerCase(t *testing.T) {
	// Testing an empty string
	input := ""
	expected := ""
	result := StringFromLowerCase(input)
	if result != expected {
		t.Errorf("Input: %s, Expected: %s, Result: %s", input, expected, result)
	}

	// Testing a string with a lowercase first letter
	input = "hello"
	expected = "hello"
	result = StringFromLowerCase(input)
	if result != expected {
		t.Errorf("Input: %s, Expected: %s, Result: %s", input, expected, result)
	}

	// Testing a string with an uppercase first letter
	input = "World"
	expected = "world"
	result = StringFromLowerCase(input)
	if result != expected {
		t.Errorf("Input: %s, Expected: %s, Result: %s", input, expected, result)
	}

	// Testing a string with special characters
	input = "Codeium"
	expected = "codeium"
	result = StringFromLowerCase(input)
	if result != expected {
		t.Errorf("Input: %s, Expected: %s, Result: %s", input, expected, result)
	}
}

func TestDeleteEndSlash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Text ends with /",
			input:    "example/",
			expected: "example",
		},
		{
			name:     "Text ends with \\",
			input:    "example\\",
			expected: "example",
		},
		{
			name:     "Text does not end with / or \\",
			input:    "example",
			expected: "example",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := DeleteEndSlash(test.input)
			if result != test.expected {
				t.Errorf("Expected %s, but got %s", test.expected, result)
			}
		})
	}
}

func TestInt64FromString(t *testing.T) {
	// Test converting a valid string to int64
	input1 := "12345"
	expected1 := int64(12345)
	result1, err1 := Int64FromString(input1)
	if err1 != nil {
		t.Errorf("Expected no error, but got: %v", err1)
	}
	if result1 != expected1 {
		t.Errorf("Expected %d, but got: %d", expected1, result1)
	}

	// Test converting an empty string to int64
	input2 := ""
	expected2 := int64(0)
	result2, err2 := Int64FromString(input2)
	if err2 == nil {
		t.Errorf("Expected error, but got: %v", err2)
	}
	if result2 != expected2 {
		t.Errorf("Expected %d, but got: %d", expected2, result2)
	}

	// Test converting an invalid string to int64
	input3 := "abc"
	expected3 := int64(0)
	result3, err3 := Int64FromString(input3)
	if err3 == nil {
		t.Error("Expected an error, but got none")
	}
	if result3 != expected3 {
		t.Errorf("Expected %d, but got: %d", expected3, result3)
	}
}

func TestFindLastPos(t *testing.T) {
	s := "Hello, World!"
	pos1 := FindLastPos(s, " ")
	if pos1 < 0 {
		t.Error("microfunctions_test.TestFindLastPos() FindLastPos()=nil !")
	}
}

func TestStringFloat64_Dimension2(t *testing.T) {
	// Testing for a positive float number
	result := StringFloat64_Dimension2(3.14159)
	if result != "3.14" {
		t.Errorf("Expected '3.14' but got %s", result)
	}

	// Testing for a negative float number
	result = StringFloat64_Dimension2(-123.456)
	if result != "-123.46" {
		t.Errorf("Expected '-123.46' but got %s", result)
	}

	// Testing for zero
	result = StringFloat64_Dimension2(0.0)
	if result != "0.00" {
		t.Errorf("Expected '0.00' but got %s", result)
	}
}

func TestStringFloat32_Dimension2(t *testing.T) {
	// Testing for a positive float number
	result := StringFloat32_Dimension2(3.14159)
	if result != "3.14" {
		t.Errorf("Expected '3.14' but got %s", result)
	}

	// Testing for a negative float number
	result = StringFloat32_Dimension2(-123.456)
	if result != "-123.46" {
		t.Errorf("Expected '-123.46' but got %s", result)
	}

	// Testing for zero
	result = StringFloat32_Dimension2(0.0)
	if result != "0.00" {
		t.Errorf("Expected '0.00' but got %s", result)
	}
}

func TestShowTimePassed(t *testing.T) {
	defer ShowTimePassed(time.Now())
}
