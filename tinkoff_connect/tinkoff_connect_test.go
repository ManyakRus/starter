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
