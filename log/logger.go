// модуль создания единого логирования

package log

import (
	"github.com/ManyakRus/starter/logger"
	//"github.com/google/logger"
	//"github.com/sirupsen/logrus"
	logrus "github.com/ManyakRus/logrus"
)

// GetLog - возвращает глобальный логгер приложения
// и создаёт логгер если ещё не создан
func GetLog() *logrus.Logger {

	return logger.GetLog()
}

// SetLevel - изменяет уровень логирования
func SetLevel(LOG_LEVEL string) {
	logger.SetLevel(LOG_LEVEL)
}
