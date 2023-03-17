package nats_connect

import (
	"github.com/nats-io/nats.go"
	"os"

	//"gitlab.aescorp.ru/dsp_dev/claim/nikitin/common/v0/micro"
	"gitlab.aescorp.ru/dsp_dev/claim/nikitin/common/v0/contextmain"
	"gitlab.aescorp.ru/dsp_dev/claim/nikitin/common/v0/logger"
	"gitlab.aescorp.ru/dsp_dev/claim/nikitin/common/v0/stopapp"
)

// Conn - соединение к серверу nats
var Conn *nats.Conn

// log - глобальный логгер
var log = logger.GetLog()

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	NATS_SERVER   string
	NATS_PORT     string
	NATS_LOGIN    string
	NATS_PASSWORD string
}

// Connect - подключается к серверу Nats
func Connect() {
	var err error

	if Settings.NATS_SERVER == "" {
		FillSettings()
	}

	sNATS_PORT := (Settings.NATS_PORT)
	URL := "nats://" + Settings.NATS_SERVER + ":" + sNATS_PORT
	UserInfo := nats.UserInfo(Settings.NATS_LOGIN, Settings.NATS_PASSWORD)
	Conn, err = nats.Connect(URL, UserInfo)
	if err != nil {
		log.Panicln("Connect() error: ", err)
	}

	nats.ManualAck()
}

// StartNats - необходимые процедуры для подключения к серверу Nats
func StartNats() {
	Connect()

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

}

// CloseConnection - закрывает соединение с сервером Nats
func CloseConnection() error {
	var err error

	if Conn == nil {
		return err
	}

	Conn.Close()

	return err
}

// WaitStop - ожидает отмену глобального контекста или сигнала завершения приложения
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. NATS_Connect.")
	}

	err := CloseConnection()
	if err != nil {
		log.Error("CloseConnection() error: ", err)
	}
	stopapp.GetWaitGroup_Main().Done()
}

// FillSettings загружает переменные окружения в структуру из файла или из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.NATS_SERVER = os.Getenv("NATS_SERVER")
	Settings.NATS_PORT = os.Getenv("NATS_PORT")
	Settings.NATS_LOGIN = os.Getenv("NATS_LOGIN")
	Settings.NATS_PASSWORD = os.Getenv("NATS_PASSWORD")

	if Settings.NATS_SERVER == "" {
		log.Panicln("Need fill NATS_SERVER ! in os.ENV ")
	}

	if Settings.NATS_PORT == "" {
		log.Panicln("Need fill NATS_PORT ! in os.ENV ")
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
