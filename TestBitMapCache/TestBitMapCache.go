package main

import (
	"flag"
	"github.com/6xiao/go/BitMapService"
	"github.com/6xiao/go/Common"
	"log"
	"net/rpc"
	"os"
)

var (
	flgRpc  = flag.String("rpc", "127.0.0.1:12321", "rpc  address of cache server")
	flgBits = flag.Int("bits", 64, "cache bits")
)

func init() {
	Common.Init(os.Stderr)
}

func RpcSet() {
	defer Common.CheckPanic()

	client, err := rpc.Dial("tcp", *flgRpc)
	if err != nil {
		log.Println("error: DialRPC ", *flgRpc, err)
		return
	}

	for i := 0; i < 256; i++ {
		skv := BitMapService.SliceKeyValue{}
		for id, last := int64(i*8192), int64((i+1)*8192); id < last; id++ {
			skv.Ids = append(skv.Ids, id)
		}

		umr := BitMapService.BitMapRequest{*flgBits, "UidTest", BitMapService.OPER_SETLAST, skv}
		out := BitMapService.SliceKeyValue{}
		if err := client.Call(BitMapService.RPC_CALL, umr, &out); err != nil {
			log.Println("error: Call ", err)
			return
		}

		log.Println("send:", len(skv.Ids), "changed:", len(out.Ids), err)
	}
}

func RpcGet() {
	defer Common.CheckPanic()

	client, err := rpc.Dial("tcp", *flgRpc)
	if err != nil {
		log.Println("error: DialRPC ", *flgRpc, err)
		return
	}

	skv := BitMapService.SliceKeyValue{}
	for id := int64(0); id < 10; id++ {
		skv.Ids = append(skv.Ids, id)
	}
	req := BitMapService.BitMapRequest{*flgBits, "UidTest", BitMapService.OPER_GET, skv}

	out := BitMapService.SliceKeyValue{}
	if err := client.Call(BitMapService.RPC_CALL, req, &out); err != nil {
		log.Println("error: Call ", err)
		return
	}
	log.Println("get :", len(out.Ids), len(out.Flags))

	for i, id := range out.Ids {
		log.Println("id:", id, "flags:", out.Flags[i])
	}
}

func RpcShift() {
	defer Common.CheckPanic()

	client, err := rpc.Dial("tcp", *flgRpc)
	if err != nil {
		log.Println("error: DialRPC ", *flgRpc, err)
		return
	}

	req := BitMapService.BitMapRequest{*flgBits, "UidTest", BitMapService.OPER_SHIFT, BitMapService.SliceKeyValue{}}
	out := BitMapService.SliceKeyValue{}
	if err := client.Call(BitMapService.RPC_CALL, req, &out); err != nil {
		log.Println("error: Call ", err)
		return
	}
	log.Println("shift end, out len:", len(out.Ids))
}

func main() {
	RpcSet()
	RpcGet()
	RpcShift()
	RpcGet()
}
