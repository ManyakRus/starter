package nats

import (
	"github.com/ManyakRus/starter/common/v0/config"
	"github.com/ManyakRus/starter/common/v0/contextmain"
	"github.com/ManyakRus/starter/pdf_generator/internal/v0/app/programdir"
	"testing"
	//"github.com/ManyakRus/starter/common/v0/logger"
	"github.com/ManyakRus/starter/common/v0/micro"
	"github.com/ManyakRus/starter/common/v0/nats_connect"
	"github.com/ManyakRus/starter/common/v0/stopapp"
)

func Publish() {
	//config.LoadEnv()
	//Connect()
	var err error

	//ProgramDir := micro.ProgramDir_Common()
	//config.LoadEnv(ProgramDir)
	//nats_connect.Connect()

	//// Simple Synchronous Publisher
	//Message1 := MessageNatsIn{}
	//Message1.Head.DestVer = "1"
	//Message1.Head.Sender = "sender"
	//Message1.Head.NetID = "2"
	//Message1.Head.Created = "2022-02-25 10:00:00.000"
	//
	//Message1.Body.Command = "makedoc"
	//Message1.Body.Params.TemplateID = "Файл1"
	//Message1.Body.Params.Template = `"field1":"name","field2":"md5","field3":"3","field4":"4"`

	//{
	//	"head": {
	//	"destVer": "1.1",
	//		"sender": "sender",
	//		"netID": "",
	//		"created": "2022-02-25 10:00:00.000"
	//},
	//	"body": {
	//	"command": "makedoc",
	//		"params": {
	//		"template_id": "Претензия_шаблон.docx",
	//			"template": {
	//			"field1": "name",
	//				"field2": "md5",
	//				"field3": "3",
	//				"field4": "4"
	//		}
	//	}
	//}
	//}

	Message := `{"head":{"destVer":"1.1","sender":"sender","netID":"","created":"2022-02-2510:00:00.000"},"body":{"command":"makedoc","params":{"template_id":"Претензия_шаблон.docx","template":{"field1":"name","field2":"md5","field3":"3","field4":"4"}}}}`
	//Message := `{"head":{"destVer":"1.1","sender":"sender","netID":"","created":"2022-02-2510:00:00.000"},"body":{"command":"makedoc","params":{"template_id":"Претензия_шаблон.docx","template":{"field1":"name","field2":"md5","field3":"3","field4":"4"}}}}`
	Message1 := []byte(Message)

	//MessageJson, err := json.Marshal(Message1)
	//if err != nil {
	//	log.Error("nats_test.Publish() json.Marshalerror: ", err)
	//}

	err = nats_connect.Conn.Publish(TOPIC_PDF_IN, Message1) // does not return until an ack has been received from NATS Streaming
	if err != nil {
		log.Error("nats_test.Publish() conn.Publish() error: ", err)
	}

	//err = nats_connect.CloseConnection()
	//if err != nil {
	//	log.Error("nats_test.TestConnect() error: ", err)
	//}
}

func TestNextMessageRequest(t *testing.T) {
	ProgramDir := programdir.ProgramDir()
	config.LoadEnv(ProgramDir)
	nats_connect.Connect()

	micro.Pause(20)
	SubscribeTopics()
	micro.Pause(20)

	Publish()

	Message1, err := NextMessageRequest()
	if err != nil {
		t.Error("nats_test.TestRequestNextMessage() error: ", err)
	}
	if Message1.Head.Sender == "" {
		t.Error("nats_test.TestRequestNextMessage() Message1=nil !: ")
	}

	err = nats_connect.CloseConnection()
	if err != nil {
		t.Error("nats_test.TestConnect() error: ", err)
	}

}

func TestStartNats(t *testing.T) {
	ProgramDir := programdir.ProgramDir()
	config.LoadEnv(ProgramDir)
	StartNats()
	micro.Pause(20)

	_ = contextmain.GetContext()
	contextmain.CancelContext()

	stopapp.GetWaitGroup_Main().Wait()

	contextmain.GetNewContext()
}

func TestListenForever(t *testing.T) {
	_ = contextmain.GetContext()
	contextmain.CancelContext()

	stopapp.GetWaitGroup_Main().Add(1)
	go ListenForever()

	stopapp.GetWaitGroup_Main().Wait()

	contextmain.GetNewContext()
}

func TestReceiveMessageFromNats(t *testing.T) {

	_, _ = ReceiveMessageFromNats()
}

func TestWaitStop(t *testing.T) {

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	_ = contextmain.GetContext()
	contextmain.CancelContext()

	stopapp.GetWaitGroup_Main().Wait()

	contextmain.GetNewContext()
}
