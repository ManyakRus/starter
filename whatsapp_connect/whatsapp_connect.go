package whatsapp_connect

//Licenses:
//Whatsapp клиент - go.mau.fi/whatsmeow" - MPL-2.0
//QR-код - github.com/mdp/qrterminal/v3" - MIT License
//SQL сервер - "github.com/mattn/go-sqlite3" - MIT License

import (
	"context"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/logger"
	"go.mau.fi/whatsmeow/types/events"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"
)

// clientWhatsApp - клиент соединения мессенджера Whatsapp
var clientWhatsApp *whatsmeow.Client

// log - глобальный логгер приложения
var log = logger.GetLog()

// filenameDB - имя файла локально базы данных sqllite
var filenameDB string

// MaxSendMessageCountIn1Second - максимальное количество сообщений в 1 секунду
var MaxSendMessageCountIn1Second float32 = 0.1

// lastSendTime - время последней отправки сообщения и мьютекс
var lastSendTime = lastSendTimeMutex{}

// lastSendTimeMutex - структура хранения времени последней отправки и мьютекс
type lastSendTimeMutex struct {
	time time.Time
	sync.Mutex
}

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	WHATSAPP_PHONE_FROM      string
	WHATSAPP_PHONE_SEND_TEST string
}

// MessageWhatsapp - сообщение из WhatsApp сокращённо
type MessageWhatsapp struct {
	Text      string
	NameFrom  string
	PhoneFrom string
	PhoneChat string
	IsFromMe  bool
	MediaType string
	//NameTo    string
	IsGroup  bool
	ID       string
	TimeSent time.Time
}

func FillMessageWhatsapp(mess *events.Message) MessageWhatsapp {
	Otvet := MessageWhatsapp{}

	Otvet.ID = mess.Info.ID
	Otvet.PhoneFrom = mess.Info.Sender.User
	Otvet.IsFromMe = mess.Info.IsFromMe
	Otvet.MediaType = mess.Info.MediaType
	Otvet.NameFrom = mess.Info.PushName
	Otvet.IsGroup = mess.Info.IsGroup
	if mess.Message != nil && mess.Message.Conversation != nil {
		//простое сообщение
		Otvet.Text = *mess.Message.Conversation
	} else if mess.Message != nil && mess.Message.ExtendedTextMessage != nil {
		//сообщение ответ
		Otvet.Text = *mess.Message.ExtendedTextMessage.Text
	}

	Otvet.PhoneChat = mess.Info.Chat.User
	Otvet.TimeSent = mess.Info.Timestamp

	return Otvet
}

// SendMessage - отправка сообщения в мессенджер Телеграм
// возвращает:
// id = id отправленного сообщения в WhatsApp
// err = error
func SendMessage(phone_send_to string, text string) (string, error) {
	var id string
	//var is_sent bool

	TimeLimit()
	log.Debug("phone_send_to: ", phone_send_to, " text: "+text)
	//

	ctxMain := context.Background()
	ctx, cancel := context.WithTimeout(ctxMain, 120*time.Second)
	defer cancel()

	recipient, ok := ParseJID(phone_send_to)
	if !ok {
		text1 := "ParseJID() invalid JID: " + phone_send_to
		log.Error(text1)
		return id, errors.New(text1)
	}

	MessageID := whatsmeow.GenerateMessageID()
	message1 := &waProto.Message{}
	message1.Conversation = &text
	_, err := clientWhatsApp.SendMessage(ctx, recipient, message1)
	if err != nil {
		text1 := "Message not sent, to: " + phone_send_to + " !"
		log.Error(text1)
		err = errors.New(text1)
		return "", err
	}

	id = string(MessageID)
	return id, nil
}

// eventHandler - получение событий из сервера whatsapp
func eventHandler_test(evt interface{}) {
	if evt == nil {
		log.Error("evt is null !")
	}
	switch v := evt.(type) {
	case *events.Message:
		mess := evt.(*events.Message)
		messW := FillMessageWhatsapp(mess)
		fmt.Println("Received a message from: ", messW.NameFrom, " phone: ", messW.PhoneFrom, "text: ", messW.Text)
		//fmt.Println("Received a message: ", mess.Message, " from: ", v.Message.GetContactMessage(), "text: ", v.Info.MediaType)
	default:
		fmt.Printf("Received: %#v \n", v)
	}
}

// Connect - создание клиента Whatsapp
func Connect(eventHandler func(evt interface{})) {
	err := Connect_err(eventHandler)
	if err != nil {
		log.Panic("WHATSAPP Connect_err() error: ", err)
	} else {
		log.Info("WHATSAPP connected. Phone from: ", Settings.WHATSAPP_PHONE_FROM)
	}
}

// Connect_err - создание клиента Whatsapp, и возвращает ошибку
func Connect_err(eventHandler func(evt interface{})) error {

	if Settings.WHATSAPP_PHONE_FROM == "" {
		FillSettings()
	}

	//ProgramDir := programdir.ProgramDir()
	ProgramDir := micro.ProgramDir_Common()
	filenameDB = ProgramDir + "whatsapp.db"

	dbLog := waLog.Stdout("Database", "WARN", true)
	container, err := sqlstore.New("sqlite3", "file:"+filenameDB+"?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Panicln(err)
	}
	clientLog := waLog.Stdout("Client", "WARN", true)
	clientWhatsApp = whatsmeow.NewClient(deviceStore, clientLog)
	clientWhatsApp.AddEventHandler(eventHandler)

	if clientWhatsApp.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := clientWhatsApp.GetQRChannel(context.Background())
		err = clientWhatsApp.Connect()
		if err != nil {
			log.Panicln(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				log.Println("QR code:", evt.Code)
			} else {
				log.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = clientWhatsApp.Connect()
		if err != nil {
			log.Panicln(err)
		}
	}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	//StopWhatsApp()
	return err
}

// StopWhatsApp - остановка работы клиента мессенджера Whatsapp
func StopWhatsApp() {
	clientWhatsApp.Disconnect()

}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled.")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("whatsapp")

	//
	StopWhatsApp()
	stopapp.GetWaitGroup_Main().Done()
}

// ParseJID parses a JID out of the given string. It supports both regular and AD JIDs.
func ParseJID(arg string) (types.JID, bool) {
	if arg == "" {
		return types.NewJID(arg, types.DefaultUserServer), false
	}

	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			_ = fmt.Errorf("Invalid JID %s: %v", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			_ = fmt.Errorf("Invalid JID %s: no server specified", arg)
			return recipient, false
		}
		return recipient, true
	}
}

// TimeLimit - пауза для ограничения количество сообщений в секунду
func TimeLimit() {
	//if MaxSendMessageCountIn1Second == 0 {
	//	return
	//}

	lastSendTime.Lock()
	defer lastSendTime.Unlock()

	if lastSendTime.time.IsZero() {
		lastSendTime.time = time.Now()
		return
	}

	t := time.Now()
	ms := int(t.Sub(lastSendTime.time).Milliseconds())
	msNeedWait := int(1000 / MaxSendMessageCountIn1Second)
	if ms < msNeedWait {
		micro.Sleep(msNeedWait - ms)
	}

	lastSendTime.time = time.Now()
}

// FillSettings загружает переменные окружения в структуру из файла или из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.WHATSAPP_PHONE_FROM = os.Getenv("WHATSAPP_PHONE_FROM")
	Settings.WHATSAPP_PHONE_SEND_TEST = os.Getenv("WHATSAPP_PHONE_SEND_TEST")

	if Settings.WHATSAPP_PHONE_FROM == "" {
		log.Panicln("Need fill WHATSAPP_PHONE_FROM ! in os.ENV ")
	}

	if Settings.WHATSAPP_PHONE_SEND_TEST == "" {
		log.Panicln("Need fill WHATSAPP_PHONE_SEND_TEST ! in os.ENV ")
	}

}

// Start - делает соединение с БД, отключение и др.
func Start(eventHandler func(evt interface{})) {
	Connect(eventHandler)

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

}

func (m MessageWhatsapp) String() string {
	Otvet := ""

	Otvet = Otvet + fmt.Sprint("Text: ", m.Text, "\n")
	Otvet = Otvet + fmt.Sprint("NameFrom: ", m.NameFrom, "\n")
	Otvet = Otvet + fmt.Sprint("PhoneFrom: ", m.PhoneFrom, "\n")
	Otvet = Otvet + fmt.Sprint("PhoneChat: ", m.PhoneChat, "\n")
	Otvet = Otvet + fmt.Sprint("IsFromMe: ", m.IsFromMe, "\n")
	Otvet = Otvet + fmt.Sprint("MediaType: ", m.MediaType, "\n")
	Otvet = Otvet + fmt.Sprint("IsGroup: ", m.IsGroup, "\n")
	Otvet = Otvet + fmt.Sprint("ID: ", m.ID, "\n")
	Otvet = Otvet + fmt.Sprint("TimeSent: ", m.TimeSent, "\n")

	return Otvet
}
