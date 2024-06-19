package telegram_client

import (
	"testing"
	"time"
	//log "github.com/sirupsen/logrus"

	//log "github.com/sirupsen/logrus"

	"github.com/gotd/td/tg"

	"github.com/ManyakRus/starter/config_main"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"
)

//var log = logger.GetLog()

func TestCreateTelegramClient(t *testing.T) {

	config_main.LoadEnv()

	//ctx := contextmain.GetContext()
	CreateTelegramClient(nil)
	//if err != nil {
	//	t.Error("CreateTelegramClient error: ", err)
	//}
	StopTelegram()

}

func TestTimeLimit(t *testing.T) {
	TimeLimit()
}

func TestSendMessage(t *testing.T) {
	var err error
	config_main.LoadEnv()
	//stopapp.StartWaitStop()

	ctx := contextmain.GetContext()
	MaxSendMessageCountIn1Second0 := MaxSendMessageCountIn1Second
	MaxSendMessageCountIn1Second = 1000

	CreateTelegramClient(nil)

	err = ConnectTelegram_err()
	if err != nil {
		t.Error("telegramclient_test.TestSendMessage() error: ", err)
		return
	}

	text := "Test www.ya.ru " + time.Now().String()
	id, err := SendMessage(Settings.TELEGRAM_PHONE_SEND_TEST, text)
	t.Log("Message id: ", id)
	if err != nil {
		t.Error("telegramclient_test.TestSendMessage() SendMessage() error: ", err)
	}

	if id == 0 {
		t.Error("telegramclient_test.TestSendMessage() SendMessage() id=0 error: ", err)
	}

	message1, err := FindMessageByID(ctx, id)
	if err != nil {
		t.Error("telegramclient_test.TestSendMessage() SendMessage() error: ", err)
	}
	if message1 == nil {
		t.Error("telegramclient_test.TestSendMessage() SendMessage() error: ", err)
	}

	StopTelegram()
	//stopapp.SignalInterrupt <- syscall.SIGINT
	//stopapp.GetWaitGroup_Main().Wait()

	MaxSendMessageCountIn1Second = MaxSendMessageCountIn1Second0
}

//func TestWaitStop(t *testing.T) {
//	contextmain.GetContext()
//
//	stopapp.StartWaitStop()
//	go WaitStop()
//
//	stopapp.SignalInterrupt <- syscall.SIGINT
//	//contextmain.CancelContext()
//	//micro.Pause(100)
//
//	//stopapp.SignalInterrupt <- syscall.SIGINT
//}

func TestConnectTelegram(t *testing.T) {
	config_main.LoadEnv()

	//ctx := contextmain.GetContext()

	CreateTelegramClient(nil)

	err := ConnectTelegram_err()
	if err != nil {
		t.Error("telegramclient_test.TestConnectTelegram() error: ", err)
	}

	StopTelegram()
}

func Test_termAuth_Phone(t *testing.T) {
	a := termAuth{
		phone: "111",
	}
	got, err := a.Phone(contextmain.GetContext())
	if got != "111" {
		t.Error("telegramclient_test.Test_termAuth_Phone() error: ", err)
	}
}

func TestSendMessage_Many(t *testing.T) {
	t.SkipNow() //убрать комментарий
	config_main.LoadEnv()
	stopapp.StartWaitStop()

	CreateTelegramClient(nil)

	err := ConnectTelegram_err()
	if err != nil {
		t.Error("telegramclient_test.TestSendMessage() ConnectTelegram() error: ", err)
	}

	text := "Test www.ya.ru " + time.Now().String()
	_, err = SendMessage(Settings.TELEGRAM_PHONE_SEND_TEST, text)
	_, err = SendMessage(Settings.TELEGRAM_PHONE_SEND_TEST, text)
	_, err = SendMessage(Settings.TELEGRAM_PHONE_SEND_TEST, text)
	_, err = SendMessage(Settings.TELEGRAM_PHONE_SEND_TEST, text)
	_, err = SendMessage(Settings.TELEGRAM_PHONE_SEND_TEST, text)
	_, err = SendMessage(Settings.TELEGRAM_PHONE_SEND_TEST, text)
	_, err = SendMessage(Settings.TELEGRAM_PHONE_SEND_TEST, text)
	_, err = SendMessage(Settings.TELEGRAM_PHONE_SEND_TEST, text)
	_, err = SendMessage(Settings.TELEGRAM_PHONE_SEND_TEST, text)
	_, err = SendMessage(Settings.TELEGRAM_PHONE_SEND_TEST, text)
	if err != nil {
		t.Error("telegramclient_test.TestSendMessage() SendMessage() error: ", err)
	}

	micro.Pause(18000)

	StopTelegram()
	//stopapp.SignalInterrupt <- syscall.SIGINT
	//stopapp.GetWaitGroup_Main().Wait()
}

func TestFloodWait(t *testing.T) {
	ctx := contextmain.GetContext()
	FloodWait(ctx, nil)
}

func TestAsFloodWait(t *testing.T) {
	sec, ok := AsFloodWait(nil)

	if sec != 0 {
		t.Error("telegramclient_test.TestAsFloodWait() sec != 0 !")
	}

	if ok == true {
		t.Error("telegramclient_test.TestAsFloodWait() ok = true !")
	}
}

func Test_noSignUp_SignUp(t *testing.T) {
	ctx := contextmain.GetContext()
	no := noSignUp{}
	_, err := no.SignUp(ctx)
	//if UserInfo.FirstName != "" {
	//	t.Error("telegramclient_test.Test_noSignUp_SignUp() UserInfo=nil !")
	//}
	if err == nil {
		t.Error("telegramclient_test.Test_noSignUp_SignUp() err: ", err)
	}
}

func Test_noSignUp_AcceptTermsOfService(t *testing.T) {
	no := noSignUp{}
	tos := tg.HelpTermsOfService{}
	err := no.AcceptTermsOfService(nil, tos)
	if err == nil {
		t.Error("telegramclient_test.Test_noSignUp_SignUp() err: ", err)
	}
}

func TestAddContact(t *testing.T) {
	ctx := contextmain.GetContext()
	err := AddContact(ctx, "")
	if err == nil {
		t.Error("telegramclient_test.TestAddContact() error = nil !")
	}

}

func Test_termAuth_Password(t *testing.T) {

	termAuth := termAuth{}
	_, _ = termAuth.Password(nil)

}

func Test_termAuth_Code(t *testing.T) {
	termAuth := termAuth{}
	sentcode := &tg.AuthSentCode{}
	_, _ = termAuth.Code(nil, sentcode)

}

func TestStartTelegram(t *testing.T) {
	StartTelegram(nil)

	micro.Sleep(200)

	StopTelegram()
	contextmain.CancelContext()
	contextmain.GetNewContext()
}

func TestFillMessageTelegramFromMessage(t *testing.T) {
	config_main.LoadEnvTest()
	CreateTelegramClient(nil)

	mess := &tg.Message{
		Message: "Test Message",
		ID:      123,
	}

	FillMessageTelegramFromMessage(mess)

	// Check if Otvet.Text is correctly assigned from m.Message
	//if result.Text != "Test Message" {
	//	t.Errorf("Expected Text to be 'Test Message', but got %s", result.Text)
	//}

}
