package chatgpt_proxy

import (
	"testing"

	//log "github.com/sirupsen/logrus"

	"github.com/ManyakRus/starter/config_main"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"

	//	logger "github.com/ManyakRus/starter/common/v0/logger"
	"github.com/ManyakRus/starter/stopapp"
)

func TestConnect_err(t *testing.T) {
	//Connect_Panic()

	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	err := Connect_err()
	if err != nil {
		t.Error("TestConnect error: ", err)
	}

	err = CloseConnection_err()
	if err != nil {
		t.Error("TestConnect() error: ", err)
	}
}

func TestWaitStop(t *testing.T) {
	stopapp.StartWaitStop()

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	micro.Pause(10)

	//stopapp.SignalInterrupt <- syscall.SIGINT
	contextmain.CancelContext()
}

func TestStartDB(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	Start()
	err := CloseConnection_err()
	if err != nil {
		t.Error("chatgpt_connect_test.TestStart() CloseConnection() error: ", err)
	}
}

func TestConnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	Connect()

	CloseConnection()
}

func TestSendMessage(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()

	Text := "Is ChatGPT enabled now ?"
	Otvet, err := SendMessage(Text, "")
	if err != nil {
		t.Error("chatgpt_connect_test.TestSendMessage() error: ", err)
	}
	t.Log("chatgpt_connect_test.TestSendMessage() Otvet: ", Otvet)
}
