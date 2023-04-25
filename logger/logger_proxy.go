// дублирует все функции логгера logrus
package logger

//
//import (
//	"context"
//	"github.com/sirupsen/logrus"
//	"io"
//	"time"
//)
//
//// WithField allocates a new entry and adds a field to it.
//// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
//// this new returned entry.
//// If you want multiple fields, use `WithFields`.
//func WithField(key string, value interface{}) *logrus.Entry {
//	return log.WithField(key, value)
//}
//
//// Adds a struct of fields to the log entry. All it does is call `WithField` for
//// each `Field`.
//func WithFields(fields logrus.Fields) *logrus.Entry {
//	GetLog()
//	return log.WithFields(fields)
//}
//
//// Add an error as single field to the log entry.  All it does is call
//// `WithError` for the given `error`.
//func WithError(err error) *logrus.Entry {
//	GetLog()
//	return log.WithError(err)
//}
//
//// Add a context to the log entry.
//func WithContext(ctx context.Context) *logrus.Entry {
//	GetLog()
//	return log.WithContext(ctx)
//}
//
//// Overrides the time of the log entry.
//func WithTime(t time.Time) *logrus.Entry {
//	GetLog()
//	return WithTime(t)
//}
//
//func Logf(level logrus.Level, format string, args ...interface{}) {
//	GetLog()
//	log.Logf(level, format, args...)
//}
//
//func Tracef(format string, args ...interface{}) {
//	GetLog()
//	log.Tracef(format, args...)
//}
//
//func Debugf(format string, args ...interface{}) {
//	GetLog()
//	log.Debugf(format, args...)
//}
//
//func Infof(format string, args ...interface{}) {
//	GetLog()
//	log.Infof(format, args...)
//}
//
//func Printf(format string, args ...interface{}) {
//	GetLog()
//	log.Printf(format, args...)
//}
//
//func Warnf(format string, args ...interface{}) {
//	GetLog()
//	log.Warnf(format, args...)
//}
//
//func Warningf(format string, args ...interface{}) {
//	GetLog()
//	log.Warningf(format, args...)
//}
//
//func Errorf(format string, args ...interface{}) {
//	GetLog()
//	log.Errorf(format, args...)
//}
//
//func Fatalf(format string, args ...interface{}) {
//	GetLog()
//	log.Fatalf(format, args...)
//}
//
//func Panicf(format string, args ...interface{}) {
//	GetLog()
//	log.Panicf(format, args...)
//}
//
//// Log will log a message at the level given as parameter.
//// Warning: using Log at Panic or Fatal level will not respectively Panic nor Exit.
//// For this behaviour Logger.Panic or Logger.Fatal should be used instead.
//func Log(level logrus.Level, args ...interface{}) {
//	GetLog()
//	log.Log(level, args...)
//}
//
//func LogFn(level logrus.Level, fn logrus.LogFunction) {
//	GetLog()
//	log.LogFn(level, fn)
//}
//
//func Trace(args ...interface{}) {
//	GetLog()
//	log.Trace(args...)
//}
//
//func Debug(args ...interface{}) {
//	GetLog()
//	log.Debug(args...)
//}
//
//func Info(args ...interface{}) {
//	GetLog()
//	log.Info(args...)
//}
//
//func Print(args ...interface{}) {
//	GetLog()
//	log.Print(args...)
//}
//
//func Warn(args ...interface{}) {
//	GetLog()
//	log.Warn(args...)
//}
//
//func Warning(args ...interface{}) {
//	GetLog()
//	log.Warning(args...)
//}
//
//func Error(args ...interface{}) {
//	GetLog()
//	log.Error(args...)
//}
//
//func Fatal(args ...interface{}) {
//	GetLog()
//	log.Fatal(args...)
//}
//
//func Panic(args ...interface{}) {
//	GetLog()
//	log.Panic(args...)
//}
//
//func TraceFn(fn logrus.LogFunction) {
//	GetLog()
//	log.TraceFn(fn)
//}
//
//func DebugFn(fn logrus.LogFunction) {
//	GetLog()
//	log.DebugFn(fn)
//}
//
//func InfoFn(fn logrus.LogFunction) {
//	GetLog()
//	log.InfoFn(fn)
//}
//
//func PrintFn(fn logrus.LogFunction) {
//	GetLog()
//	log.PrintFn(fn)
//}
//
//func WarnFn(fn logrus.LogFunction) {
//	GetLog()
//	log.WarnFn(fn)
//}
//
//func WarningFn(fn logrus.LogFunction) {
//	GetLog()
//	log.WarningFn(fn)
//}
//
//func ErrorFn(fn logrus.LogFunction) {
//	GetLog()
//	log.ErrorFn(fn)
//}
//
//func FatalFn(fn logrus.LogFunction) {
//	GetLog()
//	log.FatalFn(fn)
//}
//
//func PanicFn(fn logrus.LogFunction) {
//	GetLog()
//	log.PanicFn(fn)
//}
//
//func Logln(level logrus.Level, args ...interface{}) {
//	GetLog()
//	log.Logln(level, args...)
//}
//
//func Traceln(args ...interface{}) {
//	GetLog()
//	log.Traceln(args...)
//}
//
//func Debugln(args ...interface{}) {
//	GetLog()
//	log.Debugln(args...)
//}
//
//func Infoln(args ...interface{}) {
//	GetLog()
//	log.Infoln(args...)
//}
//
//func Println(args ...interface{}) {
//	GetLog()
//	log.Println(args...)
//}
//
//func Warnln(args ...interface{}) {
//	GetLog()
//	log.Warnln(args...)
//}
//
//func Warningln(args ...interface{}) {
//	GetLog()
//	log.Warningln(args...)
//}
//
//func Errorln(args ...interface{}) {
//	GetLog()
//	log.Errorln(args...)
//}
//
//func Fatalln(args ...interface{}) {
//	GetLog()
//	log.Fatalln(args...)
//}
//
//func Panicln(args ...interface{}) {
//	GetLog()
//	log.Panicln(args...)
//}
//
//func Exit(code int) {
//	GetLog()
//	log.Exit(code)
//}
//
////When file is opened with appending mode, it's safe to
////write concurrently to a file (within 4k message on Linux).
////In these cases user can choose to disable the lock.
//func SetNoLock() {
//	GetLog()
//	log.SetNoLock()
//}
//
////// SetLevel sets the logger level.
////func SetLevel(level logrus.Level) {
////}
//
//// GetLevel returns the logger level.
//func GetLevel() logrus.Level {
//	GetLog()
//	return log.GetLevel()
//}
//
//// AddHook adds a hook to the logger hooks.
//func AddHook(hook logrus.Hook) {
//	GetLog()
//	log.AddHook(hook)
//}
//
//// IsLevelEnabled checks if the log level of the logger is greater than the level param
//func IsLevelEnabled(level logrus.Level) bool {
//	GetLog()
//	return log.IsLevelEnabled(level)
//}
//
//// SetFormatter sets the logger formatter.
//func SetFormatter(formatter logrus.Formatter) {
//	GetLog()
//	log.SetFormatter(formatter)
//}
//
//// SetOutput sets the logger output.
//func SetOutput(output io.Writer) {
//	GetLog()
//	log.SetOutput(output)
//}
//
//func SetReportCaller(reportCaller bool) {
//	GetLog()
//	log.SetReportCaller(reportCaller)
//}
//
//// ReplaceHooks replaces the logger hooks and returns the old ones
//func ReplaceHooks(hooks logrus.LevelHooks) logrus.LevelHooks {
//	GetLog()
//	return ReplaceHooks(hooks)
//}
//
//// SetBufferPool sets the logger buffer pool.
//func SetBufferPool(pool logrus.BufferPool) {
//	GetLog()
//	log.SetBufferPool(pool)
//}
