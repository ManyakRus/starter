// модуль для работы с базой данных

package postgres_stek

import (
	"context"
	"errors"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/port_checker"
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/object_model/entities/connections"
	"time"

	"sync"
	//"time"

	//"github.com/jmoiron/sqlx"
	//_ "github.com/lib/pq"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"

	"golang.org/x/exp/maps"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

//// Conn - соединение к базе данных
//var Conn *gorm.DB

// Conn - все соединения к 10 базам данных
var MapConn = make(map[int64]*gorm.DB)

// MapConnection - все объекты Connection
var MapConnection = make(map[int64]connections.Connection)

// log - глобальный логгер
//var log = logger.GetLog()

// mutex_Connect - защита от многопоточности Reconnect()
var mutex_Connect = &sync.RWMutex{}

// mutex_ReConnect - защита от многопоточности ReConnect()
var mutex_ReConnect = &sync.RWMutex{}

// NeedReconnect - флаг необходимости переподключения
var NeedReconnect bool

//var MutexConnection sync.Mutex

// Connect_err - подключается к базе данных
func Connect(Connection connections.Connection) {

	if Connection.Server == "" {
		log.Panicln("Need fill Connection.Server")
	}

	port_checker.CheckPort(Connection.Server, Connection.Port)

	err := Connect_err(Connection)
	LogInfo_Connected(err, Connection)

}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error, Connection connections.Connection) {
	if err != nil {
		log.Panicln("POSTGRES gorm stack Connect() to database host: ", Connection.Server, ", Error: ", err)
	} else {
		log.Info("POSTGRES gorm stack Connected. host: ", Connection.Server, ", base name: ", Connection.DbName, ", schema: ", Connection.DbScheme)
	}

}

// Connect_err - подключается к базе данных
func Connect_err(Connection connections.Connection) error {

	var err error

	if Connection.Server == "" {
		log.Panicln("Need fill Connection.Server")
	}

	//ctxMain := context.Background()
	//ctxMain := contextmain.GetContext()
	//ctx, cancel := context.WithTimeout(ctxMain, 5*time.Second)
	//defer cancel()

	// get the database connection URL.
	dsn := GetDSN(Connection)

	//
	conf := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	}
	conn := postgres.Open(dsn)
	Conn, err := gorm.Open(conn, conf)
	Conn.Config.NamingStrategy = schema.NamingStrategy{TablePrefix: Connection.DbScheme + "."}
	//Conn.Config.Logger = gormlogger.Default.LogMode(gormlogger.Warn)

	if err == nil {
		DB, err := Conn.DB()
		if err != nil {
			log.Error("Conn.DB() error: ", err)
			return err
		}

		err = DB.Ping()
	}

	mutex_Connect.Lock() //race

	MapConnection[Connection.ID] = Connection
	MapConn[Connection.ID] = Conn

	mutex_Connect.Unlock()

	return err
}

// IsClosed проверка что база данных закрыта
func IsClosed(Connection connections.Connection) bool {
	var otvet bool
	Conn := MapConn[Connection.ID]
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
func Reconnect(Connection connections.Connection, err error) {
	mutex_ReConnect.Lock()
	defer mutex_ReConnect.Unlock()

	if err == nil {
		return
	}

	if errors.Is(err, context.Canceled) {
		return
	}

	Conn := MapConn[Connection.ID]
	if Conn == nil {
		log.Warn("Reconnect()")
		err := Connect_err(Connection)
		if err != nil {
			log.Error("error: ", err)
		}
		return
	}

	if IsClosed(Connection) {
		micro.Pause(1000)
		log.Warn("Reconnect()")
		err := Connect_err(Connection)
		if err != nil {
			log.Error("error: ", err)
		}
		return
	}

	sError := err.Error()
	if sError == "Conn closed" {
		micro.Pause(1000)
		log.Warn("Reconnect()")
		err := Connect_err(Connection)
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

// CloseConnectionAll - закрытие всех соединений к базам данных
func CloseConnectionAll() {

	var MapConnection_copy = make(map[int64]connections.Connection)

	maps.Copy(MapConnection_copy, MapConnection) // копия для race error

	for _, Connection := range MapConnection_copy {
		if Connection.Server == "" {
			continue
		}
		CloseConnection(Connection)
	}
}

// CloseConnection - закрытие соединения с базой данных
func CloseConnection(Connection connections.Connection) {
	Conn := MapConn[Connection.ID]
	if Conn == nil {
		return
	}

	err := CloseConnection_err(Connection)
	if err != nil {
		log.Error("Postgres gorm stack CloseConnection() error: ", err)
	} else {
		log.Info("Postgres gorm stack connection closed")
	}

	return
}

// CloseConnection_err - закрытие соединения с базой данных
func CloseConnection_err(Connection connections.Connection) error {

	Conn := MapConn[Connection.ID]
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

	mutex_Connect.Lock() //race

	delete(MapConnection, Connection.ID)
	delete(MapConn, Connection.ID)

	mutex_Connect.Unlock()

	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {
	defer stopapp.GetWaitGroup_Main().Done()

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. postgres_stek")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("postgres_stek")

	//
	CloseConnectionAll()

}

// StartDB - делает соединение с БД, отключение и др.
func StartDB(Connection connections.Connection) {
	var err error

	ctx := contextmain.GetContext()
	WaitGroup := stopapp.GetWaitGroup_Main()
	err = Start_ctx(&ctx, WaitGroup, Connection)
	LogInfo_Connected(err, Connection)

}

// Start_ctx - необходимые процедуры для подключения к серверу БД
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто StartDB()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup, Connection connections.Connection) error {
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
	err = Connect_err(Connection)
	if err != nil {
		return err
	}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	stopapp.GetWaitGroup_Main().Add(1)
	go ping_go()

	return err
}

// GetDSN - возвращает строку соединения к базе данных
func GetDSN(Connection connections.Connection) string {
	dsn := "host=" + Connection.Server + " "
	dsn += "user=" + Connection.Login + " "
	dsn += "password=" + Connection.Password + " "
	dsn += "dbname=" + Connection.DbName + " "
	dsn += "port=" + Connection.Port + " sslmode=disable TimeZone=UTC"

	return dsn
}

// GetConnection - возвращает соединение к нужной базе данных
func GetConnection(Connection connections.Connection) *gorm.DB {
	//мьютекс чтоб не подключаться одновременно
	mutex_Connect.RLock()
	defer mutex_Connect.RUnlock()

	//мьютекс чтоб не подключаться одновременно
	mutex_Connect.RLock()
	defer mutex_Connect.RUnlock()

	//
	Conn := MapConn[Connection.ID]
	if Conn == nil {
		Connect(Connection)
		Conn = MapConn[Connection.ID]
	}

	return Conn
}

// ping_go - делает пинг каждые 60 секунд, и реконнект
func ping_go() {
	var err error

	defer stopapp.GetWaitGroup_Main().Done()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	//бесконечный цикл
loop:
	for {

		for _, Connection := range MapConnection {
			if Connection.Server == "" {
				continue
			}

			addr := Connection.Server + ":" + Connection.Port

			select {
			case <-contextmain.GetContext().Done():
				log.Warn("Context app is canceled. postgres_stek.ping")
				break loop
			case <-ticker.C:
				err = port_checker.CheckPort_err(Connection.Server, Connection.Port)
				//log.Debug("ticker, ping err: ", err) //удалить
				if err != nil {
					NeedReconnect = true
					log.Warn("postgres_stek CheckPort(", addr, ") error: ", err)
				} else if NeedReconnect == true {
					log.Warn("postgres_stek CheckPort(", addr, ") OK. Start Reconnect()")
					NeedReconnect = false
					err = Connect_err(Connection)
					if err != nil {
						NeedReconnect = true
						log.Error("Connect_err() error: ", err)
					}
				}
			}
		}
	}

}
