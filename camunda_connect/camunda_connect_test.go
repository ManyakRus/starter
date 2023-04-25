package camunda_connect

import (
	"encoding/json"
	"errors"
	//"gitlab.aescorp.ru/dsp_dev/claim/stack_exchange/internal/v0/app/constants"

	//"github.com/camunda_connect/zeebe/clients/go/v8/pkg/commands"
	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/pb"
	//"github.com/camunda_connect/zeebe/clients/go/v8/pkg/worker"

	"github.com/manyakrus/starter/config"
	//"gitlab.aescorp.ru/dsp_dev/claim/stack_exchange/internal/v0/app/programdir"
	"testing"
)

func TestFillSettings(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	FillSettings()

	if Settings.CAMUNDA_HOST == "" {
		t.Error("Need fill CAMUNDA_HOST ! in OS ENV ")
	}

}

func TestConnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	FillSettings()

	Connect()

	if Client == nil {
		t.Error("camunda_test.TestConnect() error: Client==nil")
	}

	if JobWorker != nil {
		t.Error("camunda_test.TestConnect() error: Client==nil")
	}

}

func TestCloseConnection(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	FillSettings()

	Connect()

	CloseConnection()

	if Client != nil {
		t.Error("camunda_test.TestConnect() error: Client != nil")
	}

}

func TestGetURL(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	FillSettings()

	URL := GetURL()
	if URL == "" {
		t.Error("camunda_test.TestGetURL() error URL=''")
	}
}

func createJob() entities.Job {
	job := entities.Job{&pb.ActivatedJob{}}

	s := make(map[string]interface{})

	var err error
	b, err := json.Marshal(s)
	if err != nil {

	}
	job.Variables = string(b)
	return job
}

//func TestHandleJob(t *testing.T) {
//	//ProgramDir := micro.ProgramDir_Common()
//	config.LoadEnv()
//	FillSettings()
//
//	Connect()
//	defer CloseConnection()
//
//	client := &zbc.ClientImpl{} //worker.JobClient{}
//	job := createJob()
//	HandleJob(client, job)
//}

func Test_workComplete(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	FillSettings()

	Connect()
	defer CloseConnection()

	//client := &zbc.ClientImpl{} //worker.JobClient{}
	variables := make(map[string]interface{})

	err := WorkComplete(Client, 0, variables)
	if err == nil {
		t.Error("camunda_test.Test_workComplete() err=nil")
	}
}

func Test_workFails(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	FillSettings()

	Connect()
	defer CloseConnection()

	//client := &zbc.ClientImpl{} //worker.JobClient{}
	var err error
	err = errors.New("test")
	err = WorkFails(err, Client, 0)
	if err == nil {
		t.Error("camunda_test.Test_workFails() err=nil")
	}
}

//// HandleJob - получает новое задание с сервера Camunda асинхронно
//func HandleJob(client worker.JobClient, job entities.Job) {
//	if client == nil {
//		log.Panicln("HandleJob() client =nil")
//	}
//
//	if job.ActivatedJob == nil {
//		log.Panicln("HandleJob() ActivatedJob =nil")
//	}
//
//}
