package simplelogger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	// _infoLogger       *log.Logger
	// _warningLogger    *log.Logger
	// _errorLogger      *log.Logger
	_logWriter        *log.Logger
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
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalln("open log file failure:", err)
		}
		defer logFile.Close()

		_logWriter = log.New(logFile, "", log.Ldate|log.Ltime)
		// _infoLogger = log.New(logFile, "[info]", log.Ldate|log.Ltime)
		// _warningLogger = log.New(logFile, "[warning]", log.Ldate|log.Ltime)
		// _errorLogger = log.New(logFile, "[error]", log.Ldate|log.Ltime|log.Llongfile)

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
	err1 := errors.WithStack(err)
	fmt.Printf("[error]%+v", err1)
	_logWriter.Printf("[error]%+v\n--------------------------------------------------error end--------------------------------------------------", err1)
}

func LogErrorWithRemark(err error, remark string) {
	createLogFileIfNotExist()
	err1 := errors.WithStack(err)
	fmt.Printf("[error]%+v\n[error remakr]"+remark, err1)
	_logWriter.Printf("[error]%+v"+"\n[error remark]"+remark+"\n--------------------------------------------------error end--------------------------------------------------", err1)
}

func LogWarning(msg string) {
	createLogFileIfNotExist()
	fmt.Println("[warning]" + msg)
	_logWriter.Println("[warning]" + msg)
}

func LogInfo(msg string) {
	createLogFileIfNotExist()
	fmt.Println("[info]" + msg)
	_logWriter.Println("[info]" + msg)
}

func LogToConsoleWithTime(msg string, paras ...interface{}) {
	fmt.Printf("%s %s\n", time.Now().UTC().Add(8*time.Hour).Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, paras...))
}
