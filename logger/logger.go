package logger

import (
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

func InitLog(logDir string) error {
	var err error
	traceLogger, err = newLevelLogger(logDir, traceFileName, tracePrefix)
	if err != nil {
		return err
	}
	infoLogger, err = newLevelLogger(logDir, infoFileName, infoPrefix)
	if err != nil {
		return err
	}
	warningLogger, err = newLevelLogger(logDir, warningFileName, warningPrefix)
	if err != nil {
		return err
	}
	errorLogger, err = newLevelLogger(logDir, errorFileName, errorPrefix)
	if err != nil {
		return err
	}
	return nil
}

func newLevelLogger(dir, name, prefix string) (levelLogger *log.Logger, err error) {
	os.MkdirAll(dir, os.ModePerm)
	path := filepath.Join(dir, name)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		Errorln("Failure to open log file", err)
		return
	}
	levelLogger = log.New(file, tracePrefix, logFormat)
	return
}
