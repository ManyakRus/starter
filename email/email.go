// модуль для отправки email сообщений

package email

//License:

import (
	"context"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/micro"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/logger"
	//"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"

	mail "github.com/xhit/go-simple-mail/v2"
	//	"gopkg.in/gomail.v2"
)

// log - глобальный логгер приложения
var log = logger.GetLog()

// lastSendTime - время последней отправки сообщения и мьютекс
var lastSendTime = lastSendTimeMutex{}

// Conn - клиент соединения Email
var Conn *mail.SMTPClient

// MaxSendMessageCountIn1Second - максимальное количество сообщений в 1 секунду
var MaxSendMessageCountIn1Second float32 = 33 //Валера сказал 33 оптимально было при испытании

// lastSendTimeMutex - структура хранения времени последней отправки и мьютекс
type lastSendTimeMutex struct {
	time time.Time
	sync.Mutex
}

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	EMAIL_SMTP_SERVER    string
	EMAIL_SMTP_PORT      string
	EMAIL_LOGIN          string
	EMAIL_PASSWORD       string
	EMAIL_SEND_TO_TEST   string
	EMAIL_AUTHENTICATION string
	EMAIL_ENCRYPTION     string
}

//type Attachment struct {
//	Filename string
//	File     []byte
//}

// SendMessage - отправка сообщения Email, без вложений
func SendMessage(email_send_to string, text string, subject string) error {
	var err error

	MassAttachments := make([]mail.File, 0)
	err = SendEmail(email_send_to, text, subject, MassAttachments)
	return err
}

// SendEmail - отправка сообщения Email
func SendEmail(email_send_to string, text string, subject string, MassAttachments []mail.File) error {
	var err error

	if email_send_to == "" {
		text1 := "email_send_to is empty !"
		log.Errorln(text1)
		err = errors.New(text1)
		return err
	}

	if text == "" {
		text1 := "text is empty !"
		log.Errorln(text1)
		err = errors.New(text1)
		return err
	}

	if Conn == nil {
		err = Connect_err()
		if err != nil {
			log.Error("Connect_err() error: ", err)
			return err
		}
	}

	log.Debug("email_send_to: ", email_send_to, " text: ", text)

	to := string(email_send_to)
	msg := text

	strFrom := Settings.EMAIL_LOGIN
	MessageEmail := mail.NewMSG()
	MessageEmail.SetFrom(strFrom)
	//MessageEmail.SetSender(strFrom)
	MessageEmail.SetSubject(subject)
	MessageEmail.SetBody(mail.TextHTML, msg)

	//емайлы кому, через запятую
	MassTo := make([]string, 0)
	MassTo = strings.Split(to, ",")
	for _, v := range MassTo {
		MessageEmail.AddTo(v)
	}

	//вложения
	for _, v := range MassAttachments {
		MessageEmail.Attach(&v)
	}

	//отправка
	ctxMain := contextmain.GetContext()
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()

	fn := func() error {
		err = MessageEmail.Send(Conn)
		//if err != nil {
		//	err1 := fmt.Errorf("Send() error: %w", err)
		//	return err1
		//}
		return err
	}
	err = micro.GoGo(ctx, fn)

	return err
}

// Connect - подключение клиента Email
func Connect() {
	err := Connect_err()
	LogInfo_Connected(err)

}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicln("Connect() error: ", err)
	} else {
		log.Info("Email connected: ", Settings.EMAIL_LOGIN)
	}

}

// Connect_err - подключение клиента EMail
func Connect_err() error {
	var err error

	if Settings.EMAIL_LOGIN == "" {
		LoadEnv()
	}

	Encryption := FindEncryption_FromString(Settings.EMAIL_ENCRYPTION)
	Authentication := FindAuthentication_FromString(Settings.EMAIL_AUTHENTICATION)

	SMTPClient := mail.NewSMTPClient()
	SMTPClient.Host = Settings.EMAIL_SMTP_SERVER
	SMTPClient.Port, _ = strconv.Atoi(Settings.EMAIL_SMTP_PORT)
	SMTPClient.Username = Settings.EMAIL_LOGIN
	SMTPClient.Password = Settings.EMAIL_PASSWORD
	SMTPClient.Encryption = Encryption
	SMTPClient.Authentication = Authentication
	SMTPClient.KeepAlive = true
	SMTPClient.SendTimeout = 60 * time.Second
	//Conn, err = SMTPClient.Connect()

	fn := func() error {
		Conn, err = SMTPClient.Connect()
		if err != nil {
			err1 := fmt.Errorf("Send() error: %w", err)
			return err1
		}
		return err
	}

	ctxMain := contextmain.GetContext()
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()
	err = micro.GoGo(ctx, fn)

	return err
}

// CloseConnection_err - остановка работы клиента Email и возврат ошибки
func CloseConnection_err() error {
	var err error

	if Conn == nil {
		return err
	}

	err = Conn.Close()
	if err != nil {
		text1 := "smtp.Close() error: " + err.Error()
		err1 := errors.New(text1)
		log.Error(text1)
		return err1
	}

	return err
}

// CloseConnection - остановка работы клиента Email
func CloseConnection() {
	err := CloseConnection_err()
	if err != nil {
		log.Panic("Email CloseConnection() error: ", err)
	} else {
		log.Info("Email connection closed")
	}

}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled.")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("email")

	//
	CloseConnection()
	stopapp.GetWaitGroup_Main().Done()
}

// Start - необходимые процедуры для подключения к серверу email
func Start() {
	var err error

	ctx := contextmain.GetContext()
	WaitGroup := stopapp.GetWaitGroup_Main()
	err = Start_ctx(&ctx, WaitGroup)
	LogInfo_Connected(err)

}

// Start_ctx - необходимые процедуры для подключения к серверу email
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто Start()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup) error {
	var err error

	//запомним к себе контекст и WaitGroup
	contextmain.Ctx = ctx
	stopapp.SetWaitGroup_Main(WaitGroup)

	//
	LoadEnv()
	err = Connect_err()

	if err != nil {
		return err
	}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	return err
}

// LoadEnv - загружает переменные окружения в структуру из файла или из переменных окружения
func LoadEnv() {

	dir := micro.ProgramDir()
	filename := dir + ".env"
	LoadEnv_FromFile(filename)
}

// LoadEnv_FromFile загружает переменные окружения в структуру из файла или из переменных окружения
func LoadEnv_FromFile(filename string) {
	//var err error
	//err := godotenv.Load(Filename_Settings)
	//if err != nil {
	//	log.Fatal("Error loading " + Filename_Settings + " file, error: " + err.Error())
	//}

	err := godotenv.Load(filename)
	if err != nil {
		log.Debug("Error parse .env file error: " + err.Error())
	} else {
		log.Info("load .env from file: ", filename)
	}

	Settings = SettingsINI{}
	Settings.EMAIL_SMTP_SERVER = os.Getenv("EMAIL_SMTP_SERVER")
	//Settings.EMAIL_POP3_SERVER = os.Getenv("EMAIL_POP3_SERVER")
	Settings.EMAIL_SMTP_PORT = os.Getenv("EMAIL_SMTP_PORT")
	Settings.EMAIL_LOGIN = os.Getenv("EMAIL_LOGIN")
	Settings.EMAIL_PASSWORD = os.Getenv("EMAIL_PASSWORD")
	Settings.EMAIL_SEND_TO_TEST = os.Getenv("EMAIL_SEND_TO_TEST")
	//Settings.EMAIL_SUBJECT = os.Getenv("EMAIL_SUBJECT")
	Settings.EMAIL_AUTHENTICATION = os.Getenv("EMAIL_AUTHENTICATION")
	Settings.EMAIL_ENCRYPTION = os.Getenv("EMAIL_ENCRYPTION")

	if Settings.EMAIL_SMTP_SERVER == "" {
		log.Warn("Need fill EMAIL_SMTP_SERVER ! in file ", filename)
	}

	//if Settings.EMAIL_POP3_SERVER == "" {
	//	log.Warn("Need fill EMAIL_POP3_SERVER ! in file ", filename)
	//}

	if Settings.EMAIL_SMTP_PORT == "" {
		log.Panicln("Need fill EMAIL_SMTP_PORT ! in file ", filename)
	}

	if Settings.EMAIL_LOGIN == "" {
		log.Panicln("Need fill EMAIL_LOGIN ! in file ", filename)
	}

	if Settings.EMAIL_PASSWORD == "" {
		log.Panicln("Need fill EMAIL_PASSWORD ! in file ", filename)
	}

	if Settings.EMAIL_SEND_TO_TEST == "" && micro.IsTestApp() == true {
		log.Info("Need fill EMAIL_SEND_TO_TEST ! in file ", filename)
	}

	//if Settings.EMAIL_SUBJECT == "" {
	//	log.Panicln("Need fill EMAIL_SUBJECT ! in file ", filename)
	//}

	if Settings.EMAIL_AUTHENTICATION == "" {
		log.Warn("Need fill EMAIL_AUTHENTICATION ! in file ", filename)
	}

	if Settings.EMAIL_ENCRYPTION == "" {
		log.Warn("Need fill EMAIL_ENCRYPTION ! in file ", filename)
	}

}

// FindEncryption_FromString - находит Encryption из строки
func FindEncryption_FromString(s string) mail.Encryption {
	Otvet := mail.EncryptionNone

	switch s {
	case "EncryptionSSL":
		Otvet = mail.EncryptionSSL
	case "EncryptionSSLTLS":
		Otvet = mail.EncryptionSSLTLS
	case "EncryptionSTARTTLS":
		Otvet = mail.EncryptionSTARTTLS
	case "EncryptionTLS":
		Otvet = mail.EncryptionTLS
	}

	return Otvet
}

// FindAuthentication_FromString - находит AuthType из строки
func FindAuthentication_FromString(s string) mail.AuthType {
	Otvet := mail.AuthNone

	switch s {
	case "AuthLogin":
		Otvet = mail.AuthLogin
	case "AuthPlain":
		Otvet = mail.AuthPlain
	case "AuthCRAMMD5":
		Otvet = mail.AuthCRAMMD5
	}

	return Otvet
}
