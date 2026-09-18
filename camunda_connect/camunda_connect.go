// модуль для использования сервиса CAMUNDA
package camunda_connect

import (
	"context"
	"fmt"
	"github.com/ManyakRus/starter/contextmain"
	//"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/port_checker"
	"github.com/ManyakRus/starter/stopapp"
	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/pb"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"sync"

	// "gitlab.aescorp.ru/dsp_dev/claim/stack_exchange/internal/v0/app/constants"
	// "github.com/ManyakRus/starter/mssql"
	"os"
	"time"
)

// PackageName - имя текущего пакета, для логирования
const PackageName = "camunda_connect"

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// NeedReconnect - флаг необходимости переподключения
var NeedReconnect bool

// TextRPCError - текст ошибки "rpc error"
//const TextRPCError = "rpc error"

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

	//отправим ping для проверки соединения
	ctx, cancel := context.WithTimeout(contextmain.GetContext(), 60*time.Second)
	defer cancel()
	_, err = Client.NewTopologyCommand().Send(ctx)
	if err != nil {
		err = fmt.Errorf("Connect_err() проверка ping() error: %w", err)
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
func WorkComplete(jobKey int64, variables map[string]interface{}) error {

	request, err := Client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		log.Error(err)
		return err
	}

	ctx, cancel := context.WithTimeout(contextmain.GetContext(), 300*time.Second)
	defer cancel()

	_, err = request.Send(ctx)
	if err != nil {
		log.Warn("camunda_connect.WorkComplete() error: ", err)

		//вторая попытка
		//реконнект
		err = Connect_err()
		if err != nil {
			NeedReconnect = true
			log.Error("Connect_err() error: ", err)
		}

		//повтор отправки
		request, err = Client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
		if err != nil {
			log.Error(err)
			return err
		}
		_, err = request.Send(ctx)
		if err != nil {
			log.Error("camunda_connect.WorkComplete() retry error: ", err)
		}
	}

	// log.Debugf("[INFO] HandleJob, %v, complete\n", jobKey)
	return err
}

// WorkComplete_answer - отправляет статус ОК на сервер Camunda, и возвращает ответ
func WorkComplete_answer(jobKey int64, variables map[string]interface{}) (*pb.CompleteJobResponse, error) {

	request, err := Client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		log.Panicln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	Otvet, err := request.Send(ctx)
	if err != nil {
		log.Error("camunda_connect.WorkComplete_answer() error: ", err)

		//вторая попытка
		//реконнект
		err = Connect_err()
		if err != nil {
			NeedReconnect = true
			log.Error("Connect_err() error: ", err)
		}

		//повтор отправки
		request, err = Client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
		if err != nil {
			log.Error(err)
			return Otvet, err
		}
		_, err = request.Send(ctx)
		if err != nil {
			log.Error("camunda_connect.WorkComplete_answer() retry error: ", err)
		}
	}

	// log.Debugf("[INFO] HandleJob, %v, complete\n", jobKey)
	return Otvet, err
}

// WorkFails - отправляет статус ошибки на сервер Camunda
func WorkFails(job entities.Job, err0 error) error {
	var err error

	//err должен быть непустой
	if err0 == nil {
		err2 := fmt.Errorf("WorkFails() error: err=nil")
		log.Warn(err2)
		return err2
	}

	//
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	//retries - 1
	jobKey := job.GetKey()
	retries := job.GetRetries()
	retries = retries - 1
	if retries < 0 {
		retries = 0
	}

	_, err = Client.NewFailJobCommand().JobKey(jobKey).Retries(retries).ErrorMessage(err0.Error()).Send(ctx)
	if err != nil {
		log.Error("camunda_connect.WorkFails() error: ", err)

		//вторая попытка
		//реконнект
		err = Connect_err()
		if err != nil {
			NeedReconnect = true
			log.Error("Connect_err() error: ", err)
		}

		//повтор отправки
		_, err = Client.NewFailJobCommand().JobKey(jobKey).Retries(retries).ErrorMessage(err0.Error()).Send(ctx)
		if err != nil {
			log.Error("camunda_connect.WorkFails() retry error: ", err)
		}
	}

	// log.Debugf("[WARNING] HandleJob, %v, fail\n", jobKey)
	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {
	defer waitGroup_Connect.Done()

	select {
	case <-(*ctx_Connect).Done():
		log.Warn("Context app is canceled. camunda_connect")
	}

	CloseJobWorker()

	// ждём пока отправляемых сейчас сообщений будет =0
	stopapp.WaitTotalMessagesSendingNow("camunda_connect")

	// закрываем соединение
	CloseConnection()
}

// StartCamunda - необходимые процедуры для подключения к серверу Camunda
func StartCamunda(HandleJob func(client worker.JobClient, job entities.Job), CAMUNDA_JOBTYPE string, BPMN_filename string) {
	var err error

	ctx := GetContext()
	WaitGroup := GetWaitGroup()
	err = Start_ctx(ctx, WaitGroup, HandleJob, CAMUNDA_JOBTYPE, BPMN_filename)
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

	//сохраним в список подключений
	WaitGroupContext1 := stopapp.WaitGroupContext{WaitGroup: GetWaitGroup(), Ctx: GetContext(), CancelCtxFunc: cancelCtxFunc}
	stopapp.OrderedMapConnections.Put(PackageName, WaitGroupContext1)

	//
	waitGroup_Connect.Add(1)
	go WaitStop()

	waitGroup_Connect.Add(1)
	go ping_go(settings.HandleJob, settings.CAMUNDA_JOBTYPE)

	return
}

// Start_ctx - необходимые процедуры для подключения к серверу Camunda
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто StartCamunda()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup, HandleJob func(client worker.JobClient, job entities.Job), CAMUNDA_JOBTYPE string, BPMN_filename string) error {
	var err error

	//запомним к себе контекст
	if ctx == nil {
		ctx = GetContext()
	} else {
		SetContext(ctx)
	}

	//запомним к себе WaitGroup
	if WaitGroup == nil {
		WaitGroup = GetWaitGroup()
	} else {
		SetWaitGroup(WaitGroup)
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
	//stopapp.GetWaitGroup().Add(1)
	//go WaitStop()
	//
	//stopapp.GetWaitGroup().Add(1)
	//go ping_go(HandleJob, CAMUNDA_JOBTYPE)

	return err
}

// Send_BPMN_File - отправляем файл .bpmn в камунду
func Send_BPMN_File(BPMN_filename string) {
	var err error
	if BPMN_filename == "" {
		return
	}

	ctxMain := *ctx_Connect
	ctx, ctxCancelFunc := context.WithTimeout(ctxMain, time.Second*60)
	defer ctxCancelFunc()

	FileName := BPMN_filename
	//dir := micro.ProgramDir_bin()

	//FileName = dir + FileName
	log.Info("Load .bpmn file from: ", FileName)

	res, err := Client.NewDeployResourceCommand().AddResourceFile(FileName).Send(ctx)
	if err != nil {
		err = fmt.Errorf("AddResourceFile().Send() error: %w", err)
		log.Panicln(err)
	}
	log.Infof("Send .bpmn file, result: %v", res)
}

// ping_go - делает пинг каждые 60 секунд, и реконнект
func ping_go(HandleJob func(client worker.JobClient, job entities.Job), CAMUNDA_JOBTYPE string) {
	var err error

	defer waitGroup_Connect.Done()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	addr := Settings.CAMUNDA_HOST + ":" + Settings.CAMUNDA_PORT

	// бесконечный цикл
loop:
	for {
		select {
		case <-(*ctx_Connect).Done():
			log.Warn("Context app is canceled. camunda_connect.ping")
			break loop
		case <-ticker.C:
			//проверяем порт
			err_port := port_checker.CheckPort_err(Settings.CAMUNDA_HOST, Settings.CAMUNDA_PORT)
			if err != nil {
				log.Warn("CAMUNDA CheckPort(", addr, ") error: ", err_port)
				NeedReconnect = true
				continue //реконнект нужен когда не будет ошибки
			}

			//проверяем тестовый запрос в камунду
			ctx, cancelfunc := context.WithTimeout(contextmain.GetContext(), time.Second*60)
			defer cancelfunc()
			_, err2 := Client.NewTopologyCommand().Send(ctx)
			if err2 != nil {
				log.Warn("CAMUNDA Check NewTopologyCommand() error: ", err2)
				NeedReconnect = true
			}

			//
			//err = errors.Join(err_port, err2)

			//реконнект
			if NeedReconnect == true {
				log.Warn("CAMUNDA Check ", addr, " OK. Start Reconnect()")
				NeedReconnect = false
				err = Connect_err()
				if err != nil {
					NeedReconnect = true
					log.Error("Connect_err() error: ", err)
					continue
				}

				//новый JobWorker
				if JobWorker != nil {
					JobWorker.Close()
				}
				JobWorker = Client.NewJobWorker().JobType(CAMUNDA_JOBTYPE).Handler(HandleJob).Open()
			}

		}
	}

}
