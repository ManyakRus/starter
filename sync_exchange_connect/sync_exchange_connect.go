// модуль для обмена с сервисом NATS через sync_exchange
package sync_exchange_connect

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/nats_connect"
	"github.com/ManyakRus/starter/stopapp"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_types"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

// TopicNamePprof_Heap_Suffix - имя суффикса топика с профилем памяти
const TopicNamePprof_Heap_Suffix = ".heap_profile"

// TopicNamePprof_Goroutine_Suffix - имя суффикса топика с профилем памяти
const TopicNamePprof_Goroutine_Suffix = ".goroutine_profile"

// serviceName - имя сервиса который подключается
var serviceName string

// Connect - подключение к NATS SyncExchange
func Connect(ServiceName, ServiceVersion string) {
	err := Connect_err(ServiceName, ServiceVersion)
	LogInfo_Connected(err)
}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicln("Can not start NATS, server: ", nats_connect.Settings.NATS_HOST, " error: ", err)
	} else {
		log.Info("NATS connected. OK., server: ", nats_connect.Settings.NATS_HOST, ":", nats_connect.Settings.NATS_PORT)
	}

}

// Connect_err - подключение к NATS SyncExchange
func Connect_err(ServiceName, ServiceVersion string) error {
	var err error

	//запомним ServiceName
	if serviceName == "" {
		serviceName = ServiceName
	}

	//
	nats_connect.FillSettings()
	sNATS_PORT := (nats_connect.Settings.NATS_PORT)
	url := "nats://" + nats_connect.Settings.NATS_HOST + ":" + sNATS_PORT

	err = sync_exchange.InitSyncExchange(url, ServiceName, ServiceVersion)

	return err
}

// Start - необходимые процедуры для подключения к серверу Nats SyncExchange
func Start(ServiceName, ServiceVersion string) {
	var err error

	ctx := contextmain.GetContext()
	WaitGroup := stopapp.GetWaitGroup_Main()
	err = Start_ctx(&ctx, WaitGroup, ServiceName, ServiceVersion)
	LogInfo_Connected(err)

}

// Start_ctx - необходимые процедуры для подключения к NATS с библиотекой SyncExchange
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто Start()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup, ServiceName, ServiceVersion string) error {
	var err error

	//запомним к себе контекст
	if contextmain.Ctx != ctx {
		contextmain.SetContext(ctx)
	}
	//contextmain.Ctx = ctx
	if ctx == nil {
		contextmain.GetContext()
	}

	//запомним к себе WaitGroup
	stopapp.SetWaitGroup_Main(WaitGroup)
	if WaitGroup == nil {
		stopapp.StartWaitStop()
	}

	//
	err = Connect_err(ServiceName, ServiceVersion)
	if err != nil {
		return err
	}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	return err
}

// CloseConnection - закрывает соединение с NATS
func CloseConnection() {
	err := sync_exchange.DeInitSyncExchange()
	if err != nil {
		log.Warn("CloseConnection() warning: ", err)
	} else {
		log.Info("NATS stopped")
	}

	//return err
}

// WaitStop - ожидает отмену глобального контекста или сигнала завершения приложения
func WaitStop() {
	defer stopapp.GetWaitGroup_Main().Done()

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. sync_exchange_connect")
	}

	//ждём пока отправляемых сейчас сообщений будет =0
	stopapp.WaitTotalMessagesSendingNow("sync_exchange_connect")

	//закрываем соединение
	CloseConnection()
}

// SendResponseError - Отправляет ответ в NATS SyncExchange
func SendResponseError(sp *sync_types.SyncPackage, err error) {

	if err == nil {
		return
	}

	sp_otvet := sync_types.MakeSyncError("", 0, err.Error())

	err = sync_exchange.SendResponse(sp, sp_otvet)
	if err != nil {
		log.Error("SendResponse() Error: ", err)
	}

}

// Start_PprofNats - профилирование памяти отправляет в NATS, бесконечно + WaitGroup
func Start_PprofNats() {
	TextTest := TextTestOrEmpty()
	topicHeapProfile := serviceName + TextTest + TopicNamePprof_Heap_Suffix
	log.Info("Start_PprofNats(), topic: ", topicHeapProfile)

	stopapp.GetWaitGroup_Main().Add(1)
	go pprofNats_forever_go()
}

// pprofNats_forever_go - профилирование памяти отправляет в NATS, бесконечно + WaitGroup
func pprofNats_forever_go() {
	defer stopapp.GetWaitGroup_Main().Done()
	PprofNats_forever()
}

// PprofNats_forever - профилирование памяти отправляет в NATS, бесконечно
func PprofNats_forever() {
	var err error

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-contextmain.GetContext().Done():
			log.Warn("Context app is canceled. sync_exchange_connect.ping")
			break loop
		case <-ticker.C:
			err = PprofNats1()
			if err != nil {
				err = fmt.Errorf("PprofNats1(), error: %w", err)
				log.Error(err)
				time.Sleep(60 * time.Second)
			}
		}
	}
}

// PprofNats1 - профилирование памяти отправляет в NATS 1 раз
func PprofNats1() error {
	var err error

	//память
	err = PprofMemoryProfile1()
	if err != nil {
		err = fmt.Errorf("PprofMemoryProfile1(), error: %w", err)
		log.Error(err)
		return err
	}

	//горутины
	err = PprofGoroutines1()
	if err != nil {
		err = fmt.Errorf("PprofMemoryProfile1(), error: %w", err)
		log.Error(err)
		return err
	}

	return err
}

// PprofMemoryProfile1 - профилирование памяти отправляет в NATS 1 раз
func PprofMemoryProfile1() error {
	var err error

	TextTest := TextTestOrEmpty()
	topicHeapProfile := serviceName + TextTest + TopicNamePprof_Heap_Suffix
	var buf bytes.Buffer
	err = pprof.WriteHeapProfile(&buf)
	if err != nil {
		err = fmt.Errorf("pprof.WriteHeapProfile(), topic: %v, error: %w", topicHeapProfile, err)
		log.Error(err)
		time.Sleep(10 * time.Second)
		return err
	}
	err = sync_exchange.SendRawMessage(topicHeapProfile, buf.Bytes())
	if err != nil {
		err = fmt.Errorf("sync_exchange.SendRawMessage(), topic: %v, error: %w", topicHeapProfile, err)
		log.Error(err)
		time.Sleep(10 * time.Second)
		return err
	}

	return err
}

// PprofGoroutines1 - профилирование горутин отправляет в NATS 1 раз
func PprofGoroutines1() error {
	var err error

	TextTest := TextTestOrEmpty()
	topicHeapProfile := serviceName + TextTest + TopicNamePprof_Goroutine_Suffix
	var buf bytes.Buffer
	err = pprof.Lookup("goroutine").WriteTo(&buf, 2)
	if err != nil {
		err = fmt.Errorf("pprof.WriteHeapProfile(), topic: %v, error: %w", topicHeapProfile, err)
		log.Error(err)
		time.Sleep(10 * time.Second)
		return err
	}
	err = sync_exchange.SendRawMessage(topicHeapProfile, buf.Bytes())
	if err != nil {
		err = fmt.Errorf("sync_exchange.SendRawMessage(), topic: %v, error: %w", topicHeapProfile, err)
		log.Error(err)
		time.Sleep(10 * time.Second)
		return err
	}

	return err
}

// TextTestOrEmpty - возвращает "_test" или ""
func TextTestOrEmpty() string {
	Otvet := "_test"
	stage := os.Getenv("STAGE")
	switch stage {
	case "dev":
		Otvet = ""
	case "prod":
		Otvet = ""
	}

	return Otvet
}
