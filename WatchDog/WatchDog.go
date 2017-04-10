package WatchDog

import (
	"sync"
	"sync/atomic"

	"github.com/6xiao/go/Common"
	"time"
)

type WatchDog struct {
	sync.Mutex

	wait time.Duration
	hung func()
	meat *int64
}

func NewDog(duration time.Duration, meat int64, hung func()) *WatchDog {
	d := WatchDog{sync.Mutex{}, duration, hung, new(int64)}
	d.Feed(meat)
	go d.eat()
	return &d
}

func (this *WatchDog) eat() {
	defer Common.CheckPanic()

	for {
		time.Sleep(this.wait)

		if m := atomic.LoadInt64(this.meat); m == 0 && this.hung != nil {
			this.hung()
		} else {
			atomic.StoreInt64(this.meat, m/2)
		}
	}
}

func (this *WatchDog) Feed(meat int64) bool {
	defer Common.CheckPanic()
	return atomic.AddInt64(this.meat, meat) > 0
}
