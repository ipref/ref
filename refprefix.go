/* Copyright (c) 2025 Waldemar Augustyn */

package ref

import (
	"errors"
	"strconv"
	"strings"
)

// The zero value is the prefix which contains all refs
type RefPrefix struct {
	ref Ref // The bits that aren't part of the prefix must be zero
	bits int // Must be <= 128
}

func (p RefPrefix) Ref() Ref {
	return p.ref
}

func (p RefPrefix) RefUint128() Uint128 {
	return Uint128(p.ref)
}

func (p RefPrefix) Bits() int {
	return p.bits
}

func (p RefPrefix) SizeBits() int {
	return 128 - p.bits
}

func RefPrefixFrom(ref Ref, bits int) RefPrefix {

	if bits < 0 || bits > 128 {
		panic("invalid")
	}
	if bits == 128 {
		return RefPrefix{ref, bits}
	}
	val := Uint128(ref)
	val.AndNot(UINT128_1.Lsh(uint(128 - bits)).Sub(UINT128_1))
	return RefPrefix{Ref(val), bits}
}

func (p RefPrefix) String() string {
	return p.ref.StringInPrefix() + "/" + strconv.Itoa(p.bits)
}

func ParseRefPrefix(s string) (RefPrefix, error) {

	ss := strings.Split(s, "/")
	if len(ss) != 2 {
		return RefPrefix{}, errors.New("expected one slash in ref prefix")
	}
	ref, err := ParseRefInPrefix(ss[0])
	if err != nil {
		return RefPrefix{}, err
	}
	bits, err := strconv.Atoi(ss[1])
	if err != nil || bits < 0 || bits > 128 {
		return RefPrefix{}, errors.New("invalid ref prefix length")
	}
	return RefPrefixFrom(ref, bits), nil
}

func MustParseRefPrefix(s string) RefPrefix {

	p, err := ParseRefPrefix(s)
	if err != nil {
		panic("invalid ref prefix")
	}
	return p
}

func RefPrefixSingle(ref Ref) RefPrefix {
	return RefPrefixFrom(ref, 128)
}

func RefPrefixComplete() RefPrefix {
	return RefPrefix{}
}

func (p RefPrefix) Contains(ref Ref) bool {
	return RefPrefixFrom(ref, p.bits).ref == ref
}

func RefPrefixesContain(prefixes []RefPrefix, ref Ref) bool {

	for _, prefix := range prefixes {
		if prefix.Contains(ref) {
			return true
		}
	}
	return false
}
