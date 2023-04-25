package postgres_func

import "time"

func StringSQLTime(time1 time.Time) string {
	Otvet := ""

	Otvet = "'" + time1.Format(time.RFC3339Nano) + "'"

	return Otvet
}

func StringSQLTime_WithoutTimeZone(time1 time.Time) string {
	Otvet := ""

	format := "2006-01-02T15:04:05.999999999+00:00"
	Otvet = "'" + time1.Format(format) + "'"

	return Otvet
}
