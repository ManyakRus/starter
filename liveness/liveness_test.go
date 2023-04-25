package liveness

import (
	"github.com/manyakrus/starter/contextmain"
	"github.com/manyakrus/starter/fiber_connect"
	"github.com/manyakrus/starter/micro"
	"net/http"
	"testing"
)

func TestStart(t *testing.T) {
	Start()

	micro.Pause(100)

	URL := "http://" + fiber_connect.GetHost() + ":" + fiber_connect.Settings.WEBSERVER_PORT + LIVENESS_URL
	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		t.Error("liveness_test.TestStart() error: Status != 200")
	}

	contextmain.CancelContext()

}
