package mssql_stek

import (
	"github.com/ManyakRus/starter/config_main"
	"github.com/ManyakRus/starter/mssql_gorm"
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/object_model/entities/connections"
	"testing"
)

var CONNECTION = connections.Connection{ID: 3, BranchID: 2, IsLegal: true}

func TestFindDateClosedMonth(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	mssql_gorm.Connect()
	defer mssql_gorm.CloseConnection()

	Otvet, err := FindDate_ClosedMonth(CONNECTION)
	if err != nil {
		t.Error("mssql_stek_test.TestFindDateClosedMonth() FindDate_ClosedMonth() error:", err)
	} else {
		t.Log("mssql_stek_test.TestFindDateClosedMonth() FindDate_ClosedMonth() Otvet:", Otvet)
	}
}

func TestFindDateFromTo(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	mssql_gorm.Connect()
	defer mssql_gorm.CloseConnection()

	date1_balances, date2_balances, date1_doc, date2_doc, err := FindDateFromTo(CONNECTION)
	t.Logf("date1_balances, date2_balances, date1_doc, date2_doc, err: \n %v \n %v \n %v \n %v \n %v \n", date1_balances, date2_balances, date1_doc, date2_doc, err)
	if err != nil {
		t.Error("mssql_test.TestFindDateFromTo() error: ", err)
	}
}
