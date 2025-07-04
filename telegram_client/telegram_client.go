// модуль для использования Телеграмм Клиента (или бота)
package telegram_client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ManyakRus/starter/log"
	"github.com/gotd/contrib/storage"
	"github.com/gotd/td/telegram/message/peer"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"

	"github.com/gotd/td/clock"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/go-faster/errors"
	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/updates"

	pebbledb "github.com/cockroachdb/pebble"
	boltstor "github.com/gotd/contrib/bbolt"
	"github.com/gotd/contrib/pebble"
	lj "gopkg.in/natefinch/lumberjack.v2"
)

// filenameSession - имя файла сохранения сессии мессенджера Телеграм
var filenameSession string

// Client - клиент соединения мессенджера Телеграм
var Client *telegram.Client

// UserSelf - собственный юзер в Телеграм
var UserSelf *tg.User

// lastSendTime - время последней отправки сообщения и мьютекс
var lastSendTime = lastSendTimeMutex{}

//// log - глобальный логгер приложения
//var log = logger.GetLog()

// stopTelegramFunc - функция остановки соединения с мессенджером Телеграм
var stopTelegramFunc bg.StopFunc

// MAX_MESSAGE_LEN - максимальная длина сообщения
const MAX_MESSAGE_LEN = 4096

// MaxSendMessageCountIn1Second - максимальное количество сообщений в 1 секунду
var MaxSendMessageCountIn1Second float32 = 0.13 //0.13 =4 сообщения в секунду

// lastSendTimeMutex - структура хранения времени последней отправки и мьютекс
type lastSendTimeMutex struct {
	time time.Time
	sync.Mutex
}

// noSignUp can be embedded to prevent signing up.
type noSignUp struct{}

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
// TELEGRAM_APP_ID, TELEGRAM_APP_HASH - первоначально получить по ссылке: https://my.telegram.org/apps
// TELEGRAM_PHONE_FROM - номер телефона с которого отправляются сообщения
type SettingsINI struct {
	TELEGRAM_APP_ID          int
	TELEGRAM_APP_HASH        string
	TELEGRAM_PHONE_FROM      string
	TELEGRAM_PHONE_SEND_TEST string
}

// PeerDB - локальная база данных cockroachdb, для хранения кеша PeerID отправителей
var PeerDB *pebble.PeerStorage

// Func_OnNewMessage - функция из внешнего модуля, для приёма новых сообщений
var Func_OnNewMessage func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage, Peer1 storage.Peer) error

// Contacts - список контактов в Telegram
var Contacts *tg.ContactsContacts

// MessageTelegram - сообщение из Telegram сокращённо
type MessageTelegram struct {
	Text      string
	FromID    int64
	ChatID    int64
	IsFromMe  bool
	MediaType string
	//NameTo    string
	IsGroup  bool
	ID       int
	TimeSent time.Time
}

// String - возвращает строку из структуры
func (m MessageTelegram) String() string {
	Otvet := ""

	Otvet = Otvet + fmt.Sprint("Text: ", m.Text, "\n")
	Otvet = Otvet + fmt.Sprint("MediaType: ", m.MediaType, "\n")
	Otvet = Otvet + fmt.Sprint("FromID: ", m.FromID, "\n")
	Otvet = Otvet + fmt.Sprint("IsFromMe: ", m.IsFromMe, "\n")
	Otvet = Otvet + fmt.Sprint("IsGroup: ", m.IsGroup, "\n")
	Otvet = Otvet + fmt.Sprint("ID: ", m.ID, "\n")
	Otvet = Otvet + fmt.Sprint("TimeSent: ", m.TimeSent, "\n")

	return Otvet
}

// SendMessage - отправка сообщения в мессенджер Телеграм
// возвращает:
// id = id отправленного сообщения в telegram
// err = error
func SendMessage(phone_send_to string, text string) (int, error) {
	var id int

	if Client == nil {
		CreateTelegramClient(nil)
	}

	if text == "" {
		err := fmt.Errorf("SendMessage() text id empty ! ")
		log.Error(err)
		return 0, err
	}

	if phone_send_to == "" {
		err := fmt.Errorf("SendMessage() phone_send_to id empty ! ")
		log.Error(err)
		return 0, err
	}

	TimeLimit()
	log.Debug("phone_send_to: ", phone_send_to, ", text: "+text)

	text = micro.SubstringLeft(text, MAX_MESSAGE_LEN)

	//
	api := Client.API()

	//
	ctxMain := context.Background()
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()

	//
	sender := message.NewSender(api)
	target := sender.Resolve(phone_send_to)
	target.NoForwards()

	//отправка сообщения
	UpdatesClass, err := target.Text(ctx, text)
	if err != nil {
		log.Error("Text() error: ", err)
		return 0, err
	}

	//
	if UpdatesClass != nil {
		id = findIdFromUpdatesClass(UpdatesClass)
	}

	return id, err
}

// SendStyledText - отправка сообщения в мессенджер Телеграм, все слова должны быть оформлены в StyledTextOption
// возвращает:
// id = id отправленного сообщения в telegram
// err = error
func SendStyledText(phone_send_to string, texts ...message.StyledTextOption) (int, error) {
	var id int

	if Client == nil {
		CreateTelegramClient(nil)
	}

	if len(texts) == 0 {
		err := fmt.Errorf("SendMessage() text id empty ! ")
		log.Error(err)
		return 0, err
	}

	if phone_send_to == "" {
		err := fmt.Errorf("SendMessage() phone_send_to id empty ! ")
		log.Error(err)
		return 0, err
	}

	TimeLimit()
	//text := texts.
	//log.Debug("phone_send_to: ", phone_send_to, ", text: "+ text)

	//text = micro.SubstringLeft(text, MAX_MESSAGE_LEN)

	//
	api := Client.API()

	//
	ctxMain := context.Background()
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()

	//
	sender := message.NewSender(api)
	target := sender.Resolve(phone_send_to)
	target.NoForwards()

	//отправка сообщения
	UpdatesClass, err := target.StyledText(ctx, texts...)
	if err != nil {
		log.Error("Text() error: ", err)
		return 0, err
	}

	//
	if UpdatesClass != nil {
		id = findIdFromUpdatesClass(UpdatesClass)
	}

	return id, err
}

// AddContact - добавляет новый контакт в список контактов Телеграм
func AddContact(ctx context.Context, phone_send_to string) error {
	var err error

	if phone_send_to == "" {
		text1 := "phone_send_to='' !"
		err := errors.New(text1)
		log.Error(text1)
		return err
	}

	TimeLimit()

	api := Client.API()

	//var contacts []tg.InputPhoneContact
	contact := tg.InputPhoneContact{}
	contact.Phone = phone_send_to
	contact.FirstName = phone_send_to

	contacts := make([]tg.InputPhoneContact, 1)
	contacts = append(contacts, contact)

	ContactsImportedContacts, err := api.ContactsImportContacts(ctx, contacts)
	if ContactsImportedContacts == nil {
		text1 := "ContactsImportedContacts == nil. Не удалось добавить контакт !"
		err = errors.New(text1)
		log.Error(text1)
	} else if ContactsImportedContacts.Imported == nil {
		text1 := "ContactsImportedContacts.Imported =nil. Не удалось добавить контакт !"
		err = errors.New(text1)
		log.Error(text1)
	} else if len(ContactsImportedContacts.Imported) == 0 {
		text1 := "ContactsImportedContacts.Imported len=0. Не удалось добавить контакт !"
		err = errors.New(text1)
		log.Error(text1)
	}

	return err
}

// SignUp - обязательная функция клиента Телеграм
func (noSignUp) SignUp(context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

// AcceptTermsOfService - обязательная функция клиента Телеграм
func (noSignUp) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	if ctx == nil {
		text1 := "telegramclient.AcceptTermsOfService() error: Context=nil"
		return errors.New(text1)
	}

	return &auth.SignUpRequired{TermsOfService: tos}
}

// termAuth implements authentication via terminal.
type termAuth struct {
	noSignUp

	phone string
}

// Phone - обязательная функция клиента Телеграм
func (a termAuth) Phone(_ context.Context) (string, error) {
	return a.phone, nil
}

// Password - обязательная функция клиента Телеграм
// ввод пароля с терминала
func (a termAuth) Password(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", nil
	}

	fmt.Print("Enter 2FA password: ")
	bytePwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePwd)), nil
}

// Code - обязательная функция клиента Телеграм
// ввод кода CODE с терминала
func (a termAuth) Code(ctx context.Context, _ *tg.AuthSentCode) (string, error) {
	if ctx == nil {
		return "", nil
	}

	//Stdin, _ := io.Pipe() //нужен т.к. не работает в тест

	//r, _ := io.Pipe()
	//scanner := bufio.NewScanner(r)
	//msg := "Enter code from telegram: "
	//fmt.Fprintln(os.Stdout, msg)
	//
	//scanner.Scan()
	//if err := scanner.Err(); err != nil {
	//	log.Fatal(err)
	//}
	//code := scanner.Text()
	//if len(code) == 0 {
	//	log.Fatal("empty input")
	//}

	fmt.Print("Enter code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	//code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}

// memorySession implements in-memory session storage.
// Goroutine-safe.
type memorySession struct {
	mux  sync.RWMutex
	data []byte
}

// LoadSession loads session from memory.
func (s *memorySession) LoadSession(context.Context) ([]byte, error) {
	if s == nil {
		return nil, session.ErrNotFound
	}

	s.mux.RLock()
	defer s.mux.RUnlock()

	// read the whole file at once
	cpy, err := os.ReadFile(filenameSession)
	if err != nil {
		cpy = nil
		log.Error(err)
		return nil, nil
		//return nil, session.ErrNotFound
	}

	return cpy, nil
}

// StoreSession stores session to memory.
func (s *memorySession) StoreSession(ctx context.Context, data []byte) error {
	if ctx == nil {
		text1 := "telegramclient.StoreSession() error: Context=nil"
		return errors.New(text1)
	}

	s.mux.Lock()
	defer s.mux.Unlock()
	//s.data = data

	// write the whole body at once
	err := os.WriteFile(filenameSession, data, 0644)
	if err != nil {
		panic(err)
	}

	return nil
}

// CreateTelegramClient создание клиента Телеграм, без подключения
// лучше использовать Connect() - с подключением
func CreateTelegramClient(func_OnNewMessage func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage, Peer1 storage.Peer) error) {
	err := CreateTelegramClient_err(func_OnNewMessage)
	if err != nil {
		log.Panic("CreateTelegramClient() error: ", err)
	} else {
		log.Info("CreateTelegramClient() ok, phone from: ", Settings.TELEGRAM_PHONE_FROM)
	}
	return
}

// CreateTelegramClient_err создание клиента Телеграм, без подключения
// лучше использовать Connect_err() - с подключением
func CreateTelegramClient_err(func_OnNewMessage func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage, Peer1 storage.Peer) error) error {

	//загрузим настройки, если они пустые
	if Settings.TELEGRAM_APP_ID == 0 {
		FillSettings()
	}

	//заполним глобальную функцию получения нового сообщения
	Func_OnNewMessage = func_OnNewMessage

	//
	programDir := micro.ProgramDir()

	//создадим логгер zap
	sessionDir := programDir
	LogDir := filepath.Join(sessionDir, "log")
	logFilePath := filepath.Join(LogDir, "log.jsonl")
	ok, err := micro.FileExists(LogDir)
	if ok == false {
		err = fmt.Errorf("FileExists() error: need create folder: %v", LogDir)
		log.Error(err)
		err := micro.CreateFolder(LogDir, 0600) //rw------- т.к. файл секретный
		if err != nil {
			log.Error("CreateFolder() error: ", err)
			return err
		}
		//return err
	}

	if err != nil {
		err = fmt.Errorf("FileExists() error: %w", err)
		log.Debug(err)
		return err
	}

	// Log to file, so we don't interfere with prompts and messages to user.
	logWriter := zapcore.AddSync(&lj.Logger{
		Filename:   logFilePath,
		MaxBackups: 3,
		MaxSize:    1, // megabytes
		MaxAge:     7, // days
	})
	logCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		logWriter,
		zap.DebugLevel,
	)
	lg := zap.New(logCore)

	//настройка хранения текущей сессии в файле
	sessionStorage := &telegram.FileSessionStorage{
		Path: filepath.Join(programDir, "session.txt"),
	}

	// Peer storage, for resolve caching and short updates handling.
	db, err := pebbledb.Open(filepath.Join(programDir, "peers.pebble.db"), &pebbledb.Options{})
	if err != nil {
		return errors.Wrap(err, "create pebble storage")
	}
	PeerDB = pebble.NewPeerStorage(db)

	// Setting up client.
	//
	// Dispatcher is used to register handlers for events.
	dispatcher := tg.NewUpdateDispatcher()
	// Setting up update handler that will fill peer storage before
	// calling dispatcher handlers.
	updateHandler := storage.UpdateHook(dispatcher, PeerDB)

	// Setting up persistent storage for qts/pts to be able to
	// recover after restart.
	boltdb, err := bbolt.Open(filepath.Join(sessionDir, "updates.bolt.db"), 0666, nil)
	if err != nil {
		return errors.Wrap(err, "create bolt storage")
	}

	updatesRecovery := updates.New(updates.Config{
		Handler: updateHandler, // using previous handler with PeerDB
		Logger:  lg.Named("updates.recovery"),
		Storage: boltstor.NewStateStorage(boltdb),
	})

	//// Handler of FLOOD_WAIT that will automatically retry request.
	//waiter := floodwait.NewWaiter().WithCallback(func(ctx context.Context, wait floodwait.FloodWait) {
	//	// Notifying about flood wait.
	//	lg.Warn("Flood wait", zap.Duration("wait", wait.Duration))
	//	fmt.Println("Got FLOOD_WAIT. Will retry after", wait.Duration)
	//})

	// Filling client options.
	options := telegram.Options{
		Logger:         lg,              // Passing logger for observability.
		SessionStorage: sessionStorage,  // Setting up session sessionStorage to store auth data.
		UpdateHandler:  updatesRecovery, // Setting up handler for updates from server.
		//Middlewares: []telegram.Middleware{
		//	// Setting up FLOOD_WAIT handler to automatically wait and retry request.
		//	waiter,
		//	// Setting up general rate limits to less likely get flood wait errors.
		//	ratelimit.New(rate.Every(time.Millisecond*100), 5),
		//},
	}

	Client = telegram.NewClient(Settings.TELEGRAM_APP_ID, Settings.TELEGRAM_APP_HASH, options)

	if func_OnNewMessage != nil {
		dispatcher.OnNewMessage(OnNewMessage)
	}

	api := Client.API()

	// Setting up resolver cache that will use peer storage.
	//не знаю зачем, удалить
	resolver := storage.NewResolverCache(peer.Plain(api), PeerDB)
	_ = resolver

	return err
}

// OnNewMessage - функция для получения новых сообщений, и перенаправления во внешний Func_OnNewMessage()
// чтоб отправить сообщение надо:
// api := telegram_client.Client.API()
// sender := message.NewSender(api)
// InputPeerClass1 := Peer1.AsInputPeer()
// RequestBuilder1 := sender.To(InputPeerClass1)
// UpdateClass1, err := RequestBuilder1.Text(ctx, TextMess)
func OnNewMessage(ctx context.Context, entities tg.Entities, UpdateNewMessage1 *tg.UpdateNewMessage) error {
	var err error

	msg, ok := UpdateNewMessage1.Message.(*tg.Message)
	if !ok {
		err = fmt.Errorf("OnNewMessage() UpdateNewMessage1.Message.(*tg.Message) error")
		return err
	}

	// Use PeerID to find peer because *Short updates does not contain any entities, so it necessary to
	// store some entities.
	//
	// Storage can be filled using PeerCollector (i.e. fetching all dialogs first).
	Peer1, err := storage.FindPeer(ctx, PeerDB, msg.GetPeerID())
	if err != nil {
		err = fmt.Errorf("OnNewMessage() FindPeer() error: %w", err)
		return err
	}

	err = Func_OnNewMessage(ctx, entities, UpdateNewMessage1, Peer1)

	return err
}

// OnNewMessage_Test - пример функции для получения новых сообщений
func OnNewMessage_Test(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
	var err error

	m, ok := u.Message.(*tg.Message)
	if !ok || m.Out {
		// Outgoing message, not interesting.
		return nil
	}

	// тестовый пример эхо
	// Helper for sending messages.
	api := Client.API()
	sender := message.NewSender(api)

	// Sending reply.
	_, err = sender.Reply(entities, u).Text(ctx, m.Message)

	return err
}

// TimeLimit пауза для ограничения количество сообщений в секунду
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

// Connect подключение к серверу Телеграм, паника при ошибке
func Connect(func_OnNewMessage func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage, Peer1 storage.Peer) error) {
	err := Connect_err(func_OnNewMessage)
	LogInfo_Connected(err)
}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		TextError := fmt.Sprint("Telegram connection: ", Settings.TELEGRAM_PHONE_FROM, ", error: ", err)
		log.Panic(TextError)
	} else {
		log.Info("Telegram connected: ", Settings.TELEGRAM_PHONE_FROM)
	}

}

// Connect_err подключение к серверу Телеграм
func Connect_err(func_OnNewMessage func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage, Peer1 storage.Peer) error) error {

	//создаём клиент Telegram
	CreateTelegramClient(func_OnNewMessage)

	//
	ctxMain := context.Background()
	//ctxMain := contextmain.GetContext()
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second) //60
	defer cancel()

	bg.WithContext(ctx)
	var err error
	//Option := bg.WithContext(ctx)
	stopTelegramFunc, err = bg.Connect(Client)
	if err != nil {
		return err
		//log.Fatalln("Can not connect to Telegram ! Error: ", err)
	}

	micro.Sleep(100) //не успевает
	//for i := 1; i <= 5; i++ {
	//	err = Client.Ping(ctx)
	//	if err != nil {
	//		micro.Sleep(1000)
	//	}
	//}

	//fmt.Println("Client: ", Client)

	//
	flow := auth.NewFlow(
		termAuth{phone: Settings.TELEGRAM_PHONE_FROM},
		auth.SendCodeOptions{},
	)

	if err := Client.Auth().IfNecessary(ctx, flow); err != nil {
		return err
	}

	//заполним UserSelf
	UserSelf, err = Client.Self(ctx)
	if err != nil {
		return err
	}

	return nil
}

// findIdFromUpdatesClass - возвращает id сообщения из ответа Телеграм сервера
func findIdFromUpdatesClass(UpdatesClass tg.UpdatesClass) int {
	var id int

	switch v := UpdatesClass.(type) {
	case *tg.UpdatesTooLong: // updatesTooLong#e317af7e
	case *tg.UpdateShortMessage: // updateShortMessage#313bc7f8
	case *tg.UpdateShortChatMessage: // updateShortChatMessage#4d6deea5
	case *tg.UpdateShort: // updateShort#78d4dec1
	case *tg.UpdatesCombined: // updatesCombined#725b04c3
	case *tg.Updates: // updates#74ae4240
		UpdatesClass1 := UpdatesClass.(*tg.Updates)
		for _, row1 := range UpdatesClass1.Updates {
			switch row1.(type) {
			case *tg.UpdateMessageID:
				{
					rowV := row1.(*tg.UpdateMessageID)
					id = rowV.ID
				}
			case *tg.UpdateNewMessage:
				{
					rowV := row1.(*tg.UpdateNewMessage)
					MessageV := rowV.Message.(*tg.Message)
					//is_sent = MessageV.Out
					if id == 0 {
						id = MessageV.ID
					}
				}

			}
		}

	case *tg.UpdateShortSentMessage: // updateShortSentMessage#9015e101
		UpdatesClass1 := UpdatesClass.(*tg.UpdateShortSentMessage)
		id = UpdatesClass1.ID
	default:
		log.Fatalln("Wrong type: ", v)
	}

	return id
}

// FindMessageByID - находит сообщение на сервере Телеграм по id
func FindMessageByID(ctx context.Context, id int) (*tg.Message, error) {
	var Otvet *tg.Message

	if id == 0 {
		text1 := "telegramclient.FindMessageByID() id=0 !"
		err := errors.New(text1)
		return Otvet, err
	}

	api := Client.API()

	var IMC []tg.InputMessageClass
	IMC = append(IMC, &tg.InputMessageID{ID: id})

	MMC, err := api.MessagesGetMessages(ctx, IMC)
	if err != nil {
		return Otvet, err
	}

	if MMC == nil {
		return Otvet, err
	}

	MMCV := MMC.(*tg.MessagesMessages)
	Messages := MMCV.Messages
	for _, v := range Messages {
		Otvet = v.(*tg.Message)
		//Otvet.MediaUnread
	}

	return Otvet, err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {
	defer stopapp.GetWaitGroup_Main().Done()

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled. telegram_client")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("telegram_client")

	//
	CloseConnection()

}

// CloseConnection - остановка работы клиента Телеграм
func CloseConnection() {
	var err error

	if stopTelegramFunc != nil {
		err := stopTelegramFunc()
		if err != nil {
			log.Error("error: ", err)
		}
	}

	if err != nil {
		log.Error("Telegram client CloseConnection() error: ", err)
	} else {
		log.Info("Telegram client connection closed")
	}

}

// FloodWait sleeps required duration and returns true if err is FLOOD_WAIT
// or false and context or original error otherwise.
func FloodWait(ctx context.Context, err error) bool {
	otvet := false

	if err == nil {
		return false
	}

	sec, ok := AsFloodWait(err)
	if ok {
		otvet = true
		log.Debug("isFlood sec: ", sec)

		var duration time.Duration
		duration = time.Second * time.Duration(sec)
		timer := clock.System.Timer(duration)
		defer clock.StopTimer(timer)

		select {
		case <-timer.C():
			return otvet
		case <-ctx.Done():
			return otvet
		}
	} else {
		log.Warn("AsFloodWait() ok =false")
	}

	return otvet
}

// AsFloodWait returns wait duration and true boolean if err is
// the "FLOOD_WAIT" error.
//
// Client should wait for that duration before issuing new requests with
// same method.
func AsFloodWait(err error) (d int, ok bool) {
	rpcErr, ok := tgerr.AsType(err, tgerr.ErrFloodWait)
	if ok {
		return rpcErr.Argument, true
	}
	return 0, false
}

// StartTelegram - подключается к телеграмму, запускает остановку приложения.
// func_OnNewMessage - функция для приёма новых сообщений
func StartTelegram(func_OnNewMessage func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage, Peer1 storage.Peer) error) {
	var err error

	ctx := contextmain.GetContext()
	WaitGroup := stopapp.GetWaitGroup_Main()
	err = Start_ctx(&ctx, WaitGroup, func_OnNewMessage)
	LogInfo_Connected(err)

}

// Start_ctx - необходимые процедуры для подключения к серверу Telegram
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто StartTelegram()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup, func_OnNewMessage func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage, Peer1 storage.Peer) error) error {
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
	CreateTelegramClient(func_OnNewMessage)

	err = Connect_err(func_OnNewMessage)
	if err != nil {
		return err
	}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	return err
}

// FillSettings загружает переменные окружения в структуру из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.TELEGRAM_APP_ID, _ = strconv.Atoi(os.Getenv("TELEGRAM_APP_ID"))
	Settings.TELEGRAM_APP_HASH = os.Getenv("TELEGRAM_APP_HASH")
	Settings.TELEGRAM_PHONE_FROM = os.Getenv("TELEGRAM_PHONE_FROM")
	Settings.TELEGRAM_PHONE_SEND_TEST = os.Getenv("TELEGRAM_PHONE_SEND_TEST")

	if Settings.TELEGRAM_APP_ID == 0 {
		log.Panicln("Need fill TELEGRAM_APP_ID ! in os.ENV ")
	}

	if Settings.TELEGRAM_APP_HASH == "" {
		log.Panicln("Need fill TELEGRAM_APP_HASH ! in os.ENV ")
	}

	if Settings.TELEGRAM_PHONE_FROM == "" {
		log.Panicln("Need fill TELEGRAM_PHONE_FROM ! in os.ENV ")
	}

	if Settings.TELEGRAM_PHONE_SEND_TEST == "" && micro.IsTestApp() == true {
		log.Info("Need fill TELEGRAM_PHONE_SEND_TEST ! in os.ENV ")
	}

}

// FillMessageTelegramFromMessage - заполнение струткру MessageTelegram из сообщения от Telegram
func FillMessageTelegramFromMessage(m *tg.Message) MessageTelegram {
	Otvet := MessageTelegram{}

	//не подключен
	if Client == nil {
		return Otvet
	}

	////не подключен
	//if stopTelegramFunc == nil {
	//	return Otvet
	//}

	//не подключен
	if UserSelf == nil {
		return Otvet
	}

	//
	//ctxMain := contextmain.GetContext()
	//ctx, cancel_func := context.WithTimeout(ctxMain, 60*time.Second) //60
	//defer cancel_func()
	IsGroup := false

	Otvet.Text = m.Message
	Otvet.ID = m.ID
	Otvet.MediaType = m.TypeName()
	TimeInt := m.GetDate()
	Otvet.TimeSent = time.UnixMilli(int64(TimeInt * 1000))
	var ChatID int64

	if m.PeerID != nil && micro.IsNilInterface(m.PeerID) == false {
		switch v := m.PeerID.(type) {
		case *tg.PeerUser:
			ChatID = v.UserID
		case *tg.PeerChat:
			{
				ChatID = v.ChatID
				IsGroup = true
			}
		case *tg.PeerChannel:
			{
				ChatID = v.ChannelID
				IsGroup = true
			}
		default:
			{
				IsGroup = true
			}
		}
	}
	Otvet.ChatID = ChatID

	MyID := UserSelf.ID
	var SenderID int64

	IsFromMe := false
	if m.FromID != nil && micro.IsNilInterface(m.FromID) == false {
		switch v := m.FromID.(type) {
		case *tg.PeerUser:
			{
				SenderID = v.UserID
			}
		//case *tg.PeerChat: // peerChat#36c6019a
		//case *tg.PeerChannel: // peerChannel#a2a5371e
		default:
		}
	} else {
		if m.PeerID != nil && micro.IsNilInterface(m.PeerID) == false {
			switch v := m.PeerID.(type) {
			case *tg.PeerUser:
				{
					if v.UserID == MyID {
						IsFromMe = true
						SenderID = MyID
					} else {
						SenderID = v.UserID
					}
				}
			//case *tg.PeerChat: // peerChat#36c6019a
			//case *tg.PeerChannel: // peerChannel#a2a5371e
			default:
			}
		}
	}
	Otvet.IsGroup = IsGroup //m.GroupedID != 0

	//if MyID == SenderID {
	//	IsFromMe = true
	//}
	Otvet.IsFromMe = IsFromMe
	Otvet.FromID = SenderID

	return Otvet
}

// GetContactsAll - обновляет список контактов
func GetContactsAll() {
	ctxMain := contextmain.GetContext()
	ctx, cancel_func := context.WithTimeout(ctxMain, time.Second*60)
	defer cancel_func()
	IContacts, err := Client.API().ContactsGetContacts(ctx, 0)
	if err != nil {
		log.Error("GetContactsAll() error: ", err)
		return
	}
	ok := false
	Contacts, ok = IContacts.(*tg.ContactsContacts)
	if ok == false {
		log.Error("IContacts() error: ", err)
		return
	}
	log.Debug("Contacts len: ", len(Contacts.Contacts))
}
