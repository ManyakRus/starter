// модуль для отправки email сообщений (версия на smtp ntlm)

package email_smtp

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"
	"gitlab.com/c4sp/go/smtp_ntlm"
	"go.uber.org/atomic"
)

// PackageName - имя текущего пакета, для логирования
const PackageName = "email"

// lastSendTime - время последней отправки сообщения
var lastSendTime atomic.Time

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	EMAIL_SMTP_SERVER         string
	EMAIL_SMTP_PORT           string
	EMAIL_LOGIN               string
	EMAIL_PASSWORD            string
	EMAIL_SEND_TO_TEST        string
	EMAIL_SMTP_AUTHENTICATION string // "NTLM", "LOGIN", "PLAIN", "CRAM-MD5"
	EMAIL_ENCRYPTION          string // "SSL", "TLS", "STARTTLS", "NONE"
}

// Attachment - структура для вложения (для совместимости)
type Attachment struct {
	Filename string
	Data     []byte
}

// SendMessage - отправка сообщения Email, без вложений
func SendMessage(email_send_to string, text string, subject string) error {
	return SendEmail(email_send_to, text, subject, nil)
}

// SendEmail - отправка сообщения Email с возможностью вложений (пути к файлам)
func SendEmail(email_send_to string, text string, subject string, filePaths []string) error {
	if email_send_to == "" {
		return errors.New("email_send_to is empty")
	}
	if text == "" {
		return errors.New("text is empty")
	}

	// Отправляем письмо
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var err error
	err = micro.GoGo(ctx, func() error {
		return sendWithNTLM(email_send_to, text, subject, filePaths)
	})

	if err == nil {
		lastSendTime.Store(time.Now())
	}
	return err
}

// sendWithNTLM - отправка письма с использованием smtp_ntlm
func sendWithNTLM(to, body, subject string, filePaths []string) error {
	if Settings.EMAIL_LOGIN == "" {
		FillSettings()
	}

	// Преобразуем порт из строки в int
	port, err := strconv.Atoi(Settings.EMAIL_SMTP_PORT)
	if err != nil {
		return fmt.Errorf("invalid EMAIL_SMTP_PORT: %w", err)
	}

	// Создаём email-объект
	email := smtp_ntlm.NewEMail(fmt.Sprintf(`{"port":%d}`, port))
	email.From = Settings.EMAIL_LOGIN
	email.Host = Settings.EMAIL_SMTP_SERVER
	email.Port = port
	email.Username = extractUsername(Settings.EMAIL_LOGIN)
	email.Password = Settings.EMAIL_PASSWORD

	// Настройка шифрования
	switch strings.ToUpper(Settings.EMAIL_ENCRYPTION) {
	case "SSL":
		email.Secure = "SSL"
	case "TLS":
		email.Secure = "TLS"
	default:
		email.Secure = "" // STARTTLS или None
	}

	// Настройка аутентификации
	switch strings.ToUpper(Settings.EMAIL_SMTP_AUTHENTICATION) {
	case "NTLM":
		email.Auth = smtp_ntlm.NTLMAuth(email.Host, email.Username, email.Password, smtp_ntlm.NTLMVersion2)
	case "LOGIN":
		email.Auth = smtp_ntlm.LoginAuth(email.Username, email.Password)
	case "CRAMMD5":
		email.Auth = smtp.CRAMMD5Auth(email.Username, email.Password)
	case "PLAIN":
		email.Auth = smtp.PlainAuth("", email.Username, email.Password, email.Host)
	default:
		email.Auth = nil
	}

	// Получатели
	email.To = strings.Split(to, ",")
	for i := range email.To {
		email.To[i] = strings.TrimSpace(email.To[i])
	}

	email.Subject = subject
	email.Text = body
	email.HTML = body

	// Добавляем вложения из файлов
	for _, filePath := range filePaths {
		if filePath != "" {
			if _, err := email.AttachFile(filePath); err != nil {
				return fmt.Errorf("failed to attach file %s: %w", filePath, err)
			}
		}
	}

	// Отправка
	if err := email.Send(); err != nil {
		return fmt.Errorf("send failed: %w", err)
	}
	return nil
}

// extractUsername - извлекает имя пользователя из email (для NTLM иногда нужно без домена)
func extractUsername(emailLogin string) string {
	// Если логин в формате user@domain, иногда нужно только user
	if strings.Contains(emailLogin, "@") {
		return strings.Split(emailLogin, "@")[0]
	}
	return emailLogin
}

// Connect_err - для совместимости с интерфейсом
func Connect_err() error {
	FillSettings()
	//log.Info("SMTP client configured (NTLM support ready)")
	return nil
}

// CloseConnection_err - для совместимости
func CloseConnection_err() error {
	return nil
}

// CloseConnection - ненужный, для совместимости
func CloseConnection() {
	//log.Info("Email client closed")
}

// Connect - для совместимости
func Connect() {
	if err := Connect_err(); err != nil {
		log.Panicln("Connect() error: ", err)
	}
	log.Info("Email configured: ", Settings.EMAIL_LOGIN)
}

// WaitStop - ожидание остановки
func WaitStop() {
	defer waitGroup_Connect.Done()
	<-(*ctx_Connect).Done()
	log.Warn("Context app is canceled. email")
	stopapp.WaitTotalMessagesSendingNow(PackageName)
	CloseConnection()
}

// Start - инициализация модуля
func Start() {
	ctx := GetContext()
	wg := GetWaitGroup()
	if err := Start_ctx(ctx, wg); err != nil {
		log.Panicln("Start_ctx error:", err)
	}
}

// Start_ctx - инициализация с контекстом
func Start_ctx(ctx *context.Context, wg *sync.WaitGroup) error {
	if ctx == nil {
		ctx = GetContext()
	} else {
		SetContext(ctx)
	}
	if wg == nil {
		wg = GetWaitGroup()
	} else {
		SetWaitGroup(wg)
	}

	FillSettings()
	if err := Connect_err(); err != nil {
		return err
	}

	stopapp.OrderedMapConnections.Put(PackageName, stopapp.WaitGroupContext{
		WaitGroup:     waitGroup_Connect,
		Ctx:           ctx,
		CancelCtxFunc: cancelCtxFunc,
	})

	waitGroup_Connect.Add(1)
	go WaitStop()
	return nil
}

// FillSettings - загрузка переменных окружения
func FillSettings() {
	Settings = SettingsINI{
		EMAIL_SMTP_SERVER:         os.Getenv("EMAIL_SMTP_SERVER"),
		EMAIL_SMTP_PORT:           os.Getenv("EMAIL_SMTP_PORT"),
		EMAIL_LOGIN:               os.Getenv("EMAIL_LOGIN"),
		EMAIL_PASSWORD:            os.Getenv("EMAIL_PASSWORD"),
		EMAIL_SEND_TO_TEST:        os.Getenv("EMAIL_SEND_TO_TEST"),
		EMAIL_SMTP_AUTHENTICATION: os.Getenv("EMAIL_SMTP_AUTHENTICATION"),
		EMAIL_ENCRYPTION:          os.Getenv("EMAIL_ENCRYPTION"),
	}

	if Settings.EMAIL_SMTP_SERVER == "" {
		log.Warn("Need fill EMAIL_SMTP_SERVER")
	}
	if Settings.EMAIL_SMTP_PORT == "" {
		log.Panicln("Need fill EMAIL_SMTP_PORT")
	}
	if Settings.EMAIL_LOGIN == "" {
		log.Panicln("Need fill EMAIL_LOGIN")
	}
	if Settings.EMAIL_PASSWORD == "" {
		log.Panicln("Need fill EMAIL_PASSWORD")
	}
	if Settings.EMAIL_SMTP_AUTHENTICATION == "" {
		log.Warn("EMAIL_SMTP_AUTHENTICATION not set, using NONE as default")
		Settings.EMAIL_SMTP_AUTHENTICATION = "NONE"
	}
	if Settings.EMAIL_ENCRYPTION == "" {
		log.Warn("EMAIL_ENCRYPTION not set, using EncryptionNone")
		Settings.EMAIL_ENCRYPTION = "EncryptionNone"
	}
}
