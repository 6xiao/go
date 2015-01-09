package main

import (
	"flag"
	"github.com/6xiao/go/BitMapService"
	"github.com/6xiao/go/Common"
	"log"
	"net/rpc"
)

var (
	flgRpc  = flag.String("rpc", "127.0.0.1:12321", "rpc server")
	flgBits = flag.Int("bits", 64, "cache bits")
)

func RpcUidSet() {
	defer Common.CheckPanic()

	client, err := rpc.Dial("tcp", *flgRpc)
	if err != nil {
		log.Println("error: DialRPC ", *flgRpc, err)
		return
	}

	for i := 0; i < 256; i++ {
		skv := BitMapService.SliceKeyValue{}
		for uid, last := int64(i*8192), int64((i+1)*8192); uid < last; uid++ {
			skv.Uids = append(skv.Uids, uid)
		}

		umr := BitMapService.BitMapRequest{*flgBits, "UidTest", BitMapService.OPER_SET, skv}
		out := BitMapService.SliceKeyValue{}
		if err := client.Call(BitMapService.RPC_CALL, umr, &out); err != nil {
			log.Println("error: Call ", err)
			return
		}

		log.Println("send:", len(skv.Uids), "changed:", len(out.Uids), err)
	}
}

func RpcUidGet() {
	defer Common.CheckPanic()

	client, err := rpc.Dial("tcp", *flgRpc)
	if err != nil {
		log.Println("error: DialRPC ", *flgRpc, err)
		return
	}

	skv := BitMapService.SliceKeyValue{}
	for uid := int64(0); uid < 10; uid++ {
		skv.Uids = append(skv.Uids, uid)
	}
	req := BitMapService.BitMapRequest{*flgBits, "UidTest", BitMapService.OPER_GET, skv}

	out := BitMapService.SliceKeyValue{}
	if err := client.Call(BitMapService.RPC_CALL, req, &out); err != nil {
		log.Println("error: Call ", err)
		return
	}
	log.Println("get :", len(out.Uids), len(out.Flags))

	for i, uid := range out.Uids {
		log.Println("uid:", uid, "flags:", out.Flags[i])
	}
}

func RpcCacheShift() {
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
	log.Println("shift end, out len:", len(out.Uids))
}

func main() {
	RpcUidSet()
	RpcUidGet()
	RpcCacheShift()
	RpcUidGet()
}
