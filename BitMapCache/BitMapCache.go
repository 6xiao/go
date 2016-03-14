package BitMapCache

import (
	"encoding/binary"
	"io"
)

var OpBits = [8]uint8{1, 2, 4, 8, 16, 32, 64, 128}

type Cache interface {
	//how many bits of the id's flag, eg: 1,2,4,8,16,32,64
	Bits() int

	// buffer capacity
	Capacity() int

	//how much value != 0
	Count() int

	//get the max index
	MaxIndex() int

	// all id's flags <<1
	ShiftOneBit()

	//set last bit, if is the first time return true
	SetLastBit(index int) bool

	//get index value
	GetIndex(index int) uint64

	//set index value
	SetIndex(index int, value uint64)

	ReadFrom(r io.Reader) (int, error)
	WriterTo(w io.Writer) (int, error)
}

type CacheCommon struct {
	bits  int
	capa  int
	max   int
	count int
}

func (this *CacheCommon) Bits() int {
	return this.bits
}

func (this *CacheCommon) Capacity() int {
	return this.capa
}

func (this *CacheCommon) addCount() {
	this.count++
}

func (this *CacheCommon) Count() int {
	return this.count
}

func (this *CacheCommon) resetCount() {
	this.count = 0
}

func (this *CacheCommon) setMaxIndex(index int) {
	if index > this.max {
		this.max = index
	}
}

func (this *CacheCommon) MaxIndex() int {
	if this.max >= this.capa {
		return this.capa - 1
	}

	return this.max
}

//-----------------1------------------------------
type Bit1Cache struct {
	CacheCommon
	set []uint8
}

func NewBit1Cache(capa int) *Bit1Cache {
	set := make([]uint8, capa/8+8)
	return &Bit1Cache{CacheCommon{1, capa, 0, 0}, set}
}

func (this *Bit1Cache) ShiftOneBit() {
	this.resetCount()
	for i, _ := range this.set {
		this.set[i] = 0
	}
}

func (this *Bit1Cache) SetLastBit(index int) bool {
	this.setMaxIndex(index)

	pre, post := index>>3, index&7
	value, flg := this.set[pre], OpBits[post]

	if (value & flg) == 0 {
		this.set[pre] = value | flg
		this.addCount()
		return true
	}

	return false
}

func (this *Bit1Cache) GetIndex(index int) uint64 {
	pre, post := index>>3, index&7
	value, flg := this.set[pre], OpBits[post]

	if (value & flg) == 0 {
		return 0
	}

	return 1
}

func (this *Bit1Cache) ResetIndex(index int, value uint64) {
	if value > 0 {
		this.SetLastBit(index)
	} else {
		pre, post := index>>3, index&7
		value, flg := this.set[pre], OpBits[post]
		this.set[pre] = value &^ flg
	}
}

func (this *Bit1Cache) ReadFrom(r io.Reader) (int, error) {
	return r.Read(this.set)
}

func (this *Bit1Cache) WriteTo(w io.Writer) (int, error) {
	return w.Write(this.set)
}

//-----------------2------------------------------
type Bit2Cache struct {
	CacheCommon
	set []uint8
}

func NewBit2Cache(capa int) *Bit2Cache {
	set := make([]uint8, capa/4+8)
	return &Bit2Cache{CacheCommon{2, capa, 0, 0}, set}
}

func (this *Bit2Cache) ShiftOneBit() {
	end := this.capa/4 + 1
	for i := 0; i < end; i += 2 {
		this.set[i+1] = this.set[i]
		this.set[i] = 0
	}

	this.resetCount()
	for i, last := 0, this.MaxIndex(); i <= last; i++ {
		if this.GetIndex(i) > 0 {
			this.addCount()
		}
	}
}

func (this *Bit2Cache) SetLastBit(index int) bool {
	this.setMaxIndex(index)
	if this.GetIndex(index) == 0 {
		this.addCount()
	}

	pre, post := index>>2, index&7
	value, flg := this.set[pre], OpBits[post]

	if (value & flg) == 0 {
		this.set[pre] = value | flg
		return true
	}

	return false
}

func (this *Bit2Cache) GetIndex(index int) uint64 {
	pre, post := index>>2, index&7
	value0, value1, flg := this.set[pre], this.set[pre+1], OpBits[post]
	bit0, bit1 := (value0&flg)/flg, (value1&flg)/flg
	return uint64(bit0 + (bit1 << 1))
}

func (this *Bit2Cache) ResetIndex(index int, value uint64) {
	if value > 0 && this.GetIndex(index) == 0 {
		this.addCount()
	}

	this.setMaxIndex(index)

	pre, post := index>>2, index&7
	value0, value1, flg := this.set[pre], this.set[pre+1], OpBits[post]

	switch value {
	case 0:
		this.set[pre] = value0 &^ flg
		this.set[pre+1] = value1 &^ flg

	case 1:
		this.set[pre] = value0 | flg
		this.set[pre+1] = value1 &^ flg

	case 2:
		this.set[pre] = value0 &^ flg
		this.set[pre+1] = value1 | flg

	case 3:
		this.set[pre] = value0 | flg
		this.set[pre+1] = value1 | flg
	}
}

func (this *Bit2Cache) ReadFrom(r io.Reader) (int, error) {
	return r.Read(this.set)
}

func (this *Bit2Cache) WriteTo(w io.Writer) (int, error) {
	return w.Write(this.set)
}

// ---------------------------4--------------------------------
type Bit4Cache struct {
	CacheCommon
	set []uint8
}

func NewBit4Cache(capa int) *Bit4Cache {
	set := make([]uint8, capa/2+8)
	return &Bit4Cache{CacheCommon{4, capa, 0, 0}, set}
}

func (this *Bit4Cache) ShiftOneBit() {
	this.resetCount()
	for i, flags := range this.set {
		if flags&0x7 != 0 {
			this.addCount()
		}
		if flags&0x70 != 0 {
			this.addCount()
		}
		this.set[i] = (flags << 1) & 0xEE
	}
}

func (this *Bit4Cache) SetLastBit(index int) bool {
	this.setMaxIndex(index)

	idx, lastbit := index>>1, index&1
	if lastbit == 0 {
		if (this.set[idx] & 0xF) == 0 {
			this.addCount()
		}

		if (this.set[idx] & 1) == 0 {
			this.set[idx] |= 1
			return true
		}
	} else {
		if (this.set[idx] & 0xF0) == 0 {
			this.addCount()
		}

		if (this.set[idx] & 16) == 0 {
			this.set[idx] |= 16
			return true
		}
	}

	return false
}

func (this *Bit4Cache) GetIndex(index int) uint64 {
	idx, lastbit := index>>1, index&1

	if lastbit == 0 {
		return uint64(this.set[idx] & 0xF)
	}

	return uint64(this.set[idx] >> 4)
}

func (this *Bit4Cache) ResetIndex(index int, value uint64) {
	if value > 0 && this.GetIndex(index) == 0 {
		this.addCount()
	}

	this.setMaxIndex(index)

	idx, lastbit := index>>1, index&1

	if lastbit == 0 {
		this.set[idx] = (this.set[idx] & 0xF0) | uint8(value&0xF)
	}

	this.set[idx] = uint8(value<<4) | (this.set[idx] & 0xF)
}

func (this *Bit4Cache) ReadFrom(r io.Reader) (int, error) {
	return r.Read(this.set)
}

func (this *Bit4Cache) WriteTo(w io.Writer) (int, error) {
	return w.Write(this.set)
}

// ---------------------------8--------------------------------
type Bit8Cache struct {
	CacheCommon
	set []uint8
}

func NewBit8Cache(capa int) *Bit8Cache {
	return &Bit8Cache{CacheCommon{8, capa, 0, 0}, make([]uint8, capa+1)}
}

func (this *Bit8Cache) ShiftOneBit() {
	this.resetCount()
	for i, v := range this.set {
		this.set[i] = v << 1
		if this.set[i] > 0 {
			this.addCount()
		}
	}
}

func (this *Bit8Cache) SetLastBit(index int) bool {
	if index >= this.capa {
		return false
	}

	this.setMaxIndex(index)

	if this.set[index] == 0 {
		this.addCount()
	}

	if (this.set[index] & 1) == 0 {
		this.set[index]++
		return true
	}

	return false
}

func (this *Bit8Cache) GetIndex(index int) uint64 {
	return uint64(this.set[index])
}

func (this *Bit8Cache) ResetIndex(index int, value uint64) {
	if value > 0 && this.set[index] == 0 {
		this.addCount()
	}

	this.setMaxIndex(index)
	this.set[index] = uint8(value)
}

func (this *Bit8Cache) ReadFrom(r io.Reader) (int, error) {
	return r.Read(this.set)
}

func (this *Bit8Cache) WriteTo(w io.Writer) (int, error) {
	return w.Write(this.set)
}

// ---------------------------16-------------------------------------
type Bit16Cache struct {
	CacheCommon
	set []uint16
}

func NewBit16Cache(capa int) *Bit16Cache {
	return &Bit16Cache{CacheCommon{16, capa, 0, 0}, make([]uint16, capa+1)}
}

func (this *Bit16Cache) ShiftOneBit() {
	this.resetCount()
	for i, v := range this.set {
		this.set[i] = v << 1
		if this.set[i] > 0 {
			this.addCount()
		}
	}
}

func (this *Bit16Cache) SetLastBit(index int) bool {
	this.setMaxIndex(index)

	if this.set[index] == 0 {
		this.addCount()
	}

	if (this.set[index] & 1) == 0 {
		this.set[index]++
		return true
	}

	return false
}

func (this *Bit16Cache) GetIndex(index int) uint64 {
	return uint64(this.set[index])
}

func (this *Bit16Cache) ResetIndex(index int, value uint64) {
	if value > 0 && this.set[index] == 0 {
		this.addCount()
	}

	this.setMaxIndex(index)
	this.set[index] = uint16(value)
}

func (this *Bit16Cache) ReadFrom(r io.Reader) (int, error) {
	index, buf := 0, make([]byte, 2)
	for index < len(this.set) {
		if nr, err := r.Read(buf); nr != 2 || err != nil {
			break
		}
		this.set[index] = binary.LittleEndian.Uint16(buf)
		index += 1
	}
	return index, nil
}

func (this *Bit16Cache) WriteTo(w io.Writer) (int, error) {
	buf := make([]byte, 2)
	for i, v := range this.set {
		binary.LittleEndian.PutUint16(buf, v)
		if nr, err := w.Write(buf); nr != 2 || err != nil {
			return i, nil
		}
	}
	return 0, nil
}

// -----------------------------32------------------------
type Bit32Cache struct {
	CacheCommon
	set []uint32
}

func NewBit32Cache(capa int) *Bit32Cache {
	return &Bit32Cache{CacheCommon{32, capa, 0, 0}, make([]uint32, capa+1)}
}

func (this *Bit32Cache) ShiftOneBit() {
	this.resetCount()
	for i, v := range this.set {
		this.set[i] = v << 1
		if this.set[i] > 0 {
			this.addCount()
		}
	}
}

func (this *Bit32Cache) SetLastBit(index int) bool {
	this.setMaxIndex(index)

	if this.set[index] == 0 {
		this.addCount()
	}

	if (this.set[index] & 1) == 0 {
		this.set[index]++
		return true
	}

	return false
}

func (this *Bit32Cache) GetIndex(index int) uint64 {
	return uint64(this.set[index])
}

func (this *Bit32Cache) ResetIndex(index int, value uint64) {
	if value > 0 && this.set[index] == 0 {
		this.addCount()
	}

	this.setMaxIndex(index)
	this.set[index] = uint32(value)
}

func (this *Bit32Cache) ReadFrom(r io.Reader) (int, error) {
	index, buf := 0, make([]byte, 4)
	for index < len(this.set) {
		if nr, err := r.Read(buf); nr != 4 || err != nil {
			break
		}
		this.set[index] = binary.LittleEndian.Uint32(buf)
		index += 1
	}
	return index, nil
}

func (this *Bit32Cache) WriteTo(w io.Writer) (int, error) {
	buf := make([]byte, 4)
	for i, v := range this.set {
		binary.LittleEndian.PutUint32(buf, v)
		if nr, err := w.Write(buf); nr != 4 || err != nil {
			return i, nil
		}
	}
	return 0, nil
}

// ---------------------------64----------------------------
type Bit64Cache struct {
	CacheCommon
	set []uint64
}

func NewBit64Cache(capa int) *Bit64Cache {
	return &Bit64Cache{CacheCommon{64, capa, 0, 0}, make([]uint64, capa+1)}
}

func (this *Bit64Cache) ShiftOneBit() {
	this.resetCount()
	for i, v := range this.set {
		this.set[i] = v << 1
		if this.set[i] > 0 {
			this.addCount()
		}
	}
}

func (this *Bit64Cache) SetLastBit(index int) bool {
	this.setMaxIndex(index)

	if this.set[index] == 0 {
		this.addCount()
	}

	if (this.set[index] & 1) == 0 {
		this.set[index]++
		return true
	}

	return false
}

func (this *Bit64Cache) GetIndex(index int) uint64 {
	return this.set[index]
}

func (this *Bit64Cache) ResetIndex(index int, value uint64) {
	if value > 0 && this.set[index] == 0 {
		this.addCount()
	}

	this.setMaxIndex(index)
	this.set[index] = value
}

func (this *Bit64Cache) ReadFrom(r io.Reader) (int, error) {
	index, buf := 0, make([]byte, 8)
	for index < len(this.set) {
		if nr, err := r.Read(buf); nr != 8 || err != nil {
			break
		}
		this.set[index] = binary.LittleEndian.Uint64(buf)
		index += 1
	}
	return index, nil
}

func (this *Bit64Cache) WriteTo(w io.Writer) (int, error) {
	buf := make([]byte, 8)
	for i, v := range this.set {
		binary.LittleEndian.PutUint64(buf, v)
		if nr, err := w.Write(buf); nr != 8 || err != nil {
			return i, nil
		}
	}
	return 0, nil
}
