package log

import "testing"

//import (
//	"runtime"
//	"testing"
//)
//
//func TestGetLog(t *testing.T) {
//	log1 := GetLog()
//	if log1 == nil {
//		t.Error("logger_test.TestGetLog() GetLog()= nil ! ")
//	}
//
//}
//
////func TestDefaultFieldsHook_Fire(t *testing.T) {
////	df := &DefaultFieldsHook{}
////	lre := &logrus.Entry{}
////	df.Fire(lre)
////}
////
////func TestDefaultFieldsHook_Levels(t *testing.T) {
////	df := &DefaultFieldsHook{}
////	df.Levels()
////}
//
//func TestCallerPrettyfier(t *testing.T) {
//
//	frame := runtime.Frame{}
//	FunctionName, FileName := CallerPrettyfier(&frame)
//	if FunctionName == "" {
//		//t.Error("TestCallerPrettyfier error: FunctionName is empty")
//	}
//	if FileName == "" {
//		t.Error("TestCallerPrettyfier error: FileName is empty")
//	}
//}
//
//func TestSetLevel(t *testing.T) {
//	SetLevel("debug")
//}

func TestLog(t *testing.T) {
	//log := GetLog()
	Info("test")
}
