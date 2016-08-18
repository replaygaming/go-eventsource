package main

import "log"

func logMessage(prefix string, message string, v ...interface{}) {
	log.Printf("["+prefix+"] "+message, v...)
}

func Debug(message string, v ...interface{}) {
	logMessage("DEBUG", message, v...)
}

func Info(message string, v ...interface{}) {
	logMessage("INFO", message, v...)
}

func Warn(message string, v ...interface{}) {
	logMessage("WARN", message, v...)
}

func Fatal(message string, v ...interface{}) {
	log.Fatalf("[FATAL] "+message, v...)
}
