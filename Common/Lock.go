package Common

import (
	"runtime"
	"sync/atomic"
	"time"
)

const (
	unlock = 0
	locked = 1
)

type Lock struct {
	flag       int32
	unlockChan chan bool
}

func NewLock() *Lock {
	l := new(Lock)
	l.unlockChan = make(chan bool, 3)
	return l
}

func (this *Lock) Lock() {
	for {
		if atomic.CompareAndSwapInt32(&this.flag, unlock, locked) {
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
	return false
}

func (this *Lock) TryLock() bool {
	return this.SpinLock(1)
}

func (this *Lock) TimeoutSpinLock(t time.Duration) bool {
	tk := time.NewTicker(t)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			return this.SpinLock(1)

		default:
			if this.SpinLock(1024) {
				return true
			}
		}
	}
}

func (this *Lock) SchedLock() {
	for {
		if this.SpinLock(1024) {
			return
		}
		runtime.Gosched()
	}
}

func (this *Lock) TimeoutSchedLock(t time.Duration) bool {
	tk := time.NewTicker(t)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			return this.SpinLock(1)

		default:
			if this.SpinLock(1024) {
				return true
			}
		}
		runtime.Gosched()
	}
}

func (this *Lock) Unlock() {
	atomic.SwapInt32(&this.flag, unlock)
	select {
	case this.unlockChan <- true:
	default:
	}
}

type Semaphore struct {
	flag chan bool
}

func NewSemaphore(capa int) *Semaphore {
	return &Semaphore{make(chan bool, capa)}
}

func (this *Semaphore) Alloc() {
	this.flag <- true
}

func (this *Semaphore) TryAlloc() bool {
	select {
	case this.flag <- true:
		return true
	default:
	}
	return false
}

func (this *Semaphore) Free() bool {
	select {
	case <-this.flag:
		return true
	default:
	}
	return false
}
