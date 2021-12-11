package test

import (
	"github.com/go-logr/logr"
)

var NoOpLogger = logger{}

type logger struct {
	logr.Logger
}

func (l logger) Enabled() bool                                             { return false }
func (l logger) Info(msg string, keysAndValues ...interface{})             {}
func (l logger) Error(err error, msg string, keysAndValues ...interface{}) {}
func (l logger) V(level int) logr.Logger                                   { return l }
func (l logger) WithValues(keysAndValues ...interface{}) logr.Logger       { return l }
func (l logger) WithName(name string) logr.Logger                          { return l }
