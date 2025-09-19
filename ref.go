/* Copyright (c) 2025 Waldemar Augustyn */

package ref

import (
	"errors"
	"fmt"
	"strings"
)

// The zero ref is Ref{}
type Ref Uint128

func (ref Ref) IsZero() bool {
	return ref == Ref{}
}

func ParseRef(str string) (Ref, error) {
	return parse_ref(str, false)
}

func ParseRefInPrefix(str string) (Ref, error) {
	return parse_ref(str, true)
}

func parse_ref(str string, cidr bool) (Ref, error) {

	if !cidr && !strings.Contains(str, "-") {
		if val, ok := ParseUint128(str, 10); ok {
			return Ref(val), nil
		}
		return Ref{}, errors.New("invalid format")
	}
	ss := strings.Split(str, "--")
	if len(ss) == 1 {
		val, bits, err := parse_ref_comps(ss[0])
		if cidr && bits != 128 {
			return Ref{}, errors.New("ref in prefix needs '--' unless it is full-length")
		}
		return Ref(val), err
	}
	if len(ss) == 2 {
		if len(ss[0]) == 0 {
			return Ref{}, errors.New("ref cannot have leading '--'")
		}
		if len(ss[1]) == 0 {
			val, bits, err := parse_ref_comps(ss[0])
			return Ref(val.Lsh(128 - bits)), err
		}
		a, abits, err := parse_ref_comps(ss[0])
		if err != nil {
			return Ref{}, err
		}
		b, bbits, err := parse_ref_comps(ss[1])
		if err != nil {
			return Ref{}, err
		}
		if abits + bbits >= 128 {
			return Ref{}, errors.New("ref is larger than 128 bits")
		}
		return Ref(a.Lsh(128 - abits).Or(b)), nil
	}
	return Ref{}, errors.New("ref contains more than one '--'")
}

func parse_ref_comps(str string) (Uint128, uint, error) {

	var n Uint128
	var bits uint
	for _, comp := range strings.Split(str, "-") {
		if bits >= 128 {
			return Uint128{}, 0, errors.New("ref is larger than 128 bits")
		}
		if len(comp) > 4 {
			return Uint128{}, 0, errors.New("invalid format")
		}
		val, ok := ParseUint128(comp, 16)
		if !ok || val.Cmp(Uint128FromUint64(1 << 16)) >= 0 {
			return Uint128{}, 0, errors.New("invalid format")
		}
		n = n.Lsh(16).Or(val)
		bits += 16
	}
	return n, bits, nil
}

func MustParseRef(str string) Ref {

	ref, err := ParseRef(str)
	if err != nil {
		panic(fmt.Sprintf("invalid ref: %q", str))
	}
	return ref
}

func ParseIpRef(str string) (ipref IpRef, err error) {

	ip, ref, found := strings.Cut(str, "+")
	if !found {
		return IpRef{}, errors.New("invalid format (missing '+')")
	}
	ipref.IP, err = ParseIP(strings.TrimSpace(ip))
	if err != nil {
		return
	}
	ipref.Ref, err = ParseRef(strings.TrimSpace(ref))
	return
}

func MustParseIpRef(str string) IpRef {

	ipref, err := ParseIpRef(str)
	if err != nil {
		panic(fmt.Sprintf("invalid ipref: %q", str))
	}
	return ipref
}

func (ref Ref) String() string {

	val := Uint128(ref)
	if val.Cmp(Uint128FromUint64(1 << 16)) < 0 {
		return val.String()
	}
	var s string
	bits := 0
	for !val.IsZero() {
		if bits != 0 {
			s = "-" + s
		}
		s = val.And(Uint128FromUint64(0xffff)).FormatHex() + s
		val = val.Rsh(16)
		bits += 16
	}
	if bits <= 16 {
		s = "0-" + s
	}
	return s
}

func (ref Ref) StringInPrefix() string {

	val := Uint128(ref)
	if val.IsZero() {
		return "0--"
	}
	var s string
	bits := 0
	for !val.IsZero() {
		if bits != 0 {
			s += "-"
		}
		s += val.Rsh(128 - 16).FormatHex()
		val = val.Lsh(16)
		bits += 16
	}
	if bits != 128 {
		s += "--"
	}
	return s
}

func (ref Ref) AsSliceBE() []byte {
	return Uint128(ref).AsSliceBE()
}

func RefFromBytesBE(src []byte) Ref {
	return Ref(Uint128FromBytesBE(src))
}
