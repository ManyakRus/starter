package nats_liveness

import (
	bytes "bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/nats_connect"
	"github.com/ManyakRus/starter/stopapp"
	"github.com/klauspost/compress/zstd"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_types"
	"os"
	"strings"
	"sync"
	"time"
)

//var bufOut1 *bytes.Buffer

var Settings SettingsINI

var ServiceStartAt time.Time

var Ticker *time.Ticker

type SettingsINI struct {
	NATS_LIVENESS_TOPIC string
	SERVICE_NAME        string
	SERVICE_NUMBER      string
	STAGE               string
	SERVICE_NAME_FULL   string
}

type Message struct {
	ServiceName   string `json:"service_name"`   // Имя сервиса
	ServiceTime   string `json:"service_time"`   // Фактическое время сервиса
	ServiceUptime string `json:"service_uptime"` // Аптайм сервиса
	ServiceNum    string `json:"service_num"`    // Уникальный номер сервиса
	KernelVers    string `json:"kernel_version"` // Фактическая версия ядра
	KernelType    string `json:"kernel_type"`    // Фактическая версия ядра
}

// Connect - подключается к серверу Nats-sync_exchange
func Connect() {
	err := Connect_err()
	if err != nil {
		log.Panicln("Can not connect NATS: ", nats_connect.Settings.NATS_HOST, ", error: ", err)
	}
	log.Info("NATS liveness connected. host: ", nats_connect.Settings.NATS_HOST, ":", nats_connect.Settings.NATS_PORT, ", topic: ", Settings.NATS_LIVENESS_TOPIC)
}

// Connect_err - подключается к серверу Nats-sync_exchange и возвращает ошибку
func Connect_err() error {
	var err error

	err = nats_connect.Connect_err()

	nats_connect.FillSettings()

	//sNATS_PORT := (nats_connect.Settings.NATS_PORT)
	//url := "nats://" + nats_connect.Settings.NATS_HOST + ":" + sNATS_PORT
	//err = sync_exchange.InitSyncExchange(url, Settings.SERVICE_NAME)

	ServiceStartAt = time.Now()

	return err
}

// CloseConnection - закрывает соединение с сервером Nats-sync_exchange
func CloseConnection() {
	err := CloseConnection_err()
	if err != nil {
		log.Warn("Can not CloseConnection() NATS Liveness: ", nats_connect.Settings.NATS_HOST, " warning: ", err)
	}
}

// CloseConnection - закрывает соединение с сервером Nats-sync_exchange, и возвращает ошибку
func CloseConnection_err() error {
	if Ticker != nil {
		Ticker.Stop()
	}
	//err := sync_exchange.DeInitSyncExchange()
	nats_connect.CloseConnection()
	return nil
}

// SendMessage - отправляет 1 сообщение в Nats-sync_exchange
func SendMessage() {
	var err error
	//obj := sync_types.SyncObject(b)
	//pack := sync_types.MakeSyncObject(obj)
	//sync_exchange.SendMessage(Settings.NATS_LIVENESS_TOPIC, pack)

	now := time.Now().String()
	duration := time.Since(ServiceStartAt)
	sDuration := duration.String()

	//SERVICE_NAME := Settings.SERVICE_NAME
	//SERVICE_NAME = SERVICE_NAME + "_" + Settings.STAGE
	//if micro.IsTestApp() == true {
	//}

	Message1 := Message{}
	Message1.ServiceName = Settings.SERVICE_NAME_FULL
	Message1.ServiceNum = Settings.SERVICE_NUMBER
	Message1.KernelType = "nikitin"
	Message1.KernelVers = ""
	Message1.ServiceTime = now //time.Now().UTC().Format("2006-01-02 15:04:05.000")
	Message1.ServiceUptime = sDuration

	var sMessage1 string
	bytes1, err := json.Marshal(Message1)
	//sMessage0 := string(bytes1)
	sMessage1 = base64.StdEncoding.EncodeToString(bytes1)

	params := make(sync_types.SyncParams)
	params["binData"] = sMessage1 //time.Now().UTC().Format("2006-01-02 15:04:05.000")
	msg := sync_types.MakeSyncCommand("live", params)

	raw_data, err := sync_types.SyncPackageToJSON(&msg)
	if err != nil {
		log.Error("SyncPackageToJSON() error: ", err)
		return
	}
	// Пакуем сообщение
	var enc *zstd.Encoder
	var bufOut *bytes.Buffer

	enc, _ = zstd.NewWriter(bufOut, zstd.WithEncoderLevel(3))

	var msg_bin = make([]byte, 0)
	msg_bin = enc.EncodeAll([]byte(raw_data), msg_bin)
	err = nats_connect.SendMessage(Settings.NATS_LIVENESS_TOPIC, msg_bin)
	if err != nil {
		log.Error("SendMessage() error: ", err)
		return
	}

	//err = sync_exchange.SendMessage(Settings.NATS_LIVENESS_TOPIC, msg)
	//if err != nil {
	//	log.Error("SendMessage() error: ", err)
	//}

}

func FillSettings(SERVICE_NAME string) {
	Settings.SERVICE_NAME = strings.ToUpper(SERVICE_NAME)
	Settings.STAGE = os.Getenv("STAGE")
	SERVICE_NAME_FULL := Settings.SERVICE_NAME
	if Settings.STAGE != "" {
		SERVICE_NAME_FULL = SERVICE_NAME_FULL + "_" + Settings.STAGE
	}
	Settings.SERVICE_NAME_FULL = SERVICE_NAME_FULL

	NATS_LIVENESS_TOPIC := "/claim/" + Settings.SERVICE_NAME_FULL + "/live/"
	Settings.NATS_LIVENESS_TOPIC = NATS_LIVENESS_TOPIC

	Settings.SERVICE_NUMBER = os.Getenv("SERVICE_NUMBER")

}

// CheckSettingsNATS - проверяет наличие переменных окружения
func CheckSettingsNATS() error {
	var err error

	NATS_HOST := os.Getenv("NATS_HOST")
	NATS_PORT := os.Getenv("NATS_PORT")
	if NATS_HOST == "" {
		NATS_HOST = os.Getenv("BUS_LOCAL_HOST")
	}

	if NATS_PORT == "" {
		NATS_PORT = os.Getenv("BUS_LOCAL_PORT")
	}

	if NATS_HOST == "" {
		log.Error("Need fill BUS_LOCAL_HOST ! in os.ENV ")
		return err
	}
	if NATS_PORT == "" {
		log.Error("Need fill BUS_LOCAL_PORT ! in os.ENV ")
		return err
	}

	return err
}

// Start - Старт работы NATS Liveness
func Start(ServiceName string) {
	var err error

	//
	err = CheckSettingsNATS()
	if err != nil {
		return
	}

	//
	FillSettings(ServiceName)

	Connect()

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	Ticker = time.NewTicker(5 * time.Second)

	stopapp.GetWaitGroup_Main().Add(1)
	go SendMessages_go()

}

// Start_ctx - необходимые процедуры для подключения к серверу NATS
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто Start()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup, ServiceName string) error {
	var err error

	//запомним к себе контекст и WaitGroup
	contextmain.Ctx = ctx
	stopapp.SetWaitGroup_Main(WaitGroup)

	//
	err = CheckSettingsNATS()
	if err != nil {
		return err
	}

	//
	FillSettings(ServiceName)

	err = Connect_err()
	if err != nil {
		return err
	}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	Ticker = time.NewTicker(5 * time.Second)

	stopapp.GetWaitGroup_Main().Add(1)
	go SendMessages_go()

	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {
	defer stopapp.GetWaitGroup_Main().Done()

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. NATS_Liveness.")
	}

	CloseConnection()

}

// SendMessages_go - Отправляет сообщения каждые 5 секунд
func SendMessages_go() {
	defer stopapp.GetWaitGroup_Main().Done()

	for {
		select {
		case <-contextmain.GetContext().Done():
			return
		case <-Ticker.C:
			SendMessage()
		}
	}

}
