package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	tracePrefix     string = "TRACE: "
	traceFileName   string = "trace.log"
	infoPrefix      string = "INFO: "
	infoFileName    string = "info.log"
	warningPrefix   string = "WARNING: "
	warningFileName string = "warning.log"
	errorPrefix     string = "ERROR: "
	errorFileName   string = "error.log"
	logFormat       int    = log.Ldate | log.Ltime | log.Lshortfile
)

var (
	traceLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func Traceln(args ...interface{}) {
	logMsgLn(traceLogger, tracePrefix, args...)
}

func Tracef(format string, args ...interface{}) {
	logMsgF(traceLogger, tracePrefix, format, args...)
}

func Intoln(args ...interface{}) {
	logMsgLn(infoLogger, infoPrefix, args...)
}

func Intof(format string, args ...interface{}) {
	logMsgF(infoLogger, infoPrefix, format, args...)
}

func Warningln(args ...interface{}) {
	logMsgLn(warningLogger, warningPrefix, args...)
}

func Warningf(format string, args ...interface{}) {
	logMsgF(warningLogger, warningPrefix, format, args...)
}

func Errorln(args ...interface{}) {
	logMsgLn(errorLogger, errorPrefix, args...)
}

func Errorf(format string, args ...interface{}) {
	logMsgF(errorLogger, errorPrefix, format, args...)
}

func logMsgLn(logger *log.Logger, prefix string, args ...interface{}) {
	if logger != nil {
		logger.Println(args...)
	} else {
		log.SetPrefix(prefix)
		log.Println(args...)
	}
}

func logMsgF(logger *log.Logger, prefix string, f string, args ...interface{}) {
	if logger != nil {
		logger.Printf(f, args...)
	} else {
		log.SetPrefix(prefix)
		log.Printf(f, args...)
	}
}

func InitLogger(logDir string) error {
	if len(logDir) == 0 {
		return initStdLogger()
	}
	return InitFileLogger(logDir)
}

func InitFileLogger(logDir string) error {
	var err error
	trace, err := logFile(logDir, traceFileName)
	if err != nil {
		return err
	}
	warn, err := logFile(logDir, warningFileName)
	if err != nil {
		return err
	}
	info, err := logFile(logDir, infoFileName)
	if err != nil {
		return err
	}
	error, err := logFile(logDir, errorFileName)
	if err != nil {
		return err
	}
	initLogger(trace, warn, info, error)
	return nil
}

func logFile(dir, name string) (file io.Writer, err error) {
	os.MkdirAll(dir, os.ModePerm)
	path := filepath.Join(dir, name)
	file, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	return
}

func initStdLogger() error {
	initLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	return nil
}

func initLogger(trace io.Writer, info io.Writer, warn io.Writer, error io.Writer) {
	traceLogger = log.New(trace, tracePrefix, logFormat)
	infoLogger = log.New(info, tracePrefix, logFormat)
	warningLogger = log.New(warn, tracePrefix, logFormat)
	errorLogger = log.New(error, tracePrefix, logFormat)
}
