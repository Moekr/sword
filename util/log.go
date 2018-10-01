package util

import "log"

var debugEnable bool

func SetDebug(debug bool) {
	debugEnable = debug
	Debug("debug mode enable")
}

func Info(v interface{}) {
	log.Println(v)
}

func Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Debug(v interface{}) {
	if debugEnable {
		log.Println(v)
	}
}

func Debugf(format string, v ...interface{}) {
	if debugEnable {
		log.Printf(format, v...)
	}
}
