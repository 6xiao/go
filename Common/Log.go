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

func fileline(file string, line int) string {
	beg, end := len(file)-1, len(file)
	for ; beg >= 0; beg-- {
		if os.IsPathSeparator(file[beg]) {
			beg++
			break
		} else if file[beg] == '.' {
			end = beg
		}
	}
	return fmt.Sprint(file[beg:end], ":", line)
}

func offset() string {
	_, file, line, _ := runtime.Caller(2)
	return fileline(file, line)
}

func DropLog(v ...interface{}) {}

func DebugLog(v ...interface{}) {
	if logLevel > LogLevelDebug {
		return
	}

	check()

	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, NumberUTC(), offset(), "debug", v)
}

func InfoLog(v ...interface{}) {
	if logLevel > LogLevelInfo {
		return
	}

	check()

	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, NumberUTC(), offset(), "info", v)
}

func WarningLog(v ...interface{}) {
	if logLevel > LogLevelWarn {
		return
	}

	check()

	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, NumberUTC(), offset(), "warning", v)
}

func ErrorLog(v ...interface{}) {
	if logLevel > LogLevelError {
		return
	}

	check()

	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, NumberUTC(), offset(), "error", v)
}

func CustomLog(v ...interface{}) {
	check()
	logLock.Lock()
	defer logLock.Unlock()
	fmt.Fprintln(logFile, NumberUTC(), v)
}
