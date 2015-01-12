package Common

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

func init() {
	Init()
}

func Init() {
	flag.Parse()

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(os.Stderr, "#%s:%s:%s\n", f.Name, f.DefValue, f.Usage)
	})

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(os.Stderr, "-%s=%v\n", f.Name, f.Value)
	})

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func CheckPanic() {
	if err := recover(); err != nil {
		log.Println("panic : ", err)

		for skip := 1; ; skip++ {
			pc, file, line, ok := runtime.Caller(skip)
			if !ok {
				break
			}

			log.Println(skip, pc, file, line)
		}
	}
}
