package postgres_func

import (
	"testing"
	"time"
)

func TestStringSQLTime(t *testing.T) {

	var loc = time.Local
	time1 := time.Date(2023, 1, 1, 0, 0, 0, 0, loc)
	Otvet := StringSQLTime(time1)
	if Otvet != "'2023-01-01T00:00:00+03:00'" {
		t.Error("postgres_func_test.TestStringSQLTime() error")
	}

}

func TestStringSQLTime_WithoutTimeZone(t *testing.T) {
	var loc = time.Local
	time1 := time.Date(2023, 1, 1, 0, 0, 0, 0, loc)
	Otvet := StringSQLTime_WithoutTimeZone(time1)
	if Otvet != "'2023-01-01T00:00:00+00:00'" {
		t.Error("postgres_func_test.TestStringSQLTime() error")
	}
}

func TestReplaceSchemaName(t *testing.T) {

	TextSQL := "select * from schema1.table1"
	SchemaNameFrom := "schema1"
	SchemaNameTo := "schema2"
	Otvet := ReplaceSchemaName(TextSQL, SchemaNameFrom, SchemaNameTo)
	if Otvet != "select * from schema2.table1" {
		t.Error("postgres_func_test.TestReplaceSchemaName() error")
	}
}

func TestNullTime_DefaultNull(t *testing.T) {

	//дата
	Date1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
	Otvet1 := NullTime_DefaultNull(Date1)
	if Otvet1.Time != Date1 || Otvet1.Valid == false {
		t.Error("postgres_func_test.TestNullTime_DefaultNull() error")
	}

	//дата alias
	type DateAlias = time.Time
	Date2 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
	DateAlias2 := DateAlias(Date2)
	Otvet2 := NullTime_DefaultNull(DateAlias2)
	if Otvet2.Time != Date1 || Otvet2.Valid == false {
		t.Error("postgres_func_test.TestNullTime_DefaultNull() error")
	}

}
