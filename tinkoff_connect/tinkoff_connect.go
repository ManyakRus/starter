package tinkoff_connect

import (
	"context"
	"errors"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/microl"
	"github.com/ManyakRus/starter/port_checker"
	"github.com/ManyakRus/starter/stopapp"
	"github.com/tinkoff/invest-api-go-sdk/investgo"
	"os"
	"sync"
	"time"
)

// PackageName - имя текущего пакета, для логирования
const PackageName = "tinkoff_connect"

// SettingsINI - тип для хранения настроек подключени
type SettingsINI struct {
	investgo.Config
	Host string
	Port string
}

// Settings - структура для хранения настроек подключения
var Settings SettingsINI

// Conn - подключение к серверу Tinkoff-GRPC
var Client *investgo.Client

// mutex_Connect - защита от многопоточности Reconnect()
var mutex_Connect = &sync.RWMutex{}

// NeedReconnect - флаг необходимости переподключения
var NeedReconnect bool

// timeoutSeconds - время ожидания запроса в Тинькофф, в секундах
var timeoutSeconds = 60

// GetConnection - возвращает соединение
func GetConnection() *investgo.Client {
	//мьютекс чтоб не подключаться одновременно
	mutex_Connect.RLock()
	defer mutex_Connect.RUnlock()

	//
	if Client == nil {
		err := Connect_err()
		LogInfo_Connected(err)
	}

	return Client
}

// Connect - подключается к серверу Tinkoff-GRPC, при ошибке вызывает панику
func Connect() {
	var err error

	err = Connect_err()

	LogInfo_Connected_Panic(err)

}

// LogInfo_Connected - выводит сообщение в Лог
func LogInfo_Connected(err error) {
	if err != nil {
		log.Errorf("Tinkoff connection ERROR. EndPoint: %s, AccountId: %s, error: %s", Settings.EndPoint, Settings.AccountId, err)
	} else {
		log.Infof("Tinkoff connection OK. EndPoint: %s, AccountId: %s", Settings.EndPoint, Settings.AccountId)
	}

}

// LogInfo_Connected_Panic - выводит сообщение в Лог или панику
func LogInfo_Connected_Panic(err error) {
	LogInfo_Connected(err)

	if err != nil {
		panic(err)
	}
}

// Connect_err - подключается к серверу Tinkoff-GRPC, возвращает ошибку
func Connect_err() error {
	var err error

	//мьютекс чтоб не подключаться одновременно
	mutex_Connect.Lock()
	defer mutex_Connect.Unlock()

	//
	if Settings.EndPoint == "" {
		err = FillSettings()
		if err != nil {
			return err
		}
	}

	ctx := ctx_Connect

	//addr := Settings.Host + ":" + Settings.Port
	Config := Settings.Config
	Client, err = investgo.NewClient(ctx, Config, log.GetLog())
	if err != nil {
		return err
	}

	return err
}

func FillSettings() error {
	var err error

	Settings = SettingsINI{}
	INVEST_HOST := os.Getenv("INVEST_HOST")
	INVEST_PORT := os.Getenv("INVEST_PORT")

	Settings.Host = INVEST_HOST
	Settings.Port = INVEST_PORT
	EndPoint := INVEST_HOST + ":" + INVEST_PORT
	Settings.EndPoint = EndPoint

	if INVEST_HOST == "" {
		TextError := "Need fill INVEST_HOST ! in OS Environment "
		err = errors.New(TextError)
		return err
	}

	if INVEST_PORT == "" {
		TextError := "Need fill INVEST_PORT ! in OS Environment "
		err = errors.New(TextError)
		return err
	}

	Name := ""
	s := ""

	//
	Name = "INVEST_TOKEN"
	s = microl.Getenv(Name, true)
	Settings.Token = s

	//
	Name = "INVEST_ACCOUNTID"
	s = microl.Getenv(Name, false)
	Settings.AccountId = s

	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {
	defer waitGroup_Connect.Done()

	select {
	case <-ctx_Connect.Done():
		log.Warn("Context app is canceled. tinkoff_connect")
	}

	// ждём пока отправляемых сейчас сообщений будет =0
	stopapp.WaitTotalMessagesSendingNow("tinkoff_connect")

	// закрываем соединение
	CloseConnection()
}

// Start - необходимые процедуры для запуска сервера Tinkoff-GRPC
// если контекст хранится в ctx_Connect
// и есть waitGroup_Connect
// при ошибке вызывает панику
func Start() {
	Start_ctx(&ctx_Connect, waitGroup_Connect)

	//waitGroup_Connect.Add(1)
	//go WaitStop()
	//
	//waitGroup_Connect.Add(1)
	//go ping_go()

}

// Start_ctx - необходимые процедуры для запуска сервера Tinkoff-GRPC
// ctx - глобальный контекст приложения
// wg - глобальный WaitGroup приложения
func Start_ctx(ctx *context.Context, wg *sync.WaitGroup) error {
	var err error
	//	if contextmain.Ctx != ctx {
	//		contextmain.SetContext(ctx)
	//	}
	//contextmain.Ctx = ctx
	stopapp.SetWaitGroup_Main(wg)

	err = Connect_err()
	if err != nil {
		return err
	}

	//сохраним в список подключений
	WaitGroupContext1 := stopapp.WaitGroupContext{WaitGroup: waitGroup_Connect, Ctx: ctx, CancelCtxFunc: cancelCtxFunc}
	stopapp.OrderedMapConnections.Put(PackageName, WaitGroupContext1)

	//
	waitGroup_Connect.Add(1)
	go WaitStop()

	waitGroup_Connect.Add(1)
	go ping_go()

	return err
}

// CloseConnection - закрывает подключение к Tinkoff-GRPC, и пишет лог
func CloseConnection() {
	err := CloseConnection_err()
	if err != nil {
		log.Error("tinkoff_connect client CloseConnection() error: ", err)
	} else {
		log.Info("tinkoff_connect client connection closed")
	}
}

// CloseConnection - закрывает подключение к Tinkoff-GRPC, и возвращает ошибку
func CloseConnection_err() error {
	err := Client.Stop()
	return err
}

// ping_go - делает пинг каждые 60 секунд, и реконнект
func ping_go() {
	var err error

	defer waitGroup_Connect.Done()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	addr := Settings.Host + ":" + Settings.Port

	//бесконечный цикл
loop:
	for {
		select {
		case <-ctx_Connect.Done():
			log.Warn("Context app is canceled. tinkoff_connect.ping")
			break loop
		case <-ticker.C:
			err = port_checker.CheckPort_err(Settings.Host, Settings.Port)
			//log.Debug("ticker, ping err: ", err) //удалить
			if err != nil {
				NeedReconnect = true
				log.Warn("tinkoff_connect CheckPort(", addr, ") error: ", err)
			} else if NeedReconnect == true {
				log.Warn("tinkoff_connect CheckPort(", addr, ") OK. Start Reconnect()")
				NeedReconnect = false
				err = Connect_err()
				LogInfo_Connected(err)
				if err != nil {
					NeedReconnect = true
					//log.Error("tinkoff_connect Connect() error: ", err)
				}
			}
		}
	}

}

// GetTimeoutSeconds - возвращает время ожидания ответа
func GetTimeoutSeconds() int {
	Otvet := timeoutSeconds

	return Otvet
}

// SetTimeoutSeconds - устанавливает время ожидания ответа
func SetTimeoutSeconds(seconds int) {
	timeoutSeconds = seconds
}
