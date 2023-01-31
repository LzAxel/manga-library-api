package logger

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry
	WithContext(ctx context.Context) *logrus.Entry
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

type logrusLogger struct {
	logger *logrus.Logger
}

func NewLogrusLogger(level string, isDebug, isProd bool) *logrusLogger {
	l := logrus.New()

	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		panic(err)
	}
	l.SetLevel(logrusLevel)
	if isDebug {
		l.SetLevel(logrus.DebugLevel)
	}
	if isProd {
		l.SetFormatter(&logrus.JSONFormatter{})
	} else {
		l.SetOutput(colorable.NewColorableStdout())
		l.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				filename := path.Base(f.File)

				var trucatedFunc string

				splittedFunc := strings.Split(f.Function, "/")
				if len(splittedFunc) >= 3 {
					trucatedFunc = strings.Join(splittedFunc[1:], "/")
				}

				return fmt.Sprintf("%v:%v", filename, f.Line),
					fmt.Sprintf("%v()", trucatedFunc)
			},
			ForceColors:            true,
			FullTimestamp:          false,
			DisableLevelTruncation: true,
		})
	}

	return &logrusLogger{logger: l}
}

func (l *logrusLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}

func (l *logrusLogger) WithError(err error) *logrus.Entry {
	return l.logger.WithError(err)
}
func (l *logrusLogger) WithContext(ctx context.Context) *logrus.Entry {
	return l.logger.WithContext(ctx)
}
func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}
func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}
func (l *logrusLogger) Warningf(format string, args ...interface{}) {
	l.logger.Warningf(format, args...)
}
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}
func (l *logrusLogger) Debugln(args ...interface{}) {
	l.logger.Debugln(args...)
}
func (l *logrusLogger) Infoln(args ...interface{}) {
	l.logger.Infoln(args...)
}
func (l *logrusLogger) Warningln(args ...interface{}) {
	l.logger.Warningln(args...)
}
func (l *logrusLogger) Errorln(args ...interface{}) {
	l.logger.Errorln(args...)
}
func (l *logrusLogger) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args...)
}
func (l *logrusLogger) Panicln(args ...interface{}) {
	l.logger.Panicln(args...)
}
