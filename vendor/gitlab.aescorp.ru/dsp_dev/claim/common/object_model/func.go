package object_model

import (
	"fmt"
	"time"
)

func russianDate(date time.Time, long bool) string {
	months := make(map[int]string, 0)
	if long {
		months[1] = "январь"
		months[2] = "февраль"
		months[3] = "март"
		months[4] = "апрель"
		months[5] = "май"
		months[6] = "июнь"
		months[7] = "июль"
		months[8] = "август"
		months[9] = "сентябрь"
		months[10] = "октябрь"
		months[11] = "ноябрь"
		months[12] = "декабрь"
	} else {
		months[1] = "янв"
		months[2] = "февр"
		months[3] = "март"
		months[4] = "апр"
		months[5] = "май"
		months[6] = "июнь"
		months[7] = "июль"
		months[8] = "авг"
		months[9] = "сент"
		months[10] = "окт"
		months[11] = "нояб"
		months[12] = "дек"
	}

	m := date.Month()
	y := date.Year()

	return fmt.Sprintf("%v %v", months[int(m)], y)
}
