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

func (this *Calc) Run() (interface{}, error) {
	sum := 0
	for i := range this.in {
		sum += i
	}
	return sum, nil
}

func (this *Calc) Quit(_ *GoPool.Pool, res interface{}, _ error) {
	fmt.Println("get sum", res)
}

func AtQuit(_ *GoPool.Pool, res interface{}, _ error) {
	fmt.Println("get sum", res)
}

func main() {
	intChan := make(chan int)
	calc := Calc{intChan}
	pool := GoPool.FromRunner(3, &calc)
	//pool := GoPool.NewPool(3, calc.Run, AtQuit)
	pool.Keepalive()
	fmt.Println("worker number:", pool.Count())

	for i := 0; i < 1024*1024; i++ {
		intChan <- i
	}
	pool.AddExecutor()
	fmt.Println("worker number:", pool.Count())

	close(intChan)

	pool.WaitAllQuit()
	fmt.Println("worker number:", pool.Count())
}
*/

type Runner interface {
	Run() (interface{}, error)
	Quit(*Pool, interface{}, error)
}

// goroutine pool
type Pool struct {
	sync.Mutex
	sync.WaitGroup

	count int64
	capa  int64
	run   func() (interface{}, error)
	quit  func(*Pool, interface{}, error)
}

func NewPool(capa int,
	run func() (interface{}, error),
	quit func(*Pool, interface{}, error)) *Pool {
	if capa <= 0 || run == nil {
		return nil
	}
	return &Pool{sync.Mutex{}, sync.WaitGroup{}, 0, int64(capa), run, quit}
}

func FromRunner(capa int, runner Runner) *Pool {
	return NewPool(capa, runner.Run, runner.Quit)
}

// count may great then capa when user AddExecutor
func (this *Pool) Count() int {
	return int(atomic.LoadInt64(&this.count))
}

func (this *Pool) Capacity() int {
	return int(this.capa)
}

// return after all goroutine quit in pool
func (this *Pool) WaitAllQuit() {
	this.Wait()
}

// check alive
func (this *Pool) Keepalive() int {
	this.Lock()
	defer this.Unlock()

	sub := int(this.capa - atomic.LoadInt64(&this.count))
	for i := 0; i < sub; i++ {
		this.AddExecutor()
	}

	return sub
}

// add goroutine manual
func (this *Pool) AddExecutor() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go this.executor(&wg)
	wg.Wait()
}

func (this *Pool) executor(wg *sync.WaitGroup) {
	defer Common.CheckPanic()
	this.Add(1)
	defer this.Done()
	atomic.AddInt64(&this.count, 1)
	defer atomic.AddInt64(&this.count, -1)
	wg.Done()

	result, err := this.run()
	if this.quit != nil {
		this.quit(this, result, err)
	}
}
