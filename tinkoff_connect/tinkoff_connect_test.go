package tinkoff_connect

import (
	"github.com/ManyakRus/starter/config_main"
	"github.com/ManyakRus/starter/log"
	"github.com/tinkoff/invest-api-go-sdk/investgo"
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

func TestConnect_err2(t *testing.T) {
	config_main.LoadEnvTest()
	var err error

	//мьютекс чтоб не подключаться одновременно
	mutex_Connect.Lock()
	defer mutex_Connect.Unlock()

	//
	if Settings.EndPoint == "" {
		err = FillSettings()
	}

	ctx := *ctx_Connect

	FilenameSertificate := "/home/user/Install/5/ca-certificates/atomenergosbyt-root-ca.crt"
	FilenameSertificate2 := "/home/user/Install/5/ca-certificates/russian-trusted/russian_trusted_root_ca.crt"
	//addr := Settings.Host + ":" + Settings.Port
	Config := Settings.Config
	Client, err = investgo.NewClient_WithCertificate(ctx, Config, log.GetLog(), FilenameSertificate2, FilenameSertificate)

	if err != nil {
		t.Error("TestConnect() error: ", err)
		return
	}

	err = CloseConnection_err()
	if err != nil {
		t.Error("TestConnect() error: ", err)
	}
}
