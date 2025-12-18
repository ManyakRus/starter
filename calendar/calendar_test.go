package calendar

import (
	"github.com/ManyakRus/starter/constants_starter"
	"testing"
	"time"
)

func TestFindPreviousWorkDay(t *testing.T) {

	Date := time.Date(2023, 03, 13, 0, 0, 0, 0, constants_starter.Loc)
	Otvet := FindPreviousWorkDay(Date)
	if Otvet != time.Date(2023, 03, 10, 0, 0, 0, 0, constants_starter.Loc) {
		t.Error("TestFindPreviousWorkDay error")
	}

}

func TestIsWorkDay(t *testing.T) {
	Date := time.Date(2023, 03, 19, 0, 0, 0, 0, constants_starter.Loc)
	Otvet := IsWorkDay(Date)
	if Otvet != false {
		t.Error("TestIsWorkDay error")
	}
}

func TestHoursMinutesSeconds_UnmarshalByte(t *testing.T) {

	Otvet := HoursMinutesSeconds{}
	Otvet.UnmarshalJSON([]byte("01:02:03"))
	if Otvet.Hours != 1 || Otvet.Minutes != 2 || Otvet.Seconds != 3 {
		t.Error("TestHoursMinutesSeconds_UnmarshalByte error")
	}
}

func TestDiff_dates(t *testing.T) {

	Date1 := time.Date(2023, 03, 19, 0, 0, 0, 0, constants_starter.Loc)
	Date2 := time.Date(2023, 03, 20, 0, 0, 0, 0, constants_starter.Loc)
	_, _, days, _, _, _ := Diff_dates(Date1, Date2)
	if days != 1 {
		t.Error("TestDiff_dates error")
	}
}
