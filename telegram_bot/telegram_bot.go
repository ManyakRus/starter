package telegram_bot

import (
	"context"
	"errors"
	"fmt"
	//"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"sync"
)

// PackageName - имя текущего пакета, для логирования
const PackageName = "telegram_bot"

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	TELEGRAM_API_KEY      string
	TELEGRAM_CHAT_ID_TEST string
}

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// Client - клиент подклюенный к Telegram
var Client *botapi.BotAPI

// mutexReconnect - защита от многопоточности Reconnect()
var mutexReconnect = &sync.Mutex{}

// NeedReconnect - флаг необходимости переподключения
var NeedReconnect bool

// Connect - подключается к Telegram, или паника при ошибке
func Connect() {

	err := Connect_err()
	LogInfo_Connected(err)

}

// Connect_err - подключается к Telegram, возвращает ошибку
func Connect_err() error {
	var err error

	if Settings.TELEGRAM_API_KEY == "" {
		FillSettings()
	}

	Client, err = botapi.NewBotAPI(Settings.TELEGRAM_API_KEY)
	if err != nil {
		return err
	}

	Client.Debug = false

	//log.Printf("Authorized on account %s", bot.Self.UserName)
	//
	//u := botapi.NewUpdate(0)
	//u.Timeout = 60
	//
	//updates := bot.GetUpdatesChan(u)
	//var Users []User
	//Users = make([]User, 0)
	//
	//for update := range updates {
	//	if update.Message == nil { // ignore any non-Message Updates
	//		continue
	//	}
	//
	//	//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	//
	//	TextFromUser := update.Message.Text
	//
	//	Text := "Not understand. Fill name or phone number. "
	//	if valid.IsInt(TextFromUser) == true {
	//		Users = FindUser_by_Phone(TextFromUser)
	//		if len(Users) == 0 {
	//			Users = FindUser_by_CellPhone(TextFromUser)
	//		}
	//	} else if HaveAt(TextFromUser) == true {
	//		Users = FindUser_by_Email(TextFromUser)
	//	} else if valid.IsInt(TextFromUser) == false && HaveNumbers(TextFromUser) == false {
	//		Users = FindUser_by_Name(TextFromUser)
	//		if len(Users) == 0 {
	//			Users = FindUser_by_Post(TextFromUser)
	//		}
	//	} else {
	//		Users = FindUser_by_Adress(TextFromUser)
	//	}
	//
	//	if len(Users) > 0 {
	//		Text = ""
	//		for _, User1 := range Users {
	//			Text = Text + User1.String() + "\n"
	//		}
	//	}
	//
	//	if len(Text) > 2000 {
	//		Text = Text[0:2000]
	//		Text = Text + "\n" + "..."
	//		Text = Text + "\n" + "..."
	//		Text = Text + "\n" + "..."
	//	}
	//
	//	msg := botapi.NewMessage(update.Message.Chat.ID, Text)
	//	msg.ReplyToMessageID = update.Message.MessageID
	//	//msg.Entities = append(msg.Entities, )
	//	//append(msg.Entities, )
	//
	//	bot.Send(msg)
	//}

	return err
}

// SendMessageChatID - отправка сообщения в мессенджер Телеграм
// возвращает:
// id = id отправленного сообщения в telegram
// err = error
func SendMessageChatID(ChatID int64, Text string) (int, error) {
	var ID int
	var err error

	//
	if Client == nil {
		err = fmt.Errorf("error: telegram Client == nil")
		return ID, err
	}

	msg := botapi.NewMessage(ChatID, Text)
	msg.ParseMode = "HTML"

	Message, err := Client.Send(msg)
	if err != nil {
		return ID, err
	}
	ID = Message.MessageID

	return ID, err
}

// SendMessage - отправка сообщения в мессенджер Телеграм
// возвращает:
// id = id отправленного сообщения в telegram
// err = error
func SendMessage(UserName string, Text string) (int, error) {
	var ID int
	var err error

	//
	if Client == nil {
		err = fmt.Errorf("error: telegram Client == nil")
		return ID, err
	}

	//экранируем запрещённые символы
	//Text = html.EscapeString(Text)

	//
	msg := botapi.NewMessageToChannel(UserName, Text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true

	Message, err := Client.Send(msg)
	if err != nil {
		return ID, err
	}
	ID = Message.MessageID

	return ID, err
}

// Reconnect повторное подключение к Telegram, если оно отключено
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

	if Client == nil {
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

	////остановим программу т.к. она не должна работать при неработающеё БД
	//log.Error("STOP app. Error: ", err)
	//stopapp.StopApp()

}

// CloseConnection - закрытие соединения с Telegram
func CloseConnection() {
	if Client == nil {
		return
	}

	err := CloseConnection_err()
	if err != nil {
		log.Error("Telegram bot CloseConnection() error: ", err)
	} else {
		log.Info("Telegram bot connection closed")
	}

	return
}

// CloseConnection - закрытие соединения с Telegram, возвращает ошибку
func CloseConnection_err() error {
	var err error
	if Client == nil {
		return err
	}

	Client.StopReceivingUpdates()
	//Client = nil

	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {
	defer waitGroup_Connect.Done()

	select {
	case <-(*ctx_Connect).Done():
		log.Warn("Context app is canceled. telegram_bot")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("telegram_bot")

	//
	CloseConnection()

}

// Start_ctx - необходимые процедуры для подключения к серверу Telegram
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто StartDB()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup) error {
	var err error

	//запомним к себе контекст
	//	if contextmain.Ctx != ctx {
	//		contextmain.SetContext(ctx)
	//	}
	//contextmain.Ctx = ctx
	if ctx == nil {
		ctx = ctx_Connect
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

	return err
}

// Start - делает соединение с Telegram, отключение и др.
func Start() {
	err := Connect_err()
	if err != nil {
		log.Panic("telegram_bot Start() error: ", err)
	}

	//сохраним в список подключений
	ctx := ctx_Connect
	WaitGroupContext1 := stopapp.WaitGroupContext{WaitGroup: waitGroup_Connect, Ctx: ctx, CancelCtxFunc: cancelCtxFunc}
	stopapp.OrderedMapConnections.Put(PackageName, WaitGroupContext1)

	//
	waitGroup_Connect.Add(1)
	go WaitStop()

}

// FillSettings загружает переменные окружения в структуру из файла или из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}

	// заполним из переменных оуружения
	Settings.TELEGRAM_API_KEY = os.Getenv("TELEGRAM_API_KEY")

	//
	Name := ""
	s := ""

	//
	Name = "TELEGRAM_API_KEY"
	s = Getenv(Name, true)
	Settings.TELEGRAM_API_KEY = s

	Name = "TELEGRAM_CHAT_ID_TEST"
	s = Getenv(Name, false)
	Settings.TELEGRAM_CHAT_ID_TEST = s

}

// Getenv - возвращает переменную окружения
func Getenv(Name string, IsRequired bool) string {
	TextError := "Need fill OS environment variable: "
	Otvet := os.Getenv(Name)
	if IsRequired == true && Otvet == "" {
		log.Error(TextError + Name)
	}

	return Otvet
}

// GetConnection - возвращает соединение к нужной базе данных
func GetConnection() *botapi.BotAPI {
	if Client == nil {
		Connect()
	}

	return Client
}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicf("Telegram bot not connected with API KEY: %s, error: %v", Settings.TELEGRAM_API_KEY, err)
	} else {
		log.Infof("Telegram bot Connected. With API KEY: %s", Settings.TELEGRAM_API_KEY)
	}

}
