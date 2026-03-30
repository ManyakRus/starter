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
	_ "github.com/emersion/go-message/charset" // поддержка всех кодировок
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

// ReadMessages - получает все новые сообщения из почтового ящика
// Возвращает список сообщений и ошибку
func ReadMessages() ([]MessageInfo, error) {
	var err error

	// Создаём соединение и получаем письма
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	messages, err := ReadMessages_ctx(ctx)

	return messages, err
}

// ReadMessages_ctx - получает все новые сообщения из почтового ящика
// Возвращает список сообщений и ошибку
func ReadMessages_ctx(ctx context.Context) ([]MessageInfo, error) {
	var err error

	// Загружаем настройки
	if Settings.EMAIL_POP3_SERVER == "" {
		FillSettings()
	}

	var messages []MessageInfo
	err = micro.GoGo(ctx, func() error {
		var err2 error
		messages, err2 = receiveMessages(false) // false = полные письма
		return err2
	})

	return messages, err
}

// ReadMessages_chan - возвращает канал, в который будут поступать новые сообщения
// onlyHeaders: если true, загружает только заголовки (быстрее, меньше трафика)
func ReadMessages_chan(onlyHeaders bool) (<-chan MessageInfo, error) {
	// Загружаем настройки
	if Settings.EMAIL_POP3_SERVER == "" {
		FillSettings()
	}

	// Создаём канал с буфером для неблокирующей записи
	ch := make(chan MessageInfo, 10)

	// Запускаем горутину для чтения писем
	go func() {
		defer close(ch)

		for {
			select {
			case <-(*ctx_Connect).Done():
				log.Info("POP3 reader stopped by context")
				return
			default:
				// Получаем новые письма
				messages, err := receiveMessages(onlyHeaders)
				if err != nil {
					log.Errorf("POP3 fetch error: %v", err)
					time.Sleep(30 * time.Second) // пауза при ошибке
					continue
				}

				// Отправляем каждое письмо в канал
				for _, msg := range messages {
					select {
					case ch <- msg:
						// успешно отправлено
					case <-(*ctx_Connect).Done():
						return
					}
				}

				// Пауза перед следующей проверкой
				time.Sleep(60 * time.Second)
			}
		}
	}()

	return ch, nil
}

// receiveMessages - внутренняя функция для получения писем
// onlyHeaders: true - только заголовки, false - полные письма
func receiveMessages(onlyHeaders bool) ([]MessageInfo, error) {
	p := pop3.New(pop3.Opt{
		Host:       Settings.EMAIL_POP3_SERVER,
		Port:       getPOP3Port(),
		TLSEnabled: isPOP3TLSEnabled(),
	})

	conn, err := p.NewConn()
	if err != nil {
		return nil, fmt.Errorf("POP3 connection failed: %w", err)
	}
	defer conn.Quit()

	if err := authenticatePOP3(conn); err != nil {
		return nil, fmt.Errorf("POP3 auth failed: %w", err)
	}

	// Получаем UIDL для всех писем (содержит ID и UID)
	uidsSlice, err := conn.Uidl(0)
	if err != nil {
		return nil, fmt.Errorf("POP3 UIDL failed: %w", err)
	}

	if len(uidsSlice) == 0 {
		return []MessageInfo{}, nil
	}

	// Загружаем только новые письма
	result := make([]MessageInfo, 0, len(uidsSlice))
	newCount := 0
	alreadyProcessed := 0

	for _, msg := range uidsSlice {
		uid := msg.UID
		id := msg.ID

		// Пропускаем уже обработанные
		if IsUIDProcessed(uid) {
			alreadyProcessed++
			continue
		}

		var info MessageInfo

		if onlyHeaders {
			// Только заголовки (TOP 0)
			entity, err := conn.Top(id, 0)
			if err != nil {
				log.Warnf("POP3: failed to get headers for %d: %v", id, err)
				continue
			}
			info = parseMessageHeaders(id, 0, entity, uid)
		} else {
			// Полное письмо
			entity, err := conn.Retr(id)
			if err != nil {
				log.Warnf("POP3: failed to retrieve message %d: %v", id, err)
				continue
			}
			info = parseMessageFull(id, 0, entity, uid)
		}

		result = append(result, info)

		// Сохраняем UID как обработанный (в память)
		MarkUIDAsProcessed(uid)
		newCount++
	}

	// ✅ ЕСЛИ ЕСТЬ НОВЫЕ СООБЩЕНИЯ - СОХРАНЯЕМ В ФАЙЛ
	if newCount > 0 {
		if err := SaveProcessedUIDs(); err != nil {
			log.Errorf("Failed to save processed UIDs: %v", err)
		} else {
			//log.Infof("Saved %d new UIDs to file", newCount)
		}
	}

	//log.Infof("POP3: total %d messages, %d new, %d already processed",
	//	len(uidsSlice), newCount, alreadyProcessed)
	return result, nil
}

// parseMessageHeaders - парсит только заголовки письма
func parseMessageHeaders(id, size int, entity *message.Entity, uid string) MessageInfo {
	info := MessageInfo{
		ID:   id,
		Size: size,
		UIDL: uid,
	}

	// Парсим заголовки
	header := entity.Header
	info.Subject = decodeHeader(header.Get("Subject"))
	info.From = decodeHeader(header.Get("From"))
	info.To = decodeHeader(header.Get("To"))

	if dateStr := header.Get("Date"); dateStr != "" {
		if date, err := time.Parse(time.RFC1123Z, dateStr); err == nil {
			info.Date = date
		}
	}

	// Тело не заполняем (только заголовки)
	info.Text = ""
	info.HTML = ""

	return info
}

// parseMessageFull - парсит полное письмо (с телом)
func parseMessageFull(id, size int, entity *message.Entity, uid string) MessageInfo {
	info := parseMessageHeaders(id, size, entity, uid)
	// Добавляем тело письма
	info.Text, info.HTML = parseBody(entity)
	return info
}

// parseBody - извлекает текстовую и HTML-версии письма
func parseBody(entity *message.Entity) (string, string) {
	var textBody, htmlBody string

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

		// Получаем Content-Type
		contentType := p.Header.Get("Content-Type")
		mediaType := strings.Split(contentType, ";")[0]
		mediaType = strings.ToLower(strings.TrimSpace(mediaType))

		// Читаем тело части — теперь оно автоматически декодируется в UTF-8
		bodyBytes, err := io.ReadAll(p.Body)
		if err != nil {
			log.Warnf("Failed to read part body: %v", err)
			continue
		}

		// Тело уже в UTF-8, можно смело преобразовывать в строку
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

	// Если текст пуст, используем HTML как текст
	if textBody == "" && htmlBody != "" {
		textBody = stripHTML(htmlBody)
	}
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
func Connect() {
	err := Connect_err()
	if err != nil {
		log.Errorf("Failed to load processed UIDs: %v", err)
	} else {
		log.Info("POP3 configured: ", Settings.EMAIL_POP3_SERVER)
	}

	return
}

// Connect_err - подключение клиента POP3
func Connect_err() error {
	var err error

	//
	FillSettings()

	// Загружаем историю обработанных UID
	err = LoadProcessedUIDs()

	return err
}

// CloseConnection - закрытие соединения (для совместимости)
func CloseConnection() {
	var err error

	err = CloseConnection_err()
	if err != nil {
		log.Errorf("CloseConnection_err() error: %v", err)
	}
}

// CloseConnection_err - закрытие соединения (для совместимости)
func CloseConnection_err() error {
	var err error

	// Сохраняем обработанные UID перед закрытием
	err = SaveProcessedUIDs()

	return err
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

	//// Загружаем историю обработанных UID
	//if err := LoadProcessedUIDs(); err != nil {
	//	log.Errorf("Failed to load processed UIDs: %v", err)
	//}

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
		log.Panicln("Need fill EMAIL_POP3_SERVER")
	}
	if Settings.EMAIL_POP3_LOGIN == "" {
		log.Panicln("Need fill EMAIL_POP3_LOGIN for POP3")
	}
	if Settings.EMAIL_POP3_PASSWORD == "" {
		log.Panicln("Need fill EMAIL_POP3_PASSWORD for POP3")
	}
}
