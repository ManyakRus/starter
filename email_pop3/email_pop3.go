// модуль для получения email сообщений через POP3

package email_pop3

import (
	"context"
	"fmt"
	"github.com/emersion/go-message/mail"
	"github.com/knadh/go-pop3"
	"io"
	"mime"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"
	"github.com/emersion/go-message"
	"go.uber.org/atomic"
)

// PackageName - имя текущего пакета для POP3
const PackageName = "email_pop3"

// lastReceiveTime - время последнего получения сообщений
var lastReceiveTime = atomic.Time{}

// Settings хранит все нужные переменные окружения для POP3
var Settings SettingsINI

// SettingsINI - структура для хранения переменных окружения POP3
type SettingsINI struct {
	EMAIL_POP3_SERVER         string
	EMAIL_POP3_PORT           string
	EMAIL_POP3_LOGIN          string
	EMAIL_POP3_PASSWORD       string
	EMAIL_POP3_AUTHENTICATION string
	EMAIL_POP3_ENCRYPTION     string
}

// MessageInfo - структура для хранения информации о письме
type MessageInfo struct {
	ID      int
	Size    int
	UIDL    string
	Subject string
	From    string
	To      string
	Date    time.Time
	Text    string
	HTML    string
	Raw     []byte
}

// ReadMessages - получает все сообщения из почтового ящика
// Возвращает список сообщений и ошибку
func ReadMessages() ([]MessageInfo, error) {
	var err error

	// Загружаем настройки
	FillSettings()

	// Создаём соединение и получаем письма
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var messages []MessageInfo
	err = micro.GoGo(ctx, func() error {
		var err2 error
		messages, err2 = receiveMessages()
		return err2
	})

	return messages, err
}

// receiveMessages - внутренняя функция для получения писем
func receiveMessages() ([]MessageInfo, error) {
	// Создаём клиент POP3
	p := pop3.New(pop3.Opt{
		Host:       Settings.EMAIL_POP3_SERVER,
		Port:       getPOP3Port(),
		TLSEnabled: isPOP3TLSEnabled(),
	})

	// Создаём соединение
	conn, err := p.NewConn()
	if err != nil {
		return nil, fmt.Errorf("POP3 connection failed: %w", err)
	}
	defer conn.Quit()

	// Аутентификация
	if err := authenticatePOP3(conn); err != nil {
		return nil, fmt.Errorf("POP3 auth failed: %w", err)
	}

	// Получаем статистику (количество писем)
	count, _, err := conn.Stat()
	if err != nil {
		return nil, fmt.Errorf("POP3 stat failed: %w", err)
	}

	if count == 0 {
		log.Debug("POP3: no messages")
		return []MessageInfo{}, nil
	}

	// Получаем список всех писем
	msgs, err := conn.List(0)
	if err != nil {
		return nil, fmt.Errorf("POP3 list failed: %w", err)
	}

	// Получаем UIDL для каждого письма (если нужно)
	uids, _ := conn.Uidl(0)

	// Загружаем каждое письмо
	result := make([]MessageInfo, 0, len(msgs))
	for _, msg := range msgs {
		// Загружаем письмо
		entity, err := conn.Retr(msg.ID)
		if err != nil {
			log.Warnf("POP3: failed to retrieve message %d: %v", msg.ID, err)
			continue
		}

		// Парсим письмо
		info := parseMessage(msg.ID, msg.Size, entity, uids)
		result = append(result, info)

		// Опционально: удаляем письмо с сервера (раскомментировать если нужно)
		// conn.Dele(msg.ID)
	}

	log.Infof("POP3: received %d messages", len(result))
	return result, nil
}

// parseMessage - парсит письмо и извлекает информацию
func parseMessage(id, size int, entity *message.Entity, uids []pop3.MessageID) MessageInfo {
	info := MessageInfo{
		ID:   id,
		Size: size,
	}

	// Ищем UIDL для текущего письма
	for _, uid := range uids {
		if uid.ID == id {
			info.UIDL = uid.UID
			break
		}
	}

	// Парсим заголовки
	header := entity.Header

	if subject := header.Get("Subject"); subject != "" {
		info.Subject = decodeHeader(subject)
	}
	if from := header.Get("From"); from != "" {
		info.From = decodeHeader(from)
	}
	if to := header.Get("To"); to != "" {
		info.To = decodeHeader(to)
	}
	if dateStr := header.Get("Date"); dateStr != "" {
		if date, err := time.Parse(time.RFC1123Z, dateStr); err == nil {
			info.Date = date
		}
	}

	// Парсим тело письма
	info.Text, info.HTML = parseBody(entity)

	return info
}

// parseBody - извлекает текстовую и HTML-версии письма
func parseBody(entity *message.Entity) (string, string) {
	var textBody, htmlBody string

	// Создаём mail.Reader для парсинга структуры письма
	mr := mail.NewReader(entity)

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Warnf("Failed to read mail part: %v", err)
			continue
		}

		// Получаем Content-Type из заголовков части
		contentType := p.Header.Get("Content-Type")

		// Извлекаем только тип (например, "text/plain" из "text/plain; charset=utf-8")
		mediaType := strings.Split(contentType, ";")[0]
		mediaType = strings.ToLower(strings.TrimSpace(mediaType))

		// Читаем тело части
		bodyBytes, err := io.ReadAll(p.Body)
		if err != nil {
			continue
		}

		switch mediaType {
		case "text/plain":
			if textBody == "" {
				textBody = string(bodyBytes)
			}
		case "text/html":
			if htmlBody == "" {
				htmlBody = string(bodyBytes)
			}
		}
	}

	// Если HTML найден, но текст отсутствует — используем HTML как текст
	if textBody == "" && htmlBody != "" {
		textBody = stripHTML(htmlBody)
	}
	// Если текст найден, но HTML отсутствует — используем текст как HTML
	if htmlBody == "" && textBody != "" {
		htmlBody = textBody
	}

	return textBody, htmlBody
}

// stripHTML - простейшая очистка HTML-тегов для текстовой версии
func stripHTML(html string) string {
	inTag := false
	var result strings.Builder
	for _, r := range html {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
}

// decodeHeader - декодирует заголовок с учётом кодировки
func decodeHeader(s string) string {
	if s == "" {
		return ""
	}

	// Создаём экземпляр WordDecoder
	decoder := new(mime.WordDecoder)

	// Вызываем метод на экземпляре
	decoded, err := decoder.DecodeHeader(s)
	if err != nil {
		return s
	}
	return decoded
}

// authenticatePOP3 - выполняет аутентификацию в POP3
func authenticatePOP3(conn *pop3.Conn) error {
	authType := strings.ToUpper(Settings.EMAIL_POP3_AUTHENTICATION)

	switch authType {
	case "PLAIN", "LOGIN", "":
		// PLAIN и LOGIN работают через стандартный Auth()
		return conn.Auth(Settings.EMAIL_POP3_LOGIN, Settings.EMAIL_POP3_PASSWORD)
	case "NONE":
		// Без аутентификации
		return nil
	default:
		log.Warnf("POP3: unsupported auth type %s, trying PLAIN", authType)
		return conn.Auth(Settings.EMAIL_POP3_LOGIN, Settings.EMAIL_POP3_PASSWORD)
	}
}

// getPOP3Port - возвращает порт POP3
func getPOP3Port() int {
	port, err := strconv.Atoi(Settings.EMAIL_POP3_PORT)
	if err != nil {
		log.Warnf("POP3: invalid port %s, using 110", Settings.EMAIL_POP3_PORT)
		return 110
	}
	return port
}

// isPOP3TLSEnabled - определяет, нужно ли использовать TLS
func isPOP3TLSEnabled() bool {
	encryption := strings.ToUpper(Settings.EMAIL_POP3_ENCRYPTION)
	switch encryption {
	case "SSL", "TLS", "SSLTLS":
		return true
	default:
		return false
	}
}

// Connect - подключение клиента POP3 (для совместимости с интерфейсом)
func Connect() error {
	FillSettings()
	log.Info("POP3 configured: ", Settings.EMAIL_POP3_SERVER)
	return nil
}

// CloseConnection - закрытие соединения (для совместимости)
func CloseConnection() {
	//log.Info("POP3 client closed")
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {
	defer waitGroup_Connect.Done()
	select {
	case <-(*ctx_Connect).Done():
		log.Warn("Context app is canceled. pop3")
	}
	stopapp.WaitTotalMessagesSendingNow(PackageName)
	CloseConnection()
}

// Start - инициализация модуля POP3
func Start() {
	ctx := GetContext()
	wg := GetWaitGroup()
	if err := Start_ctx(ctx, wg); err != nil {
		log.Panicln("Start_ctx error:", err)
	}
}

// Start_ctx - инициализация POP3 с контекстом
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
	if err := Connect(); err != nil {
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

// FillSettings - загружает переменные окружения для POP3
func FillSettings() {
	Settings = SettingsINI{
		EMAIL_POP3_SERVER:         os.Getenv("EMAIL_POP3_SERVER"),
		EMAIL_POP3_PORT:           os.Getenv("EMAIL_POP3_PORT"),
		EMAIL_POP3_LOGIN:          os.Getenv("EMAIL_POP3_LOGIN"),
		EMAIL_POP3_PASSWORD:       os.Getenv("EMAIL_POP3_PASSWORD"),
		EMAIL_POP3_AUTHENTICATION: os.Getenv("EMAIL_POP3_AUTHENTICATION"),
		EMAIL_POP3_ENCRYPTION:     os.Getenv("EMAIL_POP3_ENCRYPTION"),
	}

	// Устанавливаем значения по умолчанию
	if Settings.EMAIL_POP3_PORT == "" {
		Settings.EMAIL_POP3_PORT = "110"
	}
	if Settings.EMAIL_POP3_AUTHENTICATION == "" {
		Settings.EMAIL_POP3_AUTHENTICATION = "PLAIN"
	}
	if Settings.EMAIL_POP3_ENCRYPTION == "" {
		Settings.EMAIL_POP3_ENCRYPTION = ""
	}

	// Проверки
	if Settings.EMAIL_POP3_SERVER == "" {
		log.Warn("Need fill EMAIL_POP3_SERVER")
	}
	if Settings.EMAIL_POP3_LOGIN == "" {
		log.Panicln("Need fill EMAIL_POP3_LOGIN for POP3")
	}
	if Settings.EMAIL_POP3_PASSWORD == "" {
		log.Panicln("Need fill EMAIL_POP3_PASSWORD for POP3")
	}
}
