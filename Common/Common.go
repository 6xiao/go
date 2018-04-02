package Common

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

// zero size, empty struct
type EmptyStruct struct{}

// parse flag and print usage/value to writer
func Init(writer io.Writer) {
	flag.Parse()

	if writer != nil {
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(writer, "-%s=%v \n", f.Name, f.Value)
		})
	}
}

// check panic when exit
func CheckPanic() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "\n%v %v\n", NumberNow(), err)

		for skip := 1; ; skip++ {
			if pc, file, line, ok := runtime.Caller(skip); ok {
				fn := runtime.FuncForPC(pc).Name()
				fmt.Fprintln(os.Stderr, NumberNow(), fn, filebase(file), line)
			} else {
				break
			}
		}
	}
}

// reload signal
func HupSignal() <-chan os.Signal {
	signals := make(chan os.Signal, 3)
	signal.Notify(signals, syscall.SIGHUP)
	return signals
}

// recive quit signal
func QuitSignal() <-chan os.Signal {
	signals := make(chan os.Signal, 3)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	return signals
}

// create a uuid string
func NewUUID() string {
	u := [16]byte{}
	rand.Read(u[:])
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func RunShell(exeStr string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", exeStr)
	bytes, err := cmd.CombinedOutput()
	return string(bytes), err
}

func RuntimeInfo(w io.Writer) {
	fmt.Fprintln(w, "runtime @", NumberNow())
	fmt.Fprintln(w, "NrCpu :", runtime.NumCPU())
	fmt.Fprintln(w, "NrGoroutine :", runtime.NumGoroutine())
	fmt.Fprintln(w, "NrCgo :", runtime.NumCgoCall())

	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)
	fmt.Fprintln(w, "NrMemHeapAlloc :", mem.HeapAlloc)
	fmt.Fprintln(w, "NrMemHeapObjects :", mem.HeapObjects)

	fmt.Fprintln(w, "NrMemHeapSys :", mem.HeapSys)
	fmt.Fprintln(w, "NrMemHeapInuse :", mem.HeapInuse)

	fmt.Fprintln(w, "NrMemStackSys :", mem.StackSys)
	fmt.Fprintln(w, "NrMemStackInuse :", mem.StackInuse)

	fmt.Fprintln(w, "NrMemMSpanSys :", mem.MSpanSys)
	fmt.Fprintln(w, "NrMemMSpanInuse :", mem.MSpanInuse)

	fmt.Fprintln(w, "NrMemMCacheSys :", mem.MCacheSys)
	fmt.Fprintln(w, "NrMemMCacheInuse :", mem.MCacheInuse)

	ps := pprof.Profiles()
	for _, ps := range ps {
		fmt.Fprintf(w, "NrProfile-%s : %d\n", ps.Name(), ps.Count())
	}
}
