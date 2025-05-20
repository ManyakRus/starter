package calendar

import (
	"fmt"
	"github.com/ManyakRus/starter/constants"
	"github.com/dromara/carbon/v2"
	"time"
)

// HoursMinutesSeconds - структура для хранения часов, минут и секунд
type HoursMinutesSeconds struct {
	Hours   int
	Minutes int
	Seconds int
}

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

// UnmarshalJSON - преобразует байты время в HoursMinutesSeconds{}
func (d *HoursMinutesSeconds) UnmarshalJSON(b []byte) error {
	str := string(b)
	err := d.UnmarshalString(str)

	return err
}

// UnmarshalString - преобразует строку время в HoursMinutesSeconds{}
func (d *HoursMinutesSeconds) UnmarshalString(str string) error {
	if str != "" && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// parse string
	t, err := time.Parse(constants.LayoutTime, str)
	if err != nil {
		err = fmt.Errorf("invalid time string: %s, error: %w", str, err)
		return err
	}

	d.Hours = t.Hour()
	d.Minutes = t.Minute()
	d.Seconds = t.Second()

	return nil
}

// Diff_dates - вычисляет разницу между двумя датами
func Diff_dates(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
