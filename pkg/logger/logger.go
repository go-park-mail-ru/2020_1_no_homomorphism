package pkg

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type MainLogger struct {
	*logrus.Logger
}

func NewLogger() *MainLogger {
	baseLogger := logrus.New()
	standardLogger := &MainLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}

	return standardLogger
}

func (l *MainLogger) LogError(err error, r http.Request) {
	l.WithField(logrus.Fields{
		"user_addr": r.RemoteAddr,
		"error": err.Error()
	}).Error(err)
}

func (l *MainLogger) LogInfo() {

}

func (l *MainLogger) LogWarning() {

}
