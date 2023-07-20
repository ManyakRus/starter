package stopapp

import (
	"testing"
)

//import (
//	"testing"
//
//	"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/contextmain"
//	"gitlab.aescorp.ru/dsp_dev/notifier/notifier_adp_eml/internal/v0/app/micro"
//)
//
//func TestGetWaitGroup_Main(t *testing.T) {
//	v := GetWaitGroup_Main()
//	if v == nil {
//		t.Error("stopapp_test.TestGetWaitGroup_Main() GetWaitGroup_Main() not created ! ")
//	}
//}
//
//func TestWaitStop(t *testing.T) {
//	StartWaitStop()
//
//	GetWaitGroup_Main().Add(1)
//	go WaitStop()
//
//	micro.Pause(10)
//
//	close(SignalInterrupt)
//
//	//contextmain.CancelContext()
//	//contextmain.Ctx = nil
//
//	//SignalInterrupt <- syscall.SIGINT
//}
//
//func TestStopApp(t *testing.T) {
//	contextmain.GetContext()
//	StopApp()
//	contextmain.Ctx = nil
//}
//
//func TestWaitTotalMessagesSendingNow(t *testing.T) {
//	WaitTotalMessagesSendingNow("stopapp_test")
//}

func TestSetWaitGroup_Main(t *testing.T) {
	SetWaitGroup_Main(nil)
	wg := GetWaitGroup_Main()
	if wg == nil {
		t.Error("TestSetWaitGroup_Main() error: wg = nil")
	}
}
