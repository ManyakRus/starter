package postgres_stek

import (
	"errors"
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/db/tables/table_connections"
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/object_model/entities/connections"
	"testing"

	//log "github.com/sirupsen/logrus"

	"github.com/ManyakRus/starter/config_main"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"

	//	logger "github.com/ManyakRus/starter/common/v0/logger"
	"github.com/ManyakRus/starter/stopapp"
)

// CONNECTION - объект Соединение, настроенный
var CONNECTION = connections.Connection{Table_Connection: table_connections.Table_Connection{ID: constants_starter.CONNECTION_ID, BranchID: constants_starter.BRANCH_ID, IsLegal: true, Server: "10.1.9.153", Port: "5432", DbName: "kol_atom_ul_uni", DbScheme: "stack", Login: "", Password: ""}}

func TestConnect_err(t *testing.T) {
	//Connect_Panic()

	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	err := Connect_err(CONNECTION)
	if err != nil {
		t.Error("TestConnect error: ", err)
	}

	err = CloseConnection_err(CONNECTION)
	if err != nil {
		t.Error("TestConnect() error: ", err)
	}
}

func TestIsClosed(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()

	err := Connect_err(CONNECTION)
	if err != nil {
		t.Error("TestIsClosed Connect() error: ", err)
	}

	isClosed := IsClosed(CONNECTION)
	if isClosed == true {
		t.Error("TestIsClosed() isClosed = true ")
	}

	err = CloseConnection_err(CONNECTION)
	if err != nil {
		t.Error("TestIsClosed() CloseConnection() error: ", err)
	}

}

func TestReconnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	err := Connect_err(CONNECTION)
	if err != nil {
		t.Error("TestIsClosed Connect() error: ", err)
	}

	//ctx := context.Background()
	Reconnect(CONNECTION, errors.New(""))

	err = CloseConnection_err(CONNECTION)
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
	config_main.LoadEnv()
	StartDB(CONNECTION)
	err := CloseConnection_err(CONNECTION)
	if err != nil {
		t.Error("db_test.TestStartDB() CloseConnection() error: ", err)
	}
}

func TestConnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	Connect(CONNECTION)

	CloseConnection(CONNECTION)
}

func TestGetConnection(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	GetConnection(CONNECTION)

	CloseConnection(CONNECTION)
}
