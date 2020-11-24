package logger

import (
	"fmt"
	"log"
)

type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}

type stdoutLogger struct{}

func NewStdoutLogger() Logger {
	return new(stdoutLogger)
}

func (*stdoutLogger) Debugf(msg string, args ...interface{}) {
	log.Printf("[DEBUG] %s", fmt.Sprintf(msg, args...))
}

func (*stdoutLogger) Infof(msg string, args ...interface{}) {
	log.Printf("[INFO] %s", fmt.Sprintf(msg, args...))
}

func (*stdoutLogger) Errorf(msg string, args ...interface{}) {
	log.Printf("[ERROR] %s", fmt.Sprintf(msg, args...))
}
