package postgres_func

import (
	"database/sql"
	"strings"
	"time"
)

// StringSQLTime - преобразует время в строку в формате SQL
func StringSQLTime(time1 time.Time) string {
	Otvet := ""

	format := "2006-01-02T15:04:05.999999Z07:00"
	Otvet = "'" + time1.Format(format) + "'"
	//Otvet = "'" + time1.Format(time.RFC3339Nano) + "'"

	return Otvet
}

// StringSQLTime - преобразует время в строку в формате SQL, без часового пояса
func StringSQLTime_WithoutTimeZone(time1 time.Time) string {
	Otvet := ""

	format := "2006-01-02T15:04:05.999999+00:00"
	Otvet = "'" + time1.Format(format) + "'"

	return Otvet
}

// ReplaceSchemaName - заменяет имя схемы в тексте SQL
func ReplaceSchemaName(TextSQL, SchemaNameFrom, SchemaNameTo string) string {
	Otvet := TextSQL

	Otvet = strings.ReplaceAll(Otvet, SchemaNameFrom+".", SchemaNameTo+".")

	return Otvet
}

// NullString_DefaultNull - преобразует строку в sql.NullString, если она пустая то Valid = false
func NullString_DefaultNull(s string) sql.NullString {
	Otvet := sql.NullString{}
	Otvet.String = s
	Otvet.Valid = true

	if s == "" {
		Otvet.Valid = false
	}
	return Otvet
}

// NullInt64_DefaultNull - преобразует значение в sql.NullInt64, если пусто то Valid = false
func NullInt64_DefaultNull(Value int64) sql.NullInt64 {
	Otvet := sql.NullInt64{}
	Otvet.Int64 = Value
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}

// NullInt32_DefaultNull - преобразует значение в sql.NullInt32, если пусто то Valid = false
func NullInt32_DefaultNull(Value int32) sql.NullInt32 {
	Otvet := sql.NullInt32{}
	Otvet.Int32 = Value
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}

// NullInt16_DefaultNull - преобразует значение в sql.NullInt16, если пусто то Valid = false
func NullInt16_DefaultNull(Value int16) sql.NullInt16 {
	Otvet := sql.NullInt16{}
	Otvet.Int16 = Value
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}

// NullTime_DefaultNull - преобразует значение в sql.NullTime, если пусто то Valid = false
func NullTime_DefaultNull(Value time.Time) sql.NullTime {
	Otvet := sql.NullTime{}
	Otvet.Time = Value
	Otvet.Valid = true

	if Value.IsZero() == true {
		Otvet.Valid = false
	}

	return Otvet
}

// NullByte_DefaultNull - преобразует значение в sql.NullByte, если пусто то Valid = false
func NullByte_DefaultNull(Value byte) sql.NullByte {
	Otvet := sql.NullByte{}
	Otvet.Byte = Value
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}

// NullFloat64_DefaultNull - преобразует значение в sql.NullFloat64, если пусто то Valid = false
func NullFloat64_DefaultNull(Value float64) sql.NullFloat64 {
	Otvet := sql.NullFloat64{}
	Otvet.Float64 = Value
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}

// NullNullBool_DefaultNull - преобразует значение в sql.NullBool, если пусто то Valid = false
func NullNullBool_DefaultNull(Value float64) sql.NullBool {
	Otvet := sql.NullBool{}
	Otvet.Bool = Value
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}
