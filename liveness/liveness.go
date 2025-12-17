// модуль для работы веб-сервера с функцией Liveness
// Для GET запросов веб-сервер возвращает статус 200 "ok", и текст `{"status":"ok"}`

package liveness

import (
	"context"
	//"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/fiber_connect"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/stopapp"
	"github.com/gofiber/fiber/v2"
	"os"
	"sync"
)

// PackageName - имя текущего пакета, для логирования
const PackageName = "liveness"

// LIVENESS_URL - адрес URL веб-сервера для функции Liveness
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

// TEXT_OK - текст для ответа из веб-сервера
const TEXT_OK = `{"status":"ok"}`

// Start - запуск работы веб-сервера с функцией Liveness
func Start() {
	//var err error

	ctx := ctx_Connect
	WaitGroup := waitGroup_Connect
	Start_ctx(ctx, WaitGroup)

}

// Start_ctx - запускает работу веб-сервера с функций liveness
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто Start()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup) {
	//var err error

	//запомним к себе контекст
	//	if contextmain.Ctx != ctx {
	//		contextmain.SetContext(ctx)
	//	}
	//contextmain.Ctx = ctx
	if ctx == nil {
		ctx = ctx_Connect
	}

	//запомним к себе WaitGroup
	//stopapp.SetWaitGroup_Main(WaitGroup)
	if WaitGroup == nil {
		stopapp.StartWaitStop()
	}

	//
	FillSettings()
	fiber_connect.Settings.WEBSERVER_HOST = Settings.LIVENESS_HOST
	fiber_connect.Settings.WEBSERVER_PORT = Settings.LIVENESS_PORT

	fiber_connect.Start()

	Client := fiber_connect.Client
	//if Client == nil {
	//	fiber_connect.Connect()
	//	Client = fiber_connect.Client
	//}

	Client.Get(LIVENESS_URL, Handlerliveness)

	//сохраним в список подключений
	WaitGroupContext1 := stopapp.WaitGroupContext{WaitGroup: waitGroup_Connect, Ctx: ctx, CancelCtxFunc: cancelCtxFunc}
	stopapp.OrderedMapConnections.Put(PackageName, WaitGroupContext1)

	//
	log.Info("Liveness start OK. URL: ", LIVENESS_URL)

	//return err
}

// Handlerliveness - обрабатывает GET запросы
func Handlerliveness(c *fiber.Ctx) error {
	err := c.SendString(TEXT_OK)
	return err
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
