package postgres_func

import (
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
