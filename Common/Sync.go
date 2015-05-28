package Common

import (
	"sync/atomic"
)

type Counter int64

func (c *Counter) Count() int64 {
	return atomic.LoadInt64((*int64)(c))
}

func (c *Counter) Enter() {
	atomic.AddInt64((*int64)(c), 1)
}

func (c *Counter) Leave() {
	atomic.AddInt64((*int64)(c), -1)
}

func (c *Counter) Put() {
	atomic.AddInt64((*int64)(c), 1)
}

func (c *Counter) Get() {
	atomic.AddInt64((*int64)(c), -1)
}
