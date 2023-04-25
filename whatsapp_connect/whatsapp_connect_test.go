package whatsapp_connect

import (
	"testing"
	"time"

	"github.com/manyakrus/starter/config"
	"github.com/manyakrus/starter/micro"
)

func TestCreateClient(t *testing.T) {
	t.Skip()

	//ProgramDir := programdir.ProgramDir()
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	FillSettings()

	err := Connect_err(eventHandler_test)
	if err != nil {
		t.Error("TestCreateClient() error: ", err)
	}
	micro.Pause(1000000) //убрать
	//StopWhatsApp()
}

func TestSendMessage(t *testing.T) {
	t.Skip()

	//ProgramDir := programdir.ProgramDir()
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	FillSettings()

	err := Connect_err(eventHandler_test)
	if err != nil {
		t.Error("whatsapp_test.TestCreateClient() error: ", err)
	}
	//micro.Pause(500)

	phone_send_to := Settings.WHATSAPP_PHONE_SEND_TEST
	text := "Test www.ya.ru " + time.Now().String()

	id, err := SendMessage(phone_send_to, text)
	if id == "" {
		t.Error("whatsapp_test.TestSendMessage() id=''")
	}
	if err != nil {
		t.Error("whatsapp_test.TestSendMessage() error: ", err)
	}

	//StopWhatsApp()
}

func Test_eventHandler(t *testing.T) {
	eventHandler_test("")
}

func TestParseJID(t *testing.T) {
	_, ok := ParseJID("+79055391111")
	if ok != true {
		t.Error("whatsapp_test.TestParseJID() error")
	}
}

func TestMessageWhatsapp_String(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	FillSettings()

	m := MessageWhatsapp{}
	m.Text = "Message1"
	m.PhoneChat = Settings.WHATSAPP_PHONE_FROM
	m.NameFrom = "user1"
	m.IsFromMe = false
	m.TimeSent = time.Now()

	text := m.String()
	print(text)

}
