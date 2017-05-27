package WatchDog

import (
	"sync/atomic"
	"time"

	"github.com/6xiao/go/Common"
)

type WatchDog struct {
	wait  time.Duration
	hung  func()
	meat  int32
	count int
}

func NewDog(duration time.Duration, meat int32, hung func()) *WatchDog {
	d := new(WatchDog)
	d.wait = duration
	d.hung = hung
	d.meat = meat

	go d.eat()
	return d
}

func (this *WatchDog) eat() {
	defer Common.CheckPanic()

	for this.hung != nil && this.count < 3 {
		time.Sleep(this.wait)

		m := atomic.LoadInt32(&this.meat)
		if m < 0 {
			return
		}

		if m == 0 {
			this.hung()
			this.count++
		} else {
			atomic.StoreInt32(&this.meat, m/2)
			this.count = 0
		}
	}
}

func (this *WatchDog) Feed(meat uint16) bool {
	defer Common.CheckPanic()
	return atomic.AddInt32(&this.meat, int32(meat)) > 0
}

func (this *WatchDog) Kill() {
	defer Common.CheckPanic()
	atomic.StoreInt32(&this.meat, -65536)
}
