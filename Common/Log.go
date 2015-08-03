package Common

// 每日滚动的LOG实现

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logDir   = os.TempDir()
	logDay   = 0
	logDebug = false
	logInfo  = true
	logError = true
	logLock  = sync.Mutex{}
	logFile  *os.File
	logInst  *log.Logger
)

// parse env ENABLE_DEBUG_LOG and DISABLE_INFO_LOG
func init() {
	if proc, err := filepath.Abs(os.Args[0]); err == nil {
		SetLogDir(filepath.Dir(proc))
	}

	if len(os.Getenv("ENABLE_DEBUG_LOG")) > 0 {
		logDebug = true
	}

	if len(os.Getenv("DISABLE_INFO_LOG")) > 0 {
		logInfo = false
	}
}

func SetLogDir(dir string) {
	logDir = dir
}

func EnableDebugLog(dbg bool) {
	logDebug = dbg
}

func EnableInfoLog(info bool) {
	logInfo = info
}

func EnableErrorLog(err bool) {
	logError = err
}

func check() {
	logLock.Lock()
	defer logLock.Unlock()

	if logDay == time.Now().Day() {
		return
	}

	logDay = time.Now().Day()
	logFile.Sync()
	logFile.Close()
	logProc := filepath.Base(os.Args[0])
	filename := filepath.Join(logDir,
		fmt.Sprintf("%s.%s.log", logProc, time.Now().Format("2006-01-02")))
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logInst = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
		logInst.Output(2, fmt.Sprintln("error", "check log file", err, "use STDOUT"))
		return
	}

	logInst = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func DebugLog(v ...interface{}) {
	if logDebug {
		check()
		logLock.Lock()
		defer logLock.Unlock()
		logInst.Output(2, fmt.Sprintln("debug", v))
	}
}

func InfoLog(v ...interface{}) {
	if logInfo {
		check()
		logLock.Lock()
		defer logLock.Unlock()
		logInst.Output(2, fmt.Sprintln("info", v))
	}
}

func ErrorLog(v ...interface{}) {
	if logError {
		check()
		logLock.Lock()
		defer logLock.Unlock()
		logInst.Output(2, fmt.Sprintln("error", v))
	}
}
