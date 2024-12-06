package calendar

import (
	"github.com/dromara/carbon/v2"
	"time"
)

// FindPreviousWorkDay - возвращает дату начала предыдущего рабочего(!) дня
// доделать БД Postgres Календарь
func FindPreviousWorkDay(DateNow time.Time) time.Time {
	var Otvet time.Time

	CarbonDate := carbon.CreateFromStdTime(DateNow).StartOfDay()

	Weekday := int(DateNow.Weekday())
	switch Weekday {
	case 0: //воскресенье
		CarbonDate = CarbonDate.AddDays(-2)
	case 1: //понедельник
		CarbonDate = CarbonDate.AddDays(-3)
	default:
		CarbonDate = CarbonDate.AddDays(-1)
	}

	Otvet = CarbonDate.StdTime()

	return Otvet
}

//// FindPreviousWorkDay - возвращает дату начала предыдущего рабочего(!) дня
//// доделать БД Postgres Календарь
//func FindPreviousWorkDay(DateNow time.Time) time.Time {
//	var Otvet time.Time
//
//	//DateNow := time.Now()
//	Otvet = carbon.CreateFromStdTime(DateNow).StartOfDay().StdTime()
//
//	Weekday := int(DateNow.Weekday())
//	switch Weekday {
//	case 0: //воскресенье
//		Otvet = carbon.CreateFromStdTime(Otvet).AddDays(-2).StdTime()
//	case 1: //понедельник
//		Otvet = carbon.CreateFromStdTime(Otvet).AddDays(-3).StdTime()
//	default:
//		Otvet = carbon.CreateFromStdTime(Otvet).AddDays(-1).StdTime()
//	}
//
//	return Otvet
//}

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
