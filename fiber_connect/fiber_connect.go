package fiber_connect

import (
	"context"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/logger"
	"github.com/ManyakRus/starter/stopapp"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

// empty - пустая структура для имени пакета
type empty struct{}

// PackageName - имя пакета golang
var PackageName = filepath.Base(reflect.TypeOf(empty{}).PkgPath())

// log - глобальный логгер
var log = logger.GetLog()

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	WEBSERVER_HOST string
	WEBSERVER_PORT string
}

// Client - клиент веб сервера
var Client *fiber.App

// WEBSERVER_PORT_DEFAULT - порт веб-сервера по умолчанию
var WEBSERVER_PORT_DEFAULT = "3000"

func Connect() {
	if Settings.WEBSERVER_PORT == "" {
		FillSettings()
	}
	Client = fiber.New()

	go Listen()

	log.Info("Fiber connected. OK. host: ", Settings.WEBSERVER_HOST, ":", Settings.WEBSERVER_PORT)

}

// FillSettings загружает переменные окружения в структуру из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.WEBSERVER_HOST = os.Getenv("WEBSERVER_HOST")
	Settings.WEBSERVER_PORT = os.Getenv("WEBSERVER_PORT")
	if Settings.WEBSERVER_HOST == "" {
		log.Debug("Need fill WEBSERVER_HOST ! in OS Environment ")
	}

	if Settings.WEBSERVER_PORT == "" {
		Settings.WEBSERVER_HOST = os.Getenv("WEB_SERVER_HOST")
		Settings.WEBSERVER_PORT = os.Getenv("WEB_SERVER_PORT")
	}

	if Settings.WEBSERVER_PORT == "" {
		log.Warn("Need fill WEBSERVER_PORT ! in OS Environment. Use default: ", WEBSERVER_PORT_DEFAULT)
		Settings.WEBSERVER_PORT = WEBSERVER_PORT_DEFAULT
	}

	//

}

// Listen_err - начинает прослушивать порт, паника при ошибке
func Listen() {
	err := Listen_err()
	if err != nil {
		log.Panic(PackageName, "Fiber Listen() error: ", err)
	}

}

// Listen_err - начинает прослушивать порт, возвращает ошибку
func Listen_err() error {
	addr := Settings.WEBSERVER_HOST + ":" + Settings.WEBSERVER_PORT //3002
	err := Client.Listen(addr)
	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled.")
	}

	//ждём пока отправляемых сейчас сообщений будет =0
	//stopapp.WaitTotalMessagesSendingNow(PackageName)

	//закрываем соединение
	CloseConnection()
	stopapp.GetWaitGroup_Main().Done()
}

func GetHost() string {
	Otvet := "127.0.0.1"

	if Settings.WEBSERVER_HOST != "" {
		Otvet = Settings.WEBSERVER_HOST
	}

	return Otvet
}

// Start - запускает работу веб-сервера
// Graceful shutdown только для тех кто пользуется этим репозиторием для старта и останова
func Start() {
	if Client == nil {
		FillSettings()
		Connect()
	}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

}

// Start_ctx - запускает работу веб-сервера
// Graceful shutdown - для тех кто передаст сюда свой контекст и WaitGroup
// для тех кто НЕ пользуется этим репозиторием для старта и останова
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup) error {
	var err error

	//запомним к себе контекст и WaitGroup
	contextmain.Ctx = ctx
	stopapp.SetWaitGroup_Main(WaitGroup)

	//
	Start()

	return err
}

// CloseConnection - закрывает соединения веб-сервера, возвращает ошибку
func CloseConnection_err() error {
	err := Client.Shutdown()

	return err
}

// CloseConnection - закрывает соединения веб-сервера
func CloseConnection() {
	err := CloseConnection_err()
	if err != nil {
		log.Error("Fiber CloseConnection() error: ", err)
	} else {
		log.Info("Fiber connection closed.")
	}
}
