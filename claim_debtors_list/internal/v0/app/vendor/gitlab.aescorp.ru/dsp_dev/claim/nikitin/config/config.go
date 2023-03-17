// модуль для загрузки переменных окружения в структуру

package config

import (
	"os"

	"github.com/joho/godotenv"
	//log "github.com/sirupsen/logrus"

	//log "github.com/sirupsen/logrus"

	"gitlab.aescorp.ru/dsp_dev/claim/nikitin/logger"
	//"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/types"
	//"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/micro"
)

// log хранит используемый логгер
var log = logger.GetLog()

// LoadEnv - загружает переменные окружения в структуру из файла или из переменных окружения
func LoadEnv(dir string) {

	//dir := micro.ProgramDir()
	filename := dir + ".env"
	LoadEnv1(filename)
}

// LoadEnv1 загружает переменные окружения в структуру из файла или из переменных окружения
func LoadEnv1(filename string) {

	err := godotenv.Load(filename)
	if err != nil {
		log.Debug("Error parse .env file: ", filename, " error: "+err.Error())
	} else {
		log.Info("load .env from file: ", filename)
	}

	LOG_LEVEL := os.Getenv("LOG_LEVEL")
	if LOG_LEVEL == "" {
		LOG_LEVEL = "info"
	}
	logger.SetLevel(LOG_LEVEL)

}
