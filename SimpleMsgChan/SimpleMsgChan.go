package main

import (
	"container/list"
	"flag"
	"github.com/6xiao/go/Common"
	"io"
	"log"
	"net"
	"sync"
	"time"
	"unsafe"
)

var (
	flgAddr = flag.String("addr", "127.0.0.1:5381", "message channel server addr")
)

const (
	HeadLen    = 8
	MsgHeadLen = 40
)

type MessageHead struct {
	msglen   uint64
	sender   uint64
	recver   uint64
	sendTime uint64
	keepTime uint64
}

//connect ...
type Connect struct {
	socket     net.Conn
	dataChan   chan []byte
	remoteAddr string
}

func NewConnect(socket net.Conn) *Connect {
	return &Connect{socket, make(chan []byte), ""}
}

func (c *Connect) Addr() string {
	return c.remoteAddr
}

func (c *Connect) Close() {
	close(c.dataChan)
}

func (c *Connect) Send(msg []byte) {
	if msg != nil {
		c.dataChan <- msg
	}
}

func (c *Connect) Sender() {
	defer Common.CheckPanic()
	defer c.socket.Close()

	for buf := range c.dataChan {
		if _, err := c.socket.Write(buf); err != nil {
			log.Println(c.Addr(), "Socket Send Error ", err)
		}
	}
}

func (c *Connect) Recver(in *Intranet) {
	defer Common.CheckPanic()
	defer c.Close()
	defer in.DelCon(c)

	c.remoteAddr = c.socket.RemoteAddr().String()

	for {
		nethead := [HeadLen]byte{}
		if _, e := io.ReadFull(c.socket, nethead[:]); e != nil {
			break
		}

		msglen := *(*uint64)(unsafe.Pointer(&nethead))
		if msglen < MsgHeadLen {
			log.Println(c.Addr(), "msg head too short")
			break
		}

		msgbuf := make([]byte, HeadLen+msglen)

		if _, e := io.ReadFull(c.socket, msgbuf[HeadLen:]); e != nil {
			log.Println(c.Addr(), "msg body too short", e)
			break
		}

		hd := (*MessageHead)(unsafe.Pointer(&msgbuf[0]))
		hd.msglen = msglen

		if hd.sender != 0 {
			in.AddCon(hd.sender, c)
		}

		if hd.recver != 0 {
			in.Publish(hd.recver, msgbuf, hd.sendTime == hd.keepTime)
		}
	}
}

type SafeList struct {
	lock sync.Mutex
	lst  *list.List
}

func NewSafeList() *SafeList {
	return &SafeList{sync.Mutex{}, list.New()}
}

func (sl *SafeList) Push(msg []byte) {
	sl.lock.Lock()
	sl.lst.PushBack(msg)
	sl.lock.Unlock()
}

func (sl *SafeList) Pop() []byte {
	sl.lock.Lock()
	f := sl.lst.Front()
	if f == nil {
		sl.lock.Unlock()
		return nil
	}

	sl.lst.Remove(f)
	sl.lock.Unlock()

	if msg, ok := f.Value.([]byte); ok {
		hd := (*MessageHead)(unsafe.Pointer(&msg[0]))
		if hd.keepTime == 0 || hd.keepTime > Common.NumberTime(time.Now()) {
			return msg
		}
	}
	return nil
}

func (sl *SafeList) Len() int {
	return sl.lst.Len()
}

func (sl *SafeList) OutTimeClear() {
	defer Common.CheckPanic()

	sl.lock.Lock()
	defer sl.lock.Unlock()

	tm := Common.NumberTime(time.Now())
	for e := sl.lst.Front(); e != nil; e = e.Next() {
		if msg, ok := e.Value.([]byte); ok {
			hd := (*MessageHead)(unsafe.Pointer(&msg[0]))
			if hd.keepTime != 0 && hd.keepTime < tm {
				e = e.Prev()
				sl.lst.Remove(e.Next())
			}
		}
	}
}

type Observer struct {
	name uint64
	set  map[string]*Connect
	list *SafeList
}

func NewObserver(name uint64) *Observer {
	return &Observer{name, make(map[string]*Connect), NewSafeList()}
}

func (o *Observer) SendBuffer() {
	defer Common.CheckPanic()

	for len(o.set) > 0 && o.list.Len() > 0 {
		if msg := o.list.Pop(); msg != nil {
			for _, v := range o.set {
				v.Send(msg)
			}
		}
	}
}

func (o *Observer) Add(c *Connect) {
	empty := (len(o.set) == 0)

	if _, ok := o.set[c.Addr()]; !ok {
		o.set[c.Addr()] = c
	}

	if empty {
		go o.SendBuffer()
	}
}

func (o *Observer) Publish(msg []byte, realtime bool) {
	if realtime || len(o.set) > 0 {
		for _, v := range o.set {
			v.Send(msg)
		}
	} else {
		o.list.Push(msg)
		l := o.list.Len()
		if (l & (l - 1)) == 0 {
			go o.list.OutTimeClear()
		}
	}
}

func (o *Observer) Delete(c *Connect) {
	delete(o.set, c.Addr())
}

type Intranet struct {
	tunnel map[uint64]*Observer
}

func NewIntranet() *Intranet {
	return &Intranet{make(map[uint64]*Observer)}
}

func (ni *Intranet) DelCon(c *Connect) {
	for _, obs := range ni.tunnel {
		obs.Delete(c)
	}
}

func (ni *Intranet) AddCon(name uint64, c *Connect) {
	if o, ok := ni.tunnel[name]; ok {
		o.Add(c)
	} else {
		no := NewObserver(name)
		no.Add(c)
		ni.tunnel[name] = no
	}
}

func (ni *Intranet) Publish(recver uint64, msg []byte, realtime bool) {
	if o, ok := ni.tunnel[recver]; ok {
		o.Publish(msg, realtime)
	} else {
		ni.tunnel[recver] = NewObserver(recver)
		ni.tunnel[recver].Publish(msg, realtime)
	}
}

func (ni *Intranet) Reactiver(c net.Conn) {
	nc := NewConnect(c)
	go nc.Sender()
	go nc.Recver(ni)
}

func main() {
	Common.Init(nil)

	ni := NewIntranet()
	err := Common.ListenSocket(*flgAddr, true, ni.Reactiver, nil)
	log.Println(err)
}
