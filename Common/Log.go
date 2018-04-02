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
	logTime   = true
	logLevel  = LogLevelInfo
	logFile   = os.Stderr
	logLock   = NewLock()
	logTicker = time.NewTicker(time.Second)
	logSlice  = make([]interface{}, 1, 1024)
)

func SetLogTime(logtime bool) {
	logTime = logtime
}

func SetLogDir(dir string) {
	logDir = dir
	newfile(time.Now())
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

	if logLock.TryLock() {
		defer logLock.Unlock()
	} else {
		return
	}

	if now := time.Now().UTC(); logDay != now.Day() {
		newfile(now)
	}
}

func newfile(now time.Time) {
	if logFile != os.Stderr {
		logFile.Close()
		logFile = os.Stderr
	}

	if len(logDir) == 0 {
		return
	}

	logDay = now.Day()
	logProc := filepath.Base(os.Args[0])
	filename := filepath.Join(logDir, fmt.Sprintf("%s.%s.log", logProc, now.Format("2006-01-02")))

	newlog, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, NumberUTC(), "open log file", filename, err, "use STDOUT")
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

func writeLog(level int, pre string, v ...interface{}) {
	if logLevel > level {
		return
	}

	check()
	logLock.Lock()
	defer logLock.Unlock()

	_, file, line, _ := runtime.Caller(2)
	if logTime {
		logSlice[0] = fmt.Sprintf("%d %s %s[%d]:", NumberUTC(), pre, filebase(file), line)
	} else {
		logSlice[0] = fmt.Sprintf("- %s %s[%d]:", pre, filebase(file), line)
	}
	fmt.Fprintln(logFile, append(logSlice, v...)...)
}

func DropLog(v ...interface{}) {}

func DebugLog(v ...interface{}) {
	writeLog(LogLevelDebug, "debug", v...)
}

func InfoLog(v ...interface{}) {
	writeLog(LogLevelInfo, "info", v...)
}

func WarningLog(v ...interface{}) {
	writeLog(LogLevelWarn, "warn", v...)
}

func ErrorLog(v ...interface{}) {
	writeLog(LogLevelError, "error", v...)
}
