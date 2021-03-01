package behaviago

import (
	"errors"
	"fmt"
	"log"
)

type Logger interface {
	Tracef(format string, params ...interface{})
	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
	Warnf(format string, params ...interface{}) error
	Errorf(format string, params ...interface{}) error
	Criticalf(format string, params ...interface{}) error

	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{}) error
	Error(v ...interface{}) error
	Critical(v ...interface{}) error

	Close()
	Flush()
	Closed() bool
}

var BTGLog Logger = &defaultLogger{}

func SetLogger(log Logger) {
	BTGLog = log
}

type defaultLogger struct {
}

func (l *defaultLogger) Tracef(format string, params ...interface{}) {
	log.Printf("[T] "+format, params...)
}

func (l *defaultLogger) Debugf(format string, params ...interface{}) {
	log.Printf("[D] "+format, params...)
}

func (l *defaultLogger) Infof(format string, params ...interface{}) {
	log.Printf("[I] "+format, params...)
}

func (l *defaultLogger) Warnf(format string, params ...interface{}) error {
	log.Printf("[W] "+format, params...)
	return errors.New(fmt.Sprintf(format, params))
}

func (l *defaultLogger) Errorf(format string, params ...interface{}) error {
	log.Printf("[E] "+format, params...)
	return errors.New(fmt.Sprintf(format, params))
}

func (l *defaultLogger) Criticalf(format string, params ...interface{}) error {
	log.Printf("[C] "+format, params...)
	return errors.New(fmt.Sprintf(format, params))
}

func (l *defaultLogger) Trace(params ...interface{}) {
	log.Println("[T] ", params)
}

func (l *defaultLogger) Debug(params ...interface{}) {
	log.Println("[D] ", params)
}

func (l *defaultLogger) Info(params ...interface{}) {
	log.Println("[I] ", params)
}

func (l *defaultLogger) Warn(params ...interface{}) error {
	log.Println("[W] ", params)
	return errors.New(fmt.Sprintf("", params))
}

func (l *defaultLogger) Error(params ...interface{}) error {
	log.Println("[E] ", params)
	return errors.New(fmt.Sprintf("", params))
}

func (l *defaultLogger) Critical(params ...interface{}) error {
	log.Println("[C] ", params)
	return errors.New(fmt.Sprintf("", params))
}

func (l *defaultLogger) Close() {
}

func (l *defaultLogger) Flush() {
}

func (l *defaultLogger) Closed() bool {
	return true
}
