package calendar

import (
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

	Otvet = CarbonDate.ToStdTime()

	return Otvet
}

//// FindPreviousWorkDay - возвращает дату начала предыдущего рабочего(!) дня
//// доделать БД Postgres Календарь
//func FindPreviousWorkDay(DateNow time.Time) time.Time {
//	var Otvet time.Time
//
//	//DateNow := time.Now()
//	Otvet = carbon.CreateFromStdTime(DateNow).StartOfDay().ToStdTime()
//
//	Weekday := int(DateNow.Weekday())
//	switch Weekday {
//	case 0: //воскресенье
//		Otvet = carbon.CreateFromStdTime(Otvet).AddDays(-2).ToStdTime()
//	case 1: //понедельник
//		Otvet = carbon.CreateFromStdTime(Otvet).AddDays(-3).ToStdTime()
//	default:
//		Otvet = carbon.CreateFromStdTime(Otvet).AddDays(-1).ToStdTime()
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
