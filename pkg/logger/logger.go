package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"time"
)

//TODO сделать под 2 конфига - консоль или в файл

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
		"pkg":      pkg,
		"pkg_func": funcName,
	}).Error(err)
}

//функция для access логов, возвращает сгенерированый id запроса
func (l *MainLogger) LogRequest(r http.Request) string {
	rand.Seed(time.Now().UnixNano())
	rid := fmt.Sprintf("%016x", rand.Int())[:5]
	l.WithFields(logrus.Fields{
		"id":        rid,
		"usr_addr":  r.RemoteAddr,
		"req_addr":  r.RequestURI,
		"method":    r.Method,
		"usr_agent": r.UserAgent(),
	}).Info("request")
	return rid
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
		"pkg":      pkg,
		"pkg_func": funcName,
	}).Warn(msg)
}
