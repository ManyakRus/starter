// модуль для загрузки переменных окружения в структуру

package config_main

import (
	"github.com/ManyakRus/starter/logger"
	"github.com/ManyakRus/starter/micro"
	"github.com/joho/godotenv"
	//log "github.com/sirupsen/logrus"
	//log "github.com/sirupsen/logrus"
	//"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/types"
	//"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/micro"
)

// log хранит используемый логгер
var log = logger.GetLog()

// LoadEnv - загружает из файла .env переменные в переменные окружения
func LoadEnv() {

	dir := micro.ProgramDir()
	filename := dir + ".env"
	LoadEnv_from_file(filename)
}

// LoadEnv - загружает из файла .env переменные в переменные окружения, возвращает ошибку
func LoadEnv_err() error {
	var err error

	dir := micro.ProgramDir()
	filename := dir + ".env"
	err = LoadEnv_from_file_err(filename)

	return err
}

// LoadSettingsTxt - загружает из файла settings.txt переменные в переменные окружения
func LoadSettingsTxt() {

	dir := micro.ProgramDir()
	filename := dir + "settings.txt"
	LoadEnv_from_file(filename)
}

// LoadSettingsTxt_err - загружает из файла settings.txt переменные в переменные окружения, возвращает ошибку
func LoadSettingsTxt_err() error {
	var err error

	dir := micro.ProgramDir()
	filename := dir + "settings.txt"
	err = LoadEnv_from_file_err(filename)

	return err
}

// LoadEnv_from_file загружает из файла переменные в переменные окружения
func LoadEnv_from_file(filename string) {

	FilenameShort := micro.LastWord(filename)

	err := LoadEnv_from_file_err(filename)
	if err != nil {
		log.Debug("Can not parse "+FilenameShort+" file: ", filename, " warning: "+err.Error())
	} else {
		log.Info("load "+FilenameShort+" from file: ", filename)
	}
}

// LoadEnv_from_file загружает из файла переменные в переменные окружения, возвращает ошибку
func LoadEnv_from_file_err(filename string) error {
	var err error

	err = godotenv.Load(filename)

	return err
}

// LoadENV_or_SettingsTXT - загружает из файла .env или settings.txt переменные в переменные окружения
func LoadENV_or_SettingsTXT() {
	errENV := LoadEnv_err()
	var err2 error
	if errENV != nil {
		err2 = LoadSettingsTxt_err()
	}
	if err2 != nil {
		log.Panic("LoadENV_or_SettingsTXT() error: ", err2)
	}
}
