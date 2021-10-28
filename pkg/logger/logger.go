package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

var lg Interface //nolint:gochecknoglobals //provide global logger

type Interface interface {
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
}

func Setup(level, srvName, srVersion string) {
	var l logrus.Level
	switch strings.ToLower(level) {
	case "debug":
		l = logrus.DebugLevel
	case "info":
		l = logrus.InfoLevel
	case "warn":
		l = logrus.WarnLevel
	case "error":
		l = logrus.ErrorLevel
	default:
		l = logrus.InfoLevel
	}
	log := logrus.New()
	log.Level = l
	format := logrus.JSONFormatter{
		TimestampFormat: "15:04:05 02/01/2006",
	}
	log.SetFormatter(&format)
	log.SetReportCaller(false)
	lg = &internalLg{
		lg: log.WithFields(map[string]interface{}{
			"service": srvName,
			"version": srVersion,
		}),
	}
}

type internalLg struct {
	lg *logrus.Entry
}

func (l *internalLg) Debugf(format string, args ...interface{}) {
	l.lg.Debugf(format, args...)
}

func (l *internalLg) Debug(args ...interface{}) {
	l.lg.Debug(args...)
}

func (l *internalLg) Infof(format string, args ...interface{}) {
	l.lg.Infof(format, args...)
}

func (l *internalLg) Info(args ...interface{}) {
	l.lg.Info(args...)
}

func (l *internalLg) Warnf(format string, args ...interface{}) {
	l.lg.Warnf(format, args...)
}

func (l *internalLg) Warn(args ...interface{}) {
	l.lg.Warn(args...)
}

func (l *internalLg) Errorf(format string, args ...interface{}) {
	l.lg.Errorf(format, args...)
}

func (l *internalLg) Error(args ...interface{}) {
	l.lg.Error(args...)
}

func Debugf(format string, args ...interface{}) {
	lg.Debugf(format, args)
}

func Debug(args ...interface{}) {
	lg.Debug(args)
}

func Infof(format string, args ...interface{}) {
	lg.Infof(format, args...)
}

func Info(args ...interface{}) {
	lg.Info(args)
}

func Warnf(format string, args ...interface{}) {
	lg.Warnf(format, args...)
}

func Warn(args ...interface{}) {
	lg.Warn(args)
}

func Errorf(format string, args ...interface{}) {
	lg.Errorf(format, args...)
}

func Error(args ...interface{}) {
	lg.Error(args)
}
