package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type MainLogger struct {
	*logrus.Logger
}

func NewLogger(writer io.Writer) *MainLogger {
	baseLogger := logrus.New()
	standardLogger := &MainLogger{baseLogger}
	Formatter := new(logrus.JSONFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	standardLogger.SetFormatter(Formatter)
	standardLogger.SetOutput(writer)
	return standardLogger
}
func (l *MainLogger) LogError(rid string, pkg string, funcName string, err error) {
	l.WithFields(logrus.Fields{
		"id":       rid,
		"package":  pkg,
		"function": funcName,
	}).Error(err)
}

//функция для access логов, возвращает сгенерированый id запроса
func (l *MainLogger) StartReq(r http.Request, rid string) {
	l.WithFields(logrus.Fields{
		"id":         rid,
		"usr_addr":   r.RemoteAddr,
		"req_URI":    r.RequestURI,
		"method":     r.Method,
		"user_agent": r.UserAgent(),
	}).Info("request started")
}

func (l *MainLogger) EndReq(start time.Time, rid string) {
	l.WithFields(logrus.Fields{
		"id":              rid,
		"elapsed_time,μs": time.Since(start).Microseconds(),
	}).Info("request ended")
}

func (l *MainLogger) HttpInfo(rid string, msg string, status int) {
	l.WithFields(logrus.Fields{
		"id":     rid,
		"status": status,
	}).Info(msg)
}

func (l *MainLogger) LogWarning(rid string, pkg string, funcName string, msg string) {
	l.WithFields(logrus.Fields{
		"id":       rid,
		"package":  pkg,
		"function": funcName,
	}).Warn(msg)
}
