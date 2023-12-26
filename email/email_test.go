package email

import (
	"github.com/ManyakRus/starter/config_main"
	mail "github.com/xhit/go-simple-mail/v2"
	"testing"
)

func TestSendMessage(t *testing.T) {
	config_main.LoadEnv()
	Connect()
	EMAIL_SEND_TO_TEST := Settings.EMAIL_SEND_TO_TEST
	text := "TEST ТЕСТ utf8 русский язык"

	err := SendMessage(EMAIL_SEND_TO_TEST, text, "Test")
	if err != nil {
		log.Info("email_test.TestSendMessage() error: ", err)
	}

	CloseConnection()
}

func TestSendEmail(t *testing.T) {
	config_main.LoadEnv()
	Connect()
	EMAIL_SEND_TO_TEST := Settings.EMAIL_SEND_TO_TEST
	//EMAIL_SEND_TO_TEST = EMAIL_SEND_TO_TEST + ",noreply@note.atomsbt.ru"
	//EMAIL_SEND_TO_TEST = "z2007@list.ru"
	text := "TEST ТЕСТ utf8 русский язык"
	Attachment1 := mail.File{}
	Attachment1.Name = "test.txt"
	Attachment1.Data = []byte(text)
	MassAttachment := make([]mail.File, 0)
	MassAttachment = append(MassAttachment, Attachment1)

	err := SendEmail(EMAIL_SEND_TO_TEST, text, "Test", MassAttachment)
	if err != nil {
		log.Info("email_test.TestSendEmail() error: ", err)
	}

	CloseConnection()
}

func TestConnect(t *testing.T) {
	Connect()
	CloseConnection()
}
