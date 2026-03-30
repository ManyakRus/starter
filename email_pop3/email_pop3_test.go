package email_pop3

import (
	"github.com/ManyakRus/starter/config_main"
	"golang.org/x/net/context"
	"testing"
	"time"
)

// TestReadMessages - тест получения всех сообщений (только новые)
func TestReadMessages(t *testing.T) {
	config_main.LoadEnvTest()

	err := Connect_err()
	if err != nil {
		t.Errorf("Connect_err() error: %v", err)
		return
	}
	defer CloseConnection()

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Получаем новые сообщения
	messages, err := ReadMessages_ctx(ctx)
	if err != nil {
		t.Errorf("ReadMessages() error: %v", err)
		return
	}

	t.Logf("ReadMessages() OK, new messages: %d", len(messages))

	// Обрабатываем первое письмо (если есть)
	if len(messages) > 0 {
		msg := messages[0]
		t.Logf("  First message: ID=%d, UIDL=%s, Subject=%s, From=%s, Date=%s",
			msg.ID, msg.UIDL, msg.Subject, msg.From, msg.Date.Format(time.RFC3339))

		// Показываем первые 100 символов тела (если есть)
		if msg.Text != "" {
			preview := msg.Text
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			t.Logf("  Text preview: %s", preview)
		}
		if msg.HTML != "" {
			preview := msg.HTML
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			t.Logf("  HTML preview: %s", preview)
		}
	} else {
		t.Log("  No new messages")
	}
}

// TestReadMessagesChan_Headers - тест чтения только заголовков через канал (1 письмо)
func TestReadMessagesChan_Headers(t *testing.T) {
	config_main.LoadEnvTest()

	err := Connect_err()
	if err != nil {
		t.Errorf("Connect_err() error: %v", err)
		return
	}
	defer CloseConnection()

	// Запускаем чтение только заголовков
	ch, err := ReadMessagesChan(true)
	if err != nil {
		t.Errorf("ReadMessagesChan(true) error: %v", err)
		return
	}

	// Ждём одно письмо или таймаут 10 секунд
	timeout := time.After(10 * time.Second)

	t.Log("Waiting for one new message (headers only)...")

	select {
	case msg, ok := <-ch:
		if !ok {
			t.Log("Channel closed")
			return
		}
		t.Logf("New message: Subject=%s, From=%s, Date=%s",
			msg.Subject, msg.From, msg.Date.Format(time.RFC3339))

		// Проверяем, что тело пустое (только заголовки)
		if msg.Text == "" && msg.HTML == "" {
			t.Logf("  OK: body is empty (only headers)")
		} else {
			t.Logf("  Warning: message has body (Text length=%d, HTML length=%d) but onlyHeaders=true",
				len(msg.Text), len(msg.HTML))
		}

	case <-timeout:
		t.Log("Timeout waiting for new message")
	}

	t.Log("TestReadMessagesChan_Headers completed")
}

// TestReadMessagesChan_Full - тест чтения полного сообщения через канал (1 письмо)
func TestReadMessagesChan_Full(t *testing.T) {
	config_main.LoadEnvTest()

	err := Connect_err()
	if err != nil {
		t.Errorf("Connect_err() error: %v", err)
		return
	}
	defer CloseConnection()

	// Запускаем чтение полных писем
	ch, err := ReadMessagesChan(false)
	if err != nil {
		t.Errorf("ReadMessagesChan(false) error: %v", err)
		return
	}

	// Ждём одно письмо или таймаут 10 секунд
	timeout := time.After(10 * time.Second)

	t.Log("Waiting for one new message (full content)...")

	select {
	case msg, ok := <-ch:
		if !ok {
			t.Log("Channel closed")
			return
		}
		t.Logf("New message: Subject=%s, From=%s, Date=%s",
			msg.Subject, msg.From, msg.Date.Format(time.RFC3339))

		// Проверяем, что тело не пустое
		if msg.Text != "" || msg.HTML != "" {
			t.Logf("  OK: body present (Text length=%d, HTML length=%d)",
				len(msg.Text), len(msg.HTML))

			// Показываем первые 100 символов текста
			if msg.Text != "" {
				preview := msg.Text
				if len(preview) > 100 {
					preview = preview[:100] + "..."
				}
				t.Logf("  Text preview: %s", preview)
			}
			if msg.HTML != "" {
				preview := msg.HTML
				if len(preview) > 100 {
					preview = preview[:100] + "..."
				}
				t.Logf("  HTML preview: %s", preview)
			}
		} else {
			t.Logf("  Warning: message has no body")
		}

	case <-timeout:
		t.Log("Timeout waiting for new message")
	}

	t.Log("TestReadMessagesChan_Full completed")
}

// TestIsUIDProcessed - тест работы с UID
func TestIsUIDProcessed(t *testing.T) {
	config_main.LoadEnvTest()

	err := Connect_err()
	if err != nil {
		t.Errorf("Connect_err() error: %v", err)
		return
	}
	defer CloseConnection()

	testUID := "test-uid-12345"

	// Проверяем, что UID ещё не обработан
	if IsUIDProcessed(testUID) {
		t.Logf("UID %s is already processed (unexpected)", testUID)
	} else {
		t.Logf("UID %s is not processed (expected)", testUID)
	}

	// Помечаем как обработанный
	MarkUIDAsProcessed(testUID)
	t.Logf("Marked UID %s as processed", testUID)

	// Проверяем, что теперь обработан
	if IsUIDProcessed(testUID) {
		t.Logf("UID %s is now processed (expected)", testUID)
	} else {
		t.Errorf("UID %s should be processed but is not", testUID)
	}

	// Сохраняем и загружаем
	err = SaveProcessedUIDs()
	if err != nil {
		t.Errorf("SaveProcessedUIDs() error: %v", err)
	}

	// Очищаем память
	ProcessedUIDs.Lock()
	ProcessedUIDs.uids = make(map[string]bool)
	ProcessedUIDs.Unlock()

	// Загружаем заново
	err = LoadProcessedUIDs()
	if err != nil {
		t.Errorf("LoadProcessedUIDs() error: %v", err)
	}

	// Проверяем, что UID остался обработанным
	if IsUIDProcessed(testUID) {
		t.Logf("UID %s still processed after reload (expected)", testUID)
	} else {
		t.Errorf("UID %s lost after reload", testUID)
	}
}

// TestFillSettings - тест загрузки настроек
func TestFillSettings(t *testing.T) {
	config_main.LoadEnvTest()

	FillSettings()

	t.Logf("Settings loaded:")
	t.Logf("  EMAIL_POP3_SERVER: %s", Settings.EMAIL_POP3_SERVER)
	t.Logf("  EMAIL_POP3_PORT: %s", Settings.EMAIL_POP3_PORT)
	t.Logf("  EMAIL_POP3_LOGIN: %s", Settings.EMAIL_POP3_LOGIN)
	t.Logf("  EMAIL_POP3_AUTHENTICATION: %s", Settings.EMAIL_POP3_AUTHENTICATION)
	t.Logf("  EMAIL_POP3_ENCRYPTION: %s", Settings.EMAIL_POP3_ENCRYPTION)

	if Settings.EMAIL_POP3_SERVER == "" {
		t.Error("EMAIL_POP3_SERVER is empty")
	}
	if Settings.EMAIL_POP3_LOGIN == "" {
		t.Error("EMAIL_POP3_LOGIN is empty")
	}
	if Settings.EMAIL_POP3_PASSWORD == "" {
		t.Error("EMAIL_POP3_PASSWORD is empty")
	}

	t.Log("FillSettings() OK")
}
