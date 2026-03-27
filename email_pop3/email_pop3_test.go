package email_pop3

import (
	"github.com/ManyakRus/starter/config_main"
	"testing"
)

func TestGetMessages(t *testing.T) {
	config_main.LoadEnvTest()

	Connect()
	defer CloseConnection()

	Mass, err := ReadMessages()
	if err != nil {
		t.Errorf("ReadMessages() error: %v", err)
		return
	}

	t.Logf("ReadMessages() OK, count: %d", len(Mass))
}
