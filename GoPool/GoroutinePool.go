package GoPool

import (
	"sync"
	"sync/atomic"

	"github.com/6xiao/go/Common"
)

/*
eg:
type Calc struct {
	in chan int
}

func (this *Calc) CalcSum() interface{} {
	sum := 0
	for i := range this.in {
		sum += i
	}
	return sum
}

func AtQuit(res interface{}) {
	fmt.Println("get sum", res)
}

func main() {
	intChan := make(chan int)
	calc := Calc{intChan}
	pool := GoPool.NewPool(3, calc.CalcSum, AtQuit)
	pool.Check()
	fmt.Println("worker number:", pool.Count())

	for i := 0; i < 1024*1024; i++ {
		intChan <- i
	}

	pool.Check()
	fmt.Println("worker number:", pool.Count())
	close(intChan)
	pool.Wait()
	fmt.Println("worker number:", pool.Count())
}
*/

type Pool struct {
	sync.Mutex
	sync.WaitGroup

	count int64
	capa  int64
	run   func() interface{}
	quit  func(interface{})
}

func NewPool(capa int, run func() interface{}, quit func(interface{})) *Pool {
	if capa <= 0 || run == nil {
		return nil
	}
	return &Pool{sync.Mutex{}, sync.WaitGroup{}, 0, int64(capa), run, quit}
}

func (this *Pool) Count() int {
	return int(atomic.LoadInt64(&this.count))
}

func (this *Pool) Check() {
	this.Lock()
	defer this.Unlock()

	if this.capa == this.count {
		return
	}

	wg := sync.WaitGroup{}
	for i := this.capa - this.count; i > 0; i-- {
		this.Add(1)
		wg.Add(1)
		go this.Run(&wg)
	}
	wg.Wait()
}

func (this *Pool) Run(wg *sync.WaitGroup) {
	defer Common.CheckPanic()
	atomic.AddInt64(&this.count, 1)
	defer atomic.AddInt64(&this.count, -1)
	defer this.Done()
	wg.Done()

	result := this.run()
	if this.quit != nil {
		this.quit(result)
	}
}
