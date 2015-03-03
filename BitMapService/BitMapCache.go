package BitMapService

import "flag"

var (
	flg10E = flag.Int64("billion", 5, "billions of cache size, billion=1<<30")
	OpBits [64]uint64
)

func init() {
	for i := 0; i < len(OpBits); i++ {
		OpBits[i] = 1 << uint(i)
	}
}

type CacheCommon struct {
	bits  int
	max   int64
	count uint64
}

func (cc *CacheCommon) Bits() int {
	return cc.bits
}

func (cc *CacheCommon) Capacity() int64 {
	return (*flg10E) * BILLION
}

func (cc *CacheCommon) AddCount() {
	cc.count++
}

func (cc *CacheCommon) Count() uint64 {
	return cc.count
}

func (cc *CacheCommon) ResetCount() {
	cc.count = 0
}

func (cc *CacheCommon) SetMaxId(id int64) {
	if id > cc.max {
		cc.max = id
	}
}

func (cc *CacheCommon) MaxId() int64 {
	if cc.max >= cc.Capacity() {
		return cc.Capacity() - 1
	}

	return cc.max
}

//-----------------1------------------------------
type Bit1Cache struct {
	CacheCommon
	set []uint64
}

func NewBit1Cache() *Bit1Cache {
	return &Bit1Cache{CacheCommon{1, 0, 0}, make([]uint64, (*flg10E)*BILLION/64)}
}

func (sc *Bit1Cache) Shift() {
	sc.ResetCount()
	for i, _ := range sc.set {
		sc.set[i] = 0
	}
}

func (sc *Bit1Cache) SetIdLastFlag(id int64) bool {
	sc.SetMaxId(id)

	pre, post := id>>6, id&077
	value, flg := sc.set[pre], OpBits[post]

	if (value & flg) == 0 {
		sc.set[pre] = value | flg
		sc.AddCount()
		return true
	}

	return false
}

func (sc *Bit1Cache) GetIdFlags(id int64) uint64 {
	pre, post := id>>6, id&077
	value, flg := sc.set[pre], OpBits[post]

	if (value & flg) == 0 {
		return 0
	}

	return 1
}

func (sc *Bit1Cache) SetIdFlags(id int64, value uint64) {
	if value > 0 {
		sc.SetIdLastFlag(id)
	} else {
		pre, post := id>>6, id&077
		value, flg := sc.set[pre], OpBits[post]
		sc.set[pre] = value &^ flg
	}
}

//-----------------2------------------------------
type Bit2Cache struct {
	CacheCommon
	set0 []uint64
	set1 []uint64
}

func NewBit2Cache() *Bit2Cache {
	return &Bit2Cache{CacheCommon{2, 0, 0}, make([]uint64, (*flg10E)*BILLION/64), make([]uint64, (*flg10E)*BILLION/64)}
}

func (sc *Bit2Cache) Shift() {
	for i, v := range sc.set0 {
		sc.set1[i] = v
		sc.set0[i] = 0
	}

	sc.ResetCount()
	for i, last := int64(0), sc.MaxId(); i <= last; i++ {
		if sc.GetIdFlags(i) > 0 {
			sc.AddCount()
		}
	}
}

func (sc *Bit2Cache) SetIdLastFlag(id int64) bool {
	sc.SetMaxId(id)
	if sc.GetIdFlags(id) == 0 {
		sc.AddCount()
	}

	pre, post := id>>6, id&077
	value0, flg := sc.set0[pre], OpBits[post]

	if (value0 & flg) == 0 {
		sc.set0[pre] = value0 | flg
		return true
	}

	return false
}

func (sc *Bit2Cache) GetIdFlags(id int64) uint64 {
	pre, post := id>>6, id&077
	value0, value1, flg := sc.set0[pre], sc.set1[pre], OpBits[post]

	if (value0 & flg) == 0 {
		if (value1 & flg) == 0 {
			return 0
		} else {
			return 2
		}
	}

	if (value1 & flg) == 0 {
		return 1
	}

	return 3
}

func (sc *Bit2Cache) SetIdFlags(id int64, value uint64) {
	if value > 0 && sc.GetIdFlags(id) == 0 {
		sc.AddCount()
	}

	sc.SetMaxId(id)

	pre, post := id>>6, id&077
	value0, value1, flg := sc.set0[pre], sc.set1[pre], OpBits[post]

	switch value {
	case 0:
		sc.set0[pre] = value0 &^ flg
		sc.set1[pre] = value1 &^ flg

	case 1:
		sc.set0[pre] = value0 | flg
		sc.set1[pre] = value1 &^ flg

	case 2:
		sc.set0[pre] = value0 &^ flg
		sc.set1[pre] = value1 | flg

	case 3:
		sc.set0[pre] = value0 | flg
		sc.set1[pre] = value1 | flg
	}
}

// ---------------------------4--------------------------------
type Bit4Cache struct {
	CacheCommon
	set []uint8
}

func NewBit4Cache() *Bit4Cache {
	return &Bit4Cache{CacheCommon{4, 0, 0}, make([]uint8, (*flg10E)*BILLION/2)}
}

func (sc *Bit4Cache) Shift() {
	sc.ResetCount()
	for i, flags := range sc.set {
		if flags&0x7 != 0 {
			sc.AddCount()
		}
		if flags&0x70 != 0 {
			sc.AddCount()
		}
		sc.set[i] = (flags << 1) & 0xEE
	}
}

func (sc *Bit4Cache) SetIdLastFlag(id int64) bool {
	sc.SetMaxId(id)

	index, lastbit := id>>1, id&1
	if lastbit == 0 {
		if (sc.set[index] & 0xF) == 0 {
			sc.AddCount()
		}

		if (sc.set[index] & 1) == 0 {
			sc.set[index] |= 1
			return true
		}
	} else {
		if (sc.set[index] & 0xF0) == 0 {
			sc.AddCount()
		}

		if (sc.set[index] & 16) == 0 {
			sc.set[index] |= 16
			return true
		}
	}

	return false
}

func (sc *Bit4Cache) GetIdFlags(id int64) uint64 {
	index, lastbit := id>>1, id&1

	if lastbit == 0 {
		return uint64(sc.set[index] & 0xF)
	}

	return uint64(sc.set[index] >> 4)
}

func (sc *Bit4Cache) SetIdFlags(id int64, value uint64) {
	if value > 0 && sc.GetIdFlags(id) == 0 {
		sc.AddCount()
	}

	sc.SetMaxId(id)

	index, lastbit := id>>1, id&1

	if lastbit == 0 {
		sc.set[index] = (sc.set[index] & 0xF0) | uint8(value&0xF)
	}

	sc.set[index] = uint8(value<<4) | (sc.set[index] & 0xF)
}

// ---------------------------8--------------------------------
type Bit8Cache struct {
	CacheCommon
	set []uint8
}

func NewBit8Cache() *Bit8Cache {
	return &Bit8Cache{CacheCommon{8, 0, 0}, make([]uint8, (*flg10E)*BILLION)}
}

func (sc *Bit8Cache) Shift() {
	sc.ResetCount()
	for i, v := range sc.set {
		sc.set[i] = v << 1
		if sc.set[i] > 0 {
			sc.AddCount()
		}
	}
}

func (sc *Bit8Cache) SetIdLastFlag(id int64) bool {
	sc.SetMaxId(id)

	if sc.set[id] == 0 {
		sc.AddCount()
	}

	if (sc.set[id] & 1) == 0 {
		sc.set[id]++
		return true
	}

	return false
}

func (sc *Bit8Cache) GetIdFlags(id int64) uint64 {
	return uint64(sc.set[id])
}

func (sc *Bit8Cache) SetIdFlags(id int64, value uint64) {
	if value > 0 && sc.set[id] == 0 {
		sc.AddCount()
	}

	sc.SetMaxId(id)
	sc.set[id] = uint8(value)
}

// ---------------------------16-------------------------------------
type Bit16Cache struct {
	CacheCommon
	set []uint16
}

func NewBit16Cache() *Bit16Cache {
	return &Bit16Cache{CacheCommon{16, 0, 0}, make([]uint16, (*flg10E)*BILLION)}
}

func (sc *Bit16Cache) Shift() {
	sc.ResetCount()
	for i, v := range sc.set {
		sc.set[i] = v << 1
		if sc.set[i] > 0 {
			sc.AddCount()
		}
	}
}

func (sc *Bit16Cache) SetIdLastFlag(id int64) bool {
	sc.SetMaxId(id)

	if sc.set[id] == 0 {
		sc.AddCount()
	}

	if (sc.set[id] & 1) == 0 {
		sc.set[id]++
		return true
	}

	return false
}

func (sc *Bit16Cache) GetIdFlags(id int64) uint64 {
	return uint64(sc.set[id])
}

func (sc *Bit16Cache) SetIdFlags(id int64, value uint64) {
	if value > 0 && sc.set[id] == 0 {
		sc.AddCount()
	}

	sc.SetMaxId(id)
	sc.set[id] = uint16(value)
}

// -----------------------------32------------------------
type Bit32Cache struct {
	CacheCommon
	set []uint32
}

func NewBit32Cache() *Bit32Cache {
	return &Bit32Cache{CacheCommon{32, 0, 0}, make([]uint32, (*flg10E)*BILLION)}
}

func (sc *Bit32Cache) Shift() {
	sc.ResetCount()
	for i, v := range sc.set {
		sc.set[i] = v << 1
		if sc.set[i] > 0 {
			sc.AddCount()
		}
	}
}

func (sc *Bit32Cache) SetIdLastFlag(id int64) bool {
	sc.SetMaxId(id)

	if sc.set[id] == 0 {
		sc.AddCount()
	}

	if (sc.set[id] & 1) == 0 {
		sc.set[id]++
		return true
	}

	return false
}

func (sc *Bit32Cache) GetIdFlags(id int64) uint64 {
	return uint64(sc.set[id])
}

func (sc *Bit32Cache) SetIdFlags(id int64, value uint64) {
	if value > 0 && sc.set[id] == 0 {
		sc.AddCount()
	}

	sc.SetMaxId(id)
	sc.set[id] = uint32(value)
}

// ---------------------------64----------------------------
type Bit64Cache struct {
	CacheCommon
	set []uint64
}

func NewBit64Cache() *Bit64Cache {
	return &Bit64Cache{CacheCommon{64, 0, 0}, make([]uint64, (*flg10E)*BILLION)}
}

func (sc *Bit64Cache) Shift() {
	sc.ResetCount()
	for i, v := range sc.set {
		sc.set[i] = v << 1
		if sc.set[i] > 0 {
			sc.AddCount()
		}
	}
}

func (sc *Bit64Cache) SetIdLastFlag(id int64) bool {
	sc.SetMaxId(id)

	if sc.set[id] == 0 {
		sc.AddCount()
	}

	if (sc.set[id] & 1) == 0 {
		sc.set[id]++
		return true
	}

	return false
}

func (sc *Bit64Cache) GetIdFlags(id int64) uint64 {
	return sc.set[id]
}

func (sc *Bit64Cache) SetIdFlags(id int64, value uint64) {
	if value > 0 && sc.set[id] == 0 {
		sc.AddCount()
	}

	sc.SetMaxId(id)
	sc.set[id] = value
}
