package sync_exchange_connect

import (
	"github.com/ManyakRus/starter/config_main"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/pkg/version"
	"github.com/ManyakRus/starter/stopapp"
	"testing"
)

var SERVICE_NAME_TEST = "test_nikitin"

func TestConnect(t *testing.T) {
	config_main.LoadEnv()
	Connect(SERVICE_NAME_TEST, version.Version)
	defer CloseConnection()

	micro.Pause(100)

}

func TestStartNats(t *testing.T) {
	config_main.LoadEnv()
	Start(SERVICE_NAME_TEST, version.Version)
	defer CloseConnection()

	micro.Pause(100)

	contextmain.CancelContext()
	contextmain.GetNewContext()
}

func TestCloseConnection(t *testing.T) {
	config_main.LoadEnv()
	Connect(SERVICE_NAME_TEST, version.Version)
	defer CloseConnection()
}

func TestWaitStop(t *testing.T) {

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	_ = contextmain.GetContext()
	contextmain.CancelContext()

	stopapp.GetWaitGroup_Main().Wait()

	contextmain.GetNewContext()
}

func TestPprofNats1(t *testing.T) {
	config_main.LoadEnvTest()
	Connect(SERVICE_NAME_TEST, version.Version)
	defer CloseConnection()

	PprofNats1()
}
