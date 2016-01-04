package Common

import (
	"compress/gzip"
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
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
				fmt.Fprintln(os.Stderr, NumberNow(), fn, fileline(file, line))
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

// hash : []byte to uint64
func Hash(mem []byte) uint64 {
	var hash uint64 = 5381
	for _, b := range mem {
		hash = (hash << 5) + hash + uint64(b)
	}
	return hash
}

// create a uuid string
func NewUUID() string {
	u := [16]byte{}
	rand.Read(u[:])
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func Ternary(cond bool, valTrue, valFlase interface{}) interface{} {
	if cond {
		return valTrue
	}
	return valFlase
}

type Int64Slice []int64
func (this Int64Slice) Len() int           { return len(this) }
func (this Int64Slice) Less(i, j int) bool { return this[i] < this[j] }
func (this Int64Slice) Swap(i, j int)      { this[i], this[j] = this[j], this[i] }

// compress data use gzip
func Gzip(in []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	if wt, err := gzip.NewWriterLevel(buf, gzip.BestCompression); err != nil {
		return nil, err
	} else {
		if _, err := wt.Write(in); err != nil {
			return nil, err
		}
		wt.Close()
	}
	return buf.Bytes(), nil
}