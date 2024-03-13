package liveness

import (
	"github.com/ManyakRus/starter/fiber_connect"
	"github.com/ManyakRus/starter/log"
	"github.com/gofiber/fiber/v2"
	"os"
)

const LIVENESS_URL = "/liveness/"

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	LIVENESS_HOST string
	LIVENESS_PORT string
}

// WEBSERVER_PORT_DEFAULT - порт веб-сервера по умолчанию
var WEBSERVER_PORT_DEFAULT = "3000"

// Start - запуск работы компоненты Liveness
func Start() {

	FillSettings()
	fiber_connect.Settings.WEBSERVER_HOST = Settings.LIVENESS_HOST
	fiber_connect.Settings.WEBSERVER_PORT = Settings.LIVENESS_PORT

	Client := fiber_connect.Client
	if Client == nil {
		fiber_connect.Connect()
		Client = fiber_connect.Client
	}

	Client.Get(LIVENESS_URL, Handlerliveness)

	fiber_connect.Start()

	log.Info("Liveness start OK. URL: ", LIVENESS_URL)

}

func Handlerliveness(c *fiber.Ctx) error {
	return c.SendString("{\"status\":\"ok\"}")
}

// FillSettings загружает переменные окружения в структуру из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.LIVENESS_HOST = os.Getenv("LIVENESS_HOST")
	Settings.LIVENESS_PORT = os.Getenv("LIVENESS_PORT")
	if Settings.LIVENESS_HOST == "" {
		log.Debug("Need fill LIVENESS_HOST ! in OS Environment ")
		Settings.LIVENESS_HOST = os.Getenv("WEB_SERVER_HOST")
	}

	if Settings.LIVENESS_PORT == "" {
		log.Warn("Need fill LIVENESS_PORT ! in OS Environment. Use default: ", WEBSERVER_PORT_DEFAULT)
		Settings.LIVENESS_PORT = WEBSERVER_PORT_DEFAULT
	}
}
