// модуль для обмена с сервисом NATS через sync_exchange
package sync_exchange_connect

import (
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/nats_connect"
	"github.com/ManyakRus/starter/stopapp"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_types"
)

// Connect - подключение к NATS SyncExchange
func Connect(ServiceName string) {
	var err error

	nats_connect.FillSettings()
	sNATS_PORT := (nats_connect.Settings.NATS_PORT)
	url := "nats://" + nats_connect.Settings.NATS_HOST + ":" + sNATS_PORT

	err = sync_exchange.InitSyncExchange(url, ServiceName)
	if err != nil {
		log.Panicln("Can not start NATS, server: ", nats_connect.Settings.NATS_HOST, " error: ", err)
	}

	//return err
}

// Start - необходимые процедуры для подключения к серверу Nats SyncExchange
func Start(ServiceName string) {
	//var err error

	Connect(ServiceName)

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

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
