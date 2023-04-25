package ping

import (
	"testing"
	"time"
)

func TestPingPort(t *testing.T) {

	IP := "10.1.9.151"
	Port := "26500"
	TimeStart := time.Now()

	err := Ping_err(IP, Port)

	t.Log("Прошло время: ", time.Since(TimeStart))

	if err != nil {
		t.Error("PingPort() error: ", err)
	} else {
		t.Log("PingPort() OK.")
	}

}
