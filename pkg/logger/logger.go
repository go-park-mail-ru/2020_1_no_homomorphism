package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

const requestId int = 1

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

func (l *MainLogger) LogError(ctx context.Context, pkg string, funcName string, err error) {
	l.WithFields(logrus.Fields{
		"id":       l.GetIdFromContext(ctx),
		"package":  pkg,
		"function": funcName,
	}).Error(err)
}

func (l *MainLogger) GetIdFromContext(ctx context.Context) string {
	rid, ok := ctx.Value(requestId).(string)
	if !ok {
		l.WithFields(logrus.Fields{
			"id":       "NO_ID",
			"package":  "logger",
			"function": "GetIdFromContext",
		}).Warn("can't get request id from context")
		return ""
	}
	return rid
}

func (l *MainLogger) StartReq(r http.Request, rid string) {
	l.WithFields(logrus.Fields{
		"id":         rid,
		"usr_addr":   r.RemoteAddr,
		"req_URI":    r.RequestURI,
		"method":     r.Method,
		"user_agent": r.UserAgent(),
	}).Info("request started")
}

func (l *MainLogger) EndReq(start time.Time, ctx context.Context) {
	l.WithFields(logrus.Fields{
		"id":              l.GetIdFromContext(ctx),
		"elapsed_time,μs": time.Since(start).Microseconds(),
	}).Info("request ended")
}

func (l *MainLogger) HttpInfo(ctx context.Context, msg string, status int) {
	l.WithFields(logrus.Fields{
		"id":     l.GetIdFromContext(ctx),
		"status": status,
	}).Info(msg)
}

func (l *MainLogger) LogWarning(ctx context.Context, pkg string, funcName string, msg string) {
	l.WithFields(logrus.Fields{
		"id":       l.GetIdFromContext(ctx),
		"package":  pkg,
		"function": funcName,
	}).Warn(msg)
}
