package logs

import "log"

var debugEnable bool

func SetDebug(debug bool) {
	debugEnable = debug
	Debug("debug mode enable")
}

func Debug(format string, v ...interface{}) {
	if debugEnable {
		log.Printf(format, v...)
	}
}

func Info(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Warn(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}
