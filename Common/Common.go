package Common

import (
	"flag"
	"fmt"
	"io"
	"log"
	"path"
	"runtime"
	"time"
)

// zero size, empty struct
type Empty struct{}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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
		fmt.Fprintf(writer, "-%s=%v \\\n", f.Name, f.Value)
	})
}

// check panic when exit
func CheckPanic() {
	if err := recover(); err != nil {
		log.Println("panic : ", err)

		for skip := 1; ; skip++ {
			if _, file, line, ok := runtime.Caller(skip); ok {
				log.Println(skip, path.Base(file), line)
			} else {
				break
			}
		}
	}
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
