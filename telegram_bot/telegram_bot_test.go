package telegram_bot

import (
	"errors"
	"testing"
	"time"

	"github.com/ManyakRus/starter/config_main"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"

	//	logger "github.com/ManyakRus/starter/common/v0/logger"
	"github.com/ManyakRus/starter/stopapp"
)

func TestConnect_err(t *testing.T) {
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

func TestReconnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	err := Connect_err()
	if err != nil {
		t.Error("TestIsClosed Connect() error: ", err)
	}

	//ctx := context.Background()
	Reconnect(errors.New(""))

	err = CloseConnection_err()
	if err != nil {
		t.Error("TestReconnect() CloseConnection() error: ", err)
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

func TestConnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	Connect()
	defer CloseConnection()
}

func TestGetConnection(t *testing.T) {
	config_main.LoadEnv()
	GetConnection()
	defer CloseConnection()
}

func TestSendMessageChatID(t *testing.T) {
	t.SkipNow()

	config_main.LoadEnv()
	GetConnection()
	defer CloseConnection()

	s := Settings.TELEGRAM_CHAT_ID_TEST
	i, err := micro.Int64FromString(s)
	if err != nil {
		t.Errorf("TELEGRAM_CHAT_ID_TEST: %s, error: %v", s, err)
	}

	Text := "test " + time.Now().String()
	_, err = SendMessageChatID(i, Text)
	if err != nil {
		t.Error("TestSendMessage() error: ", err)
	}

}

func TestSendMessage(t *testing.T) {
	config_main.LoadEnv()
	GetConnection()
	defer CloseConnection()

	Text := "test " + time.Now().String()
	ID := Settings.TELEGRAM_CHAT_ID_TEST
	_, err := SendMessage(ID, Text)
	if err != nil {
		t.Error("TestSendMessage() error: ", err)
	}

}
