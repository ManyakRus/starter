// модуль для корректной остановки работы приложения

package stopapp

import (
	"github.com/ManyakRus/starter/logger"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	//"github.com/sirupsen/logrus"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	//	"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/micro"
	//"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/db"
	//"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/grpcserver"
	//"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/logger"
	//"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/micro"
)

// log - глобальный логгер
var log = logger.GetLog()

// SignalInterrupt - канал для ожидания сигнала остановки приложения
var SignalInterrupt chan os.Signal

// wgMain - группа ожидания завершения всех частей программы
var wgMain *sync.WaitGroup

// lockWGMain - гарантирует получение WGMain с учётом многопоточности
var lockWGMain = &sync.Mutex{}

// onceWGMain - гарантирует создание WGMain один раз
var onceWGMain sync.Once

// TotalMessagesSendingNow - количество сообщений отправляющихся прям сейчас
var TotalMessagesSendingNow int32

// SecondsWaitTotalMessagesSendingNow - количество секунд ожидания для отправки последнего сообщения
const SecondsWaitTotalMessagesSendingNow = 10

// SetWaitGroup_Main - присваивает внешний WaitGroup
func SetWaitGroup_Main(wg *sync.WaitGroup) {
	wgMain = wg
}

// GetWaitGroup_Main - возвращает группу ожидания завершения всех частей программы
func GetWaitGroup_Main() *sync.WaitGroup {
	lockWGMain.Lock()
	defer lockWGMain.Unlock()
	//
	//if wgMain == nil {
	//	wgMain = &sync.WaitGroup{}
	//}

	if wgMain == nil {
		//onceWGMain.Do(func() {
		wgMain = &sync.WaitGroup{}
		//})
	}

	return wgMain
}

// StartWaitStop - запускает ожидание сигнала завершения приложения
func StartWaitStop() {
	SignalInterrupt = make(chan os.Signal, 1)

	fnWait := func() {
		signal.Notify(SignalInterrupt, os.Interrupt, syscall.SIGTERM)
	}
	go fnWait()

	GetWaitGroup_Main().Add(1)
	go WaitStop()

}

// StopApp - отмена глобального контекста для остановки работы приложения
func StopApp() {
	if contextmain.CancelContext != nil {
		contextmain.CancelContext()
	} else {
		//os.Exit(0)
		log.Warn("Context = nil")
	}

}

// StopApp - отмена глобального контекста для остановки работы приложения
func StopAppAndWait() {
	if contextmain.CancelContext != nil {
		contextmain.CancelContext()
	} else {
		os.Exit(0)
	}

	GetWaitGroup_Main().Wait()
}

// WaitStop - ожидает отмену глобального контекста или сигнала завершения приложения
func WaitStop() {

	select {
	case <-SignalInterrupt:
		log.Warn("Interrupt clean shutdown.")
		contextmain.CancelContext()
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled.")
	}

	GetWaitGroup_Main().Done()
}

// ожидает чтоб прям щас ничего не отправлялось
func WaitTotalMessagesSendingNow(filename string) {
	for f := 0; f < SecondsWaitTotalMessagesSendingNow; f++ {
		TotalMessages := atomic.LoadInt32(&TotalMessagesSendingNow)
		if TotalMessages == 0 {
			break
		}
		log.Warn("TotalMessagesSendingNow =", TotalMessages, " !=0 sleep(1000), filename: ", filename)
		micro.Sleep(1000)
	}
}
