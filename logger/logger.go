package logger

import (
	"fmt"
	"log"
)

type Logger interface {
	GetLogger(topic string) Logger
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}

type stdoutLogger struct {
	prefix string
}

func NewStdoutLogger() Logger {
	return &stdoutLogger{
		prefix: " ",
	}
}

func (l *stdoutLogger) GetLogger(topic string) Logger {
	return &stdoutLogger{
		prefix: fmt.Sprintf("%s[%s] ", l.prefix, topic),
	}
}

func (l *stdoutLogger) Debugf(msg string, args ...interface{}) {
	log.Printf("[DEBUG]%s%s", l.prefix, fmt.Sprintf(msg, args...))
}

func (l *stdoutLogger) Infof(msg string, args ...interface{}) {
	log.Printf("[INFO]%s%s", l.prefix, fmt.Sprintf(msg, args...))
}

func (l *stdoutLogger) Errorf(msg string, args ...interface{}) {
	log.Printf("[ERROR]%s%s", l.prefix, fmt.Sprintf(msg, args...))
}
