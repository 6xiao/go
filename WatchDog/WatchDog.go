package WatchDog

import (
	"math"
	"sync/atomic"
	"time"

	"github.com/6xiao/go/Common"
)

type WatchDog struct {
	wait time.Duration
	hung func()
	meat int64
}

func NewDog(duration time.Duration, hung func()) *WatchDog {
	d := new(WatchDog)
	d.wait = duration
	d.hung = hung
	d.meat = math.MaxInt32

	go d.eat()
	return d
}

func (this *WatchDog) eat() {
	defer Common.CheckPanic()

	for this.hung != nil {
		time.Sleep(this.wait)

		m := atomic.LoadInt64(&this.meat)
		if m < 0 {
			return
		}

		if m == 0 {
			this.hung()
		} else {
			atomic.StoreInt64(&this.meat, m/2)
		}
	}
}

func (this *WatchDog) Feed(meat int64) bool {
	defer Common.CheckPanic()

	if meat > math.MaxInt32 {
		return false
	}
	return atomic.AddInt64(&this.meat, meat) > 0
}

func (this *WatchDog) Kill() {
	defer Common.CheckPanic()
	atomic.StoreInt64(&this.meat, math.MinInt32)
}
