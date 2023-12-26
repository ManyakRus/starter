package nats_connect

import (
	"context"
	"github.com/ManyakRus/starter/config_main"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"testing"
	"time"

	//"github.com/ManyakRus/starter/common/v0/logger"
	"github.com/ManyakRus/starter/stopapp"
)

func TestConnect_err(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	err := Connect_err()
	if err != nil {
		t.Error("nats_connect.TestConnect_err() error: ", err)
	}
	CloseConnection()
}

func TestCloseConnection(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	Connect()
	CloseConnection()
}

func TestStartNats(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
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
	config_main.LoadEnv()
	Connect()
	defer CloseConnection()
}

func TestSendMessageCtx(t *testing.T) {

	config_main.LoadEnv()
	Connect()
	defer CloseConnection()

	ctxMain := contextmain.GetContext()
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()

	Subject := "TESTING"
	Data := []byte("testing")
	err := SendMessageCtx(ctx, Subject, Data)
	if err != nil {
		t.Error("TestSendMessageCtx() error: ", err)
	}

}
