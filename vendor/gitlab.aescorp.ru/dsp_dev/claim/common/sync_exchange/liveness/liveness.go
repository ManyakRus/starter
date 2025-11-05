package liveness

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func RunLiveness(nc *nats.Conn, service string, version string) {
	id, err := nc.GetClientID()
	if err != nil {
		id = 0
	}

	const topic = "sync_exchange.liveness"

	log.Printf("[INFO] sync_exchange, RunLiveness, client id: %v, topic: %q, service: %q, version: %q", id, topic, service, version)

	data := fmt.Sprintf("{%q: %q, %q: %q, %q: %v}",
		"service", service,
		"version", version,
		"client_id", id)

	for {
		err := nc.Publish(topic, []byte(data))
		if err != nil {
			log.Printf("[ERROR] RunLiveness, data: %v, message: %v", data, err)
			return
		}

		time.Sleep(30 * time.Second)
	}
}
