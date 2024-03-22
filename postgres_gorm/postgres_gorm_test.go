package postgres_gorm

import (
	"errors"
	"testing"
	"time"

	//log "github.com/sirupsen/logrus"

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

func TestIsClosed(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()

	err := Connect_err()
	if err != nil {
		t.Error("TestIsClosed Connect() error: ", err)
	}

	isClosed := IsClosed()
	if isClosed == true {
		t.Error("TestIsClosed() isClosed = true ")
	}

	err = CloseConnection_err()
	if err != nil {
		t.Error("TestIsClosed() CloseConnection() error: ", err)
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

func TestStartDB(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	StartDB()
	err := CloseConnection_err()
	if err != nil {
		t.Error("db_test.TestStartDB() CloseConnection() error: ", err)
	}
}

func TestConnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	Connect()

	CloseConnection()
}

func TestGetConnection(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	GetConnection()

	CloseConnection()
}

func TestConnect_WithApplicationName_err(t *testing.T) {

	config_main.LoadEnv()
	err := Connect_WithApplicationName_err("starter test")
	if err != nil {
		t.Error("TestConnect_WithApplicationName_err error: ", err)
	}

	//micro.Pause(60 * 1000)

	err = CloseConnection_err()
	if err != nil {
		t.Error("TestConnect_WithApplicationName_err() error: ", err)
	}
}

func TestRawMultipleSQL(t *testing.T) {
	config_main.LoadEnv()
	GetConnection()
	defer CloseConnection()

	TimeStart := time.Now()

	TextSQL := `
drop table if exists temp_TestRawMultipleSQL2; 
CREATE TEMPORARY TABLE temp_TestRawMultipleSQL2 (id int);

insert into temp_TestRawMultipleSQL2
select 1;

SELECT * FROM temp_TestRawMultipleSQL2
`
	//TextSQL := "SELECT 1; SELECT 2"
	tx := RawMultipleSQL(Conn, TextSQL)
	err := tx.Error
	if err != nil {
		t.Error("TestRawMultipleSQL() error: ", err)
	}

	t.Log("Прошло время: ", time.Since(TimeStart))
}

func TestRawMultipleSQL2(t *testing.T) {
	config_main.LoadEnv()
	GetConnection()
	defer CloseConnection()

	TimeStart := time.Now()

	TextSQL := `
drop table if exists temp_TestRawMultipleSQL2; 
CREATE TEMPORARY TABLE temp_TestRawMultipleSQL2 (id int);

insert into temp_TestRawMultipleSQL2
select 1;

SELECT * FROM temp_TestRawMultipleSQL2
`
	tx := RawMultipleSQL(Conn, TextSQL)
	err := tx.Error
	if err != nil {
		t.Error("TestRawMultipleSQL() error: ", err)
	}

	if tx.RowsAffected != 1 {
		t.Error("TestRawMultipleSQL() RowsAffected = ", tx.RowsAffected)
	}

	t.Log("Прошло время: ", time.Since(TimeStart))
}

func TestRawMultipleSQL3(t *testing.T) {
	config_main.LoadEnv()
	GetConnection()
	defer CloseConnection()

	TimeStart := time.Now()

	TextSQL := `
drop table if exists temp_TestRawMultipleSQL2; 
CREATE TEMPORARY TABLE temp_TestRawMultipleSQL2 (id int);

insert into temp_TestRawMultipleSQL2
select 1;

SELECT * FROM temp_TestRawMultipleSQL2
`
	f := func(t *testing.T) {
		tx := RawMultipleSQL(Conn, TextSQL)
		err := tx.Error
		if err != nil {
			t.Error("TestRawMultipleSQL3() error: ", err)
		}

		if tx.RowsAffected != 1 {
			t.Error("TestRawMultipleSQL3() RowsAffected = ", tx.RowsAffected)
		}

	}

	//запустим 100 потоков
	for i := 0; i < 100; i++ {
		stopapp.GetWaitGroup_Main().Add(1)
		go f(t)
		stopapp.GetWaitGroup_Main().Done()
	}

	stopapp.GetWaitGroup_Main().Wait()

	t.Log("Прошло время: ", time.Since(TimeStart))
}
