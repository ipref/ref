/* Copyright (c) 2025 Waldemar Augustyn */

package ref

import (
	"errors"
	"net/netip"
)

type IPPrefix netip.Prefix // .Addr().Zone() must be "", and must be .Masked()

func (p IPPrefix) Addr() IP {
	return IP(netip.Prefix(p).Addr())
}

func (p IPPrefix) Bits() int {
	return netip.Prefix(p).Bits()
}

func (p IPPrefix) SizeBits() int {
	return p.Addr().Len() * 8 - p.Bits()
}

func IPPrefixFrom(ip IP, bits int) IPPrefix {
	return IPPrefix(netip.PrefixFrom(netip.Addr(ip), bits).Masked())
}

func IPPrefixAllVer(ipver int) IPPrefix {
	return IPPrefixFrom(IPZero(IPVerToLen(ipver)), 0)
}

func (p IPPrefix) String() string {
	return netip.Prefix(p).String()
}

func ParseIPPrefix(s string) (IPPrefix, error) {

	p, err := netip.ParsePrefix(s)
	if err != nil {
		return IPPrefix{}, err
	}
	if p.Addr().Zone() != "" {
		return IPPrefix{}, errors.New("IP address prefix may not have zone")
	}
	return IPPrefix(p.Masked()), nil
}

func MustParseIPPrefix(s string) IPPrefix {

	p, err := ParseIPPrefix(s)
	if err != nil {
		panic("invalid IP address prefix")
	}
	return p
}

func IPPrefixSingle(ip IP) IPPrefix {
	return IPPrefixFrom(ip, ip.Len() * 8)
}

func IPPrefixComplete(ipver int) IPPrefix {
	return IPPrefixFrom(IPZero(IPVerToLen(ipver)), 0)
}

func (p IPPrefix) Contains(ip IP) bool {
	return netip.Prefix(p).Contains(netip.Addr(ip))
}

func IPPrefixesContain(prefixes []IPPrefix, ip IP) bool {

	for _, prefix := range prefixes {
		if prefix.Contains(ip) {
			return true
		}
	}
	return false
}

// Returns the 2^l subnets of prefix length 'a.Bits() + l' within a, in order.
// If l is invalid, then nil is returned.
func (a IPPrefix) Subnets(l int) []IPPrefix {

	ip := a.Addr().AsUint128Cast()
	alen := a.Bits()
	if a.Addr().Ver() == 4 {
		alen += 128 - 32
	}
	blen := alen + l
	if l < 0 || l >= 64 || blen > 128 {
		return nil
	}
	prefixes := []IPPrefix {}
	for i := uint64(0); i < 1 << uint64(l); i++ {
		x := ip.Or(Uint128FromUint64(i).Lsh(uint(128 - blen)))
		y := IPFromUint128(x)
		bits := blen
		if a.Addr().Ver() == 4 {
			y = IPFromUint32(x.Uint32())
			bits -= 128 - 32
		}
		prefixes = append(prefixes, IPPrefixFrom(y, bits))
	}
	return prefixes
}
