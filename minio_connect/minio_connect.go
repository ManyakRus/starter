package minio_connect

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/port_checker"
	"github.com/ManyakRus/starter/stopapp"

	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Conn - соединение к Minio
var Conn *miniogo.Client

// PackageName - имя текущего пакета, для логирования
const PackageName = "minio_connect"

// mutexReconnect - защита от многопоточности Reconnect()
var mutexReconnect = &sync.Mutex{}

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// NeedReconnect - флаг необходимости переподключения
var NeedReconnect bool

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	MINIO_HOST                 string
	MINIO_PORT                 string
	MINIO_KEY                  string
	MINIO_SECRET_KEY           string
	MINIO_USE_SSL              bool
	MINIO_INSECURE_SKIP_VERIFY bool
}

// Connect_err - подключается к Minio
func Connect() {

	if Settings.MINIO_HOST == "" {
		FillSettings()
	}

	//ping.CheckPort(Settings.MINIO_HOST, Settings.MINIO_PORT)

	err := Connect_err()
	LogInfo_Connected(err)

}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicln("Minio Connect_err() host: ", Settings.MINIO_HOST, ", Error: ", err)
	} else {
		log.Info("Minio connected. host: ", Settings.MINIO_HOST, ", port: ", Settings.MINIO_PORT)
	}

}

// Connect_err - подключается к Minio
func Connect_err() error {
	var err error

	if Settings.MINIO_HOST == "" {
		FillSettings()
	}

	//ctxMain := context.Background()
	//ctxMain := contextmain.GetContext()
	//ctx, cancel := context.WithTimeout(ctxMain, 5*time.Second)
	//defer cancel()

	addr := Settings.MINIO_HOST + ":" + Settings.MINIO_PORT
	options := &miniogo.Options{
		Creds:     credentials.NewStaticV4(Settings.MINIO_KEY, Settings.MINIO_SECRET_KEY, ""),
		Secure:    Settings.MINIO_USE_SSL,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: Settings.MINIO_INSECURE_SKIP_VERIFY}},
	}
	Conn, err = miniogo.New(addr, options)
	if err == nil {
		port_checker.CheckPort(Settings.MINIO_HOST, Settings.MINIO_PORT)
	}

	return err
}

// IsClosed проверка что Minio закрыто
func IsClosed() bool {
	var Otvet bool
	if Conn == nil {
		return true
	}

	Otvet = Conn.IsOffline()
	return Otvet
}

// Reconnect повторное подключение к Minio, если оно отключено
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

// CloseConnection - закрытие соединения с Minio
func CloseConnection() error {
	if Conn == nil {
		return nil
	}

	err := CloseConnection_err()
	if err != nil {
		log.Error("Minio CloseConnection() error: ", err)
	} else {
		log.Info("Minio connection closed")
	}

	return err
}

// CloseConnection - закрытие соединения с Minio
func CloseConnection_err() error {
	var err error
	if Conn == nil {
		return nil
	}

	//ctx := contextmain.GetContext()
	//ctx := context.Background()
	//err := Conn.Close()

	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {
	defer waitGroup_Connect.Done()

	select {
	case <-ctx_Connect.Done():
		log.Warn("Context app is canceled. minio_connect")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("Minio")

	//
	err := CloseConnection()
	if err != nil {
		log.Error("CloseConnection() error: ", err)
	}
}

// StartMinio - необходимые процедуры для подключения к серверу Minio
func StartMinio() {
	var err error

	ctx := ctx_Connect
	WaitGroup := waitGroup_Connect
	err = Start_ctx(&ctx, WaitGroup)
	LogInfo_Connected(err)

}

// Start_ctx - необходимые процедуры для подключения к серверу Minio
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто StartMinio()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup) error {
	var err error

	//запомним к себе контекст
	//	if contextmain.Ctx != ctx {
	//		contextmain.SetContext(ctx)
	//	}
	//contextmain.Ctx = ctx
	if ctx == nil {
		ctx = &ctx_Connect
	}

	//запомним к себе WaitGroup
	//stopapp.SetWaitGroup_Main(WaitGroup)
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

// FillSettings загружает переменные окружения в структуру из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.MINIO_HOST = os.Getenv("MINIO_HOST")
	Settings.MINIO_PORT = os.Getenv("MINIO_PORT")
	Settings.MINIO_KEY = os.Getenv("MINIO_KEY")
	Settings.MINIO_SECRET_KEY = os.Getenv("MINIO_SECRET_KEY")

	//
	sMINIO_USE_SSL := os.Getenv("MINIO_USE_SSL")
	MINIO_USE_SSL := micro.BoolFromString(sMINIO_USE_SSL)
	if sMINIO_USE_SSL == "" {
		MINIO_USE_SSL = true
	}
	Settings.MINIO_USE_SSL = MINIO_USE_SSL

	//
	sINSECURE_SKIP_VERIFY := os.Getenv("MINIO_INSECURE_SKIP_VERIFY")
	INSECURE_SKIP_VERIFY := micro.BoolFromString(sINSECURE_SKIP_VERIFY)
	if sINSECURE_SKIP_VERIFY == "" {
		INSECURE_SKIP_VERIFY = true
	}
	Settings.MINIO_INSECURE_SKIP_VERIFY = INSECURE_SKIP_VERIFY

	if Settings.MINIO_HOST == "" {
		log.Panicln("Need fill MINIO_HOST ! in os.ENV ")
	}

	if Settings.MINIO_PORT == "" {
		log.Panicln("Need fill MINIO_PORT ! in os.ENV ")
	}

	if Settings.MINIO_KEY == "" {
		log.Panicln("Need fill MINIO_KEY ! in os.ENV ")
	}

	if Settings.MINIO_SECRET_KEY == "" {
		log.Panicln("Need fill MINIO_SECRET_KEY ! in os.ENV ")
	}
	//

}

// ping_go - делает пинг каждые 60 секунд, и реконнект
func ping_go() {
	var err error

	defer waitGroup_Connect.Done()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	addr := Settings.MINIO_HOST + ":" + Settings.MINIO_PORT

	//бесконечный цикл
loop:
	for {
		select {
		case <-ctx_Connect.Done():
			log.Warn("Context app is canceled. minio_connect.ping")
			break loop
		case <-ticker.C:
			err = port_checker.CheckPort_err(Settings.MINIO_HOST, Settings.MINIO_PORT)
			//log.Debug("ticker, ping err: ", err) //удалить
			if err != nil {
				NeedReconnect = true
				log.Warn("minio_connect CheckPort(", addr, ") error: ", err)
			} else if NeedReconnect == true {
				log.Warn("minio_connect CheckPort(", addr, ") OK. Start Reconnect()")
				NeedReconnect = false
				Connect()
			}
		}
	}

}

// CreateBucketCtx_err -создание бакета (раздела) хранения файлов
// bucketName - имя бакета (раздела)
// location - локация (moscow)
func CreateBucketCtx_err(ctx context.Context, bucketName string, location string) error {
	var err error

	ctxMain := ctx_Connect
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()

	err = Conn.MakeBucket(ctx, bucketName, miniogo.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := Conn.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			// log.Printf("[INFO] BucketExists, %q", t.bucketName)
		} else {
			return fmt.Errorf("MakeBucket, error: %w", err)
		}
	}
	return nil
}

// CreateBucketCtx -создание бакета (раздела) хранения файлов
// bucketName - имя бакета (раздела)
// location - локация (moscow)
func CreateBucketCtx(ctx context.Context, bucketName string, location string) {
	err := CreateBucketCtx_err(ctx, bucketName, location)
	if err != nil {
		log.Panic("CreateBucketCtx() bucketName: ", bucketName, " error: ", err)
	} else {
		log.Debug("CreateBucketCtx() bucketName: ", bucketName, " OK")
	}
}

// UploadFileCtx_err - загружает файл на сервер MinIO
// возвращает ошибку
func UploadFileCtx_err(ctx context.Context, bucketName, objectName, filePath string) (string, error) {
	contentType := "application/pdf"
	info, err := Conn.FPutObject(ctx, bucketName, objectName, filePath, miniogo.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", fmt.Errorf("UploadFile, error: %w", err)
	}

	log.Debugf("UploadFile, Successfully upload file: %q, size: %d, tag: %s\n", objectName, info.Size, info.ETag)
	return info.ETag, nil
}

// UploadFileCtx - загружает файл на сервер MinIO, при ошибке паника
// возвращаю ETag и ошибку
func UploadFileCtx(ctx context.Context, bucketName, objectName, filePath string) string {
	Otvet, err := UploadFileCtx_err(ctx, bucketName, objectName, filePath)

	if err != nil {
		log.Panicf("UploadFileCtx() objectName: %s, bucketName: %s, error: %v", objectName, bucketName, err)
	} else {
		log.Debugf("UploadFileCtx() objectName: %s, bucketName: %s, OK", objectName, bucketName)
	}

	return Otvet
}

// DownloadFileCtx - загружает файл на сервер MinIO, при ошибке паника
// возвращаю файл
func DownloadFileCtx(ctx context.Context, bucketName, objectName string) []byte {
	Otvet, err := DownloadFileCtx_err(ctx, bucketName, objectName)

	if err != nil {
		log.Panic("UploadFileCtx() objectName: %s, bucketName: %s, error: %v", objectName, bucketName, err)
	} else {
		log.Debug("UploadFileCtx() objectName: %s, bucketName: %s, OK", objectName, bucketName)
	}

	return Otvet
}

// DownloadFileCtx_err - загружает файл из сервера MinIO
// возвращает файл и ошибку
func DownloadFileCtx_err(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	Otvet := make([]byte, 100)
	var err error

	Object, err := Conn.GetObject(ctx, bucketName, objectName, miniogo.GetObjectOptions{})
	if err != nil {
		//log.Panic("GetObject() error: ", err)
		//return Otvet, fmt.Errorf("DownloadFileCtx_err(), error: %w", err)
		return Otvet, err
	}
	defer Object.Close()

	//count, err := Object.Read(Otvet)
	Otvet, err = io.ReadAll(Object)
	if err != nil {
		//log.Panic("minio Read() error: ", err)
		//return Otvet, err
		return Otvet, err
	}
	if len(Otvet) == 0 {
		TextError := "minio Read() error: len=0"
		log.Error(TextError)
		err := errors.New(TextError)
		return Otvet, err
	}

	log.Debug("DownloadFileCtx_err() OK, objectName: ", objectName)

	return Otvet, nil
}
