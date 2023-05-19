// дублирует все функции логгера logrus
package log

import (
	"context"
	"github.com/ManyakRus/logrus"
	"github.com/ManyakRus/starter/logger"
	"io"
	"time"
)

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func WithField(key string, value interface{}) *logrus.Entry {
	return logger.GetLog().WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLog().WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func WithError(err error) *logrus.Entry {
	return GetLog().WithError(err)
}

// Add a context to the log entry.
func WithContext(ctx context.Context) *logrus.Entry {
	return GetLog().WithContext(ctx)
}

// Overrides the time of the log entry.
func WithTime(t time.Time) *logrus.Entry {

	return GetLog().WithTime(t)
}

func Logf(level logrus.Level, format string, args ...interface{}) {
	GetLog().Logf(level, format, args...)
}

func Tracef(format string, args ...interface{}) {
	GetLog().Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	GetLog().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	GetLog().Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	GetLog().Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	GetLog().Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	GetLog().Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	GetLog().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	GetLog().Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	GetLog().Panicf(format, args...)
}

// Log will log a message at the level given as parameter.
// Warning: using Log at Panic or Fatal level will not respectively Panic nor Exit.
// For this behaviour Logger.Panic or Logger.Fatal should be used instead.
func Log(level logrus.Level, args ...interface{}) {
	GetLog().Log(level, args...)
}

func LogFn(level logrus.Level, fn logrus.LogFunction) {
	GetLog().LogFn(level, fn)
}

func Trace(args ...interface{}) {
	GetLog().Trace(args...)
}

func Debug(args ...interface{}) {
	GetLog().Debug(args...)
}

func Info(args ...interface{}) {
	GetLog().Info(args...)
}

func Print(args ...interface{}) {
	GetLog().Print(args...)
}

func Warn(args ...interface{}) {
	GetLog().Warn(args...)
}

func Warning(args ...interface{}) {
	GetLog().Warning(args...)
}

func Error(args ...interface{}) {
	GetLog().Error(args...)
}

func Fatal(args ...interface{}) {
	GetLog().Fatal(args...)
}

func Panic(args ...interface{}) {
	GetLog().Panic(args...)
}

func TraceFn(fn logrus.LogFunction) {
	GetLog().TraceFn(fn)
}

func DebugFn(fn logrus.LogFunction) {
	GetLog().DebugFn(fn)
}

func InfoFn(fn logrus.LogFunction) {
	GetLog().InfoFn(fn)
}

func PrintFn(fn logrus.LogFunction) {
	GetLog().PrintFn(fn)
}

func WarnFn(fn logrus.LogFunction) {
	GetLog().WarnFn(fn)
}

func WarningFn(fn logrus.LogFunction) {
	GetLog().WarningFn(fn)
}

func ErrorFn(fn logrus.LogFunction) {
	GetLog().ErrorFn(fn)
}

func FatalFn(fn logrus.LogFunction) {
	GetLog().FatalFn(fn)
}

func PanicFn(fn logrus.LogFunction) {
	GetLog().PanicFn(fn)
}

func Logln(level logrus.Level, args ...interface{}) {
	GetLog().Logln(level, args...)
}

func Traceln(args ...interface{}) {
	GetLog().Traceln(args...)
}

func Debugln(args ...interface{}) {
	GetLog().Debugln(args...)
}

func Infoln(args ...interface{}) {
	GetLog().Infoln(args...)
}

func Println(args ...interface{}) {
	GetLog().Println(args...)
}

func Warnln(args ...interface{}) {
	GetLog().Warnln(args...)
}

func Warningln(args ...interface{}) {
	GetLog().Warningln(args...)
}

func Errorln(args ...interface{}) {
	GetLog().Errorln(args...)
}

func Fatalln(args ...interface{}) {
	GetLog().Fatalln(args...)
}

func Panicln(args ...interface{}) {
	GetLog().Panicln(args...)
}

func Exit(code int) {
	GetLog().Exit(code)
}

// When file is opened with appending mode, it's safe to
// write concurrently to a file (within 4k message on Linux).
// In these cases user can choose to disable the lock.
func SetNoLock() {
	GetLog().SetNoLock()
}

//// SetLevel sets the logger level.
//func SetLevel(level logrus.Level) {
//}

// GetLevel returns the logger level.
func GetLevel() logrus.Level {
	return GetLog().GetLevel()
}

// AddHook adds a hook to the logger hooks.
func AddHook(hook logrus.Hook) {
	GetLog().AddHook(hook)
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func IsLevelEnabled(level logrus.Level) bool {
	return GetLog().IsLevelEnabled(level)
}

// SetFormatter sets the logger formatter.
func SetFormatter(formatter logrus.Formatter) {
	GetLog().SetFormatter(formatter)
}

// SetOutput sets the logger output.
func SetOutput(output io.Writer) {
	GetLog().SetOutput(output)
}

func SetReportCaller(reportCaller bool) {
	GetLog().SetReportCaller(reportCaller)
}

// ReplaceHooks replaces the logger hooks and returns the old ones
func ReplaceHooks(hooks logrus.LevelHooks) logrus.LevelHooks {
	GetLog()
	return ReplaceHooks(hooks)
}

// SetBufferPool sets the logger buffer pool.
func SetBufferPool(pool logrus.BufferPool) {
	GetLog().SetBufferPool(pool)
}
