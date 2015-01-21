package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/6xiao/go/BitMapService"
	"github.com/6xiao/go/Common"
	"log"
	"net/rpc"
	"strings"
	"sync"
	"time"
)

var (
	flgRpc  = flag.String("rpc", "127.0.0.1:12321", "rpc address of this server")
	flgSync = flag.String("sync", "127.0.0.1:12321", "sync rpc address of remote server")
	lock    = sync.Mutex{}
	caches  = make(map[string]BitMapService.Cache)
)

func init() {
	Common.Init()
}

func GetCache(name string, bits int) (BitMapService.Cache, error) {
	cache, ok := caches[name]
	if !ok {
		switch bits {
		case 1:
			cache = BitMapService.NewBit1Cache()

		case 2:
			cache = BitMapService.NewBit2Cache()

		case 4:
			cache = BitMapService.NewBit4Cache()

		case 8:
			cache = BitMapService.NewBit8Cache()

		case 16:
			cache = BitMapService.NewBit16Cache()

		case 32:
			cache = BitMapService.NewBit32Cache()

		case 64:
			cache = BitMapService.NewBit64Cache()

		default:
			return nil, errors.New(fmt.Sprint("unsupport bits", bits))
		}

		caches[name] = cache
		log.Println("create new cache:", name, " bits:", bits)
	}

	if cache.Bits() != bits {
		return nil, errors.New(fmt.Sprint("can't use", cache.Bits(), "bitmap as", bits, "bitmap named", name))
	}

	return cache, nil
}

func Request(in BitMapService.BitMapRequest, out *BitMapService.SliceKeyValue) error {
	defer Common.CheckPanic()
	defer lock.Unlock()
	lock.Lock()

	if in.Operator == BitMapService.OPER_INFO {
		for cachename, cache := range caches {
			log.Println("infomation:", cachename,
				" bits:", cache.Bits(),
				" max uid:", cache.MaxUid(),
				" total count:", cache.Count())
		}

		return nil
	}

	cache, err := GetCache(in.Name, in.Bits)
	if err != nil {
		return err
	}

	maxuid := cache.Capacity() - 1
	switch in.Operator {
	case BitMapService.OPER_SET:
		for _, uid := range in.Uids {
			if uid > maxuid {
				log.Println("too big uid :", uid)
			} else if cache.SetUidLastFlag(uid) {
				out.Uids = append(out.Uids, uid)
			}
		}

	case BitMapService.OPER_GET:
		for _, uid := range in.Uids {
			if uid > maxuid {
				log.Println("too big uid :", uid)
			} else {
				out.Uids = append(out.Uids, uid)
				out.Flags = append(out.Flags, cache.GetUidFlags(uid))
			}
		}

	case BitMapService.OPER_SHIFT:
		cache.Shift()
		log.Println("shift:", in.Name, " max uid:", cache.MaxUid(), " total count:", cache.Count())

	case BitMapService.OPER_SYNC:
		if len(in.Uids) != len(in.Flags) {
			return errors.New("uids != flags")
		}

		for index, uid := range in.Uids {
			if uid > maxuid {
				log.Println("too big uid :", uid)
			} else {
				cache.SetUidFlags(uid, in.Flags[index]|cache.GetUidFlags(uid))
			}
		}

	default:
		return errors.New("unsupport operator: " + in.Operator)
	}

	return nil
}

func ReverseSync(in string, out *bool) error {
	defer Common.CheckPanic()
	log.Println("sync request from :", in)

	conn, err := rpc.Dial("tcp", in)
	if err != nil {
		log.Println("error: DialRPC ", in, err)
		return err
	}

	defer conn.Close()

	for cachename, cache := range caches {
		log.Println("sync:", cachename, "total count :", cache.Count())

		skv := BitMapService.SliceKeyValue{}
		for uid, max := int64(0), int64(cache.MaxUid()); uid <= max; uid++ {
			if flags := cache.GetUidFlags(uid); flags != 0 {
				skv.Uids = append(skv.Uids, uid)
				skv.Flags = append(skv.Flags, flags)
			}

			if len(skv.Uids) == 1024*1024 || uid == max {
				umr := BitMapService.BitMapRequest{cache.Bits(), cachename, BitMapService.OPER_SYNC, skv}

				if err := conn.Call(BitMapService.RPC_CALL, umr, nil); err != nil {
					log.Println("error: Call ", err)
					return err
				}

				skv = BitMapService.SliceKeyValue{}
			}
		}
	}

	log.Println("sync over :", in)

	*out = true
	return nil
}

func RpcSync() {
	defer Common.CheckPanic()
	time.Sleep(time.Second * 3)

	rpcs := strings.Split(*flgSync, ",")
	for _, uri := range rpcs {
		client, err := rpc.Dial("tcp", uri)
		if err != nil {
			log.Println("error: DialRPC ", uri, err)
			continue
		}
		defer client.Close()

		suc := false
		err = client.Call(BitMapService.RPC_SYNC, *flgRpc, &suc)

		log.Println("sync result :", uri, suc, err)

		for cachename, cache := range caches {
			log.Println("infomation:", cachename,
				" bits:", cache.Bits(),
				" max uid:", cache.MaxUid(),
				" total count:", cache.Count())
		}
	}
}

func main() {
	go RpcSync()

	trpc := BitMapService.NewBitMapRpc(*flgRpc, ReverseSync, Request)
	BitMapService.ListenBitMapRpc(trpc)
	log.Fatal("exit ...")
}
