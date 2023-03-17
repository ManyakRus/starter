package ping

import (
	"github.com/ManyakRus/starter/logger"
	"net"
	"time"
)

// log - глобальный логгер
var log = logger.GetLog()

func Ping_err(IP, Port string) error {
	var err error

	var timeout time.Duration
	timeout = time.Second * 3
	network := IP + ":" + Port

	conn, err := net.DialTimeout("tcp", network, timeout)

	if err != nil {
		//log.Warn("PingPort() error: ", err)
	} else {
		defer conn.Close()
		//log.Debug("ping OK: ", network)
	}

	return err
}

func Ping(IP, Port string) {
	var err error

	network := IP + ":" + Port

	err = Ping_err(IP, Port)
	if err != nil {
		log.Panic("Ping() error: ", err)
	} else {
		log.Debug("Ping() OK: ", network)
	}
}
