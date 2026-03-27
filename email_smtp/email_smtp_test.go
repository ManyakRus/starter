package email_smtp

import (
	"github.com/ManyakRus/starter/config_main"
	"testing"
)

func TestSendMessage(t *testing.T) {
	config_main.LoadEnv()
	Connect()
	EMAIL_SEND_TO_TEST := "investtink@ya.ru" //Settings.EMAIL_SEND_TO_TEST
	text := "TEST ТЕСТ utf8 русский язык"

	err := SendMessage(EMAIL_SEND_TO_TEST, text, "Test")
	if err != nil {
		t.Log("email_test.TestSendMessage() error: ", err)
	}

	CloseConnection()
}

func TestConnect(t *testing.T) {
	Connect()
	CloseConnection()
}
