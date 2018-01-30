package Common

// 每日滚动的LOG实现
import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	LogLevelDrop = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelNone
)

var (
	logDir    = ""
	logDay    = 0
	logLevel  = LogLevelInfo
	logFile   = os.Stderr
	logLock   = NewLock()
	logTicker = time.NewTicker(time.Second)
	logSlice  = make([]interface{}, 1, 1024)
)

func SetLogDir(dir string) {
	logDir = dir
}

func SetLogLevel(level int) {
	logLevel = level
}

func check() {
	select {
	case <-logTicker.C:
	default:
		return
	}

	if len(logDir) == 0 {
		return
	}

	if logLock.TryLock() {
		defer logLock.Unlock()
	} else {
		return
	}

	now := time.Now().UTC()
	if logDay == now.Day() {
		return
	}

	if logFile != os.Stderr {
		logFile.Close()
	}

	logDay = now.Day()
	logProc := filepath.Base(os.Args[0])
	filename := filepath.Join(logDir, fmt.Sprintf("%s.%s.log", logProc, now.Format("2006-01-02")))

	newlog, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, NumberUTC(), "open log file", err, "use STDOUT")
		logFile = os.Stderr
	} else {
		logFile = newlog
	}
}

func filebase(file string) string {
	beg, end := len(file)-1, len(file)
	for ; beg >= 0; beg-- {
		if os.IsPathSeparator(file[beg]) {
			beg++
			break
		} else if file[beg] == '.' {
			end = beg
		}
	}
	return file[beg:end]
}

func DropLog(v ...interface{}) {}

func DebugLog(v ...interface{}) {
	if logLevel > LogLevelDebug {
		return
	}

	check()

	_, file, line, _ := runtime.Caller(1)
	logSlice[0] = fmt.Sprintf("%d debug %s[%d]:", NumberUTC(), filebase(file), line)
	out := append(logSlice, v...)
	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, out...)
}

func InfoLog(v ...interface{}) {
	if logLevel > LogLevelInfo {
		return
	}

	check()

	_, file, line, _ := runtime.Caller(1)
	logSlice[0] = fmt.Sprintf("%d info %s[%d]:", NumberUTC(), filebase(file), line)
	out := append(logSlice, v...)
	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, out...)
}

func WarningLog(v ...interface{}) {
	if logLevel > LogLevelWarn {
		return
	}

	check()

	_, file, line, _ := runtime.Caller(1)
	logSlice[0] = fmt.Sprintf("%d warn %s[%d]:", NumberUTC(), filebase(file), line)
	out := append(logSlice, v...)
	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, out...)
}

func ErrorLog(v ...interface{}) {
	if logLevel > LogLevelError {
		return
	}

	check()

	_, file, line, _ := runtime.Caller(1)
	logSlice[0] = fmt.Sprintf("%d error %s[%d]:", NumberUTC(), filebase(file), line)
	out := append(logSlice, v...)
	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, out...)
}
