package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// levels
const (
	debugLevel   = 0
	releaseLevel = 1
	errorLevel   = 2
	fatalLevel   = 3
)

const (
	printDebugLevel   = "[debug  ] "
	printReleaseLevel = "[release] "
	printErrorLevel   = "[error  ] "
	printFatalLevel   = "[fatal  ] "
)

type Logger struct {
	level      int
	baseLogger *log.Logger
	baseFile   *os.File
	logName    string
	pathName   string
}

func New(strLevel string, pathname string, flag int, logName string) (*Logger, error) {
	// level
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
	case "release":
		level = releaseLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		return nil, errors.New("unknown level: " + strLevel)
	}

	// logger
	var baseLogger *log.Logger
	var baseFile *os.File

	//if pathname != "" {
	//	now := time.Now()
	//
	//	//filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
	//	//	now.Year(),
	//	//	now.Month(),
	//	//	now.Day(),
	//	//	now.Hour(),
	//	//	now.Minute(),
	//	//	now.Second())
	//	filename := fmt.Sprintf("%d%02d%02d.log",
	//		now.Year(),
	//		now.Month(),
	//		now.Day())
	//
	//	file, err := os.Create(path.Join(pathname, filename))
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	baseLogger = log.New(file, "", flag)
	//	baseFile = file
	//} else {
	//	baseLogger = log.New(os.Stdout, "", flag)
	//}

	file := createLogFile(pathname, logName)
	if file != nil {
		baseLogger = log.New(file, "", flag)
		baseFile = file
	} else {
		baseLogger = log.New(os.Stdout, "", flag)
	}

	// new
	logger := new(Logger)
	logger.level = level
	logger.baseLogger = baseLogger
	logger.baseFile = baseFile
	logger.logName = logName
	logger.pathName = pathname

	return logger, nil
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	logger.baseLogger.Output(3, fmt.Sprintf(format, a...))

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Release(format string, a ...interface{}) {
	logger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func (logger *Logger) ResetLogFile() {
	if logger.pathName != "" {
		// 创建日志文件
		file := createLogFile(logger.pathName, logger.logName)
		if file != nil {
			logger.baseLogger.SetOutput(file)
		}
	}
}

var gLogger, _ = New("debug", "", log.LstdFlags, "")

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
	}
}

func Debug(format string, a ...interface{}) {
	gLogger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func Release(format string, a ...interface{}) {
	gLogger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func Close() {
	gLogger.Close()
}

func createLogFile(pathName string, logName string) *os.File {
	if len(pathName) == 0 || len(logName) == 0 {
		return nil
	}

	dayStr := getTodayDay()
	path := pathName + "/" + dayStr
	if exists, _ := pathExists(path); !exists {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("ResetLogFile Mkdir path : %v error :", path, err)
		}
	}

	logPath := path + "/" + logName
	file, err := os.Create(logPath)
	if err != nil {
		fmt.Println("ResetLogFile path : %v error : ", logPath, err)
		return nil
	}
	return file
}

func ResetLog() {
	if gLogger != nil {
		gLogger.ResetLogFile()
	}
}

func getTodayDay() string {
	now := time.Now()
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		filename := fmt.Sprintf("%d%02d%02d",
			now.Year(),
			now.Month(),
			now.Day())
		return filename
	}
	return now.In(loc).Format("20060102")
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
