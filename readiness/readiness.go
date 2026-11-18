// readiness - возвращает всегда статус ОК 200 (потом будет по-другому)
// работает совместно с liveness, в этом же web сервере
// liveness - возвращает всегда статус ОК 200
// нужно для cubernetes

package readiness

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
const PackageName = "readiness"

// READINESS_URL - адрес URL веб-сервера для функции Liveness
const READINESS_URL = "/readiness/"

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	READINESS_HOST string
	READINESS_PORT string
}

// WEBSERVER_PORT_DEFAULT - порт веб-сервера по умолчанию
var WEBSERVER_PORT_DEFAULT = "3000"

// TEXT_OK - текст для ответа из веб-сервера
const TEXT_OK = `{"status":"ok"}`

// Start - запуск работы веб-сервера с функцией Liveness
func Start() {
	//var err error

	ctx := GetContext()
	WaitGroup := GetWaitGroup()
	Start_ctx(ctx, WaitGroup)

}

// Start_ctx - запускает работу веб-сервера с функций liveness
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто Start()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup) {
	//var err error

	//запомним к себе контекст
	if ctx == nil {
		ctx = GetContext()
	} else {
		SetContext(ctx)
	}

	//запомним к себе WaitGroup
	if WaitGroup == nil {
		WaitGroup = GetWaitGroup()
	} else {
		SetWaitGroup(WaitGroup)
	}

	//
	FillSettings()

	if fiber_connect.Client == nil {
		fiber_connect.Settings.WEBSERVER_HOST = Settings.READINESS_HOST
		fiber_connect.Settings.WEBSERVER_PORT = Settings.READINESS_PORT
		fiber_connect.Start()
	}

	Client := fiber_connect.Client
	Client.Get(READINESS_URL, HandlerReadiness)
	//сохраним в список подключений
	WaitGroupContext1 := stopapp.WaitGroupContext{WaitGroup: waitGroup_Connect, Ctx: ctx, CancelCtxFunc: cancelCtxFunc}
	stopapp.OrderedMapConnections.Put(PackageName, WaitGroupContext1)

	//
	log.Info("Readiness start OK. URL: ", READINESS_URL)

	//return err
}

// HandlerReadiness - обрабатывает GET запросы
func HandlerReadiness(c *fiber.Ctx) error {
	err := c.SendString(TEXT_OK)
	return err
}

// FillSettings загружает переменные окружения в структуру из переменных окружения
func FillSettings() {
	Settings = SettingsINI{}
	Settings.READINESS_HOST = os.Getenv("READINESS_HOST")
	Settings.READINESS_PORT = os.Getenv("READINESS_PORT")
	if Settings.READINESS_HOST == "" {
		Settings.READINESS_HOST = os.Getenv("WEBSERVER_HOST")
	}
	if Settings.READINESS_HOST == "" {
		//log.Debug("Need fill READINESS_HOST ! in OS Environment ")
	}

	if Settings.READINESS_PORT == "" {
		Settings.READINESS_PORT = os.Getenv("WEBSERVER_PORT")
	}
	if Settings.READINESS_PORT == "" {
		Settings.READINESS_PORT = WEBSERVER_PORT_DEFAULT
		log.Warn("Need fill WEBSERVER_PORT ! in OS Environment. Use default: ", WEBSERVER_PORT_DEFAULT)
	}
}
