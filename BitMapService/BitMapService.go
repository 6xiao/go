package BitMapService

import (
	"errors"
)

const (
	OPER_GET     = "GET"
	OPER_SET     = "SET"
	OPER_SETLAST = "SETLAST"
	OPER_SHIFT   = "SHIFT"
	BILLION      = 1024 * 1024 * 1024
	RPC_CALL     = "BitMapServer.RequestCache"
	RPC_SYNC     = "BitMapServer.ReverseSync"
)

type Cache interface {
	//how many bits of the id's flag, such as 1,2,4,8,16,32,64
	Bits() int

	Capacity() int64

	//how much uid in cache
	Count() uint64

	//get the max id
	MaxId() int64

	// all id's flags <<1
	Shift()

	//set id's last bit, if is the first time return true
	SetIdLastFlag(id int64) bool

	//get id's flags
	GetIdFlags(id int64) uint64

	//set id's all flag in once time
	SetIdFlags(id int64, value uint64)
}

type SliceKeyValue struct {
	Ids   []int64
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

type BitMapServer struct {
	SyncFunc SyncInvoke
	ReqFunc  ReqInvoke
}

func NewBitMapServer(sync SyncInvoke, req ReqInvoke) *BitMapServer {
	return &BitMapServer{sync, req}
}

//1, client give a rpc address(in) to server
//2, server connect to client
//3, server push all datas to client use the RequestCache interface
//in : address
//out : success
func (s *BitMapServer) ReverseSync(in string, out *bool) error {
	if s.SyncFunc != nil {
		return s.SyncFunc(in, out)
	}

	return errors.New("can't ReverseSync")
}

//out: first inserted id when Operator == SET
func (s *BitMapServer) RequestCache(in BitMapRequest, out *SliceKeyValue) error {
	if s.ReqFunc != nil {
		*out = SliceKeyValue{}
		return s.ReqFunc(in, out)
	}

	return errors.New("can't RequestCache")
}
