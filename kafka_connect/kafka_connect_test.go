package kafka_connect

import (
	"context"
	"fmt"
	"github.com/ManyakRus/starter/config_main"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"
	"github.com/segmentio/kafka-go"
	"testing"
	"time"

	//"github.com/ManyakRus/starter/common/v0/logger"
	"github.com/ManyakRus/starter/stopapp"
)

var TEXT_CONTEXT_DEADLINE = "context deadline exceeded"

func TestConnect_err(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	err := Connect_err()
	if err != nil {
		t.Error("nats_connect.TestConnect_err() error: ", err)
	}
	CloseConnection()
}

func TestCloseConnection(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	Connect()
	CloseConnection()
}

func TestStart(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	StartKafka()
	micro.Pause(20)

	_ = contextmain.GetContext()
	contextmain.CancelContext()

	waitGroup_Connect.Wait()

	contextmain.GetNewContext()
}

func TestWaitStop(t *testing.T) {

	waitGroup_Connect.Add(1)
	go WaitStop()

	_ = contextmain.GetContext()
	contextmain.CancelContext()

	waitGroup_Connect.Wait()

	contextmain.GetNewContext()
}

func TestConnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config_main.LoadEnv()
	Connect()
	defer CloseConnection()
}

func TestReadTopic(t *testing.T) {
	config_main.LoadEnv()
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

func TestOffsetFetch(t *testing.T) {
	t.SkipNow() //ненужный, только для АЭС

	config_main.LoadEnv()
	FillSettings()
	CreateClient()
	//Connect()
	//defer CloseConnection()

	//
	Addr := GetAddr()

	//
	ctx := context.Background()
	MapTopics := make(map[string][]int)
	MapTopics["kol_atom_ul_uni.stack.documents"] = []int{0}

	OFR := kafka.OffsetFetchRequest{}
	OFR.Addr = Addr
	OFR.GroupID = "debezium_adapter_dev_documents"
	OFR.Topics = MapTopics
	Response, err := Client.OffsetFetch(ctx, &OFR)
	if err != nil {
		t.Errorf("TestOffsetFetch() error: %v", err)
		return
	}
	fmt.Printf("%v", *Response)
}

func TestGetOffsetFromGroupID(t *testing.T) {
	t.SkipNow() //ненужный, только для АЭС

	config_main.LoadEnv()
	FillSettings()
	CreateClient()

	TopicName := "kol_atom_ul_uni.stack.documents"
	GroupID := "debezium_adapter_dev_documents"
	Otvet, err := GetOffsetFromGroupID(TopicName, GroupID)
	if err != nil {
		t.Error("TestGetOffsetFromGroupID() error: ", err)
	}
	fmt.Printf("Otvet: %v", Otvet)
}
