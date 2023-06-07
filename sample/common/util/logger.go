package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var fileDate string
var basePath string

// InitLogger 初始化日志对象
func InitLogger(path string) (err error) {
	basePath = path
	if warnFile != nil {
		warnFile.Close()
	}
	if errFile != nil {
		errFile.Close()
	}
	if infoFile != nil {
		infoFile.Close()
	}
	if perfFile != nil {
		perfFile.Close()
	}
	warnLogger, warnFile, err = newLogger(path, "warn")
	if err != nil {
		return err
	}
	errLogger, errFile, err = newLogger(path, "error")
	if err != nil {
		return err
	}
	infoLogger, infoFile, err = newLogger(path, "info")
	if err != nil {
		return err
	}
	perfLogger, perfFile, err = newLogger(path, "perf")
	if err != nil {
		return err
	}
	return nil
}

// CloseLogger 释放日志对象
func CloseLogger() {
	closeLoggerFile()
}

// LogWarn 警告日志
func LogWarn(content ...interface{}) {
	if time.Now().Format("2006-01-02") != fileDate {
		InitLogger(basePath)
	}
	go func(content ...interface{}) {
		fmt.Println(content)
		warnLogger.Println(content)
	}(content...)
}

// LogError 错误日志
func LogError(content ...interface{}) {
	if time.Now().Format("2006-01-02") != fileDate {
		InitLogger(basePath)
	}
	go func(content ...interface{}) {
		fmt.Println(content)
		errLogger.Println(content)
	}(content...)
}

// LogErrorInterface 错误日志
func LogErrorInterface(format string, a ...interface{}) {
	if time.Now().Format("2006-01-02") != fileDate {
		InitLogger(basePath)
	}
	content := fmt.Sprintf(format, a...)
	go func(content string) {
		fmt.Println(content)
		errLogger.Println(content)
	}(content)
}

// LogInfo 信息日志
func LogInfo(content ...interface{}) {
	if time.Now().Format("2006-01-02") != fileDate {
		InitLogger(basePath)
	}
	go func(content ...interface{}) {
		fmt.Println(content)
		infoLogger.Println(content)
	}(content...)
}

// LogInfo 信息日志
func LogPerf(content ...interface{}) {
	if time.Now().Format("2006-01-02") != fileDate {
		InitLogger(basePath)
	}
	go func(content ...interface{}) {
		fmt.Println(content)
		infoLogger.Println(content)
	}(content...)
}

var (
	warnLogger *log.Logger
	warnFile   *os.File
	errLogger  *log.Logger
	errFile    *os.File
	infoLogger *log.Logger
	infoFile   *os.File
	perfLogger *log.Logger
	perfFile   *os.File

)

// newLogger 创建日志对象
func newLogger(basePath, prefix string) (*log.Logger, *os.File, error) {
	currentDate := time.Now().Format("2006-01-02")
	fileDate = currentDate
	name := prefix + "-" + currentDate + ".txt"
	filePath := filepath.Join(basePath, prefix)
	fileName := filepath.Join(filePath, name)
	var logFile *os.File
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			// 创建目录
			if err := os.MkdirAll(filePath, 0777); err != nil {
				return nil, nil, err
			}
			// 创建文件
			logFile, err = os.Create(fileName)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	} else {
		// 打开文件
		logFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		// logFile, err = os.Open(fileName)
		if err != nil {
			return nil, nil, err
		}
	}

	logger := log.New(logFile, fmt.Sprintf("[%s]", prefix), log.LstdFlags)
	return logger, logFile, nil
}

// closeLoggerFile 释放日志文件对象
func closeLoggerFile() {
	warnFile.Close()
	errFile.Close()
	infoFile.Close()
	perfFile.Close()
}
