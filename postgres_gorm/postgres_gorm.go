// модуль для работы с базой данных

package postgres_gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/logger"
	"github.com/ManyakRus/starter/port_checker"
	"strings"
	"time"

	//"github.com/jackc/pgconn"
	"os"
	"sync"
	//"time"

	//"github.com/jmoiron/sqlx"
	//_ "github.com/lib/pq"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"

	"gorm.io/driver/postgres"
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

	port_checker.CheckPort(Settings.DB_HOST, Settings.DB_PORT)

	err := Connect_err()
	if err != nil {
		log.Panicln("POSTGRES gorm Connect() to database host: ", Settings.DB_HOST, ", Error: ", err)
	} else {
		log.Info("POSTGRES gorm Connected. host: ", Settings.DB_HOST, ", base name: ", Settings.DB_NAME, ", schema: ", Settings.DB_SCHEMA)
	}

}

// Connect_err - подключается к базе данных
func Connect_err() error {
	var err error
	err = Connect_WithApplicationName_err("")

	return err
}

// Connect_WithApplicationName_err - подключается к базе данных, с указанием имени приложения
func Connect_WithApplicationName_err(ApplicationName string) error {
	var err error

	if Settings.DB_HOST == "" {
		FillSettings()
	}

	//ctxMain := context.Background()
	//ctxMain := contextmain.GetContext()
	//ctx, cancel := context.WithTimeout(ctxMain, 5*time.Second)
	//defer cancel()

	// get the database connection URL.
	dsn := GetDSN(ApplicationName)

	//
	conf := &gorm.Config{}
	//conn := postgres.Open(dsn)

	dialect := postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, //для запуска мультизапросов
	})
	Conn, err = gorm.Open(dialect, conf)

	//Conn, err = gorm.Open(conn, conf)
	Conn.Config.NamingStrategy = schema.NamingStrategy{TablePrefix: Settings.DB_SCHEMA + "."}
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
		log.Error("DB.CheckPort() error: ", err)
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
func CloseConnection() {
	if Conn == nil {
		return
	}

	err := CloseConnection_err()
	if err != nil {
		log.Error("Postgres gorm CloseConnection() error: ", err)
	} else {
		log.Info("Postgres gorm connection closed")
	}

	return
}

// CloseConnection - закрытие соединения с базой данных
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

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. Postgres gorm.")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("Postgres gorm")

	//
	CloseConnection()

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

// Start - делает соединение с БД, отключение и др.
func Start(ApplicationName string) {
	Connect_WithApplicationName_err(ApplicationName)

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	stopapp.GetWaitGroup_Main().Add(1)
	go ping_go()

}

// FillSettings загружает переменные окружения в структуру из файла или из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	// заполним из переменных оуружения
	Settings.DB_HOST = os.Getenv("DB_HOST")
	Settings.DB_PORT = os.Getenv("DB_PORT")
	Settings.DB_NAME = os.Getenv("DB_NAME")
	Settings.DB_SCHEMA = os.Getenv("DB_SCHEME")
	Settings.DB_USER = os.Getenv("DB_USER")
	Settings.DB_PASSWORD = os.Getenv("DB_PASSWORD")

	//// заполним из переменных оуружения как у Нечаева
	//if Settings.DB_HOST == "" {
	//	Settings.DB_HOST = os.Getenv("STORE_HOST")
	//	Settings.DB_PORT = os.Getenv("STORE_PORT")
	//	Settings.DB_NAME = os.Getenv("STORE_NAME")
	//	Settings.DB_SCHEMA = os.Getenv("STORE_SCHEME")
	//	Settings.DB_USER = os.Getenv("STORE_LOGIN")
	//	Settings.DB_PASSWORD = os.Getenv("STORE_PASSWORD")
	//}

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
		log.Panicln("Need fill DB_SCHEMA ! in os.ENV ")
	}

	if Settings.DB_USER == "" {
		log.Panicln("Need fill DB_USER ! in os.ENV ")
	}

	if Settings.DB_PASSWORD == "" {
		log.Panicln("Need fill DB_PASSWORD ! in os.ENV ")
	}

	//
}

// GetDSN - возвращает строку соединения к базе данных
func GetDSN(ApplicationName string) string {
	ApplicationName = strings.ReplaceAll(ApplicationName, " ", "_")

	dsn := "host=" + Settings.DB_HOST + " "
	dsn += "user=" + Settings.DB_USER + " "
	dsn += "password=" + Settings.DB_PASSWORD + " "
	dsn += "dbname=" + Settings.DB_NAME + " "
	dsn += "port=" + Settings.DB_PORT + " sslmode=disable TimeZone=UTC "
	dsn += "application_name=" + ApplicationName

	return dsn
}

// GetConnection - возвращает соединение к нужной базе данных
func GetConnection() *gorm.DB {
	if Conn == nil {
		Connect()
	}

	return Conn
}

// GetConnection_WithApplicationName - возвращает соединение к нужной базе данных, с указанием имени приложения
func GetConnection_WithApplicationName(ApplicationName string) *gorm.DB {
	if Conn == nil {
		Connect_WithApplicationName_err(ApplicationName)
	}

	return Conn
}

// ping_go - делает пинг каждые 60 секунд, и реконнект
func ping_go() {

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	addr := Settings.DB_HOST + ":" + Settings.DB_PORT

	//бесконечный цикл
loop:
	for {
		select {
		case <-contextmain.GetContext().Done():
			log.Warn("Context app is canceled. postgres_gorm.ping")
			break loop
		case <-ticker.C:
			err := port_checker.CheckPort_err(Settings.DB_HOST, Settings.DB_PORT)
			//log.Debug("ticker, ping err: ", err) //удалить
			if err != nil {
				NeedReconnect = true
				log.Warn("postgres_gorm CheckPort(", addr, ") error: ", err)
			} else if NeedReconnect == true {
				log.Warn("postgres_gorm CheckPort(", addr, ") OK. Start Reconnect()")
				NeedReconnect = false
				Connect()
			}
		}
	}

	stopapp.GetWaitGroup_Main().Done()
}

// RawMultipleSQL - выполняет текст запроса, отдельно для каждого запроса
func RawMultipleSQL(db *gorm.DB, TextSQL string) *gorm.DB {
	var tx *gorm.DB
	var err error
	tx = db

	// запустим все запросы отдельно
	sqlSlice := strings.Split(TextSQL, ";")
	len1 := len(sqlSlice)
	for i, v := range sqlSlice {
		if i == len1-1 {
			tx = tx.Raw(v)
			err = tx.Error
		} else {
			tx = tx.Exec(v)
			err = tx.Error
		}
		if err != nil {
			TextError := fmt.Sprint("db.Raw() error: ", err, ", TextSQL: \n", v)
			err = errors.New(TextError)
			break
		}
	}

	if tx == nil {
		log.Panic("db.Raw() error: rows =nil")
	}

	return tx
}
