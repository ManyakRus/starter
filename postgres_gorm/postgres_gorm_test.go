package postgres_gorm

import (
	"errors"
	"sync"
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

	waitGroup_Connect.Add(1)
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
	wg := sync.WaitGroup{}
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

		wg.Done()
	}

	//запустим 100 потоков
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go f(t)
	}

	//micro.Pause(100)
	wg.Wait()

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

func TestReplaceTemporaryTableNamesToUnique(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "One temporary table",
			input: "CREATE TEMPORARY TABLE temp_TableName (id int); SELECT * FROM temp_TableName",
		},
		{
			name:  "Multiple temporary tables",
			input: "CREATE TEMPORARY TABLE temp_TableName1 (id int); CREATE TEMPORARY TABLE temp_TableName2 (name varchar); SELECT * FROM temp_TableName1; SELECT * FROM temp_TableName2",
		},
		{
			name:  "Temporary table with different cases",
			input: "CREATE TEMPORARY TABLE temp_tableName (id int); SELECT * FROM temp_tableName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceTemporaryTableNamesToUnique(tt.input)
			if got == tt.input {
				t.Errorf("ReplaceTemporaryTableNamesToUnique() error: no changes")
			}
		})
	}
}

func TestReplaceTemporaryTableNamesToUnique_EmptyInput(t *testing.T) {
	input := ""
	expected := ""
	result := ReplaceTemporaryTableNamesToUnique(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestReplaceTemporaryTableNamesToUnique_NoTempTable(t *testing.T) {
	input := "SELECT * FROM TableName"
	expected := "SELECT * FROM TableName"
	result := ReplaceTemporaryTableNamesToUnique(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestReplaceTemporaryTableNamesToUnique_OneTempTable(t *testing.T) {
	input := "CREATE TEMPORARY TABLE temp_Test (id int); SELECT * FROM temp_Test "
	result := ReplaceTemporaryTableNamesToUnique(input)
	if result == input {
		t.Errorf("TestReplaceTemporaryTableNamesToUnique_OneTempTable() error: no changes")
	}
}

func TestReplaceTemporaryTableNamesToUnique_MultipleTempTables(t *testing.T) {
	input := "CREATE TEMPORARY TABLE temp1 (id int); CREATE TEMPORARY TABLE temp2 (name varchar); SELECT * FROM temp1;"
	result := ReplaceTemporaryTableNamesToUnique(input)
	if result == input {
		t.Errorf("TestReplaceTemporaryTableNamesToUnique_OneTempTable() error: no changes")
	}
}

func TestSetSingularTableNames(t *testing.T) {
	// Test case 1: IsSingular is true
	IsSingular := true
	SetSingularTableNames(IsSingular)
	if NamingStrategy.SingularTable != IsSingular {
		t.Errorf("Expected SingularTable to be %v, but got %v", IsSingular, NamingStrategy.SingularTable)
	}

	// Test case 2: IsSingular is false
	IsSingular = false
	SetSingularTableNames(IsSingular)
	if NamingStrategy.SingularTable != IsSingular {
		t.Errorf("Expected SingularTable to be %v, but got %v", IsSingular, NamingStrategy.SingularTable)
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
