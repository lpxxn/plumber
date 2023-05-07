package log

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/google/wire"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func provideLog(logLevel string) Loger {
	logrusLog := logrus.New()
	logrusLog.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true, CallerPrettyfier: getCallerPrettyfier()}
	logrusLog.SetReportCaller(true)
	jackOut := &lumberjack.Logger{
		Filename:   "logs/plumber.log",
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30, //days
	}
	logrusLog.SetOutput(io.MultiWriter(jackOut, os.Stdout))
	switch logLevel {
	case "debug":
		logrusLog.SetLevel(logrus.DebugLevel)
	case "info":
		logrusLog.SetLevel(logrus.InfoLevel)
	case "warn":
		logrusLog.SetLevel(logrus.WarnLevel)
	case "error":
		logrusLog.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrusLog.SetLevel(logrus.FatalLevel)
	case "trace":
		logrusLog.SetLevel(logrus.TraceLevel)
	default:
		logrusLog.SetLevel(logrus.InfoLevel)
	}
	logrusLog.Debugf("log init success")
	return logrusLog
}

func getCallerPrettyfier() func(f *runtime.Frame) (string, string) {
	return func(f *runtime.Frame) (string, string) {
		// https://github.com/sirupsen/logrus/blob/v1.9.0/example_custom_caller_test.go
		// https://github.com/kubernetes/klog/blob/v2.90.1/klog.go#L644
		_, file, line, ok := runtime.Caller(10)
		if !ok {
			file = "???"
			line = 1
		} else {
			//if slash := strings.LastIndex(file, "/"); slash >= 0 {
			//	file = file[slash+1:]
			//}
		}
		return "", fmt.Sprintf("%s:%d", file, line)
	}
}

func logLevel() string {
	return os.Getenv("LOG_LEVEL")
}

var logSet = wire.NewSet(logLevel, provideLog)
