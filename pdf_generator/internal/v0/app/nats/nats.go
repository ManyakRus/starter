package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/ManyakRus/starter/common/v0/nats_connect"
	"github.com/ManyakRus/starter/pdf_generator/internal/v0/app/create_file"
	"github.com/ManyakRus/starter/pdf_generator/internal/v0/app/types"
	"time"

	"github.com/ManyakRus/starter/common/v0/contextmain"
	"github.com/ManyakRus/starter/common/v0/logger"
	"github.com/ManyakRus/starter/common/v0/micro"
	"github.com/ManyakRus/starter/common/v0/stopapp"
)

// log - глобальный логгер
var log = logger.GetLog()

// NatsSubscriptionIn - подписка на нужный топик
var NatsSubscriptionIn *nats.Subscription

// PAUSE_WAIT_NEW_MESSAGE - количество миллисекунд паузы обращения к серверу NATS
const PAUSE_WAIT_NEW_MESSAGE = 1000

var TOPIC_PDF_IN = "/claim/pdf_generator/in/"

// StartNats - необходимые процедуры для подключения к серверу Nats
func StartNats() {
	SubscribeTopics()

	stopapp.GetWaitGroup_Main().Add(1)
	go ListenForever()

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

}

// ListenForever - бесконечный цикл получения данных с сервера Nats
func ListenForever() {
	var err error

loop:
	for {
		select {
		case <-contextmain.GetContext().Done():
			log.Warn("Context app is canceled.")
			break loop
		default:
			Message1, haveNew := ReceiveMessageFromNats()
			if haveNew == false {
				micro.Pause(PAUSE_WAIT_NEW_MESSAGE)
				continue
			}

			err = create_file.StartCreateFile(Message1)
			if err != nil {
				log.Error(err)
				continue
			}

		}
	}

	stopapp.GetWaitGroup_Main().Done()
}

// ReceiveMessageFromNats - получение 1 сообщения с сервера Nats
// возвращает true если есть сообщение
func ReceiveMessageFromNats() (types.MessageNatsIn, bool) {
	//Otvet := false

	Message1, err := NextMessageRequest()
	if err != nil {
		log.Error(err)
		return Message1, false
	}

	if Message1.Head.Sender == "" {
		ErrorText := "Message.Head.Sender ='' !"
		log.Error(ErrorText)
		return Message1, false
	}

	return Message1, true
}

// NextMessageRequest - получает следующее 1 сообщение с сервера Nats
func NextMessageRequest() (types.MessageNatsIn, error) {
	var err error
	Message1 := types.MessageNatsIn{}

	if NatsSubscriptionIn == nil {
		SubscribeTopics()
	}

	ctxMain := contextmain.GetContext()
	ctx, cancel := context.WithTimeout(ctxMain, 20*time.Second)
	defer cancel()

	MessageNats1, err := NatsSubscriptionIn.NextMsgWithContext(ctx)
	if err == context.DeadlineExceeded {
		return Message1, err
	} else if err != nil {
		text1 := fmt.Sprintln("NextMsgWithContext() error: ", err)
		log.Error(text1)
		err = errors.New(text1)
		return Message1, err
	}

	//err = MessageNats1.InProgress()
	//if err != nil {
	//	text1 := fmt.Sprintln("MessageNats1.InProgress() error: ", err)
	//	log.Error(text1)
	//	err = errors.New(text1)
	//	return Message1, err
	//}

	//var objmap map[string]json.RawMessage
	//err = json.Unmarshal(MessageNats1.Data, &objmap)

	err = json.Unmarshal(MessageNats1.Data, &Message1)
	if err != nil {
		text1 := fmt.Sprintln("Unmarshal() error: ", err)
		log.Error(text1)
		err = errors.New(text1)
		return Message1, err
	}

	//err = MessageNats1.AckSync()
	//if err != nil {
	//	text1 := fmt.Sprintln("MessageNats1.AckSync() error: ", err)
	//	log.Error(text1)
	//	err = errors.New(text1)
	//	return Message1, err
	//}

	return Message1, err
}

// SubscribeTopics - подписывается на топики nats
func SubscribeTopics() {
	var err error

	//topic pdf/in
	NatsSubscriptionIn, err = nats_connect.Conn.SubscribeSync(TOPIC_PDF_IN)

	if err != nil {
		log.Panicln("SubscribeTopics() TOPIC_PDF_IN: ", TOPIC_PDF_IN, "error: ", err)
	} else {
		log.Info("SubscribeTopics() OK. URL: ", nats_connect.Settings.NATS_SERVER, " Subject:", TOPIC_PDF_IN)
	}

}

// WaitStop - ожидает отмену глобального контекста или сигнала завершения приложения
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. NATS.")
	}

	//err := CloseConnection()
	//if err != nil {
	//	log.Error("CloseConnection() error: ", err)
	//}
	stopapp.GetWaitGroup_Main().Done()
}
