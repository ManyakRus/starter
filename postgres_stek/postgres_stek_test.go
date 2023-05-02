package postgres_stek

import (
	"errors"
	"github.com/ManyakRus/starter/constants"
	"testing"

	//log "github.com/sirupsen/logrus"

	"github.com/ManyakRus/starter/config"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"

	//	logger "github.com/ManyakRus/starter/common/v0/logger"
	"github.com/ManyakRus/starter/stopapp"
)

func TestConnect_err(t *testing.T) {
	//Connect_Panic()

	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	err := Connect_err(constants.CONNECTION)
	if err != nil {
		t.Error("TestConnect error: ", err)
	}

	err = CloseConnection_err(constants.CONNECTION)
	if err != nil {
		t.Error("TestConnect() error: ", err)
	}
}

func TestIsClosed(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()

	err := Connect_err(constants.CONNECTION)
	if err != nil {
		t.Error("TestIsClosed Connect() error: ", err)
	}

	isClosed := IsClosed(constants.CONNECTION)
	if isClosed == true {
		t.Error("TestIsClosed() isClosed = true ")
	}

	err = CloseConnection_err(constants.CONNECTION)
	if err != nil {
		t.Error("TestIsClosed() CloseConnection() error: ", err)
	}

}

func TestReconnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	err := Connect_err(constants.CONNECTION)
	if err != nil {
		t.Error("TestIsClosed Connect() error: ", err)
	}

	//ctx := context.Background()
	Reconnect(constants.CONNECTION, errors.New(""))

	err = CloseConnection_err(constants.CONNECTION)
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

func TestStartDB(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	StartDB(constants.CONNECTION)
	err := CloseConnection_err(constants.CONNECTION)
	if err != nil {
		t.Error("db_test.TestStartDB() CloseConnection() error: ", err)
	}
}

func TestConnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	Connect(constants.CONNECTION)

	CloseConnection(constants.CONNECTION)
}

func TestGetConnection(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	GetConnection(constants.CONNECTION)

	CloseConnection(constants.CONNECTION)
}
