package Common

import (
	"log"
	"net"
	"net/rpc"
)

func ListenRpc(addr string, obj interface{}, atexit func()) {
	defer CheckPanic()
	if atexit != nil {
		defer atexit()
	}

	if obj == nil {
		log.Println("rpc object is nil")
		return
	}

	server := rpc.NewServer()
	server.Register(obj)

	if listener, err := net.Listen("tcp", addr); err != nil {
		log.Println("rpc listen error", addr, err)
	} else {
		log.Println("rpc running @", addr)
		server.Accept(listener)
	}
}

func ListenSocket(addr string, keepalive bool, reactiver func(net.Conn), atexit func()) {
	defer CheckPanic()
	if atexit != nil {
		defer atexit()
	}

	if reactiver == nil {
		log.Println("socket reactiver is nil")
		return
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Println("can't resolve addr", addr, err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println("can't listen tcp ", addr, err)
		return
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("accept ", err)
			continue
		}

		conn.SetNoDelay(true)
		conn.SetKeepAlive(keepalive)
		conn.SetLinger(-1)
		go reactiver(conn)
	}
}
