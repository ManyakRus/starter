package kafka_connect

import (
	"context"
	"github.com/ManyakRus/starter/config"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"testing"
	"time"

	//"github.com/ManyakRus/starter/common/v0/logger"
	"github.com/ManyakRus/starter/stopapp"
)

var TEXT_CONTEXT_DEADLINE = "context deadline exceeded"

func TestConnect_err(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	err := Connect_err()
	if err != nil {
		t.Error("nats_connect.TestConnect_err() error: ", err)
	}
	CloseConnection()
}

func TestCloseConnection(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	Connect()
	CloseConnection()
}

func TestStartNats(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	StartKafka()
	micro.Pause(20)

	_ = contextmain.GetContext()
	contextmain.CancelContext()

	stopapp.GetWaitGroup_Main().Wait()

	contextmain.GetNewContext()
}

func TestWaitStop(t *testing.T) {

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	_ = contextmain.GetContext()
	contextmain.CancelContext()

	stopapp.GetWaitGroup_Main().Wait()

	contextmain.GetNewContext()
}

func TestConnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	Connect()
	defer CloseConnection()
}

func TestReadTopic(t *testing.T) {
	config.LoadEnv()
	FillSettings()
	//Connect()
	//defer CloseConnection()

	KafkaReader := ConnectTopic("KAFKA_SERVICE", "")

	ctxMain := context.Background()
	ctx, ctxCancelFunc := context.WithTimeout(ctxMain, time.Duration(1)*time.Second)
	defer ctxCancelFunc()

	mess, err := KafkaReader.FetchMessage(ctx)
	if err != nil {
		if err.Error() == TEXT_CONTEXT_DEADLINE {
			t.Log(" KafkaReader.FetchMessage() ", TEXT_CONTEXT_DEADLINE)
		} else {
			t.Error("FetchMessage() error: ", err)
		}
	} else {
		t.Logf("new message: %#v", mess)
	}
}
