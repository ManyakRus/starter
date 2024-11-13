// модуль для работы с базой данных

package postgres_pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/logger"
	"github.com/ManyakRus/starter/port_checker"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"

	//"github.com/jackc/pgconn"
	"os"
	"sync"
	//"time"

	//_ "github.com/jackc/pgconn"
	//_ "github.com/jackc/pgx/v4"
	//"github.com/jmoiron/sqlx"
	//_ "github.com/lib/pq"
	//log "github.com/sirupsen/logrus"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"
)

// Conn - соединение к базе данных
var Conn *pgx.Conn

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
	LogInfo_Connected(err)

}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicln("POSTGRES pgx Connect() to database host: ", Settings.DB_HOST, ", Error: ", err)
	} else {
		log.Info("POSTGRES pgx Connected. host: ", Settings.DB_HOST, ", base name: ", Settings.DB_NAME, ", schema: ", Settings.DB_SCHEMA)
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
	ctxMain := contextmain.GetContext()
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()

	// get the database connection URL.
	//databaseUrl := "postgres://" + Settings.DB_USER + ":" + Settings.DB_PASSWORD
	//databaseUrl += "@" + Settings.DB_HOST + ":" + Settings.DB_PORT + "/" + Settings.DB_NAME

	databaseUrl := GetConnectionString(ApplicationName)

	//
	config, err := pgx.ParseConfig(databaseUrl)
	//config.PreferSimpleProtocol = true //для мульти-запросов
	Conn = nil
	Conn, err = pgx.ConnectConfig(ctx, config)

	if err != nil {
		err = fmt.Errorf("ConnectConfig() error: %w", err)
		//log.Panicln("Unable to connect to database host: ", Settings.DB_HOST, " Error: ", err)
	}

	if err == nil {
		err = Conn.Ping(ctx)
	}

	return err
}

// GetConnectionString - возвращает строку соединения к базе данных
func GetConnectionString(ApplicationName string) string {
	ApplicationName = strings.ReplaceAll(ApplicationName, " ", "_")

	dsn := "host=" + Settings.DB_HOST + " "
	dsn += "user=" + Settings.DB_USER + " "
	dsn += "password=" + Settings.DB_PASSWORD + " "
	dsn += "dbname=" + Settings.DB_NAME + " "
	dsn += "port=" + Settings.DB_PORT + " sslmode=disable TimeZone=UTC "
	dsn += "application_name=" + ApplicationName

	return dsn
}

// IsClosed проверка что база данных закрыта
func IsClosed() bool {
	var otvet bool
	if Conn == nil {
		return true
	}

	ctx := contextmain.GetContext()
	err := Conn.Ping(ctx)
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
		log.Error("Postgres pgx CloseConnection() error: ", err)
	} else {
		log.Info("Postgres pgx connection closed")
	}

}

// CloseConnection - закрытие соединения с базой данных
func CloseConnection_err() error {
	if Conn == nil {
		return nil
	}

	ctxMain := contextmain.GetContext()
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()
	err := Conn.Close(ctx)

	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled.")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("Postgres pgx")

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

	stopapp.GetWaitGroup_Main().Add(1)
	go ping_go()

	return err
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

// GetConnection - возвращает соединение к нужной базе данных
func GetConnection() *pgx.Conn {
	if Conn == nil || Conn.IsClosed() {
		Connect()
	}

	return Conn
}

// GetConnection_WithApplicationName - возвращает соединение к нужной базе данных, с указанием имени приложения
func GetConnection_WithApplicationName(ApplicationName string) *pgx.Conn {
	if Conn == nil {
		Connect_WithApplicationName_err(ApplicationName)
	}

	return Conn
}

// RawMultipleSQL - выполняет текст запроса, отдельно для каждого запроса
func RawMultipleSQL(tx pgx.Tx, TextSQL string) (pgx.Rows, error) {
	var Rows pgx.Rows
	var err error

	if tx == nil {
		TextError := "RawMultipleSQL() error: tx =nil"
		log.Error(TextError)
		err = errors.New(TextError)
		return Rows, err
	}

	//if tx.IsClosed() {
	//	TextError := "RawMultipleSQL() error: tx is closed"
	//	log.Error(TextError)
	//	err = errors.New(TextError)
	//	return Rows, err
	//}

	ctx := contextmain.GetContext()

	//запустим транзакцию
	//tx, err := tx.Begin(ctx)
	//if err != nil {
	//	log.Error(err)
	//	return Rows, err
	//}
	//defer tx.Commit()

	//
	TextSQL1 := ""
	TextSQL2 := TextSQL

	//запустим все запросы, кроме последнего
	pos1 := strings.LastIndex(TextSQL, ";")
	if pos1 > 0 {
		TextSQL1 = TextSQL[0:pos1]
		TextSQL2 = TextSQL[pos1:]
		_, err := tx.Exec(ctx, TextSQL1)
		if err != nil {
			TextError := fmt.Sprint("tx.Exec() error: ", err, ", TextSQL: \n", TextSQL1)
			err = errors.New(TextError)
			log.Error(err)
			return Rows, err
		}
	}

	//запустим последний запрос, с возвратом результата
	Rows, err = tx.Query(ctx, TextSQL2)
	if err != nil {
		TextError := fmt.Sprint("tx.Raw() error: ", err, ", TextSQL: \n", TextSQL2)
		err = errors.New(TextError)
		return Rows, err
	}

	return Rows, err
}
