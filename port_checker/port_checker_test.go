package port_checker

import (
	"testing"
	"time"
)

func TestCheckPort(t *testing.T) {

	IP := "10.1.9.151"
	Port := "26500"
	TimeStart := time.Now()

	err := CheckPort_err(IP, Port)

	t.Log("Прошло время: ", time.Since(TimeStart))

	if err != nil {
		t.Error("TestCheckPort() error: ", err)
	} else {
		t.Log("TestCheckPort() OK.")
	}

}
