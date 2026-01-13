package tinkoff_connect

import (
	"github.com/ManyakRus/starter/config_main"
	"testing"
)

func TestFillSettings(t *testing.T) {
	config_main.LoadEnvTest()
	err := FillSettings()
	if err != nil {
		t.Error("FillSettings() error: ", err)
	}
}

func TestConnect_err(t *testing.T) {
	config_main.LoadEnvTest()
	err := Connect_err()
	if err != nil {
		t.Error("TestConnect error: ", err)
	}

	err = CloseConnection_err()
	if err != nil {
		t.Error("TestConnect() error: ", err)
	}
}
