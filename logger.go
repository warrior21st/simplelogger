package logger

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

var (
	_infoLogger       *log.Logger
	_warningLogger    *log.Logger
	_errorLogger      *log.Logger
	_lastLogFileDate  string
	_fileCreateLocker sync.Mutex
)

func createLogFileIfNotExist() {
	nowStr := time.Now().Format("20060102")
	if _lastLogFileDate != nowStr {
		_fileCreateLocker.Lock()

		programPath, _ := os.Getwd()
		logFilePath := programPath + string(os.PathSeparator) + "logs"
		if !isExist(logFilePath) {
			os.Mkdir(logFilePath, os.ModePerm)
		}
		logFilePath += string(os.PathSeparator) + nowStr + ".log"
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("open log file failure:", err)
		}

		_infoLogger = log.New(logFile, "[Info]", log.Ldate|log.Ltime)
		_warningLogger = log.New(logFile, "[Warning]", log.Ldate|log.Ltime)
		_errorLogger = log.New(logFile, "[Error]", log.Ldate|log.Ltime|log.Llongfile)

		_lastLogFileDate = nowStr

		_fileCreateLocker.Unlock()
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func LogError(err error) {
	createLogFileIfNotExist()
	msg := err.Error()
	msg += "\n"
	msg += string(debug.Stack())
	fmt.Println(msg)
	_errorLogger.Println(msg)
}

func LogWarning(msg string) {
	createLogFileIfNotExist()
	fmt.Println(msg)
	_warningLogger.Println(msg)
}

func LogInfo(msg string) {
	createLogFileIfNotExist()
	fmt.Println(msg)
	_infoLogger.Println(msg)
}
