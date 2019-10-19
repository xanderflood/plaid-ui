package tools

import (
	"os"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
)

//Logger exposes a minimal logging interface
//go:generate counterfeiter . Logger
type Logger interface {
	Errorf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Warningf(message string, args ...interface{})
	Tracef(message string, args ...interface{})
	Debugf(message string, args ...interface{})
}

//NewStdoutLogger returns a Logger that writes colorized
//output to stdout.
func NewStdoutLogger() Logger {
	loggo.RegisterWriter("stdout", loggocolor.NewWriter(os.Stdout)) //nolint:errcheck
	return loggo.GetLogger("stdout")
}
