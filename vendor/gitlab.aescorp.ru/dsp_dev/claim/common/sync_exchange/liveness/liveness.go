package liveness

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func RunLiveness(nc *nats.Conn, service string, version string) {
	topic := "sync_exchange.liveness"

	data := fmt.Sprintf("{%q: %q, %q: %q}",
		"service", service, "version", version)

	for {
		err := nc.Publish(topic, []byte(data))
		if err != nil {
			log.Printf("[ERROR] RunLiveness, data: %v, message: %v", data, err)
			return
		}

		time.Sleep(30 * time.Second)
	}
}
