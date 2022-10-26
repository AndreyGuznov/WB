package log

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	dbg "runtime/debug"
	"strings"
	"sync"
)

const (
	errorLevel string = "ERROR"
	warnLevel  string = "WARN "
	infoLevel  string = "INFO "
	debugLevel string = "DEBUG"
)

const (
	// debug flag
	debug bool = true
)

// Logger stores all known loggers
type Logger struct {
	loggers map[string]*log.Logger
}

var (
	instance *Logger

	initLog sync.Once
)

// GetLogger creates the log file if not present and initializes the various loggers for use
func GetLogger() *Logger {
	initLog.Do(func() {
		loggers := make(map[string]*log.Logger)

		loggers[errorLevel] = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)
		loggers[warnLevel] = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)
		loggers[infoLevel] = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)
		if debug {
			loggers[debugLevel] = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)
		} else {
			//discard the messages
			loggers[debugLevel] = log.New(ioutil.Discard, "", log.LstdFlags|log.Lmicroseconds)
		}

		instance = &Logger{loggers: loggers}
	})

	return instance
}

func (rlog *Logger) logInt(level, msg string, err error) {
	// To get the file of the caller to log that info in the log
	// skip 1 returns logger.go and hence increasing it to 2 assumming only 1 level of call to logger
	_, file, line, ok := runtime.Caller(2)
	if ok {
		fileParts := strings.Split(file, "/")
		msg = fmt.Sprintf("%s %s:%d %s", level, fileParts[len(fileParts)-1], line, msg)
	}

	errMsg := ""
	if err != nil {
		errMsg = fmt.Sprintf("[%s]", err)
	}
	if level == errorLevel {
		errMsg += fmt.Sprintf("\nThe stacktrace is: %s", string(dbg.Stack()))
	}

	logger := rlog.loggers[level]
	logger.Println(msg, errMsg)
}

// Err logs the message & error as an error
func Err(msg string, err error) {
	GetLogger().logInt(errorLevel, msg, err)
}

// Warn logs the message & error as a warning
func Warn(msg string, err error) {
	GetLogger().logInt(warnLevel, msg, err)
}

// Info logs the message as an info
func Info(msg string) {
	GetLogger().logInt(infoLevel, msg, nil)
}

// Debug logs the message only when the debug flag is turned on
func Debug(msg string) {
	GetLogger().logInt(debugLevel, msg, nil)
}
