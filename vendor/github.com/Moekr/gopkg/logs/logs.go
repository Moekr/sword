package logs

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func InitLogs(logsPath string) {
	logger = logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(newFormatter())
	if len(logsPath) == 0 {
		logger.SetOutput(os.Stdout)
	} else {
		logger.AddHook(newRotateHook(logsPath))
	}
}

func Trace(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

func Debug(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Panic(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}
