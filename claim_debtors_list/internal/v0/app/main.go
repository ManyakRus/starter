// главный модуль программы

package main

import (
	"github.com/ManyakRus/starter/config"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/logger"
	mssql "github.com/ManyakRus/starter/mssql_connect"
	"github.com/ManyakRus/starter/stopapp"
	//	// "github.com/ManyakRus/starter/claim_debtors_list/internal/v0/app/config"

	"github.com/ManyakRus/starter/claim_debtors_list/db"
	"github.com/ManyakRus/starter/claim_debtors_list/programdir"
)

// // log - глобальный логгер
var log = logger.GetLog()

// main - старт приложения
func main() {
	StartApp()
}

// StartApp - выполнение всех операций для старта приложения
func StartApp() {
	ProgramDir := programdir.ProgramDir()
	config.LoadEnv(ProgramDir)

	stopapp.StartWaitStop()

	contextmain.GetContext()

	mssql.StartDB()

	stopapp.GetWaitGroup_Main().Add(1)
	go StartForever()

	stopapp.GetWaitGroup_Main().Wait()

	IsClosed := mssql.IsClosed()
	log.Info("DB closed: ", IsClosed)

}

func StartForever() {
	Mass, err := db.FindDebtorsList()
	if err != nil {
		log.Error("FindDebtorsList() error: ", err)
	}
	log.Debugf("DebtorsList: %#v", Mass)
	stopapp.GetWaitGroup_Main().Done()
	contextmain.CancelContext()
}
