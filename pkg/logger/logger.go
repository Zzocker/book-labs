package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"

	"github.com/Zzocker/book-labs/pkg/errors"
)

var lg *logger //nolint:gochecknoglobals //it's ok

// Logger : interface represent a logger.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type logger struct {
	lg *logrus.Entry
}

func NewServiceLogger(level, srName, version string) {
	var l logrus.Level
	switch strings.ToLower((level)) {
	case "error":
		l = logrus.ErrorLevel
	case "warn":
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	default:
		l = logrus.InfoLevel
	}
	log := logrus.New()
	log.SetLevel(l)
	format := &logrus.JSONFormatter{
		TimestampFormat: "15:04:05 02/01/2006",
	}
	log.SetFormatter(format)
	log.SetReportCaller(false)
	lg = &logger{log.WithFields(map[string]interface{}{
		"service": srName,
		"version": version,
	})}
}

func (l *logger) debugf(format string, args ...interface{}) {
	l.lg.Debugf(format, args...)
}

func (l *logger) infof(format string, args ...interface{}) {
	l.lg.Infof(format, args...)
}

func (l *logger) warnf(format string, args ...interface{}) {
	l.lg.Warnf(format, args...)
}

func (l *logger) errorf(format string, args ...interface{}) {
	l.lg.Errorf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	lg.debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	lg.infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	lg.warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	lg.errorf(format, args...)
}

func WithFields(fields map[string]interface{}) Logger {
	return lg.lg.WithFields(fields)
}

func SystemErr(err error) {
	_, ok := err.(*errors.Error) //nolint:errorlint //this Error is custom implemnation not using errors
	if !ok {
		lg.lg.Error(err)

		return
	}

	entry := lg.lg.WithFields(map[string]interface{}{
		"operations": errors.Ops(err),
		"code":       codes.Code(errors.ErrCode(err)).String(),
		"request_id": errors.ErrReqID(err),
	})

	//nolint:exhaustive //it's ok
	switch errors.ErrSeverity(err) {
	case errors.SeverityWarn:
		entry.Warnf("%v", err)
	case errors.SeverityInfo:
		entry.Infof("%v", err)
	case errors.SeverityDebug:
		entry.Debugf("%v", err)
	default:
		entry.Errorf("%v", err)
	}
}
