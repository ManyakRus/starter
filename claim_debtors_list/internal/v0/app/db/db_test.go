package db

import (
	"testing"

	//log "github.com/sirupsen/logrus"

	"github.com/manyakrus/starter/config"

	mssql "github.com/manyakrus/starter/mssql_connect"
	//logger "github.com/manyakrus/starter/common/v0/logger"
	//stopapp "github.com/manyakrus/starter/common/v0/stopapp"

	"github.com/manyakrus/starter/claim_debtors_list/programdir"
)

//func recoveryFunction() {
//	if recoveryMessage := recover(); recoveryMessage != nil {
//	}
//}

func TestFindDebtorsList(t *testing.T) {

	ProgramDir := programdir.ProgramDir()
	config.LoadEnv(ProgramDir)
	err := mssql.Connect_err()
	if err != nil {
		t.Error("TestConnect error: ", err)
	}

	Otvet, err := FindDebtorsList()
	if err != nil {
		t.Error("db_test.TestFindDebtorsList() error: ", err)
	}
	if Otvet == nil {
		t.Error("db_test.TestFindDebtorsList() Otvet = nil ")
	}

	if len(Otvet) == 0 {
		t.Error("db_test.TestFindDebtorsList() len(Otvet) = 0 ")
	}

	err = mssql.CloseConnection_err()
	if err != nil {
		t.Error("TestConnect() error: ", err)
	}

}
