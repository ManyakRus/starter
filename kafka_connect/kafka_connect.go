package kafka_connect

import (
	"context"
	"github.com/ManyakRus/starter/logger"
	"os"
	"sync"

	//"github.com/ManyakRus/starter/common/v0/micro"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/stopapp"

	"github.com/segmentio/kafka-go"
)

// Conn - соединение к серверу nats
var Conn *kafka.Conn

// log - глобальный логгер
var log = logger.GetLog()

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	KAFKA_HOST     string
	KAFKA_PORT     string
	KAFKA_LOGIN    string
	KAFKA_PASSWORD string
}

// Connect - подключается к серверу Nats
func Connect() {
	var err error

	err = Connect_err()
	if err != nil {
		log.Panicln("KAFKA Connect() host: ", Settings.KAFKA_HOST, " error: ", err)
	} else {
		log.Info("KAFKA Connect() OK, host: ", Settings.KAFKA_HOST)
	}
}

// Connect_err - подключается к серверу Nats и возвращает ошибку
func Connect_err() error {
	var err error

	if Settings.KAFKA_HOST == "" {
		FillSettings()
	}

	//sKAFKA_PORT := (Settings.KAFKA_PORT)
	//URL := "nats://" + Settings.KAFKA_HOST + ":" + sKAFKA_PORT
	//UserInfo := nats.UserInfo(Settings.KAFKA_LOGIN, Settings.KAFKA_PASSWORD)
	Conn, err = kafka.Dial("tcp", Settings.KAFKA_HOST+":"+Settings.KAFKA_PORT)

	//nats.ManualAck()
	return err
}

// StartKafka - необходимые процедуры для подключения к серверу Kafka
func StartKafka() {
	Connect()

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

}

// Start_ctx - необходимые процедуры для подключения к серверу Kafka
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто StartKafka()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup) error {
	var err error

	//запомним к себе контекст и WaitGroup
	contextmain.Ctx = ctx
	stopapp.SetWaitGroup_Main(WaitGroup)

	//
	err = Connect_err()
	if err != nil {
		return err
	}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	return err
}

// CloseConnection - закрывает соединение с сервером Nats
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

	return
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled.")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("KAFKA_connect")

	//
	CloseConnection()

	stopapp.GetWaitGroup_Main().Done()
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
