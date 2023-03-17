package calendar

import (
	"github.com/golang-module/carbon/v2"
	"time"
)

// FindPreviousWorkDay - возвращает дату начала предыдущего рабочего(!) дня
// доделать БД Postgres Календарь
func FindPreviousWorkDay(DateNow time.Time) time.Time {
	var Otvet time.Time

	//DateNow := time.Now()
	Otvet = carbon.Time2Carbon(DateNow).StartOfDay().Carbon2Time()

	Weekday := int(DateNow.Weekday())
	switch Weekday {
	case 0: //воскресенье
		Otvet = carbon.Time2Carbon(Otvet).AddDays(-2).Carbon2Time()
	case 1: //понедельник
		Otvet = carbon.Time2Carbon(Otvet).AddDays(-3).Carbon2Time()
	default:
		Otvet = carbon.Time2Carbon(Otvet).AddDays(-1).Carbon2Time()
	}

	return Otvet
}

func IsWorkDay(Date time.Time) bool {
	Otvet := false

	Weekday := int(Date.Weekday())
	switch Weekday {
	case 6, 0: //суббота+воскресенье
	default:
		Otvet = true
	}

	return Otvet
}
