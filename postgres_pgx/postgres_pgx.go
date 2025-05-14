// модуль для работы с базой данных

package postgres_pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/log"
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

// mutex_Connect - защита от многопоточности Connect()
var mutex_Connect = &sync.RWMutex{}

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

// TextConnBusy - текст ошибки "conn busy"
const TextConnBusy = "conn busy"

// timeOutSeconds - время ожидания для Ping()
const timeOutSeconds = 60

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

	//
	if contextmain.GetContext().Err() != nil {
		return contextmain.GetContext().Err()
	}

	//
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
	err := GetConnection().Ping(ctx)
	if err != nil {
		otvet = true
	}
	return otvet
}

// Reconnect повторное подключение к базе данных, если оно отключено
// или полная остановка программы
func Reconnect(err error) {
	mutex_Connect.Lock()
	defer mutex_Connect.Unlock()

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

	mutex_Connect.Lock()
	defer mutex_Connect.Unlock()

	err := GetConnection().Close(ctx)

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
	err := Connect_WithApplicationName_err(ApplicationName)
	LogInfo_Connected(err)

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
	var err error

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	addr := Settings.DB_HOST + ":" + Settings.DB_PORT

	var ctx context.Context
	//бесконечный цикл
loop:
	for {
		ctx = contextmain.GetContext()

		select {
		case <-ctx.Done():
			log.Warn("Context app is canceled. postgres_pgx.ping")
			break loop
		case <-ticker.C:

			//ping в базе данных
			mutex_Connect.RLock() //race
			//err = GetConnection().Ping(ctx) //ping делать нельзя т.к. data race
			err = Ping_err(ctx)
			mutex_Connect.RUnlock()
			if err != nil {
				switch err.Error() {
				case TextConnBusy:
					{
						log.Warn("postgres_pgx Ping() warning: ", err)
					}
				default:
					{
						NeedReconnect = true
						log.Error("postgres_pgx Ping() error: ", err)
					}
				}

			} else {
				//IsClosed
				if GetConnection().IsClosed() == true {
					NeedReconnect = true
					log.Error("postgres_pgx error: IsClosed() = true")
				}
			}

			//ping порта
			err = port_checker.CheckPort_err(Settings.DB_HOST, Settings.DB_PORT)
			if err != nil {
				NeedReconnect = true
				log.Warn("postgres_pgx CheckPort(", addr, ") error: ", err)
			} else if NeedReconnect == true {
				log.Warn("postgres_pgx CheckPort(", addr, ") OK. Start Reconnect()")
				NeedReconnect = false
				err = Connect_err()
				if err != nil {
					NeedReconnect = true
					log.Error("Connect_err() error: ", err)
				}
			}
		}
	}

	stopapp.GetWaitGroup_Main().Done()
}

// GetConnection - возвращает соединение к нужной базе данных
func GetConnection() *pgx.Conn {
	//мьютекс чтоб не подключаться одновременно
	mutex_Connect.RLock()
	defer mutex_Connect.RUnlock()

	//
	if Conn == nil || Conn.IsClosed() {
		err := Connect_err()
		if err != nil {
			log.Error("POSTGRES pgx Connect() to database host: ", Settings.DB_HOST, ", error: ", err)
		} else {
			log.Info("POSTGRES pgx Connected. host: ", Settings.DB_HOST, ", base name: ", Settings.DB_NAME, ", schema: ", Settings.DB_SCHEMA)
		}
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
// после вызова, в конце необходимо закрыть rows!
// if err != nil {
// }
// defer rows.Close()

func RawMultipleSQL(tx pgx.Tx, TextSQL string) (pgx.Rows, error) {
	var rows pgx.Rows
	var err error

	if tx == nil {
		TextError := "RawMultipleSQL() error: tx =nil"
		log.Error(TextError)
		err = errors.New(TextError)
		return rows, err
	}

	//if tx.IsClosed() {
	//	TextError := "RawMultipleSQL() error: tx is closed"
	//	log.Error(TextError)
	//	err = errors.New(TextError)
	//	return rows, err
	//}

	ctx := contextmain.GetContext()

	//запустим транзакцию
	//tx, err := tx.Begin(ctx)
	//if err != nil {
	//	log.Error(err)
	//	return rows, err
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
			return rows, err
		}
	}

	//запустим последний запрос, с возвратом результата
	rows, err = tx.Query(ctx, TextSQL2)
	if err != nil {
		TextError := fmt.Sprint("tx.Raw() error: ", err, ", TextSQL: \n", TextSQL2)
		err = errors.New(TextError)
		return rows, err
	}
	//defer rows.Close()

	return rows, err
}

// Ping_err - выполняет пустой запрос для теста соединения
func Ping_err(ctxMain context.Context) error {
	var err error

	ctx, cancelFunc := context.WithTimeout(ctxMain, timeOutSeconds*time.Second)
	defer cancelFunc()

	_, err = GetConnection().Exec(ctx, ";")
	return err
}
