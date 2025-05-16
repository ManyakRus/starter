package kafka_connect

import (
	"context"
	"fmt"
	"github.com/ManyakRus/starter/log"
	"net"
	"os"
	"sync"
	"time"

	//"github.com/ManyakRus/starter/common/v0/micro"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/stopapp"

	"github.com/segmentio/kafka-go"
)

// Conn - соединение к серверу nats
var Conn *kafka.Conn

// log - глобальный логгер
//var log = logger.GetLog()

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	KAFKA_HOST     string
	KAFKA_PORT     string
	KAFKA_LOGIN    string
	KAFKA_PASSWORD string
}

// Client - клиент для Kafka
var Client *kafka.Client

// Connect - подключается к серверу Kafka
func Connect() {
	var err error

	err = Connect_err()
	LogInfo_Connected(err)
}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicln("KAFKA Connect() host: ", Settings.KAFKA_HOST, " error: ", err)
	} else {
		log.Info("KAFKA Connect() OK, host: ", Settings.KAFKA_HOST)
	}

}

// Connect_err - подключается к серверу Kafka и возвращает ошибку
func Connect_err() error {
	var err error

	if Settings.KAFKA_HOST == "" {
		FillSettings()
	}

	//sKAFKA_PORT := (Settings.KAFKA_PORT)
	//URL := "nats://" + Settings.KAFKA_HOST + ":" + sKAFKA_PORT
	//UserInfo := nats.UserInfo(Settings.KAFKA_LOGIN, Settings.KAFKA_PASSWORD)
	Conn, err = kafka.Dial("tcp", Settings.KAFKA_HOST+":"+Settings.KAFKA_PORT)

	//
	err = CreateClient()
	if err != nil {
		return err
	}

	//nats.ManualAck()
	return err
}

// CreateClient - создаёт клиент для Kafka
func CreateClient() error {
	var err error

	Client = &kafka.Client{}
	Client.Addr = GetAddr()

	return err
}

// GetAddr - создаёт Addr
func GetAddr() net.Addr {
	URL := Settings.KAFKA_HOST + ":" + Settings.KAFKA_PORT
	Otvet := kafka.TCP(URL)

	return Otvet
}

// StartKafka - необходимые процедуры для подключения к серверу Kafka
func StartKafka() {
	var err error

	ctx := contextmain.GetContext()
	WaitGroup := stopapp.GetWaitGroup_Main()
	err = Start_ctx(&ctx, WaitGroup)
	LogInfo_Connected(err)

}

// Start_ctx - необходимые процедуры для подключения к серверу Kafka
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто StartKafka()
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

	return err
}

// CloseConnection - закрывает соединение с сервером Kafka
func CloseConnection() {
	var err error

	if Conn == nil {
		return
	}

	err = Conn.Close()
	if err != nil {
		log.Error("KAFKA CloseConnection() error: ", err)
	} else {
		log.Info("KAFKA stopped")
	}

	//
	Client = nil

	return
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {
	defer stopapp.GetWaitGroup_Main().Done()

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. kafka_connect")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("KAFKA_connect")

	//
	CloseConnection()

}

// FillSettings загружает переменные окружения в структуру из файла или из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.KAFKA_HOST = os.Getenv("KAFKA_HOST")
	Settings.KAFKA_PORT = os.Getenv("KAFKA_PORT")
	Settings.KAFKA_LOGIN = os.Getenv("KAFKA_LOGIN")
	Settings.KAFKA_PASSWORD = os.Getenv("KAFKA_PASSWORD")

	if Settings.KAFKA_HOST == "" {
		log.Panicln("Need fill KAFKA_HOST ! in os.ENV ")
	}

	if Settings.KAFKA_PORT == "" {
		log.Panicln("Need fill KAFKA_PORT ! in os.ENV ")
	}

	//if Settings.KAFKA_LOGIN == "" {
	//	log.Panicln("Need fill KAFKA_LOGIN ! in os.ENV ")
	//}
	//
	//if Settings.KAFKA_PASSWORD == "" {
	//	log.Panicln("Need fill KAFKA_PASSWORD ! in os.ENV ")
	//}

	//
}

// ConnectTopic - подключает кафку к нужному топику
func ConnectTopic(TopicName, GroupID string) *kafka.Reader {

	// make a new reader that consumes from topic
	KafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{Settings.KAFKA_HOST + ":" + Settings.KAFKA_PORT},
		GroupID:  GroupID,
		Topic:    TopicName,
		MinBytes: 10,   // 10KB 10e3
		MaxBytes: 10e6, // 10MB
	})

	return KafkaReader
}

// GetOffsetFromGroupID - получает оффсет группы для конкретного топика, партиция 0
func GetOffsetFromGroupID(TopicName, GroupID string) (int64, error) {
	var Otvet int64 = 0
	var err error

	//
	ctxMain := contextmain.GetContext()
	ctx, ctxCancelFunc := context.WithTimeout(ctxMain, time.Duration(60)*time.Second)
	defer ctxCancelFunc()

	//
	PartitionNumber := 0
	MapTopics := make(map[string][]int)
	MapTopics[TopicName] = []int{PartitionNumber}

	//
	Addr := GetAddr()
	OFR := kafka.OffsetFetchRequest{}
	OFR.Addr = Addr
	OFR.GroupID = GroupID
	OFR.Topics = MapTopics

	//
	Response, err := Client.OffsetFetch(ctx, &OFR)
	if err != nil {
		err = fmt.Errorf("OffsetFetch() error: %w", err)
		return Otvet, err
	}

	//
	MassOffset := Response.Topics[TopicName]
	if len(MassOffset) != 1 {
		err = fmt.Errorf("GetOffsetFromGroupID() error: len(MassOffset) != 1")
		return Otvet, err
	}

	//
	Otvet = MassOffset[0].CommittedOffset

	return Otvet, err
}
