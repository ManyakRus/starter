package nats_connect

import (
	"context"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/port_checker"
	"github.com/nats-io/nats.go"
	"os"
	"sync"

	//"github.com/ManyakRus/starter/common/v0/micro"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/stopapp"
)

// Conn - соединение к серверу nats
var Conn *nats.Conn

// log - глобальный логгер
//var log = logger.GetLog()

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	NATS_HOST     string
	NATS_PORT     string
	NATS_LOGIN    string
	NATS_PASSWORD string
}

// Connect - подключается к серверу Nats
func Connect(ServiceName string) {
	var err error

	err = Connect_err(ServiceName)
	LogInfo_Connected(err)
}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicln("NATS Connect() error: ", err)
	} else {
		log.Info("NATS Connect() ok ")
	}

}

// Connect_err - подключается к серверу Nats и возвращает ошибку
func Connect_err(ServiceName string) error {
	var err error

	if Settings.NATS_HOST == "" {
		FillSettings()
	}

	port_checker.CheckPort(Settings.NATS_HOST, Settings.NATS_PORT)

	sNATS_PORT := Settings.NATS_PORT
	URL := "nats://" + Settings.NATS_HOST + ":" + sNATS_PORT
	UserInfo := nats.UserInfo(Settings.NATS_LOGIN, Settings.NATS_PASSWORD)
	Options := nats.GetDefaultOptions()
	Options.Name = ServiceName
	Conn, err = nats.Connect(URL, UserInfo)

	//nats.ManualAck()
	return err
}

// StartNats - необходимые процедуры для подключения к серверу Nats
func StartNats(ServiceName string) {
	var err error

	ctx := contextmain.GetContext()
	WaitGroup := stopapp.GetWaitGroup_Main()
	err = Start_ctx(&ctx, WaitGroup, ServiceName)
	LogInfo_Connected(err)

}

// Start_ctx - необходимые процедуры для подключения к серверу NATS
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто StartNats()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup, ServiceName string) error {
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
	err = Connect_err(ServiceName)
	if err != nil {
		return err
	}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	return err
}

// CloseConnection - закрывает соединение с сервером Nats
func CloseConnection() {
	//var err error

	if Conn == nil {
		return
	}

	Conn.Close()

	log.Info("NATS stopped")

	return
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {
	defer stopapp.GetWaitGroup_Main().Done()

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. nats_connect.")
	}

	CloseConnection()

}

// FillSettings загружает переменные окружения в структуру из файла или из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.NATS_HOST = os.Getenv("NATS_HOST")
	Settings.NATS_PORT = os.Getenv("NATS_PORT")
	Settings.NATS_LOGIN = os.Getenv("NATS_LOGIN")
	Settings.NATS_PASSWORD = os.Getenv("NATS_PASSWORD")

	// как у всех
	if Settings.NATS_HOST == "" {
		Settings.NATS_HOST = os.Getenv("BUS_LOCAL_HOST")
	}

	if Settings.NATS_PORT == "" {
		Settings.NATS_PORT = os.Getenv("BUS_LOCAL_PORT")
	}

	if Settings.NATS_HOST == "" {
		log.Error("Need fill NATS_HOST ! in os.ENV ")
	}

	//
	if Settings.NATS_PORT == "" {
		log.Error("Need fill NATS_PORT ! in os.ENV ")
	}

	//if Settings.NATS_LOGIN == "" {
	//	log.Panicln("Need fill NATS_LOGIN ! in os.ENV ")
	//}
	//
	//if Settings.NATS_PASSWORD == "" {
	//	log.Panicln("Need fill NATS_PASSWORD ! in os.ENV ")
	//}

	//
}

// SendMessageCtx - Отправляет сообщение, учитывает таймаут контекста
func SendMessageCtx(ctx context.Context, subject string, data []byte) error {
	var err error

	fn := func() error {
		err = SendMessage(subject, data)
		return err
	}
	err = micro.GoGo(ctx, fn)
	return err
}

// SendMessage - Отправляет сообщение
func SendMessage(subject string, data []byte) error {
	var err error

	err = Conn.Publish(subject, data)

	return err
}
