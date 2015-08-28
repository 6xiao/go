package Common

// 每日滚动的LOG实现

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	logDir  = os.TempDir()
	logDay  = 0
	logLock = sync.Mutex{}
	logFile *os.File
)

// parse env ENABLE_DEBUG_LOG and DISABLE_INFO_LOG
func init() {
	if proc, err := filepath.Abs(os.Args[0]); err == nil {
		SetLogDir(filepath.Dir(proc))
	}
}

func SetLogDir(dir string) {
	logDir = dir
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
	var err error
	logFile, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logFile = os.Stderr
		fmt.Fprintln(os.Stderr, NumberNow(), "check log file", err, "use STDOUT")
	}
}

func DropLog(v ...interface{}) {}

func DebugLog(v ...interface{}) {
	check()
	logLock.Lock()
	defer logLock.Unlock()
	_, file, line, _ := runtime.Caller(1)
	fmt.Fprintln(logFile, NumberNow(), filepath.Base(file), line, "debug", v)
}

func InfoLog(v ...interface{}) {
	check()
	logLock.Lock()
	defer logLock.Unlock()
	_, file, line, _ := runtime.Caller(1)
	fmt.Fprintln(logFile, NumberNow(), filepath.Base(file), line, "info", v)
}

func ErrorLog(v ...interface{}) {
	check()
	logLock.Lock()
	defer logLock.Unlock()
	_, file, line, _ := runtime.Caller(1)
	fmt.Fprintln(logFile, NumberNow(), filepath.Base(file), line, "error", v)
}
