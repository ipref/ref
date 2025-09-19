/* Copyright (c) 2025 Waldemar Augustyn */

package ref

import (
	"fmt"
	"math/big"
	"math/bits"
)

var UINT128_0 = Uint128FromUint64(0)
var UINT128_1 = Uint128FromUint64(1)
var UINT128_2_127 = UINT128_1.Lsh(127)
var UINT128_MAX = UINT128_0.Sub(UINT128_1)

type Uint128 struct {
	L, H uint64
}

func Uint128FromUint8(x uint8) Uint128 {
	return Uint128{uint64(x), 0}
}

func Uint128FromUint16(x uint16) Uint128 {
	return Uint128{uint64(x), 0}
}

func Uint128FromUint32(x uint32) Uint128 {
	return Uint128{uint64(x), 0}
}

func Uint128FromUint64(x uint64) Uint128 {
	return Uint128{x, 0}
}

func Uint128FromBig(i *big.Int) (Uint128, bool) {

	if i.Sign() < 0 || i.BitLen() > 128 {
		return Uint128{}, false
	}
	return Uint128{i.Uint64(), i.Rsh(i, 64).Uint64()}, true
}

func (x Uint128) Uint8Check() (uint8, bool) {

	if x.H != 0 || x.L >> 8 != 0 {
		return 0, false
	}
	return uint8(x.L), true
}

func (x Uint128) Uint16Check() (uint16, bool) {

	if x.H != 0 || x.L >> 16 != 0 {
		return 0, false
	}
	return uint16(x.L), true
}

func (x Uint128) Uint32Check() (uint32, bool) {

	if x.H != 0 || x.L >> 32 != 0 {
		return 0, false
	}
	return uint32(x.L), true
}

func (x Uint128) Uint64Check() (uint64, bool) {

	if x.H != 0 {
		return 0, false
	}
	return uint64(x.L), true
}

func (x Uint128) IntCheck() (int, bool) {

	if x.H != 0 || uint64(int(x.L)) != x.L {
		return 0, false
	}
	return int(x.L), true
}

func (x Uint128) Uint8() uint8 {
	return uint8(x.L)
}

func (x Uint128) Uint16() uint16 {
	return uint16(x.L)
}

func (x Uint128) Uint32() uint32 {
	return uint32(x.L)
}

func (x Uint128) Uint64() uint64 {
	return x.L
}

func (x Uint128) Int() int {
	return int(x.L)
}

func (x Uint128) Big() *big.Int {

	i := new(big.Int).SetUint64(x.H)
	i = i.Lsh(i, 64)
	return i.Or(i, new(big.Int).SetUint64(x.L))
}

func (x Uint128) IsZero() bool {
	return x == Uint128{}
}

func (x Uint128) Cmp(y Uint128) int {

	switch {
	case x.H < y.H: return -1
	case x.H > y.H: return 1
	case x.L < y.L: return -1
	case x.L > y.L: return 1
	}
	return 0
}

func (x Uint128) And(y Uint128) Uint128 {
	return Uint128{x.L & y.L, x.H & y.H}
}

func (x Uint128) Or(y Uint128) Uint128 {
	return Uint128{x.L | y.L, x.H | y.H}
}

func (x Uint128) Xor(y Uint128) Uint128 {
	return Uint128{x.L ^ y.L, x.H ^ y.H}
}

func (x Uint128) AndNot(y Uint128) Uint128 {
	return Uint128{x.L &^ y.L, x.H &^ y.H}
}

func (x Uint128) Compl() Uint128 {
	return Uint128{^x.L, ^x.H}
}

func (x Uint128) Lsh(n uint) Uint128 {

	if n >= 64 {
		return Uint128{0, x.L << (n - 64)}
	} else {
		return Uint128{x.L << n, (x.H << n) | (x.L >> (64 - n))}
	}
}

func (x Uint128) Rsh(n uint) Uint128 {

	if n >= 64 {
		return Uint128{x.H >> (n - 64), 0}
	} else {
		return Uint128{x.H << (64 - n) | (x.L >> n), x.H >> n}
	}
}

func (x Uint128) LeadingZeros() int {

	if x.H != 0 {
		return bits.LeadingZeros64(x.H)
	}
	return bits.LeadingZeros64(x.L) + 64
}

func (x Uint128) TrailingZeros() int {

	if x.L != 0 {
		return bits.TrailingZeros64(x.L)
	}
	return bits.TrailingZeros64(x.H) + 64
}

func (x Uint128) Bit(n int) uint {

	if n < 64 {
		return uint((x.L >> n) & 1)
	} else if n < 128 {
		return uint((x.H >> (n - 64)) & 1)
	}
	return 0
}

func (x Uint128) BitLen() int {
	return 128 - x.LeadingZeros()
}

func (x Uint128) Add(y Uint128) (z Uint128) {
	z, _ = x.AddCarry(y, 0)
	return
}

func (x Uint128) AddCarry(y Uint128, cin uint64) (z Uint128, cout uint64) {

	z.L, cout = bits.Add64(x.L, y.L, cin)
	z.H, cout = bits.Add64(x.H, y.H, cout)
	return
}

func (x Uint128) Sub(y Uint128) (z Uint128) {
	z, _ = x.SubBorrow(y, 0)
	return
}

func (x Uint128) SubBorrow(y Uint128, bin uint64) (z Uint128, bout uint64) {

	z.L, bout = bits.Sub64(x.L, y.L, bin)
	z.H, bout = bits.Sub64(x.H, y.H, bout)
	return
}

func (x Uint128) Mul(y Uint128) (z Uint128) {

	z.H, z.L = bits.Mul64(x.L, y.L)
	z.H += x.H * y.L + x.L * y.H
	return
}

func (x Uint128) String() string {
	return x.Format(10)
}

func (x Uint128) Format(base int) string {
	return x.Big().Text(base)
}

func (x Uint128) FormatHex() string {

	if x.IsZero() {
		return "0"
	}
	if x.H == 0 {
		return fmt.Sprintf("%x", x.L)
	}
	return fmt.Sprintf("%x%016x", x.H, x.L)
}

func ParseUint128(s string, base int) (Uint128, bool) {

	n, ok := new(big.Int).SetString(s, base)
	if !ok {
		return Uint128{}, false
	}
	return Uint128FromBig(n)
}

func MustParseUint128(s string, base int) Uint128 {

	val, ok := ParseUint128(s, base)
	if !ok {
		panic("invalid")
	}
	return val
}

func (x Uint128) PutBytesLE(dst []byte) {
	le.PutUint64(dst[:8], x.L)
	le.PutUint64(dst[8:], x.H)
}

func (x Uint128) PutBytesBE(dst []byte) {
	be.PutUint64(dst[:8], x.H)
	be.PutUint64(dst[8:], x.L)
}

func (x Uint128) AsBytesLE() (dst [16]byte) {
	x.PutBytesLE(dst[:])
	return
}

func (x Uint128) AsBytesBE() (dst [16]byte) {
	x.PutBytesBE(dst[:])
	return
}

func (x Uint128) AsSliceLE() []byte {
	bs := x.AsBytesLE()
	return bs[:]
}

func (x Uint128) AsSliceBE() []byte {
	bs := x.AsBytesBE()
	return bs[:]
}

func Uint128FromBytesLE(src []byte) Uint128 {
	return Uint128{le.Uint64(src[:8]), le.Uint64(src[8:])}
}

func Uint128FromBytesBE(src []byte) Uint128 {
	return Uint128{be.Uint64(src[8:]), be.Uint64(src[:8])}
}
