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

func TestReplaceSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		schema   string
		expected string
	}{
		{
			name:     "No schema",
			input:    "SELECT * FROM public.users",
			schema:   "",
			expected: "SELECT * FROM public.users",
		},
		{
			name:     "Schema with tabs and newlines",
			input:    "\tSELECT * FROM public.users\n",
			schema:   "myschema",
			expected: "\tSELECT * FROM myschema.users\n",
		},
		{
			name:     "Schema with spaces",
			input:    "SELECT * FROM public.users ",
			schema:   "myschema",
			expected: "SELECT * FROM myschema.users ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Settings.DB_SCHEMA = tt.schema
			got := ReplaceSchema(tt.input)
			if got != tt.expected {
				t.Errorf("ReplaceSchema() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestReplaceTableNamesToUnique1(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		uuid     string
		expected string
	}{
		{
			name:     "No public schema",
			input:    "SELECT * FROM TableName",
			uuid:     "1234567890",
			expected: "SELECT * FROM TableName",
		},
		{
			name:     "Public schema with no spaces",
			input:    "SELECT * FROM public.TableName",
			uuid:     "1234567890",
			expected: "SELECT * FROM public.TableName",
		},
		{
			name:     "Public schema with spaces",
			input:    "SELECT * FROM public.Table Name",
			uuid:     "1234567890",
			expected: "SELECT * FROM public.Table_1234567890 Name",
		},
		{
			name:     "Public schema with tabs",
			input:    "SELECT * FROM public.\tTableName\t",
			uuid:     "1234567890",
			expected: "SELECT * FROM public.\tTableName_1234567890\t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceTableNamesToUnique1(tt.input, tt.uuid)
			if got != tt.expected {
				t.Errorf("ReplaceTableNamesToUnique1() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
