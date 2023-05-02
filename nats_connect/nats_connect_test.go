package nats_connect

import (
	"github.com/ManyakRus/starter/config"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"testing"
	//"github.com/ManyakRus/starter/common/v0/logger"
	"github.com/ManyakRus/starter/stopapp"
)

func TestConnect_err(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	err := Connect_err()
	if err != nil {
		t.Error("nats_connect.TestConnect_err() error: ", err)
	}
	CloseConnection()
}

func TestCloseConnection(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	Connect()
	CloseConnection()
}

func TestStartNats(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	StartNats()
	micro.Pause(20)

	_ = contextmain.GetContext()
	contextmain.CancelContext()

	stopapp.GetWaitGroup_Main().Wait()

	contextmain.GetNewContext()
}

func TestWaitStop(t *testing.T) {

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	_ = contextmain.GetContext()
	contextmain.CancelContext()

	stopapp.GetWaitGroup_Main().Wait()

	contextmain.GetNewContext()
}

func TestConnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	Connect()
	CloseConnection()
}
