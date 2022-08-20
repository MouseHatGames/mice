package stdout

import (
	"fmt"
	"log"

	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/options"
)

type stdoutLogger struct {
	prefix string
}

func Logger() options.Option {
	return func(o *options.Options) {
		o.Logger = NewStdoutLogger(" ")
	}
}

func LoggerPrefix(prefix string) options.Option {
	return func(o *options.Options) {
		o.Logger = NewStdoutLogger(prefix + " ")
	}
}

func NewStdoutLogger(prefix string) logger.Logger {
	return &stdoutLogger{
		prefix: prefix,
	}
}

func (l *stdoutLogger) GetLogger(topic string) logger.Logger {
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
