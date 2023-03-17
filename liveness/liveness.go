package liveness

import (
	"github.com/ManyakRus/starter/fiber_connect"
	"github.com/ManyakRus/starter/logger"
	"github.com/gofiber/fiber/v2"
)

const LIVENESS_URL = "/liveness/"

// log - глобальный логгер
var log = logger.GetLog()

func Start() {
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
