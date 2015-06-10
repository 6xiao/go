package Common

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"
)

// zero size, empty struct
type EmptyStruct struct{}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("stderr ")
}

// parse flag and print usage/value to writer
func Init(writer io.Writer) {
	flag.Parse()

	if writer == nil {
		return
	}

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(writer, "#%s : %s\n", f.Name, f.Usage)
	})

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(writer, "-%s=%v \n", f.Name, f.Value)
	})
}

// check panic when exit
func CheckPanic() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "\n%v %v\n", FormatNow(), err)

		for skip := 1; ; skip++ {
			if pc, file, line, ok := runtime.Caller(skip); ok {
				fn := runtime.FuncForPC(pc).Name()
				fmt.Fprintf(os.Stderr, "%v %v %v:%v\n",
					FormatNow(), fn, path.Base(file), line)
			} else {
				break
			}
		}
	}
}

// recive quit signal
func QuitSignal() <-chan os.Signal {
	signals := make(chan os.Signal, 3)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	return signals
}

// hash : []byte to uint64
func Hash(mem []byte) uint64 {
	var hash uint64 = 5381

	for _, b := range mem {
		hash = (hash << 5) + hash + uint64(b)
	}

	return hash
}

// format a time.Time to string as 2006-01-02 15:04:05.999
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.999")
}

// format time.Now() use FormatTime
func FormatNow() string {
	return FormatTime(time.Now())
}

// parse a string as "2006-01-02 15:04:05.999" to time.Time
func ParseTime(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05.999", s, time.Local)
}

// format a time.Time to number as 20060102150405999
func NumberTime(t time.Time) uint64 {
	return uint64(t.UnixNano()/1000000)%1000 +
		uint64(t.Second())*1000 +
		uint64(t.Minute())*100000 +
		uint64(t.Hour())*10000000 +
		uint64(t.Day())*1000000000 +
		uint64(t.Month())*100000000000 +
		uint64(t.Year())*10000000000000
}

// create a uuid string
func NewUUID() string {
	u := [16]byte{}
	rand.Read(u[:])
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func Ternary(cond bool, valTrue, valFlase interface{}) interface{} {
	if cond {
		return valTrue
	}

	return valFlase
}
