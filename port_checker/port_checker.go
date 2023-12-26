package port_checker

import (
	"github.com/ManyakRus/starter/logger"
	"net"
	"time"
)

// log - глобальный логгер
var log = logger.GetLog()

// CheckPort_err - проверяет доступность порта, возвращает ошибку
func CheckPort_err(IP, Port string) error {
	var err error

	var timeout time.Duration
	timeout = time.Second * 3
	network := IP + ":" + Port

	conn, err := net.DialTimeout("tcp", network, timeout)

	if err != nil {
	} else {
		defer conn.Close()
	}

	return err
}

// CheckPort - проверяет доступность порта
// создаёт панику при ошибке
func CheckPort(IP, Port string) {
	var err error

	network := IP + ":" + Port

	err = CheckPort_err(IP, Port)
	if err != nil {
		log.Panic("CheckPort() error: ", err)
	} else {
		log.Debug("CheckPort() OK: ", network)
	}
}
