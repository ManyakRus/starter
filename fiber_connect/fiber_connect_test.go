package fiber_connect

import (
	"github.com/manyakrus/starter/contextmain"
	"github.com/manyakrus/starter/micro"
	"testing"
)

func TestConnect(t *testing.T) {

	FillSettings()
	Connect()
	CloseConnection()

}

func TestStart(t *testing.T) {
	Start()

	micro.Pause(200)
	contextmain.CancelContext()
	contextmain.GetNewContext()
}

func TestGetHost(t *testing.T) {
	Otvet := GetHost()
	if Otvet == "" {
		t.Error("fiber_connect_test.TestGetHost() error: Otvet=''")
	}
}
