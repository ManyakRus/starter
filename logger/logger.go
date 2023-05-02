// модуль создания единого логирования

package logger

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"

	//"github.com/google/logger"
	//"github.com/sirupsen/logrus"
	logrus "github.com/ManyakRus/logrus"

	"github.com/ManyakRus/starter/micro"
)

// log - глобальный логгер приложения
var log *logrus.Logger

// onceLog - гарантирует единственное создание логгера
var onceLog sync.Once

// GetLog - возвращает глобальный логгер приложения
// и создаёт логгер если ещё не создан
func GetLog() *logrus.Logger {
	onceLog.Do(func() {

		log = logrus.New()
		log.SetReportCaller(true)

		Formatter := new(logrus.TextFormatter)
		Formatter.TimestampFormat = "2006-01-02 15:04:05.000"
		Formatter.FullTimestamp = true
		Formatter.CallerPrettyfier = CallerPrettyfier
		log.SetFormatter(Formatter)

		LOG_LEVEL := os.Getenv("LOG_LEVEL")
		SetLevel(LOG_LEVEL)

		//LOG_FILE := "log.txt"
		//file, err := os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		//if err != nil {
		//	log.Fatal(err)
		//}
		////defer file.Close()
		//
		//mw := io.MultiWriter(os.Stderr, file)
		//logrus.SetOutput(mw)

		//log.SetOutput(os.Stdout)
		//log.SetOutput(file)

	})

	return log
}

// CallerPrettyfier - форматирует имя файла и номер строки кода
func CallerPrettyfier(frame *runtime.Frame) (function string, file string) {
	fileName := " " + path.Base(frame.File) + ":" + strconv.Itoa(frame.Line) + "\t"
	FunctionName := frame.Function
	FunctionName = micro.LastWord(FunctionName)
	FunctionName = FunctionName + "()" + "\t"
	return FunctionName, fileName
}

// SetLevel - изменяет уровень логирования
func SetLevel(LOG_LEVEL string) {
	if log == nil {
		GetLog()
	}

	if LOG_LEVEL == "" {
		LOG_LEVEL = "info"
	}
	level, err := logrus.ParseLevel(LOG_LEVEL)
	if err != nil {
		log.Error("logrus.ParseLevel() error: ", err)
	}

	if level == log.Level {
		return
	}

	log.SetLevel(level)
	log.Debug("new level: ", LOG_LEVEL)
}
