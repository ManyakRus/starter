package mssql_stek

import (
	model "gitlab.aescorp.ru/dsp_dev/claim/common/object_model"
	"github.com/manyakrus/starter/config"
	"github.com/manyakrus/starter/mssql_gorm"
	"testing"
)

var CONNECTION = model.Connection{ID: 3, BranchId: 2, IsLegal: true}

func TestFindDateClosedMonth(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
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
	config.LoadEnv()
	mssql_gorm.Connect()
	defer mssql_gorm.CloseConnection()

	date1_balances, date2_balances, date1_doc, date2_doc, err := FindDateFromTo(CONNECTION)
	t.Logf("date1_balances, date2_balances, date1_doc, date2_doc, err: \n %v \n %v \n %v \n %v \n %v \n", date1_balances, date2_balances, date1_doc, date2_doc, err)
	if err != nil {
		t.Error("mssql_test.TestFindDateFromTo() error: ", err)
	}
}
