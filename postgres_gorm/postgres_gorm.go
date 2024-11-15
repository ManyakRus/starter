// модуль для работы с базой данных

package postgres_gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/constants"
	"github.com/ManyakRus/starter/logger"
	"github.com/ManyakRus/starter/port_checker"
	"strings"
	"time"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"
	"os"
	"sync"

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

// NeedReconnect - флаг необходимости переподключения
var NeedReconnect bool

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_SCHEMA   string
	DB_USER     string
	DB_PASSWORD string
}

// NamingStrategy - структура для хранения настроек наименования таблиц
var NamingStrategy = schema.NamingStrategy{}

// Connect - подключается к базе данных
func Connect() {

	if Settings.DB_HOST == "" {
		FillSettings()
	}

	port_checker.CheckPort(Settings.DB_HOST, Settings.DB_PORT)

	err := Connect_err()
	LogInfo_Connected(err)

}

// Connect_err - подключается к базе данных
func Connect_err() error {
	var err error
	err = Connect_WithApplicationName_err("")

	return err
}

// Connect_WithApplicationName_SingularTableName - подключается к базе данных, с указанием имени приложения, без переименования имени таблиц
func Connect_WithApplicationName_SingularTableName(ApplicationName string) {
	err := Connect_WithApplicationName_SingularTableName_err(ApplicationName)
	if err != nil {
		log.Panicln("POSTGRES gorm Connect_WithApplicationName_SingularTableName_err() to database host: ", Settings.DB_HOST, ", Error: ", err)
	} else {
		log.Info("POSTGRES gorm Connected. host: ", Settings.DB_HOST, ", base name: ", Settings.DB_NAME, ", schema: ", Settings.DB_SCHEMA)
	}
}

// Connect_WithApplicationName_SingularTableName_err - подключается к базе данных, с указанием имени приложения, без переименования имени таблиц
func Connect_WithApplicationName_SingularTableName_err(ApplicationName string) error {
	SetSingularTableNames(true)
	err := Connect_WithApplicationName_err(ApplicationName)
	return err
}

// Connect_WithApplicationName - подключается к базе данных, с указанием имени приложения
func Connect_WithApplicationName(ApplicationName string) {
	err := Connect_WithApplicationName_err(ApplicationName)
	if err != nil {
		log.Panicln("POSTGRES gorm Connect_WithApplicationName_err() to database host: ", Settings.DB_HOST, ", Error: ", err)
	} else {
		log.Info("POSTGRES gorm Connected. host: ", Settings.DB_HOST, ", base name: ", Settings.DB_NAME, ", schema: ", Settings.DB_SCHEMA)
	}
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

	//get the database connection URL.
	dsn := GetDSN(ApplicationName)

	//
	conf := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	}
	conf.NamingStrategy = NamingStrategy

	//
	config := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, //для запуска мультизапросов
	}

	//
	dialect := postgres.New(config)
	Conn, err = gorm.Open(dialect, conf)

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
	err := Connect_WithApplicationName_err(ApplicationName)
	LogInfo_Connected(err)
	//if err != nil {
	//	log.Panic("Postgres gorm Start() error: ", err)
	//}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	stopapp.GetWaitGroup_Main().Add(1)
	go ping_go()

}

// Start_SingularTableName - делает соединение с БД, отключение и др. Без переименования имени таблиц на множественное число
func Start_SingularTableName(ApplicationName string) {
	SetSingularTableNames(true)
	Start(ApplicationName)

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
	NamingStrategy.SchemaName(Settings.DB_SCHEMA)
	NamingStrategy.TablePrefix = Settings.DB_SCHEMA + "."
}

// GetDSN - возвращает строку соединения к базе данных
func GetDSN(ApplicationName string) string {
	ApplicationName = strings.ReplaceAll(ApplicationName, " ", "_")

	dsn := "host=" + Settings.DB_HOST + " "
	dsn += "user=" + Settings.DB_USER + " "
	dsn += "password=" + Settings.DB_PASSWORD + " "
	dsn += "dbname=" + Settings.DB_NAME + " "
	dsn += "port=" + Settings.DB_PORT + " "
	dsn += "sslmode=disable TimeZone=" + constants.TIME_ZONE + " "
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
		err := Connect_WithApplicationName_err(ApplicationName)
		if err != nil {
			log.Panic("GetConnection_WithApplicationName() error: ", err)
		}
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

//// RawMultipleSQL - выполняет текст запроса, отдельно для каждого запроса
//func RawMultipleSQL(db *gorm.DB, TextSQL string) *gorm.DB {
//	var tx *gorm.DB
//	var err error
//	tx = db
//
//	// запустим все запросы отдельно
//	sqlSlice := strings.Split(TextSQL, ";")
//	len1 := len(sqlSlice)
//	for i, v := range sqlSlice {
//		if i == len1-1 {
//			tx = tx.Raw(v)
//			err = tx.Error
//		} else {
//			tx = tx.Exec(v)
//			err = tx.Error
//		}
//		if err != nil {
//			TextError := fmt.Sprint("db.Raw() error: ", err, ", TextSQL: \n", v)
//			err = errors.New(TextError)
//			break
//		}
//	}
//
//	if tx == nil {
//		log.Panic("db.Raw() error: rows =nil")
//	}
//
//	return tx
//}

// RawMultipleSQL - выполняет текст запроса, отдельно для каждого запроса
func RawMultipleSQL(db *gorm.DB, TextSQL string) *gorm.DB {
	var tx *gorm.DB
	var err error
	tx = db

	if tx == nil {
		log.Error("RawMultipleSQL() error: db =nil")
		return tx
	}

	//запустим транзакцию
	//tx0 := tx.Begin()
	//defer tx0.Commit()

	//
	TextSQL1 := ""
	TextSQL2 := TextSQL

	//запустим все запросы, кроме последнего
	pos1 := strings.LastIndex(TextSQL, ";")
	if pos1 > 0 {
		TextSQL1 = TextSQL[0:pos1]
		TextSQL2 = TextSQL[pos1:]
		tx = tx.Exec(TextSQL1)
		err = tx.Error
		if err != nil {
			TextError := fmt.Sprint("db.Exec() error: ", err, ", TextSQL: \n", TextSQL1)
			err = errors.New(TextError)
			return tx
		}
	}

	//запустим последний запрос, с возвратом результата
	tx = tx.Raw(TextSQL2)
	err = tx.Error
	if err != nil {
		TextError := fmt.Sprint("db.Raw() error: ", err, ", TextSQL: \n", TextSQL2)
		err = errors.New(TextError)
		return tx
	}

	return tx
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

// ReplaceTemporaryTableNamesToUnique - заменяет "public.TableName" на "public.TableName_UUID"
func ReplaceTemporaryTableNamesToUnique(TextSQL string) string {
	Otvet := TextSQL

	sUUID := micro.StringIdentifierFromUUID()
	map1 := make(map[string]int)

	//найдём список всех временных таблиц, и заполним в map1
	s0 := Otvet
	for {
		sFind := "CREATE TEMPORARY TABLE "
		sFind2 := "create temporary table "
		pos1 := micro.IndexSubstringMin2(s0, sFind, sFind2)
		if pos1 < 0 {
			break
		}
		s2 := s0[pos1+len(sFind):]
		pos2 := micro.IndexSubstringMin2(s2, " ", "\n")
		if pos2 <= 0 {
			break
		}
		name1 := s0[pos1+len(sFind) : pos1+len(sFind)+pos2]
		if name1 == "" {
			break
		}

		s0 = s0[pos1+len(sFind)+pos2:]
		map1[name1] = len(name1)
	}

	//заменим все временные таблицы на уникальные
	MassNames := micro.SortMapStringInt_Desc(map1)
	for _, v := range MassNames {
		sFirst := " "
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+" ", " "+v+"_"+sUUID+" ")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+"\n", " "+v+"_"+sUUID+"\n")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+"\t", " "+v+"_"+sUUID+"\t")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+";", " "+v+"_"+sUUID+";")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+"(", " "+v+"_"+sUUID+"(")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+")", " "+v+"_"+sUUID+")")

		sFirst = "\t"
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+" ", " "+v+"_"+sUUID+" ")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+"\n", " "+v+"_"+sUUID+"\n")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+"\t", " "+v+"_"+sUUID+"\t")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+";", " "+v+"_"+sUUID+";")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+"(", " "+v+"_"+sUUID+"(")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+")", " "+v+"_"+sUUID+")")

		sFirst = "\n"
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+" ", " "+v+"_"+sUUID+" ")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+"\n", " "+v+"_"+sUUID+"\n")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+"\t", " "+v+"_"+sUUID+"\t")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+";", " "+v+"_"+sUUID+";")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+"(", " "+v+"_"+sUUID+"(")
		Otvet = strings.ReplaceAll(Otvet, sFirst+v+")", " "+v+"_"+sUUID+")")

	}

	return Otvet
}

// SetSingularTableNames - меняет настройку "SingularTable" - надо ли НЕ переименовывать имя таблиц во вножественное число
// true = не переименовывать
func SetSingularTableNames(IsSingular bool) {
	NamingStrategy.SingularTable = IsSingular
}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicln("POSTGRES gorm Connect() to database host: ", Settings.DB_HOST, ", Error: ", err)
	} else {
		log.Info("POSTGRES gorm Connected. host: ", Settings.DB_HOST, ", base name: ", Settings.DB_NAME, ", schema: ", Settings.DB_SCHEMA)
	}

}
