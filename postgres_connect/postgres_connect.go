// модуль для работы с базой данных

package postgres_connect

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ManyakRus/starter/logger"
	"github.com/ManyakRus/starter/ping"
	"os"
	"sync"
	"time"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"
)

// Conn - соединение к базе данных
var Conn *sqlx.DB

// log - глобальный логгер
var log = logger.GetLog()

// mutexReconnect - защита от многопоточности Reconnect()
var mutexReconnect = &sync.Mutex{}

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// NeedReconnect - флаг необходимости переподключения
var NeedReconnect bool

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_SCHEMA   string
	DB_USER     string
	DB_PASSWORD string
}

// Connect_err - подключается к базе данных
func Connect() {

	if Settings.DB_HOST == "" {
		FillSettings()
	}

	ping.Ping(Settings.DB_HOST, Settings.DB_PORT)

	err := Connect_err()
	if err != nil {
		log.Panicln("POSTGRES sqlx Connect_err() host: ", Settings.DB_HOST, ", Error: ", err)
	} else {
		log.Info("POSTGRES sqlx connected. host: ", Settings.DB_HOST, ", base name: ", Settings.DB_NAME, ", schema: ", Settings.DB_SCHEMA)
	}

}

// Connect_err - подключается к базе данных
func Connect_err() error {
	var err error

	if Settings.DB_HOST == "" {
		FillSettings()
	}

	//ctxMain := context.Background()
	//ctxMain := contextmain.GetContext()
	//ctx, cancel := context.WithTimeout(ctxMain, 5*time.Second)
	//defer cancel()

	// get the database connection URL.
	databaseUrl := "postgres://" + Settings.DB_USER + ":" + Settings.DB_PASSWORD
	databaseUrl += "@" + Settings.DB_HOST + ":5432/" + Settings.DB_NAME + "?sslmode=disable"

	//
	//Conn, err = pgx.Connect(ctx, databaseUrl)
	Conn, err = sqlx.Connect(
		"postgres",
		databaseUrl,
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

	//PgError, ok := err.(*pgconn.PgError)
	//if ok {
	//	if PgError.Code == "P0001" { // Class P0 — PL/pgSQL Error, RaiseException
	//		return //нужен
	//	}
	//}

	//остановим программу т.к. она не должна работать при неработающеё БД
	log.Error("STOP app. Error: ", err)
	stopapp.StopApp()

}

// CloseConnection - закрытие соединения с базой данных
func CloseConnection() error {
	if Conn == nil {
		return nil
	}

	err := CloseConnection_err()
	if err != nil {
		log.Error("Postgres sqlx CloseConnection() error: ", err)
	} else {
		log.Info("Postgres sqlx connection closed")
	}

	return err
}

// CloseConnection - закрытие соединения с базой данных
func CloseConnection_err() error {
	if Conn == nil {
		return nil
	}

	//ctx := contextmain.GetContext()
	//ctx := context.Background()
	err := Conn.Close()

	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled.")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("Postgres sqlx")

	//
	err := CloseConnection()
	if err != nil {
		log.Error("CloseConnection() error: ", err)
	}
	stopapp.GetWaitGroup_Main().Done()
}

// StartDB - делает соединение с БД, отключение и др.
func StartDB() {
	Connect()

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	stopapp.GetWaitGroup_Main().Add(1)
	go ping_go()

}

// FillSettings загружает переменные окружения в структуру из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.DB_HOST = os.Getenv("DB_HOST")
	Settings.DB_PORT = os.Getenv("DB_PORT")
	Settings.DB_NAME = os.Getenv("DB_NAME")
	Settings.DB_SCHEMA = os.Getenv("DB_SCHEME")
	Settings.DB_USER = os.Getenv("DB_USER")
	Settings.DB_PASSWORD = os.Getenv("DB_PASSWORD")

	if Settings.DB_HOST == "" {
		log.Panicln("Need fill DB_HOST ! in os.ENV ")
	}

	if Settings.DB_PORT == "" {
		log.Panicln("Need fill DB_PORT ! in os.ENV ")
	}

	if Settings.DB_NAME == "" {
		log.Panicln("Need fill DB_NAME ! in os.ENV ")
	}

	if Settings.DB_SCHEMA == "" {
		log.Panicln("Need fill DB_SCHEME ! in os.ENV ")
	}
	if Settings.DB_USER == "" {
		log.Panicln("Need fill DB_USER ! in os.ENV ")
	}

	if Settings.DB_PASSWORD == "" {
		log.Panicln("Need fill DB_PASSWORD ! in os.ENV ")
	}

	//
}

// ping_go - делает пинг каждые 60 секунд, и реконнект
func ping_go() {

	ticker := time.NewTicker(60 * time.Second)

	addr := Settings.DB_HOST + ":" + Settings.DB_PORT

	//бесконечный цикл
loop:
	for {
		select {
		case <-contextmain.GetContext().Done():
			log.Warn("Context app is canceled. postgres_connect.ping")
			break loop
		case <-ticker.C:
			err := ping.Ping_err(Settings.DB_HOST, Settings.DB_PORT)
			//log.Debug("ticker, ping err: ", err) //удалить
			if err != nil {
				NeedReconnect = true
				log.Warn("postgres_connect Ping(", addr, ") error: ", err)
			} else if NeedReconnect == true {
				log.Warn("postgres_connect Ping(", addr, ") OK. Start Reconnect()")
				NeedReconnect = false
				Connect()
			}
		}
	}

	stopapp.GetWaitGroup_Main().Done()
}
