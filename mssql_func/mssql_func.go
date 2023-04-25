package postgres_func

import "time"

func StringSQLTime(time1 time.Time) string {
	Otvet := ""

	//Otvet = "'" + time1.Format(time.RFC3339Nano) + "'"
	Otvet = "'" + time1.Format("20060102 15:04:05") + "'"

	return Otvet
}
