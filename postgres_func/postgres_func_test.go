package postgres_func

import (
	"testing"
	"time"
)

func TestStringSQLTime(t *testing.T) {

	var loc = time.Local
	time1 := time.Date(2023, 1, 1, 0, 0, 0, 0, loc)
	Otvet := StringSQLTime(time1)
	if Otvet != "'2023-01-01T00:00:00+03:00'" {
		t.Error("postgres_func_test.TestStringSQLTime() error")
	}

}
