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
func NullString_DefaultNull[T ~string](s T) sql.NullString {
	Otvet := sql.NullString{}
	Otvet.String = string(s)
	Otvet.Valid = true

	if s == "" {
		Otvet.Valid = false
	}
	return Otvet
}

// NullInt64_DefaultNull - преобразует значение в sql.NullInt64, если пусто то Valid = false
func NullInt64_DefaultNull[T ~int64](Value T) sql.NullInt64 {
	Otvet := sql.NullInt64{}
	Otvet.Int64 = int64(Value)
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}

// NullInt32_DefaultNull - преобразует значение в sql.NullInt32, если пусто то Valid = false
func NullInt32_DefaultNull[T ~int32](Value T) sql.NullInt32 {
	Otvet := sql.NullInt32{}
	Otvet.Int32 = int32(Value)
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}

// NullInt16_DefaultNull - преобразует значение в sql.NullInt16, если пусто то Valid = false
func NullInt16_DefaultNull[T ~int16](Value T) sql.NullInt16 {
	Otvet := sql.NullInt16{}
	Otvet.Int16 = int16(Value)
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}

// NullTime_DefaultNull - преобразует значение в sql.NullTime, если пусто то Valid = false
func NullTime_DefaultNull[T time.Time | int64](Value T) sql.NullTime {
	Otvet := sql.NullTime{}

	switch any(Value).(type) {
	case time.Time:
		v := any(Value).(time.Time)
		Otvet.Time = v
		Otvet.Valid = true

		if v.IsZero() == true {
			Otvet.Valid = false
		}
	case int64:
		v := any(Value).(int64)
		Time1 := time.Unix(v, 0)
		Otvet.Time = Time1
		Otvet.Valid = true

		if Time1.IsZero() == true {
			Otvet.Valid = false
		}
	}

	return Otvet
}

// NullByte_DefaultNull - преобразует значение в sql.NullByte, если пусто то Valid = false
func NullByte_DefaultNull[T ~byte](Value T) sql.NullByte {
	Otvet := sql.NullByte{}
	Otvet.Byte = byte(Value)
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}

// NullFloat64_DefaultNull - преобразует значение в sql.NullFloat64, если пусто то Valid = false
func NullFloat64_DefaultNull[T ~float64](Value T) sql.NullFloat64 {
	Otvet := sql.NullFloat64{}
	Otvet.Float64 = float64(Value)
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}

// NullFloat32_DefaultNull - преобразует значение в sql.NullFloat64, если пусто то Valid = false
func NullFloat32_DefaultNull[T ~float32](Value T) sql.NullFloat64 {
	Otvet := sql.NullFloat64{}
	Otvet.Float64 = float64(Value)
	Otvet.Valid = true

	if Value == 0 {
		Otvet.Valid = false
	}

	return Otvet
}

// NullBool_DefaultNull - преобразует значение в sql.NullBool, если пусто то Valid = false
func NullBool_DefaultNull[T ~bool](Value T) sql.NullBool {
	Otvet := sql.NullBool{}
	Otvet.Bool = bool(Value)
	Otvet.Valid = true

	if Value == false {
		Otvet.Valid = false
	}

	return Otvet
}
