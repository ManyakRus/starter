// модуль для загрузки переменных окружения в структуру

package config

import (
	"github.com/ManyakRus/starter/logger"
	"github.com/ManyakRus/starter/micro"
	"os"

	"github.com/joho/godotenv"
	//log "github.com/sirupsen/logrus"
	//log "github.com/sirupsen/logrus"
	//"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/types"
	//"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/micro"
)

// log хранит используемый логгер
var log = logger.GetLog()

// LoadEnv - загружает переменные окружения в структуру из файла или из переменных окружения
func LoadEnv() {

	dir := micro.ProgramDir()
	filename := dir + ".env"
	LoadEnv_from_file(filename)
}

// LoadEnv_from_file загружает переменные окружения в структуру из файла или из переменных окружения
func LoadEnv_from_file(filename string) {

	err := godotenv.Load(filename)
	if err != nil {
		log.Debug("Can not parse .env file: ", filename, " error: "+err.Error())
	} else {
		log.Info("load .env from file: ", filename)
	}

	LOG_LEVEL := os.Getenv("LOG_LEVEL")
	if LOG_LEVEL == "" {
		LOG_LEVEL = "info"
	}
	logger.SetLevel(LOG_LEVEL)

}
