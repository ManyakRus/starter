// модуль для корректной остановки работы приложения

package stopapp

import (
	"github.com/ManyakRus/starter/log"
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
//var log = logger.GetLog()

// SignalInterrupt - канал для ожидания сигнала остановки приложения
var SignalInterrupt chan os.Signal

// onceWGMain - гарантирует создание WGMain один раз
var onceWGMain sync.Once

// TotalMessagesSendingNow - количество сообщений отправляющихся прям сейчас
var TotalMessagesSendingNow int32

// SecondsWaitTotalMessagesSendingNow - количество секунд ожидания для отправки последнего сообщения
const SecondsWaitTotalMessagesSendingNow = 10

// StartWaitStop - запускает ожидание сигнала завершения приложения
func StartWaitStop() {
	//создадим контекст, т.к. попозже уже гонка данных
	contextmain.GetContext()

	//
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
		//os.Exit(0)
		log.Warn("Context = nil")
	}

	GetWaitGroup_Main().Wait()
}

// WaitStop - ожидает отмену глобального контекста или сигнала завершения приложения
func WaitStop() {
	defer GetWaitGroup_Main().Done()

	select {
	case <-SignalInterrupt:
		log.Warn("Interrupt clean shutdown.")
		if contextmain.CancelContext != nil {
			contextmain.CancelContext()
		}
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. stopapp")
	}

	////ожидаем закрытие всех подключений
	////создаём массив для обратной сортировки
	//MassWait := make([]KeyValueWaitGroupContext, 0)
	//OrderedMapConnections.OrderedRange(func(key string, value WaitGroupContext) {
	//	KeyValueWaitGroupContext1 := KeyValueWaitGroupContext{Key: key, Value: value}
	//	MassWait = append(MassWait, KeyValueWaitGroupContext1)
	//})
	//
	////запускаем Wait() в обратном порядке
	//for i := len(MassWait) - 1; i >= 0; i-- {
	//	key := MassWait[i].Key
	//	value := MassWait[i].Value
	//	log.Debugf("Ожидаем закрытия соединения: %s", key)
	//	if value.CancelCtxFunc != nil {
	//		value.CancelCtxFunc()
	//	}
	//	if value.WaitGroup != nil {
	//		WaitGroup1 := value.WaitGroup
	//		WaitGroup1.Wait()
	//	}
	//
	//	//
	//	OrderedMapConnections.Delete(key)
	//}

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

// KeyValueWaitGroupContext - структура ключ-значение
type KeyValueWaitGroupContext struct {
	Key   string
	Value WaitGroupContext
}

// Wait_GracefulShutdown - ожидает завершения всех горутин программы, а потом ожидает закрытие всех подключений
func Wait_GracefulShutdown() {

	//ожидаем отмену контекста
	select {
	case <-contextmain.GetContext().Done():

	}

	//ожидаем завершения всех горутин программы
	GetWaitGroup_Main().Wait()

	//ожидаем закрытие всех подключений
	//создаём массив для обратной сортировки
	MassWait := make([]KeyValueWaitGroupContext, 0)
	OrderedMapConnections.OrderedRange(func(key string, value WaitGroupContext) {
		KeyValueWaitGroupContext1 := KeyValueWaitGroupContext{Key: key, Value: value}
		MassWait = append(MassWait, KeyValueWaitGroupContext1)
	})

	//запускаем Wait() в обратном порядке
	for i := len(MassWait) - 1; i >= 0; i-- {
		key := MassWait[i].Key
		value := MassWait[i].Value
		log.Debugf("Ожидаем закрытия соединения: %s", key)
		CancelCtxFunc1 := value.CancelCtxFunc
		if CancelCtxFunc1 != nil {
			CancelCtxFunc1()
		}
		WaitGroup1 := value.WaitGroup
		if WaitGroup1 == nil {
			WaitGroup1.Wait()
		}

		//
		OrderedMapConnections.Delete(key)
	}
}
