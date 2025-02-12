package calendar

import (
	"github.com/ManyakRus/starter/constants"
	"testing"
	"time"
)

func TestFindPreviousWorkDay(t *testing.T) {

	Date := time.Date(2023, 03, 13, 0, 0, 0, 0, constants.Loc)
	Otvet := FindPreviousWorkDay(Date)
	if Otvet != time.Date(2023, 03, 10, 0, 0, 0, 0, constants.Loc) {
		t.Error("TestFindPreviousWorkDay error")
	}

}

func TestIsWorkDay(t *testing.T) {
	Date := time.Date(2023, 03, 19, 0, 0, 0, 0, constants.Loc)
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
