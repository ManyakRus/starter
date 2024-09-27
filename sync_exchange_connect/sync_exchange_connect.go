// модуль для обмена с сервисом NATS через sync_exchange
package sync_exchange_connect

import (
	"context"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/nats_connect"
	"github.com/ManyakRus/starter/stopapp"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_types"
	"sync"
)

// Connect - подключение к NATS SyncExchange
func Connect(ServiceName string) {
	err := Connect_err(ServiceName)
	LogInfo_Connected(err)
}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicln("Can not start NATS, server: ", nats_connect.Settings.NATS_HOST, " error: ", err)
	} else {
		log.Info("NATS connected. OK., server: ", nats_connect.Settings.NATS_HOST, ":", nats_connect.Settings.NATS_PORT, " error: ", err)
	}

}

// Connect_err - подключение к NATS SyncExchange
func Connect_err(ServiceName string) error {
	var err error

	nats_connect.FillSettings()
	sNATS_PORT := (nats_connect.Settings.NATS_PORT)
	url := "nats://" + nats_connect.Settings.NATS_HOST + ":" + sNATS_PORT

	err = sync_exchange.InitSyncExchange(url, ServiceName, "")

	return err
}

// Start - необходимые процедуры для подключения к серверу Nats SyncExchange
func Start(ServiceName string) {
	var err error

	ctx := contextmain.GetContext()
	WaitGroup := stopapp.GetWaitGroup_Main()
	err = Start_ctx(&ctx, WaitGroup, ServiceName)
	LogInfo_Connected(err)

}

// Start_ctx - необходимые процедуры для подключения к NATS с библиотекой SyncExchange
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто Start()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup, ServiceName string) error {
	var err error

	//запомним к себе контекст и WaitGroup
	contextmain.Ctx = ctx
	stopapp.SetWaitGroup_Main(WaitGroup)

	//
	err = Connect_err(ServiceName)
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
		log.Error("CloseConnection() error: ", err)
	} else {
		log.Info("NATS stopped")
	}

	//return err
}

// WaitStop - ожидает отмену глобального контекста или сигнала завершения приложения
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. NATS.")
	}

	//ждём пока отправляемых сейчас сообщений будет =0
	stopapp.WaitTotalMessagesSendingNow("nats SyncExchange")

	//закрываем соединение
	CloseConnection()
	stopapp.GetWaitGroup_Main().Done()
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
