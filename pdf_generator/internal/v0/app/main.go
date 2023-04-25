// главный модуль программы

package main

import (
	"github.com/manyakrus/starter/common/v0/config"
	"github.com/manyakrus/starter/common/v0/contextmain"
	logger "github.com/manyakrus/starter/common/v0/logger"
	stopapp "github.com/manyakrus/starter/common/v0/stopapp"
	//	// "github.com/manyakrus/starter/claim_debtors_list/internal/v0/app/config"

	"github.com/manyakrus/starter/pdf_generator/internal/v0/app/programdir"
)

//// log - глобальный логгер
var log = logger.GetLog()

//
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

	stopapp.GetWaitGroup_Main().Wait()

}
