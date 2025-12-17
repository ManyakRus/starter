// модуль для работы с базой данных
// нужен для подключения к второй БД, одновременно с первой БД

package postgres_pgxpool2

import (
	"context"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/port_checker"
	"github.com/ManyakRus/starter/postgres_pgtype"
	"github.com/ManyakRus/starter/postgres_pgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"

	"github.com/ManyakRus/starter/constants_starter"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"
	"os"
	"sync"
)

// PackageName - имя текущего пакета, для логирования
const PackageName = "postgres_pgxpool2"

// PgxPool - пул соединений к базе данных
var PgxPool *pgxpool.Pool

//// mutex_Connect - защита от многопоточности Connect()
//var mutex_Connect = &sync.RWMutex{}

// mutex_ReConnect - защита от многопоточности ReConnect()
var mutex_ReConnect = &sync.RWMutex{}

// Settings хранит все нужные переменные окружения
var Settings = SettingsINI{}

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
	NoNUll      bool
}

// TextConnBusy - текст ошибки "conn busy"
const TextConnBusy = "conn busy"

// timeOutSeconds - время ожидания для Ping()
const timeOutSeconds = 1

// Connect - подключается к базе данных
func Connect() {

	if Settings.DB_HOST == "" {
		FillSettings()
	}

	port_checker.CheckPort(Settings.DB_HOST, Settings.DB_PORT)

	err := Connect_err()
	LogInfo_Connected(err)

}

// Connect_err - подключается к базе данных, возвращает ошибку
func Connect_err() error {
	var err error
	err = Connect_WithApplicationName_err("")

	return err
}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicln("POSTGRES pgxpool Connect() to database host: ", Settings.DB_HOST, ", Error: ", err)
	} else {
		log.Info("POSTGRES pgxpool Connected. host: ", Settings.DB_HOST, ", base name: ", Settings.DB_NAME, ", schema: ", Settings.DB_SCHEMA)
	}

}

// Connect_NoNull - подключается к базе данных, с указанием имени приложения
func Connect_NoNull(ApplicationName string) {
	Settings.NoNUll = true

	err := Connect_WithApplicationName_err(ApplicationName)
	LogInfo_Connected(err)
}

// Connect_NoNull_err - подключается к базе данных, с указанием имени приложения
func Connect_NoNull_err(ApplicationName string) error {
	Settings.NoNUll = true

	err := Connect_WithApplicationName_err(ApplicationName)
	return err
}

// Connect_WithApplicationName - подключается к базе данных, с указанием имени приложения
func Connect_WithApplicationName(ApplicationName string) {
	err := Connect_WithApplicationName_err(ApplicationName)
	LogInfo_Connected(err)
}

// Connect_WithApplicationName_err - подключается к базе данных, с указанием имени приложения
func Connect_WithApplicationName_err(ApplicationName string) error {
	var err error

	if Settings.DB_HOST == "" {
		FillSettings()
	}

	//
	if ctx_Connect.Err() != nil {
		return ctx_Connect.Err()
	}

	//
	ctxMain := ctx_Connect
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()

	databaseUrl := GetConnectionString(ApplicationName)

	//
	config, err := pgxpool.ParseConfig(databaseUrl)
	if Settings.NoNUll == true {
		config.AfterConnect = AfterConnect_NoNull
	}
	//config.PreferSimpleProtocol = true //для мульти-запросов
	PgxPool = nil
	PgxPool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		err = fmt.Errorf("pgxpool.NewWithConfig() error: %w", err)
	}

	if err == nil {
		err = PgxPool.Ping(ctx)
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
	dsn += "port=" + Settings.DB_PORT + " sslmode=disable TimeZone=" + constants_starter.TIME_ZONE + " "
	dsn += "application_name=" + ApplicationName + " "
	dsn += "search_path=" + Settings.DB_SCHEMA + " "

	return dsn
}

// IsClosed проверка что база данных закрыта
func IsClosed() bool {
	var Otvet bool

	if PgxPool == nil {
		return true
	}

	ctxMain := ctx_Connect
	ctx, cancelFunc := context.WithTimeout(ctxMain, timeOutSeconds*time.Second)
	defer cancelFunc()

	err := GetConnection().Ping(ctx)
	if err != nil {
		return true
	}

	return Otvet
}

// Reconnect повторное подключение к базе данных, если оно отключено
// или полная остановка программы
func Reconnect(err error) {
	mutex_ReConnect.Lock()
	defer mutex_ReConnect.Unlock()

	if err == nil {
		return
	}

	if errors.Is(err, context.Canceled) {
		return
	}

	if PgxPool == nil {
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
	if sError == "PgxPool closed" {
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
	if PgxPool == nil {
		return
	}

	err := CloseConnection_err()
	if err != nil {
		log.Error("POSTGRES pgxpool CloseConnection() error: ", err)
	} else {
		log.Info("POSTGRES pgxpool connection closed")
	}

}

// CloseConnection_err - закрытие соединения с базой данных
func CloseConnection_err() error {
	var err error

	if PgxPool == nil {
		return err
	}

	//ctxMain := ctx_Connect
	//ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	//defer cancel()

	//mutex_Connect.Lock()
	//defer mutex_Connect.Unlock()

	//PgxPool.Reset()
	PgxPool.Close()

	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {
	defer waitGroup_Connect.Done()

	select {
	case <-ctx_Connect.Done():
		log.Warn("Context app is canceled. postgres_pgxpool2")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("postgres_pgxpool2")

	//
	CloseConnection()
}

// StartDB - делает соединение с БД, отключение и др.
func StartDB() {
	var err error

	ctx := ctx_Connect
	WaitGroup := waitGroup_Connect
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
		ctx = &ctx_Connect
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

// Start_NoNull - делает соединение с БД, отключение и др.
// запросы вместо null возвращают значение по умолчанию (пока только дата)
func Start_NoNull(ApplicationName string) {
	Settings.NoNUll = true

	err := Connect_WithApplicationName_err(ApplicationName)
	LogInfo_Connected(err)

	waitGroup_Connect.Add(1)
	go WaitStop()

	waitGroup_Connect.Add(1)
	go ping_go()

}

// Start - делает соединение с БД, отключение и др.
func Start(ApplicationName string) {
	err := Connect_WithApplicationName_err(ApplicationName)
	LogInfo_Connected(err)

	waitGroup_Connect.Add(1)
	go WaitStop()

	waitGroup_Connect.Add(1)
	go ping_go()

}

// FillSettings загружает переменные окружения в структуру из файла или из переменных окружения
func FillSettings() {
	//Settings = SettingsINI{}
	Settings.DB_HOST = os.Getenv("DB_HOST2")
	Settings.DB_PORT = os.Getenv("DB_PORT2")
	Settings.DB_NAME = os.Getenv("DB_NAME2")
	Settings.DB_SCHEMA = os.Getenv("DB_SCHEME2")
	Settings.DB_USER = os.Getenv("DB_USER2")
	Settings.DB_PASSWORD = os.Getenv("DB_PASSWORD2")

	if Settings.DB_HOST == "" {
		log.Panicln("Need fill DB_HOST2 ! in os.ENV ")
	}

	if Settings.DB_PORT == "" {
		log.Panicln("Need fill DB_PORT2 ! in os.ENV ")
	}

	if Settings.DB_NAME == "" {
		log.Panicln("Need fill DB_NAME2 ! in os.ENV ")
	}

	if Settings.DB_SCHEMA == "" {
		log.Panicln("Need fill DB_SCHEME2 ! in os.ENV ")
	}
	if Settings.DB_USER == "" {
		log.Panicln("Need fill DB_USER2 ! in os.ENV ")
	}

	if Settings.DB_PASSWORD == "" {
		log.Panicln("Need fill DB_PASSWORD2 ! in os.ENV ")
	}

	//
}

// ping_go - делает пинг каждые 60 секунд, и реконнект
func ping_go() {
	var err error

	defer waitGroup_Connect.Done()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	addr := Settings.DB_HOST + ":" + Settings.DB_PORT

	var ctx context.Context
	//бесконечный цикл
loop:
	for {
		ctx = ctx_Connect

		select {
		case <-ctx.Done():
			log.Warn("Context app is canceled. postgres_pgxpool2.ping")
			break loop
		case <-ticker.C:

			//ping в базе данных
			//mutex_Connect.RLock() //race
			//err = GetConnection().Ping(ctx) //ping делать нельзя т.к. data race
			err = Ping_err(ctx)
			//mutex_Connect.RUnlock()
			if err != nil {
				switch err.Error() {
				case TextConnBusy:
					{
						log.Warn("postgres_pgxpool2 Ping() warning: ", err)
					}
				default:
					{
						NeedReconnect = true
						log.Error("postgres_pgxpool2 Ping() error: ", err)
					}
				}

			}

			//ping порта
			err = port_checker.CheckPort_err(Settings.DB_HOST, Settings.DB_PORT)
			if err != nil {
				NeedReconnect = true
				log.Warn("postgres_pgxpool2 CheckPort(", addr, ") error: ", err)
			} else if NeedReconnect == true {
				log.Warn("postgres_pgxpool2 CheckPort(", addr, ") OK. Start Reconnect()")
				NeedReconnect = false
				err = Connect_err()
				if err != nil {
					NeedReconnect = true
					log.Error("Connect_err() error: ", err)
				}
			}
		}
	}

}

// GetConnection - возвращает соединение к нужной базе данных
func GetConnection() *pgxpool.Pool {
	//мьютекс чтоб не подключаться одновременно
	//mutex_Connect.RLock()
	//defer mutex_Connect.RUnlock()

	//
	if PgxPool == nil {
		err := Connect_err()
		if err != nil {
			log.Error("POSTGRES pgxpool Connect_err() to database host: ", Settings.DB_HOST, ", error: ", err)
		} else {
			log.Info("POSTGRES pgxpool Connected. host: ", Settings.DB_HOST, ", base name: ", Settings.DB_NAME, ", schema: ", Settings.DB_SCHEMA)
		}
	}

	return PgxPool
}

// GetConnection_NoNull - возвращает соединение к нужной базе данных
func GetConnection_NoNull(ApplicationName string) *pgxpool.Pool {
	//мьютекс чтоб не подключаться одновременно
	//mutex_Connect.RLock()
	//defer mutex_Connect.RUnlock()

	//
	if PgxPool == nil {
		err := Connect_NoNull_err(ApplicationName)
		if err != nil {
			log.Error("POSTGRES pgxpool Connect_NoNull_err() to database host: ", Settings.DB_HOST, ", error: ", err)
		} else {
			log.Info("POSTGRES pgxpool Connected. host: ", Settings.DB_HOST, ", base name: ", Settings.DB_NAME, ", schema: ", Settings.DB_SCHEMA)
		}
	}

	return PgxPool
}

// GetConnection_WithApplicationName - возвращает соединение к нужной базе данных, с указанием имени приложения
func GetConnection_WithApplicationName(ApplicationName string) *pgxpool.Pool {
	if PgxPool == nil {
		err := Connect_WithApplicationName_err(ApplicationName)
		LogInfo_Connected(err)
	}

	return PgxPool
}

// RawMultipleSQL - выполняет текст запроса, отдельно для каждого запроса
// после вызова, в конце необходимо закрыть rows!
// if err != nil {
// }
// defer rows.Close()

func RawMultipleSQL(tx postgres_pgx.IConnectionTransaction, TextSQL string) (pgx.Rows, error) {
	var rows pgx.Rows
	var err error

	if tx == nil {
		TextError := "RawMultipleSQL() error: tx =nil"
		log.Error(TextError)
		err = errors.New(TextError)
		return rows, err
	}

	ctx := ctx_Connect

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
	//ctx, cancelFunc := context.WithTimeout(ctxMain, 1*time.Second)
	defer cancelFunc()

	//mutex_Connect.Lock() //убрал т.к. зависает всё
	//defer mutex_Connect.Unlock()

	_, err = GetConnection().Exec(ctx, ";")
	if err != nil {
		err = fmt.Errorf("Ping_err() Exec() error: %w", err)
		return err
	}
	return err
}

// ReplaceSchema - заменяет "public." на Settings.DB_SCHEMA
func ReplaceSchema(TextSQL string) string {
	Otvet := TextSQL

	if Settings.DB_SCHEMA == "" {
		return Otvet
	}

	Otvet = strings.ReplaceAll(Otvet, "\tpublic.", "\t"+Settings.DB_SCHEMA+".")
	Otvet = strings.ReplaceAll(Otvet, "\npublic.", "\n"+Settings.DB_SCHEMA+".")
	Otvet = strings.ReplaceAll(Otvet, " public.", " "+Settings.DB_SCHEMA+".")

	return Otvet
}

// ReplaceSchemaName - заменяет имя схемы в тексте SQL
func ReplaceSchemaName(TextSQL, SchemaNameFrom string) string {
	Otvet := TextSQL

	Otvet = strings.ReplaceAll(Otvet, SchemaNameFrom+".", Settings.DB_SCHEMA+".")

	return Otvet
}

// AfterConnect_NoNull - регистрирует обработчики для нужных типов
// чтобы NULL=default value
func AfterConnect_NoNull(ctx context.Context, conn *pgx.Conn) error {
	// Регистрируем обработчики для нужных типов,

	//timestamptz
	conn.TypeMap().RegisterType(&pgtype.Type{
		Name:  "timestamptz",
		OID:   pgtype.TimestamptzOID,
		Codec: &postgres_pgtype.TimestamptzCodec{},
	})

	//timestamp
	conn.TypeMap().RegisterType(&pgtype.Type{
		Name:  "timestamp",
		OID:   pgtype.TimestampOID,
		Codec: &postgres_pgtype.TimestampCodec{},
	})

	//timetz
	conn.TypeMap().RegisterType(&pgtype.Type{
		Name:  "timetz",
		OID:   pgtype.TimetzOID,
		Codec: &postgres_pgtype.TimeCodec{},
	})

	//time
	conn.TypeMap().RegisterType(&pgtype.Type{
		Name:  "time",
		OID:   pgtype.TimeOID,
		Codec: &postgres_pgtype.TimeCodec{},
	})

	//date
	conn.TypeMap().RegisterType(&pgtype.Type{
		Name:  "date",
		OID:   pgtype.DateOID,
		Codec: &postgres_pgtype.DateCodec{},
	})

	return nil
}
