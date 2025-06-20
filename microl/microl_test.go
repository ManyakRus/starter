package microl

import (
	"github.com/ManyakRus/starter/constants"
	"os"
	"testing"
	"time"
)

func TestSet_FieldFromEnv_String(t *testing.T) {
	type Struct1 struct {
		TestSet_FieldFromEnv string
	}

	Name := "TestSet_FieldFromEnv"
	os.Setenv(Name, Name)

	Struct := Struct1{}
	Set_FieldFromEnv_String(&Struct, Name, true)

	if Struct.TestSet_FieldFromEnv != Name {
		t.Error("Set_FieldFromEnv_String() error")
	}

}

func TestSet_FieldFromEnv_Int(t *testing.T) {
	type Struct1 struct {
		TestSet_FieldFromEnv int
	}

	Name := "TestSet_FieldFromEnv"
	os.Setenv(Name, "1")

	Struct := Struct1{}
	Set_FieldFromEnv_Int(&Struct, Name, true)

	if Struct.TestSet_FieldFromEnv != 1 {
		t.Error("TestSet_FieldFromEnv_Int() error")
	}

}

func TestSet_FieldFromEnv_Int64(t *testing.T) {
	type Struct1 struct {
		TestSet_FieldFromEnv int64
	}

	Name := "TestSet_FieldFromEnv"
	os.Setenv(Name, "1")

	Struct := Struct1{}
	Set_FieldFromEnv_Int64(&Struct, Name, true)

	if Struct.TestSet_FieldFromEnv != 1 {
		t.Error("TestSet_FieldFromEnv_Int64() error")
	}

}
func TestSet_FieldFromEnv_Int32(t *testing.T) {
	type Struct1 struct {
		TestSet_FieldFromEnv int32
	}

	Name := "TestSet_FieldFromEnv"
	os.Setenv(Name, "1")

	Struct := Struct1{}
	Set_FieldFromEnv_Int32(&Struct, Name, true)

	if Struct.TestSet_FieldFromEnv != 1 {
		t.Error("TestSet_FieldFromEnv_Int32() error")
	}

}

func TestSet_FieldFromEnv_Time(t *testing.T) {
	type Struct1 struct {
		TestSet_FieldFromEnv time.Time
	}

	sTime := "02.01.2000 00:00:00"
	Time1, err := time.Parse(constants.LayoutDateTimeRus, sTime)
	if err != nil {
		t.Error("TestSet_FieldFromEnv_Time() error")
	}
	//sTime := Time1.GoString()
	Name := "TestSet_FieldFromEnv"
	os.Setenv(Name, sTime)

	Struct := Struct1{}
	Set_FieldFromEnv_Time(&Struct, Name, true)

	if Struct.TestSet_FieldFromEnv != Time1 {
		t.Error("TestSet_FieldFromEnv_Time() error")
	}

}

func TestSet_FieldFromEnv_Date(t *testing.T) {
	type Struct1 struct {
		TestSet_FieldFromEnv time.Time
	}

	sTime := "02.01.2000"
	Time1, err := time.Parse(constants.LayoutDateRus, sTime)
	if err != nil {
		t.Error("TestSet_FieldFromEnv_Date() error")
	}
	//sTime := Time1.GoString()
	Name := "TestSet_FieldFromEnv"
	os.Setenv(Name, sTime)

	Struct := Struct1{}
	Set_FieldFromEnv_Date(&Struct, Name, true)

	if Struct.TestSet_FieldFromEnv != Time1 {
		t.Error("TestSet_FieldFromEnv_Date() error")
	}

}

func TestSet_FieldFromEnv_Bool(t *testing.T) {
	type Struct1 struct {
		TestSet_FieldFromEnv bool
	}

	Name := "TestSet_FieldFromEnv"
	os.Setenv(Name, "1")

	Struct := Struct1{}
	Set_FieldFromEnv_Bool(&Struct, Name, true)

	if Struct.TestSet_FieldFromEnv != true {
		t.Error("TestSet_FieldFromEnv_Bool() error")
	}

}

func TestShowTimePassed_FormatText(t *testing.T) {
	ShowTimePassed_FormatText("start offer Download_and_Save_All(), time passed: %s", time.Now())
}

func TestShow_Stage(t *testing.T) {
	Show_Stage()
}
