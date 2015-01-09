package Common

import (
	"flag"
	"log"
	"runtime"
)

func init() {
	flag.Parse()
	flag.Usage()
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func CheckPanic() {
	if err := recover(); err != nil {
		log.Println("panic : ", err)
		for skip := 1; ; skip++ {
			if pc, file, line, ok := runtime.Caller(skip); ok {
				log.Println("caller : ", skip, pc, file, line)
			} else {
				break
			}
		}
	}
}
