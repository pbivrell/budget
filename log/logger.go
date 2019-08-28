package log

import (
	"github.com/sirupsen/logrus"
)

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var Logger = initialize()

type logger struct {
	logger *logrus.Logger
}

func (l *logger) ActualLogger() *logrus.Logger {
	return l.logger
}

func initialize() *logger {

	logger := &logger{
		logger: logrus.New(),
	}
	logger.logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(LogLevelDebug)
	return logger
}

func (l *logger) SetLevel(level int) {
	switch level {
	case LogLevelDebug:
		l.logger.SetLevel(logrus.DebugLevel)
	case LogLevelInfo:
		l.logger.SetLevel(logrus.InfoLevel)
	case LogLevelWarn:
		l.logger.SetLevel(logrus.WarnLevel)
	case LogLevelError:
		l.logger.SetLevel(logrus.ErrorLevel)
	default:
		l.logger.SetLevel(logrus.InfoLevel)
	}
}

func (l *logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
