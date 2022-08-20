package logger

type Logger interface {
	GetLogger(topic string) Logger
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}
