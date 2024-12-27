package micro

import (
	"context"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"reflect"
	"strings"
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

func TestStringFromInt32(t *testing.T) {
	Otvet := StringFromInt32(0)
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

func TestShowTimePassed_FormatText(t *testing.T) {
	defer ShowTimePassed_FormatText("time passed: %s", time.Now())
}

func TestShowTimePassedSeconds(t *testing.T) {
	defer ShowTimePassedSeconds(time.Now())
	Pause(1)
}

func TestShowTimePassedMilliSeconds(t *testing.T) {
	defer ShowTimePassedMilliSeconds(time.Now())
	Pause(1)
}

func TestStructDeepCopy(t *testing.T) {
	type NestedStruct struct {
		Number int
		Text   string
	}

	type TestStruct struct {
		ID      int
		Name    string
		Nested  NestedStruct
		Numbers []int
	}

	src := TestStruct{
		ID:   1,
		Name: "Test",
		Nested: NestedStruct{
			Number: 100,
			Text:   "Nested",
		},
		Numbers: []int{1, 2, 3},
	}

	var dist TestStruct

	err := StructDeepCopy(src, &dist)
	if err != nil {
		t.Fatalf("error copying struct: %v", err)
	}

	if !reflect.DeepEqual(src, dist) {
		t.Errorf("copied struct does not match original struct")
	}
}

func TestIsEmptyValue(t *testing.T) {
	// Testing for integer zero value
	if !IsEmptyValue(0) {
		t.Error("Expected true for integer zero value")
	}

	// Testing for empty string value
	if !IsEmptyValue("") {
		t.Error("Expected true for empty string value")
	}

	// Testing for empty uuid value
	uuid1 := uuid.Nil
	if !IsEmptyValue(uuid1) {
		t.Error("Expected true for empty uuid value")
	}
}

func TestStringIdentifierFromUUID(t *testing.T) {
	// Test that the function returns a non-empty string
	result := StringIdentifierFromUUID()
	if result == "" {
		t.Error("Expected non-empty string, but got empty string")
	}

	// Test that the function returns a string of the correct length
	expectedLength := 32
	if len(result) != expectedLength {
		t.Errorf("Expected string of length %d, but got %d", expectedLength, len(result))
	}

	// Test that the function returns a string with no hyphens
	if strings.Contains(result, "-") {
		t.Error("Expected string with no hyphens, but got hyphen")
	}
}

func TestIndexSubstringMin(t *testing.T) {
	// Test case 1: empty input string and no substrings provided
	s1 := ""
	Otvet1 := -1
	if IndexSubstringMin(s1) != Otvet1 {
		t.Errorf("IndexSubstringMin(%q) = %d; want %d", s1, IndexSubstringMin(s1), Otvet1)
	}

	// Test case 2: non-empty input string and no substrings provided
	s2 := "Hello, world!"
	Otvet2 := -1
	if IndexSubstringMin(s2) != Otvet2 {
		t.Errorf("IndexSubstringMin(%q) = %d; want %d", s2, IndexSubstringMin(s2), Otvet2)
	}

	// Test case 3: input string contains one of the substrings
	s3 := "Hello, world!"
	substrings3 := []string{"world"}
	Otvet3 := 7
	if IndexSubstringMin(s3, substrings3...) != Otvet3 {
		t.Errorf("IndexSubstringMin(%q, %v...) = %d; want %d", s3, substrings3, IndexSubstringMin(s3, substrings3...), Otvet3)
	}

	// Test case 4: input string contains multiple occurrences of the same substring
	s4 := "Hello, world! Hello, world!"
	substrings4 := []string{"world"}
	Otvet4 := 7
	if IndexSubstringMin(s4, substrings4...) != Otvet4 {
		t.Errorf("IndexSubstringMin(%q, %v...) = %d; want %d", s4, substrings4, IndexSubstringMin(s4, substrings4...), Otvet4)
	}

	// Test case 5: input string contains multiple different substrings
	s5 := "Hello, world! How are you?"
	substrings5 := []string{"world", "are"}
	Otvet5 := 7
	if IndexSubstringMin(s5, substrings5...) != Otvet5 {
		t.Errorf("IndexSubstringMin(%q, %v...) = %d; want %d", s5, substrings5, IndexSubstringMin(s5, substrings5...), Otvet5)
	}

	// Test case 6: input string contains a mix of substrings that overlap and don't overlap
	s6 := "Hello, world! How are you?"
	substrings6 := []string{"world", "orl"}
	Otvet6 := 7
	if IndexSubstringMin(s6, substrings6...) != Otvet6 {
		t.Errorf("IndexSubstringMin(%q, %v...) = %d; want %d", s6, substrings6, IndexSubstringMin(s6, substrings6...), Otvet6)
	}

	// Test case 6: input string contains a mix of substrings that overlap and don't overlap
	s7 := "Hello, world! How are you?"
	substring7 := "world"
	substring8 := "How"
	Otvet7 := 7
	if IndexSubstringMin(s7, substring7, substring8) != Otvet7 {
		t.Errorf("IndexSubstringMin(%q, %v...) = %d; want %d", s7, substring7, IndexSubstringMin(s7, substring7, substring8), Otvet7)
	}
}

func TestIndexSubstringMin2(t *testing.T) {
	tests := []struct {
		s        string
		substr1  string
		substr2  string
		expected int
	}{
		{s: "hello world", substr1: "world", substr2: "test", expected: 6},
		{s: "hello world", substr1: "", substr2: "test", expected: -1},
		{s: "hello world", substr1: "test", substr2: "", expected: -1},
		{s: "hello world", substr1: "", substr2: "", expected: -1},
		{s: "hello world", substr1: "world", substr2: "hello", expected: 0},
		{s: "hello world", substr1: "test", substr2: "world", expected: 6},
		{s: "hello world", substr1: "test", substr2: "test", expected: -1},
	}

	for _, test := range tests {
		result := IndexSubstringMin2(test.s, test.substr1, test.substr2)
		if result != test.expected {
			t.Errorf("IndexSubstringMin2(%q, %q, %q) = %d, expected %d", test.s, test.substr1, test.substr2, result, test.expected)
		}
	}
}

func TestSortkMapStringInt(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []string
	}{
		{
			name:     "Empty map",
			input:    map[string]int{},
			expected: []string{},
		},
		{
			name:     "Single element map",
			input:    map[string]int{"a": 1},
			expected: []string{"a"},
		},
		{
			name:     "Multiple element map",
			input:    map[string]int{"a": 1, "b": 2, "c": 3},
			expected: []string{"c", "b", "a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortMapStringInt_Desc(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestIsNilInterface(t *testing.T) {
	// Testing for a nil pointer interface
	var ptr *int
	if !IsNilInterface(ptr) {
		t.Error("Expected true for nil pointer interface")
	}

	// Testing for a nil slice interface
	var slice []int
	if !IsNilInterface(slice) {
		t.Error("Expected true for nil slice interface")
	}

	// Testing for a non-nil map interface
	m := make(map[string]int)
	if IsNilInterface(m) {
		t.Error("Expected false for non-nil map interface")
	}

	// Testing for a nil function interface
	var fn func()
	if !IsNilInterface(fn) {
		t.Error("Expected true for nil function interface")
	}

	// Testing for a nil interface
	var i interface{}
	if !IsNilInterface(i) {
		t.Error("Expected true for nil interface")
	}
}

func TestStringFromMassInt64(t *testing.T) {
	// Test with an empty array
	emptyArray := []int64{}
	emptyResult := StringFromMassInt64(emptyArray, ",")
	if emptyResult != "" {
		t.Errorf("Expected empty string, but got: %s", emptyResult)
	}

	// Test with an array of single element
	singleArray := []int64{42}
	singleResult := StringFromMassInt64(singleArray, ",")
	if singleResult != "42" {
		t.Errorf("Expected '42', but got: %s", singleResult)
	}

	// Test with an array of multiple elements
	multipleArray := []int64{1, 2, 3}
	multipleResult := StringFromMassInt64(multipleArray, "-")
	expectedResult := "1-2-3"
	if multipleResult != expectedResult {
		t.Errorf("Expected '%s', but got: %s", expectedResult, multipleResult)
	}
}

func TestIsInt(t *testing.T) {
	// Test with an empty string
	emptyResult := IsInt("")
	if emptyResult != false {
		t.Errorf("Expected false for empty string, but got: %v", emptyResult)
	}

	// Test with a string containing only digits
	digitResult := IsInt("12345")
	if digitResult != true {
		t.Errorf("Expected true for string containing only digits, but got: %v", digitResult)
	}

	// Test with a string containing non-digit characters
	nonDigitResult := IsInt("abc123")
	if nonDigitResult != false {
		t.Errorf("Expected false for string containing non-digit characters, but got: %v", nonDigitResult)
	}
}

func TestInt32FromString(t *testing.T) {
	// Test converting a valid string to int32
	input1 := "12345"
	expected1 := int32(12345)
	result1, err1 := Int32FromString(input1)
	if err1 != nil {
		t.Errorf("Expected no error, but got: %v", err1)
	}
	if result1 != expected1 {
		t.Errorf("Expected %d, but got: %d", expected1, result1)
	}

	// Test converting an empty string to int32
	input2 := ""
	expected2 := int32(0)
	result2, err2 := Int32FromString(input2)
	if err2 == nil {
		t.Errorf("Expected error, but got: %v", err2)
	}
	if result2 != expected2 {
		t.Errorf("Expected %d, but got: %d", expected2, result2)
	}

	// Test converting an invalid string to int32
	input3 := "abc"
	expected3 := int32(0)
	result3, err3 := Int32FromString(input3)
	if err3 == nil {
		t.Error("Expected an error, but got none")
	}
	if result3 != expected3 {
		t.Errorf("Expected %d, but got: %d", expected3, result3)
	}
}

func TestPause_ctx(t *testing.T) {
	t.Run("Context canceled before pause duration elapses", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		start := time.Now()
		Pause_ctx(ctx, 100)
		duration := time.Since(start)

		if duration > 50*time.Millisecond {
			t.Errorf("Pause_ctx did not return in a timely manner when context was canceled")
		}
	})

	t.Run("Context not canceled before pause duration elapses", func(t *testing.T) {
		ctx := context.Background()

		start := time.Now()
		Pause_ctx(ctx, 100)
		duration := time.Since(start)

		if duration < 100*time.Millisecond {
			t.Errorf("Pause_ctx did not wait for the specified duration when context was not canceled")
		}
	})
}

func TestExecuteShellCommand_EmptyCommand(t *testing.T) {
	expected := ""
	result, err := ExecuteShellCommand("")
	if result != expected || err == nil {
		t.Errorf("Expected empty result and non-nil error, but got result: %s, error: %v", result, err)
	}
}

func TestExecuteShellCommand_ValidCommand(t *testing.T) {
	expectedOutput := "Hello, World!"
	cmd := "echo"
	arg := "Hello, World!"

	output, err := ExecuteShellCommand(cmd, arg)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if output != (expectedOutput + "\n") {
		t.Errorf("Expected: %s, but got: %s", expectedOutput, output)
	}
}

func TestExecuteShellCommand_InvalidCommand(t *testing.T) {
	cmd := "invalidcommand"

	_, err := ExecuteShellCommand(cmd)
	if err == nil {
		t.Error("Expected an error for invalid command, but got nil")
	}
}

func TestExecuteShellCommand_Git(t *testing.T) {
	//найдём список Хэшей коммитов
	cmd := "git"
	arg := make([]string, 0)
	arg = append(arg, "rev-list")
	arg = append(arg, "--all")
	arg = append(arg, "--max-count=1")

	output, err := ExecuteShellCommand(cmd, arg...)
	if err != nil {
		t.Errorf("TestExecuteShellCommand_Git() error: %v", err)
		return
	}

	if output == "" {
		t.Error("TestExecuteShellCommand_Git() error: Output=''")
		return
	}
	MassHash0 := strings.Split(output, "\n")
	MassHash := make([]string, 0)

	for _, v := range MassHash0 {
		if v == "" {
			continue
		}
		MassHash = append(MassHash, v)
	}

	//найдём версии их Хэшей
	cmd = "git"
	arg = make([]string, 0)
	arg = append(arg, "describe")
	arg = append(arg, "--always")
	arg = append(arg, "--tags")
	arg = append(arg, MassHash...)

	output, err = ExecuteShellCommand(cmd, arg...)
	if err != nil {
		t.Errorf("TestExecuteShellCommand_Git() error: %v", err)
	}

	if output == "" {
		t.Error("TestExecuteShellCommand_Git() error: Output=''")
	}
}

func TestDeleteEndEndline(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Text ends with \\n",
			input:    "example\n",
			expected: "example",
		},
		{
			name:     "Text does not end with \\n",
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
			result := DeleteEndEndline(test.input)
			if result != test.expected {
				t.Errorf("Expected %s, but got %s", test.expected, result)
			}
		})
	}
}

func TestFind_ModifiedTime_Internal(t *testing.T) {
	_, err := Find_Directory_ModifiedTime("vendor")
	if err != nil {
		t.Error(err)
	}
}

func TestTimeMax(t *testing.T) {
	// Testing when the first date is the maximum
	x := time.Date(2023, 03, 13, 0, 0, 0, 0, time.UTC)
	y1 := time.Date(2023, 03, 10, 0, 0, 0, 0, time.UTC)
	y2 := time.Date(2023, 03, 11, 0, 0, 0, 0, time.UTC)
	result := TimeMax(x, y1, y2)
	if result != x {
		t.Error("Expected first date as the maximum, but got: ", result)
	}

	// Testing when the last date is the maximum
	x = time.Date(2023, 03, 10, 0, 0, 0, 0, time.UTC)
	y1 = time.Date(2023, 03, 11, 0, 0, 0, 0, time.UTC)
	y2 = time.Date(2023, 03, 12, 0, 0, 0, 0, time.UTC)
	result = TimeMax(x, y1, y2)
	if result != y2 {
		t.Error("Expected last date as the maximum, but got: ", result)
	}

	// Testing when the middle date is the maximum
	x = time.Date(2023, 03, 10, 0, 0, 0, 0, time.UTC)
	y1 = time.Date(2023, 03, 12, 0, 0, 0, 0, time.UTC)
	y2 = time.Date(2023, 03, 11, 0, 0, 0, 0, time.UTC)
	result = TimeMax(x, y1, y2)
	if result != y1 {
		t.Error("Expected middle date as the maximum, but got: ", result)
	}

	// Testing when all dates are the same
	x = time.Date(2023, 03, 10, 0, 0, 0, 0, time.UTC)
	y1 = time.Date(2023, 03, 10, 0, 0, 0, 0, time.UTC)
	y2 = time.Date(2023, 03, 10, 0, 0, 0, 0, time.UTC)
	result = TimeMax(x, y1, y2)
	if result != x {
		t.Error("Expected all dates the same, should return the first date, but got: ", result)
	}
}

func TestTimeMin(t *testing.T) {
	// Case 1: Test when x is the only time provided
	x := time.Date(2023, 3, 14, 12, 0, 0, 0, time.UTC)
	result := TimeMin(x)
	if result != x {
		t.Errorf("Expected %v, but got %v", x, result)
	}

	// Case 2: Test when y is the only time provided
	y := time.Date(2023, 3, 15, 12, 0, 0, 0, time.UTC)
	result = TimeMin(y)
	if result != y {
		t.Errorf("Expected %v, but got %v", y, result)
	}

	// Case 3: Test when multiple times are provided and the minimum is in the middle
	x = time.Date(2023, 3, 14, 12, 0, 0, 0, time.UTC)
	y = time.Date(2023, 3, 15, 12, 0, 0, 0, time.UTC)
	z := time.Date(2023, 3, 13, 12, 0, 0, 0, time.UTC)
	result = TimeMin(x, y, z)
	if result != z {
		t.Errorf("Expected %v, but got %v", z, result)
	}

	// Case 4: Test when the minimum time is the first one provided
	x = time.Date(2023, 3, 13, 12, 0, 0, 0, time.UTC)
	y = time.Date(2023, 3, 15, 12, 0, 0, 0, time.UTC)
	z = time.Date(2023, 3, 14, 12, 0, 0, 0, time.UTC)
	result = TimeMin(x, y, z)
	if result != x {
		t.Errorf("Expected %v, but got %v", x, result)
	}

	// Case 5: Test when the minimum time is the last one provided
	x = time.Date(2023, 3, 15, 12, 0, 0, 0, time.UTC)
	y = time.Date(2023, 3, 14, 12, 0, 0, 0, time.UTC)
	z = time.Date(2023, 3, 13, 12, 0, 0, 0, time.UTC)
	result = TimeMin(x, y, z)
	if result != z {
		t.Errorf("Expected %v, but got %v", z, result)
	}
}

func TestFind_Repository_Code_ModifiedTime(t *testing.T) {

	_, err := Find_Repository_Code_ModifiedTime()
	if err != nil {
		t.Error(err)
	}
}

func TestShow_Repository_Code_ModifiedTime(t *testing.T) {
	Show_Repository_Code_ModifiedTime()
}

func TestShow_Version(t *testing.T) {
	Show_Version("1")
}

func TestMassFrom_MapString(t *testing.T) {
	// Test case for sorting a map by column names
	t.Run("Sort map by keys", func(t *testing.T) {
		input := map[string]int{"c": 3, "a": 1, "b": 2}
		expected := []int{1, 2, 3}
		result := MassFrom_MapString(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	// Test case for handling a key not found in the map
	t.Run("Key not found", func(t *testing.T) {
		input := map[string]string{"a": "apple", "b": "banana"}
		_ = MassFrom_MapString(input) // This call should print an error message for the missing key
		// You can capture the output of fmt.Printf using https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	})
}

func TestMassFrom_Map(t *testing.T) {
	// Test case for sorting a map by column names
	t.Run("Sort map by keys", func(t *testing.T) {
		input := map[string]int{"c": 3, "a": 1, "b": 2}
		expected := []int{1, 2, 3}
		result := MassFrom_Map(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	// Test case for handling a key not found in the map
	t.Run("Key not found", func(t *testing.T) {
		input := map[string]string{"a": "apple", "b": "banana"}
		_ = MassFrom_Map(input) // This call should print an error message for the missing key
		// You can capture the output of fmt.Printf using https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	})
}

func TestMassFrom_Map_DESC(t *testing.T) {
	// Test case for sorting a map by column names
	t.Run("Sort map by keys", func(t *testing.T) {
		input := map[string]int{"c": 3, "a": 1, "b": 2}
		expected := []int{3, 2, 1}
		result := MassFrom_Map_DESC(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	// Test case for handling a key not found in the map
	t.Run("Key not found", func(t *testing.T) {
		input := map[string]string{"a": "apple", "b": "banana"}
		_ = MassFrom_Map(input) // This call should print an error message for the missing key
		// You can capture the output of fmt.Printf using https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	})
}

func TestStringDateTime(t *testing.T) {
	// Test case 1: Testing the output format for a specific date and time
	expectedResult1 := "31.12.2022 23:59:59"
	testTime1 := time.Date(2022, time.December, 31, 23, 59, 59, 0, time.UTC)
	result1 := StringDateTime(testTime1)
	if result1 != expectedResult1 {
		t.Errorf("Test case 1 failed. Expected: %s, Got: %s", expectedResult1, result1)
	}

	// Test case 2: Testing the output format for a different date and time
	expectedResult2 := "01.01.2023 00:00:00"
	testTime2 := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	result2 := StringDateTime(testTime2)
	if result2 != expectedResult2 {
		t.Errorf("Test case 2 failed. Expected: %s, Got: %s", expectedResult2, result2)
	}
}

func TestSubstring(t *testing.T) {
	tests := []struct {
		input       string
		startIndex  int
		length      int
		expected    string
		description string
	}{
		{"hello", 0, 2, "he", "Getting the first 2 characters"},
		{"world", 2, 3, "rld", "Getting characters from index 2 to 4"},
		{"test", 5, 2, "", "Start index beyond the length of the string"},
		{"", 0, 2, "", "Empty input string"},
	}

	for _, test := range tests {
		result := Substring(test.input, test.startIndex, test.length)
		if result != test.expected {
			t.Errorf("Test case failed: %s. Expected: %s, but got: %s", test.description, test.expected, result)
		}
	}
}

func TestIntNot0(t *testing.T) {
	// Testing when all input integers are 0
	result1 := IntNot0(0, 0, 0)
	if result1 != 0 {
		t.Errorf("Expected 0, but got %d", result1)
	}

	// Testing when some input integers are 0 and some are non-zero
	result2 := IntNot0(0, 5, 0, 10)
	if result2 != 5 {
		t.Errorf("Expected 5, but got %d", result2)
	}

	// Testing when all input integers are non-zero
	result3 := IntNot0(3, 7, 2)
	if result3 != 3 {
		t.Errorf("Expected 3, but got %d", result3)
	}
}

func TestInsertTextFrom_EndOfString(t *testing.T) {
	expected := "Hello, World!"
	result := InsertTextFrom("Hello, ", "World!", 7)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestInsertTextFrom_BeginningOfString(t *testing.T) {
	expected := "Hi, there"
	result := InsertTextFrom(", there", "Hi", 0)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestInsertTextFrom_MiddleOfString(t *testing.T) {
	expected := "The quick brown fox jumps over the lazy dog."
	result := InsertTextFrom("The quick fox jumps over the lazy dog.", "brown ", 10)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestDate_from_TimestampReference(t *testing.T) {
	var Timestamp1 *timestamppb.Timestamp

	result := Date_from_TimestampReference(Timestamp1)
	if result.IsZero() == false {
		t.Errorf("Expected time.Time{}, but got %s", result)
	}
}

func TestSetFieldValue(t *testing.T) {
	type Struct1 struct {
		Field1 string
	}

	Struct := &Struct1{}
	SetFieldValue(Struct, "Field1", "Test")

	if Struct.Field1 != "Test" {
		t.Errorf("Expected 'Test', but got %s", Struct.Field1)
	}
}

func TestSetFieldValue2(t *testing.T) {
	type Struct1 struct {
		Field1 int
	}

	Struct := &Struct1{}
	SetFieldValue(Struct, "Field1", "1")

	if Struct.Field1 != 1 {
		t.Errorf("Expected 'Test', but got %v", Struct.Field1)
	}
}
