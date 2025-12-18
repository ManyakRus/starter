package postgres_pgxpool2

import (
	"errors"
	"github.com/ManyakRus/starter/microl"
	"golang.org/x/net/context"
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

	waitGroup_Connect.Add(1)
	go WaitStop()

	micro.Pause(10)

	//stopapp.SignalInterrupt <- syscall.SIGINT
	contextmain.CancelContext()
	contextmain.GetNewContext()
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

func TestConnect_WithApplicationName_err(t *testing.T) {

	config_main.LoadEnv()
	err := Connect_WithApplicationName_err("test_starter_postgres_pgx")
	if err != nil {
		t.Error("TestConnect_WithApplicationName_err error: ", err)
	}

	CloseConnection()

}

func TestRawMultipleSQL(t *testing.T) {
	config_main.LoadEnv()
	connection := GetConnection()
	//defer connection.Release()
	defer CloseConnection()

	TimeStart := time.Now()

	TextSQL := `
drop table if exists temp_TestRawMultipleSQL2; 
CREATE TEMPORARY TABLE temp_TestRawMultipleSQL2 (id int);

insert into temp_TestRawMultipleSQL2
select 1;

SELECT * FROM temp_TestRawMultipleSQL2
`
	ctx := context.Background()
	tx, _ := connection.Begin(ctx)
	defer tx.Commit(ctx)
	//TextSQL := "SELECT 1; SELECT 2"
	Rows, err := RawMultipleSQL(tx, TextSQL)
	if err != nil {
		t.Error("TestRawMultipleSQL() error: ", err)
		return
	}
	if Rows == nil {
		t.Error("TestRawMultipleSQL() error: Rows == nil")
		return
	}
	defer Rows.Close()

	Otvet := 0
	for Rows.Next() {
		err := Rows.Scan(&Otvet)
		if err != nil {
			t.Error("TestRawMultipleSQL() Scan() error: ", err)
		}
	}

	t.Log("Прошло время: ", time.Since(TimeStart))
}

// TestRawMultipleSQL2 - negative test, with error
func TestRawMultipleSQL2(t *testing.T) {
	config_main.LoadEnv()
	GetConnection()
	defer CloseConnection()

	defer microl.ShowTimePassed(time.Now())

	TextSQL := `
drop table if exists temp_TestRawMultipleSQL2; 
CREATE TEMPORARY TABLE temp_TestRawMultipleSQL2 (id int);

insert into temp_TestRawMultipleSQL2
select 1;

SELECT * FROM temp_TestRawMultipleSQL2
`

	ctx := contextmain.GetContext()
	Rows, err := PgxPool.Query(ctx, TextSQL)
	if err == nil {
		t.Error("TestRawMultipleSQL2() Query() error: ", err)
		return
	}
	if Rows == nil {
		t.Error("TestRawMultipleSQL2() error: Rows == nil")
		return
	}

}

func TestReplaceSchemaName(t *testing.T) {
	TextSQL := "SELECT * FROM public.users"
	Settings.DB_SCHEMA = "myschema"
	ExpectedSQL := "SELECT * FROM myschema.users"
	ActualSQL := ReplaceSchemaName(TextSQL, "public")
	if ActualSQL != ExpectedSQL {
		t.Errorf("Expected %v, but got %v", ExpectedSQL, ActualSQL)
	}
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
