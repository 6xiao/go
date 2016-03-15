package BloomFilter

import (
	"github.com/6xiao/go/BitMapCache"
)

type BloomFilter struct {
	haskKey []uint64
	capa    uint64
	cache   *BitMapCache.Bit1Cache
}

func NewBloomFilter(keys []uint64, cache *BitMapCache.Bit1Cache) *BloomFilter {
	return &BloomFilter{keys, uint64(cache.Capacity()), cache}
}

func NewDefaultBloomFilter() *BloomFilter {
	var primes = []uint64{61, 163, 233, 347, 401, 521, 661, 809}
	return NewBloomFilter(primes, BitMapCache.NewBit1Cache(1<<32))
}

func (this *BloomFilter) Keys() []uint64 {
	return this.haskKey
}

func (this *BloomFilter) KeySize() int {
	return len(this.haskKey)
}

func (this *BloomFilter) Capacity() int {
	return int(this.capa)
}

func (this *BloomFilter) Set(data []byte) {
	for _, key := range this.haskKey {
		hash := key
		for _, b := range data {
			hash = (hash << 5) + hash + uint64(b)
		}
		this.cache.SetLastBit(int(hash % this.capa))
	}
}

func (this *BloomFilter) Hits(data []byte) int {
	res := uint64(0)
	for _, key := range this.haskKey {
		hash := key
		for _, b := range data {
			hash = (hash << 5) + hash + uint64(b)
		}
		res += this.cache.GetIndex(int(hash % this.capa))
	}
	return int(res)
}

func (this *BloomFilter) IsSet(data []byte) bool {
	return this.Hits(data) == len(this.haskKey)
}
