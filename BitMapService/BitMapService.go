package BitMapService

import (
	"errors"
	"github.com/6xiao/go/Common"
	"log"
	"net"
	"net/rpc"
)

const (
	OPER_GET   = "GET"
	OPER_SET   = "SET"
	OPER_SYNC  = "SYNC"
	OPER_SHIFT = "SHIFT"
	OPER_INFO  = "INFOMATION"
	BILLION    = 1024 * 1024 * 1024
	RPC_CALL   = "BitMapRPC.RequestCache"
	RPC_SYNC   = "BitMapRPC.ReverseSync"
)

type Cache interface {
	//how many bits of the uid's flag, such as 1,2,4,8,16,32,64
	Bits() int

	Capacity() int64

	//how much uid in cache
	Count() uint64

	//get the max uid
	MaxUid() int64

	// all uid's flags <<1
	Shift()

	//set uid's last bit, if is the first time return true
	SetUidLastFlag(uid int64) bool

	//get uid's flags
	GetUidFlags(uid int64) uint64

	//set uid's all flag in once time
	SetUidFlags(uid int64, value uint64)
}

type SliceKeyValue struct {
	Uids  []int64
	Flags []uint64
}

type BitMapRequest struct {
	Bits int

	// cache name
	Name string

	//OPER_GET / OPER_SET / OPER_SHIFT / OPER_SYNC / OPER_INFO
	Operator string

	SliceKeyValue
}

type SyncInvoke func(in string, out *bool) error
type ReqInvoke func(in BitMapRequest, out *SliceKeyValue) error

type BitMapRPC struct {
	Address  string
	SyncFunc SyncInvoke
	ReqFunc  ReqInvoke
}

func NewBitMapRpc(addr string, sync SyncInvoke, req ReqInvoke) *BitMapRPC {
	return &BitMapRPC{addr, sync, req}
}

//1, client give a rpc address(in) to server
//2, server connect to client
//3, server push all datas to client use the RequestCache interface
//in : address
//out : success
func (s *BitMapRPC) ReverseSync(in string, out *bool) error {
	if s.SyncFunc != nil {
		return s.SyncFunc(in, out)
	}

	return errors.New("can't ReverseSync")
}

//out: first inserted uid when Operator == SET
func (s *BitMapRPC) RequestCache(in BitMapRequest, out *SliceKeyValue) error {
	if s.ReqFunc != nil {
		*out = SliceKeyValue{}
		return s.ReqFunc(in, out)
	}

	return errors.New("can't RequestCache")
}

func ListenBitMapRpc(ssrpc *BitMapRPC) {
	defer log.Fatal("Big Error : Rpc Server quit")
	defer Common.CheckPanic()

	server := rpc.NewServer()
	server.Register(ssrpc)

	if listener, err := net.Listen("tcp", ssrpc.Address); err != nil {
		log.Println("Big Error : Rpc can't start", err)
	} else {
		log.Println("Task Rpc running @", ssrpc.Address)
		server.Accept(listener)
	}
}
