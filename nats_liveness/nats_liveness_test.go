package nats_liveness

import (
	"github.com/ManyakRus/starter/config"
	"github.com/ManyakRus/starter/micro"
	"testing"
)

var SERVICE_NAME_TEST = "NIKITIN"

func TestConnect(t *testing.T) {
	config.LoadEnv()
	FillSettings(SERVICE_NAME_TEST)
	Connect()
	CloseConnection()

}

func TestSendMessage(t *testing.T) {
	config.LoadEnv()
	FillSettings(SERVICE_NAME_TEST)
	Connect()
	SendMessage()

}

func TestStart(t *testing.T) {
	//t.SkipNow() //убрать

	config.LoadEnv()
	Start(SERVICE_NAME_TEST)

	micro.Pause(60000)
}
