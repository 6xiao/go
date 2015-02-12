package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/6xiao/go/BitMapService"
	"github.com/6xiao/go/Common"
	"log"
	"net/rpc"
	"os"
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
	Common.Init(os.Stderr)
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

	cache, err := GetCache(in.Name, in.Bits)
	if err != nil {
		return err
	}

	maxid := cache.Capacity() - 1
	switch in.Operator {
	case BitMapService.OPER_SETLAST:
		for _, id := range in.Ids {
			if id > maxid {
				log.Println("too big id :", id)
			} else if cache.SetIdLastFlag(id) {
				out.Ids = append(out.Ids, id)
			}
		}

	case BitMapService.OPER_GET:
		for _, id := range in.Ids {
			if id > maxid {
				log.Println("too big id :", id)
			} else {
				out.Ids = append(out.Ids, id)
				out.Flags = append(out.Flags, cache.GetIdFlags(id))
			}
		}

	case BitMapService.OPER_SHIFT:
		cache.Shift()
		log.Println("shift:", in.Name, " max id:", cache.MaxId(), " total count:", cache.Count())

	case BitMapService.OPER_SET:
		if len(in.Ids) != len(in.Flags) {
			return errors.New("length of ids != flags")
		}

		for index, id := range in.Ids {
			if id > maxid {
				log.Println("too big id :", id)
			} else {
				cache.SetIdFlags(id, in.Flags[index])
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
		for id, max := int64(0), int64(cache.MaxId()); id <= max; id++ {
			if flags := cache.GetIdFlags(id); flags != 0 {
				skv.Ids = append(skv.Ids, id)
				skv.Flags = append(skv.Flags, flags)
			}

			if len(skv.Ids) == 1024*1024 || id == max {
				umr := BitMapService.BitMapRequest{cache.Bits(), cachename, BitMapService.OPER_SET, skv}

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
	time.Sleep(time.Second)

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
				" max id:", cache.MaxId(),
				" total count:", cache.Count())
		}
	}
}

func main() {
	go RpcSync()

	trpc := BitMapService.NewBitMapServer(ReverseSync, Request)
	Common.ListenRpc(*flgRpc, trpc, nil)
	log.Fatal("exit ...")
}
