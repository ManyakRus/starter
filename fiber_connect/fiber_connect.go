package fiber_connect

import (
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/logger"
	"github.com/ManyakRus/starter/stopapp"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"reflect"
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

var Client *fiber.App

func Connect() {
	if Settings.WEBSERVER_PORT == "" {
		FillSettings()
	}
	Client = fiber.New()

	log.Info("Fiber connected. OK. host: ", Settings.WEBSERVER_HOST, ":", Settings.WEBSERVER_PORT)

}

// FillSettings загружает переменные окружения в структуру из файла или из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.WEBSERVER_HOST = os.Getenv("WEBSERVER_HOST")
	Settings.WEBSERVER_PORT = os.Getenv("WEBSERVER_PORT")
	if Settings.WEBSERVER_HOST == "" {
		log.Info("Need fill WEBSERVER_HOST ! in OS Environment ")
	}

	if Settings.WEBSERVER_PORT == "" {
		log.Info("Need fill WEBSERVER_PORT ! in OS Environment ")
		Settings.WEBSERVER_PORT = "3002"
	}

	//

}

func CloseConnection_err() error {
	err := Client.Shutdown()

	return err
}

func CloseConnection() {
	err := CloseConnection_err()
	if err != nil {
		log.Error("Fiber CloseConnection() error: ", err)
	} else {
		log.Info("Fiber connection closed.")
	}
}

// StartLiveness - делает соединение с БД, отключение и др.
func Start() {
	if Client == nil {
		FillSettings()
		Connect()
	}

	go Listen()

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

}

func Listen() {
	addr := Settings.WEBSERVER_HOST + ":" + Settings.WEBSERVER_PORT //3002
	err := Client.Listen(addr)
	if err != nil {
		log.Panic(PackageName, "Listen() error: ", err)
	}

}

// WaitStop - ожидает отмену глобального контекста или сигнала завершения приложения
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
