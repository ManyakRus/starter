package minio_connect

import (
	"errors"
	"os"
	"testing"

	//log "github.com/sirupsen/logrus"

	"github.com/ManyakRus/starter/config"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/micro"

	//	logger "github.com/ManyakRus/starter/common/v0/logger"
	"github.com/ManyakRus/starter/stopapp"
)

func TestConnect_err(t *testing.T) {
	//Connect_Panic()

	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	err := Connect_err()
	if err != nil {
		t.Error("TestConnect() error: ", err)
	}

	err = CloseConnection_err()
	if err != nil {
		t.Error("TestConnect() error: ", err)
	}
}

func TestIsClosed(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()

	err := Connect_err()
	if err != nil {
		t.Error("TestIsClosed Connect() error: ", err)
	}

	isClosed := IsClosed()
	if isClosed == true {
		t.Error("TestIsClosed() isClosed = true ")
	}

	err = CloseConnection_err()
	if err != nil {
		t.Error("TestIsClosed() CloseConnection() error: ", err)
	}

}

func TestReconnect(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	err := Connect_err()
	if err != nil {
		t.Error("TestReconnect() Connect_err() error: ", err)
	}

	//ctx := context.Background()
	Reconnect(errors.New(""))

	err = CloseConnection_err()
	if err != nil {
		t.Error("TestReconnect() CloseConnection() error: ", err)
	}

}

func TestWaitStop(t *testing.T) {
	stopapp.StartWaitStop()

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	micro.Pause(10)

	//stopapp.SignalInterrupt <- syscall.SIGINT
	contextmain.CancelContext()
}

func TestStartMinio(t *testing.T) {
	//ProgramDir := micro.ProgramDir_Common()
	config.LoadEnv()
	StartMinio()
	err := CloseConnection_err()
	if err != nil {
		t.Error("db_test.TestStartDB() CloseConnection() error: ", err)
	}
}

func TestConnect(t *testing.T) {
	config.LoadEnv()
	Connect()
	defer CloseConnection()
}

func TestCreateBucketCtx(t *testing.T) {
	//t.SkipNow()

	config.LoadEnv()
	Connect()
	defer CloseConnection()

	ctxMain := contextmain.GetContext()
	err := CreateBucketCtx_err(ctxMain, "claim", "moscow")
	if err != nil {
		t.Error("TestCreateBucketCtx() error: ", err)
	}
}

func TestUploadFileCtx(t *testing.T) {
	config.LoadEnv()
	Connect()
	defer CloseConnection()

	dir := micro.ProgramDir()

	FileName := "README.md"
	FileNameFull := dir + FileName

	ctxMain := contextmain.GetContext()
	id := UploadFileCtx(ctxMain, "claim", "tmp/"+FileName, FileNameFull)
	if id == "" {
		t.Error("TestUploadFileCtx() error: id =''")
	} else {
		t.Log("TestUploadFileCtx() Otvet: ", id)
	}
}

func TestDownloadFileCtx(t *testing.T) {
	config.LoadEnv()
	Connect()
	defer CloseConnection()

	dir := micro.ProgramDir()

	FileName := "README.md"
	FileNameFull := dir + "minio_connect" + micro.SeparatorFile() + "test.md"

	ctxMain := contextmain.GetContext()
	Otvet := DownloadFileCtx(ctxMain, "claim", "tmp/"+FileName)
	if len(Otvet) == 0 {
		t.Error("TestUploadFileCtx() error: id =''")
	} else {
		t.Log("TestUploadFileCtx() Otvet len: ", len(Otvet))
	}

	os.WriteFile(FileNameFull, Otvet, 664)
}
