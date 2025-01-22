// модуль для использования сервиса CAMUNDA
package camunda_connect

import (
	"context"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/logger"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/port_checker"
	"github.com/ManyakRus/starter/stopapp"
	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"sync"

	// "gitlab.aescorp.ru/dsp_dev/claim/stack_exchange/internal/v0/app/constants"
	// "github.com/ManyakRus/starter/mssql"
	"os"
	"time"
)

// log - глобальный логгер
var log = logger.GetLog()

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// NeedReconnect - флаг необходимости переподключения
var NeedReconnect bool

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	CAMUNDA_HOST string
	CAMUNDA_PORT string
}

// Client - клиент подключения к CAMUNDA_ID
var Client zbc.Client

// JobWorker - worker который выполняет подключение к приему сообщений от CAMUNDA
var JobWorker worker.JobWorker

// StartSettings - параметры для запуска камунды
type StartSettings struct {
	HandleJob       func(client worker.JobClient, job entities.Job)
	CAMUNDA_JOBTYPE string
	BPMN_filename   string
	TimeOut         time.Duration
	MaxJobsActive   int
	MaxConcurrency  int
}

// FillSettings загружает переменные окружения в структуру из файла или из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.CAMUNDA_HOST = os.Getenv("CAMUNDA_HOST")
	Settings.CAMUNDA_PORT = os.Getenv("CAMUNDA_PORT")
	if Settings.CAMUNDA_HOST == "" {
		log.Panic("Need fill CAMUNDA_HOST ! in OS Environment ")
	}

	if Settings.CAMUNDA_PORT == "" {
		log.Panic("Need fill CAMUNDA_PORT ! in OS Environment ")
	}

}

// Connect_err - подключается к серверу Camunda, паника при ошибке
func Connect() {
	err := Connect_err()
	LogInfo_Connected(err)
}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panic("CAMUNDA Connect_err() ip: ", Settings.CAMUNDA_HOST, " port: ", Settings.CAMUNDA_PORT, " error: ", err)
	} else {
		log.Info("CAMUNDA connected, ip: ", Settings.CAMUNDA_HOST, " port: ", Settings.CAMUNDA_PORT)
	}

}

// Connect_err - подключается к серверу Camunda, возвращает ошибку
func Connect_err() error {
	var err error

	if Settings.CAMUNDA_HOST == "" {
		FillSettings()
	}

	err = port_checker.CheckPort_err(Settings.CAMUNDA_HOST, Settings.CAMUNDA_PORT)
	if err != nil {
		return err
	}

	Client, err = zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         GetURL(),
		UsePlaintextConnection: true,
	})
	if err != nil {
		return err
	}

	//log.Infoln("CAMUNDA connected. ip: ", Settings.CAMUNDA_HOST)

	// JobWorker = Client.NewJobWorker().JobType(CAMUNDA_ID).Handler(HandleJob).Open()
	return err
}

// CloseConnection - отключается от сервера Camunda
func CloseJobWorker() {
	if JobWorker != nil {
		JobWorker.Close()
		JobWorker.AwaitClose()
	}
	JobWorker = nil
}

// CloseConnection - отключается от сервера Camunda
func CloseConnection() {
	if JobWorker != nil {
		JobWorker.Close()
		JobWorker.AwaitClose()
	}
	err := Client.Close()
	if err != nil {
		log.Panicln("Client.Close() error: ", err)
	}

	log.Infoln("CAMUNDA stopped")
	Client = nil
	JobWorker = nil
}

// GetURL - возврашает строку соединения к серверу Camunda
func GetURL() string {
	Otvet := ""

	if Settings.CAMUNDA_HOST == "" {
		log.Panicln("CAMUNDA_HOST = ''")
	}

	if Settings.CAMUNDA_PORT == "" {
		log.Panicln("CAMUNDA_PORT = ''")
	}

	Otvet = Settings.CAMUNDA_HOST + ":" + Settings.CAMUNDA_PORT

	return Otvet
}

// WorkComplete - отправляет статус ОК на сервер Camunda
func WorkComplete(client worker.JobClient, jobKey int64, variables map[string]interface{}) error {

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		log.Panicln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	_, err = request.Send(ctx)
	if err != nil {
		log.Error("camunda_connect.WorkComplete() error: ", err)
	}

	// log.Debugf("[INFO] HandleJob, %v, complete\n", jobKey)
	return err
}

// WorkFails - отправляет статус ошибки на сервер Camunda
func WorkFails(err error, client worker.JobClient, job entities.Job) error {
	if err == nil {
		log.Panicln("err =nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	jobKey := job.GetKey()
	retries := job.GetRetries()

	_, err1 := client.NewFailJobCommand().JobKey(jobKey).Retries(retries - 1).ErrorMessage(err.Error()).Send(ctx)
	if err1 != nil {
		log.Error("camunda_connect.WorkFails() error: ", err1)
	}

	// log.Debugf("[WARNING] HandleJob, %v, fail\n", jobKey)
	return err1
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. camunda_connect")
	}

	CloseJobWorker()

	// ждём пока отправляемых сейчас сообщений будет =0
	stopapp.WaitTotalMessagesSendingNow("camunda_connect")

	// закрываем соединение
	CloseConnection()
	stopapp.GetWaitGroup_Main().Done()
}

// StartCamunda - необходимые процедуры для подключения к серверу Camunda
func StartCamunda(HandleJob func(client worker.JobClient, job entities.Job), CAMUNDA_JOBTYPE string, BPMN_filename string) {
	var err error

	ctx := contextmain.GetContext()
	WaitGroup := stopapp.GetWaitGroup_Main()
	err = Start_ctx(&ctx, WaitGroup, HandleJob, CAMUNDA_JOBTYPE, BPMN_filename)
	LogInfo_Connected(err)
}

// Start_WithSettings - необходимые процедуры для подключения к серверу Camunda
func Start_WithSettings(settings StartSettings) {
	var err error

	//
	err = Connect_err()
	LogInfo_Connected(err)

	Send_BPMN_File(settings.BPMN_filename)

	JobWorkerStep3 := Client.NewJobWorker().JobType(settings.CAMUNDA_JOBTYPE).Handler(settings.HandleJob)
	if settings.MaxConcurrency > 0 {
		JobWorkerStep3 = JobWorkerStep3.Concurrency(settings.MaxConcurrency)
	}
	if settings.MaxJobsActive > 0 {
		JobWorkerStep3 = JobWorkerStep3.MaxJobsActive(settings.MaxJobsActive)
	}
	if settings.TimeOut.Seconds() > 0 {
		JobWorkerStep3 = JobWorkerStep3.Timeout(settings.TimeOut)
	}
	JobWorker = JobWorkerStep3.Open()

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	stopapp.GetWaitGroup_Main().Add(1)
	go ping_go(settings.HandleJob, settings.CAMUNDA_JOBTYPE)

	return
}

// Start_ctx - необходимые процедуры для подключения к серверу Camunda
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто StartCamunda()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup, HandleJob func(client worker.JobClient, job entities.Job), CAMUNDA_JOBTYPE string, BPMN_filename string) error {
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

	settings := StartSettings{}
	settings.HandleJob = HandleJob
	settings.CAMUNDA_JOBTYPE = CAMUNDA_JOBTYPE
	settings.BPMN_filename = BPMN_filename
	Start_WithSettings(settings)

	////
	//err = Connect_err()
	//if err != nil {
	//	return err
	//}
	//
	//Send_BPMN_File(BPMN_filename)
	//
	//JobWorker = Client.NewJobWorker().JobType(CAMUNDA_JOBTYPE).Handler(HandleJob).Open()
	//
	//stopapp.GetWaitGroup_Main().Add(1)
	//go WaitStop()
	//
	//stopapp.GetWaitGroup_Main().Add(1)
	//go ping_go(HandleJob, CAMUNDA_JOBTYPE)

	return err
}

// Send_BPMN_File - отправляем файл .bpmn в камунду
func Send_BPMN_File(BPMN_filename string) {
	var err error
	if BPMN_filename == "" {
		return
	}

	ctxMain := contextmain.GetContext()
	ctx, ctxCancelFunc := context.WithTimeout(ctxMain, time.Second*60)
	defer ctxCancelFunc()

	FileName := BPMN_filename
	dir := micro.ProgramDir_Common()
	FlagFind, err := micro.FileExists(dir + "bin")
	if FlagFind == true {
		dir = dir + "bin" + micro.SeparatorFile()
	}

	FileName = dir + FileName
	log.Info("Load .bpmn file from: ", FileName)

	res, err := Client.NewDeployResourceCommand().AddResourceFile(FileName).Send(ctx)
	if err != nil {
		log.Panicln(err)
	}
	log.Infof("Send .bpmn file, result: %v", res)
}

// ping_go - делает пинг каждые 60 секунд, и реконнект
func ping_go(HandleJob func(client worker.JobClient, job entities.Job), CAMUNDA_JOBTYPE string) {

	ticker := time.NewTicker(60 * time.Second)

	addr := Settings.CAMUNDA_HOST + ":" + Settings.CAMUNDA_PORT

	// бесконечный цикл
loop:
	for {
		select {
		case <-contextmain.GetContext().Done():
			log.Warn("Context app is canceled. camunda_connect.ping")
			break loop
		case <-ticker.C:
			err := port_checker.CheckPort_err(Settings.CAMUNDA_HOST, Settings.CAMUNDA_PORT)
			// log.Debug("ticker, ping err: ", err) //удалить
			if err != nil {
				NeedReconnect = true
				log.Warn("CAMUNDA CheckPort(", addr, ") error: ", err)
			} else if NeedReconnect == true {
				log.Warn("CAMUNDA CheckPort(", addr, ") OK. Start Reconnect()")
				NeedReconnect = false
				Connect()
				JobWorker = Client.NewJobWorker().JobType(CAMUNDA_JOBTYPE).Handler(HandleJob).Open()
			}
		}
	}

	stopapp.GetWaitGroup_Main().Done()
}
