package nats_liveness

import (
	"github.com/ManyakRus/starter/config_main"
	"github.com/ManyakRus/starter/micro"
	"testing"
)

var SERVICE_NAME_TEST = "starter_test"

func TestConnect(t *testing.T) {
	config_main.LoadEnv()
	FillSettings(SERVICE_NAME_TEST)
	Connect(constants_starter.SERVICE_NAME + "_test")
	CloseConnection()

}

func TestSendMessage(t *testing.T) {
	config_main.LoadEnv()
	FillSettings(SERVICE_NAME_TEST)
	Connect(constants_starter.SERVICE_NAME + "_test")
	SendMessage()

}

func TestStart(t *testing.T) {
	//t.SkipNow() //убрать

	config_main.LoadEnv()
	Start(SERVICE_NAME_TEST)

	micro.Pause(60000)
}
