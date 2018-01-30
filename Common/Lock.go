package Common

import (
	"sync/atomic"
	"time"
)

const (
	unlock = 1
	locked = 0
	waited = -1
)

type Lock struct {
	flag       int32
	unlockChan chan bool
}

func NewLock() *Lock {
	l := new(Lock)
	l.flag = unlock
	l.unlockChan = make(chan bool, 3)
	return l
}

func (this *Lock) Lock() {
	if atomic.AddInt32(&this.flag, waited) == locked {
		return
	}

	for {
		if atomic.LoadInt32(&this.flag) >= locked && atomic.SwapInt32(&this.flag, waited) == unlock {
			return
		}
		<-this.unlockChan
	}
}

func (this *Lock) SpinLock(count int) bool {
	for i := 0; i < count; i++ {
		if atomic.CompareAndSwapInt32(&this.flag, unlock, locked) {
			return true
		}
	}
	return atomic.CompareAndSwapInt32(&this.flag, unlock, locked)
}

func (this *Lock) TryLock() bool {
	return this.SpinLock(0)
}

func (this *Lock) TimeoutSpinLock(t time.Duration) bool {
	tk := time.NewTicker(t)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			return this.SpinLock(0)

		default:
			if this.SpinLock(8192) {
				return true
			}
		}
	}
}

func (this *Lock) Unlock() {
	if atomic.SwapInt32(&this.flag, unlock) == locked {
		return
	}

	select {
	case this.unlockChan <- true:
	default:
	}
}
