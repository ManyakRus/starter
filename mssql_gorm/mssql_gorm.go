// модуль для работы с базой данных

package mssql_gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/logger"
	"net/url"
	"os"
	"sync"
	//_ "github.com/denisenkom/go-mssqldb"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Conn - соединение к базе данных
var Conn *gorm.DB

// log - глобальный логгер
var log = logger.GetLog()

// mutexReconnect - защита от многопоточности Reconnect()
var mutexReconnect = &sync.Mutex{}

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	MSSQL_ADDRESS  string
	MSSQL_BASENAME string
	MSSQL_SCHEMA   string
	MSSQL_PORT     string
	MSSQL_LOGIN    string
	MSSQL_PASSWORD string
}

// Connect - подключается к базе данных
func Connect() {
	var err error

	err = Connect_err()
	if err != nil {
		log.Panicln("MSSQL GORM unable connect, host: ", Settings.MSSQL_ADDRESS, ", Error: ", err)
	} else {
		log.Info("MSSQL GORM connected, host: ", Settings.MSSQL_ADDRESS, ", port: ", Settings.MSSQL_PORT)
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

	conf := &gorm.Config{}

	dsn := GetDSN()
	conn := sqlserver.Open(dsn)
	Conn, err = gorm.Open(conn, conf)
	Conn.Config.NamingStrategy = schema.NamingStrategy{TablePrefix: Settings.MSSQL_SCHEMA + "."}
	Conn.Config.Logger = gormlogger.Default.LogMode(gormlogger.Warn)

	if err == nil {
		DB, err := Conn.DB()
		if err != nil {
			log.Error("Conn.DB() error: ", err)
			return err
		}

		err = DB.Ping()
	}

	return err
}

// IsClosed проверка что база данных закрыта
func IsClosed() bool {
	var otvet bool
	if Conn == nil {
		return true
	}

	DB, err := Conn.DB()
	if err != nil {
		log.Error("Conn.DB() error: ", err)
		return true
	}

	err = DB.Ping()
	if err != nil {
		log.Error("DB.Ping() error: ", err)
		return true
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
		log.Error("MSSQL gorm CloseConnection() error: ", err)
	} else {
		log.Info("MSSQL gorm stopped")
	}
}

// CloseConnection_err - закрытие соединения с базой данных
func CloseConnection_err() error {
	if Conn == nil {
		return nil
	}

	DB, err := Conn.DB()
	if err != nil {
		log.Error("Conn.DB() error: ", err)
		return err
	}
	err = DB.Close()
	if err != nil {
		log.Error("DB.Close() error: ", err)
	}
	Conn = nil

	return err
}

// WaitStop - ожидает отмену глобального контекста или сигнала завершения приложения
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled.")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("MSSQL gorm")

	//
	CloseConnection()

	stopapp.GetWaitGroup_Main().Done()
}

// StartDB - делает соединение с БД, отключение и др.
func StartDB() {
	Connect()

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

}

// FillSettings загружает переменные окружения в структуру из файла или из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.MSSQL_LOGIN = os.Getenv("MSSQL_LOGIN")
	Settings.MSSQL_PASSWORD = os.Getenv("MSSQL_PASSWORD")
	Settings.MSSQL_ADDRESS = os.Getenv("MSSQL_ADDRESS")
	Settings.MSSQL_BASENAME = os.Getenv("MSSQL_BASENAME")
	Settings.MSSQL_PORT = os.Getenv("MSSQL_PORT")
	Settings.MSSQL_SCHEMA = os.Getenv("MSSQL_SCHEMA")
	if Settings.MSSQL_LOGIN == "" {
		log.Panicln("Need fill MSSQL_LOGIN ! in os.ENV ")
	}

	if Settings.MSSQL_PASSWORD == "" {
		log.Panicln("Need fill MSSQL_PASSWORD ! in os.ENV ")
	}

	if Settings.MSSQL_ADDRESS == "" {
		log.Panicln("Need fill MSSQL_ADDRESS ! in os.ENV ")
	}

	if Settings.MSSQL_BASENAME == "" {
		log.Panicln("Need fill MSSQL_BASENAME ! in os.ENV ")
	}

	if Settings.MSSQL_PORT == "" {
		log.Panicln("Need fill MSSQL_PORT ! in os.ENV ")
	}

	if Settings.MSSQL_SCHEMA == "" {
		log.Panicln("Need fill MSSQL_SCHEMA ! in os.ENV ")
	}
	//
}

// GetConnection - возвращает соединение к нужной базе данных
func GetConnection(connection_id int64) *gorm.DB {
	if connection_id == 0 {
		log.Panicln("mssql_gorm.GetConnection() error: connection_id =0")
	}

	if Conn == nil {
		Connect()
	}

	return Conn
}

// GetDSN - возвращает строку соединения к базе данных
func GetDSN() string {
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

	//loc, err := time.LoadLocation("Russian Standard Time")
	//if err != nil {
	//	log.Panicln("time.LoadLocation() error: ", err)
	//}

	ConnectionString = ConnectionString + `&parseTime=True&loc="Russian Standard Time"`
	//ConnectionString = ConnectionString + "&loc=Local"

	return ConnectionString
}
