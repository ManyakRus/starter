// модуль для работы с базой данных

package mssql_connect

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/url"
	"os"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"
)

// Conn - соединение к базе данных
var Conn *sqlx.DB

// log - глобальный логгер
//var log = logger.GetLog()

// mutexReconnect - защита от многопоточности Reconnect()
var mutexReconnect = &sync.Mutex{}

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	MSSQL_LOGIN    string
	MSSQL_PASSWORD string
	MSSQL_ADDRESS  string
	MSSQL_BASENAME string
	MSSQL_PORT     string
}

// Connect - подключается к базе данных
func Connect() {
	var err error

	err = Connect_err()
	LogInfo_Connected(err)

}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicln("MSSQL unable connect, host: ", Settings.MSSQL_ADDRESS, ", Error: ", err)
	} else {
		log.Info("MSSQL connected, host: ", Settings.MSSQL_ADDRESS, ", port: ", Settings.MSSQL_PORT)
	}

}

// Connect_err - подключается к базе данных и возвращает ошибку
func Connect_err() error {
	var err error

	if Settings.MSSQL_ADDRESS == "" {
		FillSettings()
	}

	//ctxMain := context.Background()
	////ctxMain := contextmain.GetContext()
	//ctx, cancel := context.WithTimeout(ctxMain, 5*time.Second)
	//defer cancel()

	query := url.Values{}
	query.Add("database", Settings.MSSQL_BASENAME)

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(Settings.MSSQL_LOGIN, Settings.MSSQL_PASSWORD),
		Host:   fmt.Sprintf("%s:%s", Settings.MSSQL_ADDRESS, Settings.MSSQL_PORT),
		// Path:  instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	ConnectionString := u.String()
	//Conn, err = sql.Open("sqlserver", ConnectionString)

	Conn, err = sqlx.Connect(
		"mssql",
		ConnectionString,
	)

	if err == nil {
		err = Conn.Ping()
	}

	return err
}

// IsClosed проверка что база данных закрыта
func IsClosed() bool {
	var otvet bool
	if Conn == nil {
		return true
	}

	ctx := contextmain.GetContext()
	err := Conn.PingContext(ctx)
	if err != nil {
		otvet = true
	}
	return otvet
}

// Reconnect повторное подключение к базе данных, если оно отключено
// или полная остановка программы
func Reconnect(err error) {
	mutexReconnect.Lock()
	defer mutexReconnect.Unlock()

	if err == nil {
		return
	}

	if errors.Is(err, context.Canceled) {
		return
	}

	if Conn == nil {
		log.Warn("Reconnect()")
		err := Connect_err()
		if err != nil {
			log.Error("error: ", err)
		}
		return
	}

	if IsClosed() {
		micro.Pause(1000)
		log.Warn("Reconnect()")
		err := Connect_err()
		if err != nil {
			log.Error("error: ", err)
		}
		return
	}

	sError := err.Error()
	if sError == "Conn closed" {
		micro.Pause(1000)
		log.Warn("Reconnect()")
		err := Connect_err()
		if err != nil {
			log.Error("error: ", err)
		}
		return
	}

	//остановим программу т.к. она не должна работать при неработающеё БД
	//log.Error("STOP app. Error: ", err)
	//stopapp.StopApp()

}

// CloseConnection - закрытие соединения с базой данных
func CloseConnection() {
	if Conn == nil {
		return
	}

	err := CloseConnection_err()
	if err != nil {
		log.Error("MSSQL CloseConnection() error: ", err)
	} else {
		log.Info("MSSQL stopped")
	}
}

// CloseConnection_err - закрытие соединения с базой данных
func CloseConnection_err() error {
	if Conn == nil {
		return nil
	}

	err := Conn.Close()
	Conn = nil

	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled.")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("MSSQL sqlx")

	//
	CloseConnection()

	stopapp.GetWaitGroup_Main().Done()
}

// StartDB - делает соединение с БД, отключение и др.
func StartDB() {
	var err error

	ctx := contextmain.GetContext()
	WaitGroup := stopapp.GetWaitGroup_Main()
	err = Start_ctx(&ctx, WaitGroup)
	LogInfo_Connected(err)

}

// Start_ctx - необходимые процедуры для подключения к серверу БД
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто StartDB()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup) error {
	var err error

	//запомним к себе контекст
	contextmain.Ctx = ctx
	if ctx == nil {
		contextmain.GetContext()
	}

	//запомним к себе WaitGroup
	stopapp.SetWaitGroup_Main(WaitGroup)
	if WaitGroup == nil {
		stopapp.StartWaitStop()
	}

	//
	err = Connect_err()
	if err != nil {
		return err
	}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	return err
}

// FillSettings загружает переменные окружения в структуру из файла или из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.MSSQL_LOGIN = os.Getenv("MSSQL_LOGIN")
	Settings.MSSQL_PASSWORD = os.Getenv("MSSQL_PASSWORD")
	Settings.MSSQL_ADDRESS = os.Getenv("MSSQL_ADDRESS")
	Settings.MSSQL_BASENAME = os.Getenv("MSSQL_BASENAME")
	Settings.MSSQL_PORT = os.Getenv("MSSQL_PORT")

	if Settings.MSSQL_ADDRESS == "" {
		log.Panicln("Need fill MSSQL_ADDRESS ! in os.ENV ")
	}

	if Settings.MSSQL_BASENAME == "" {
		log.Panicln("Need fill MSSQL_BASENAME ! in os.ENV ")
	}

	if Settings.MSSQL_PORT == "" {
		log.Panicln("Need fill MSSQL_PORT ! in os.ENV ")
	}

	if Settings.MSSQL_LOGIN == "" {
		log.Panicln("Need fill MSSQL_LOGIN ! in os.ENV ")
	}

	if Settings.MSSQL_PASSWORD == "" {
		log.Panicln("Need fill MSSQL_PASSWORD ! in os.ENV ")
	}

	//
}

func GetConnection(connection_id int) *sqlx.DB {
	if Conn == nil {
		Connect()
	}

	return Conn
}
