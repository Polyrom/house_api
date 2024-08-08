package logging

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

// TODO: maybe figure out WithField - then add to loggin mw
type Logger interface {
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}

var logger *logrus.Logger
var once sync.Once

func New() *logrus.Logger {
	once.Do(func() {
		l := logrus.New()
		l.SetOutput(os.Stdout)
		l.SetReportCaller(true)
		l.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
			DisableColors: true,
			FullTimestamp: true,
		}
		l.SetLevel(logrus.TraceLevel)
		logger = l
	})
	return logger
}
