package gmx

import (
	"fmt"
	"sync"
)

// ---------------------------------------------------------------------------
// Logging integration.

// Avoid importing the log type information unnecessarily.  There's a small cost
// associated with using an interface rather than the type.  Depending on how
// often the logger is plugged in, it would be worth using the type instead.
type log_Logger interface {
	Output(calldepth int, s string) error
}

var (
	globalLogger log_Logger
	globalMutex  sync.Mutex
	logPrefix    = "[gmx]"
)

// Specify the *log.Logger object where log messages should be sent to.
func SetLogger(logger log_Logger) {
	globalLogger = logger
}

// Enable the delivery of debug messages to the logger.  Only meaningful
// if a logger is also set.
func SetDebug(debug bool) {
	Cfg.Debug = debug
}

func log(v ...interface{}) {
	if globalLogger != nil {
		str := logPrefix + fmt.Sprint(v...)
		globalLogger.Output(2, str)
	}
}

func logln(v ...interface{}) {
	if globalLogger != nil {
		str := logPrefix + fmt.Sprintln(v...)
		globalLogger.Output(2, str)
	}
}

func logf(format string, v ...interface{}) {
	if globalLogger != nil {
		str := logPrefix + fmt.Sprintf(format, v...)
		globalLogger.Output(2, str)
	}
}

func debug(v ...interface{}) {
	if Cfg.Debug && globalLogger != nil {
		str := logPrefix + fmt.Sprint(v...)
		globalLogger.Output(2, str)
	}
}

func debugln(v ...interface{}) {
	if Cfg.Debug && globalLogger != nil {
		str := logPrefix + fmt.Sprintln(v...)
		globalLogger.Output(2, str)
	}
}

func debugf(format string, v ...interface{}) {
	if Cfg.Debug && globalLogger != nil {
		str := logPrefix + fmt.Sprintf(format, v...)
		globalLogger.Output(2, str)
	}
}
