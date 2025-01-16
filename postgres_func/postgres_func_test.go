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
