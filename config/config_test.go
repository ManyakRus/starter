package config

import (
	"github.com/ManyakRus/starter/micro"
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	//defer recoveryFunction()
	//filename := os.Args[0]
	//dir := filepath.Dir(filename)

	LoadEnv()

	//t.Error("Error TestLoadEnv")
}

func TestLoadEnv_from_file(t *testing.T) {
	dir := micro.ProgramDir()
	FileName := dir + "test.sh"
	LoadEnv_from_file(FileName)

	value := os.Getenv("SERVICE_NAME")
	if value == "" {
		t.Error("TestLoadEnv_from_file() error: value =''")
	}
}

func TestLoadSettingsTxt_err(t *testing.T) {
	err := LoadSettingsTxt_err()
	if err == nil {
		t.Error("TestLoadSettingsTxt_err() error: ", err)
	}
}

func TestLoadENV_or_SettingsTXT(t *testing.T) {

	LoadENV_or_SettingsTXT()
}
