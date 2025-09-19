/* Copyright (c) 2018-2021, 2025 Waldemar Augustyn */

package ref

import (
	"errors"
	"golang.org/x/sys/unix"
	"net/netip"
)

type IP netip.Addr // IPv4 or IPv6 address; Zone() must be ""

// Tests if the IP is equal to the zero-initialized value. This is distinct from
// the zero IP address (eg. 0.0.0.0 or ::).
func (ip IP) IsZero() bool {
	return ip == IP{}
}

func (ip IP) IsZeroAddr() bool {

	if ip.IsZero() {
		panic("uninitialized")
	}
	s := ip.AsSlice()
	var b byte
	for i := 0; i < len(s); i++ {
		b |= s[i]
	}
	return b == 0
}

func (ip IP) String() string {

	if ip.IsZero() {
		return "(uninitialized)"
	}
	return netip.Addr(ip).String()
}

func ParseIP(s string) (IP, error) {

	ip, err := netip.ParseAddr(s)
	if err != nil {
		return IP{}, err
	}
	if ip.Zone() != "" {
		return IP{}, errors.New("IP address may not have zone")
	}
	return IP(ip), nil
}

func MustParseIP(s string) IP {

	ip, err := ParseIP(s)
	if err != nil {
		panic("invalid IP address")
	}
	return ip
}

// The slice must be 4 or 16 bytes
func IPFromSlice(ip []byte) IP {

	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		panic("invalid IP address")
	}
	return IP(addr)
}

func IPFromUint32(ip uint32) IP {

	var ipb [4]byte
	be.PutUint32(ipb[:], uint32(ip))
	return IP(netip.AddrFrom4(ipb))
}

func IPFromUint128(ip Uint128) IP {

	var ipb [16]byte
	ip.PutBytesBE(ipb[:])
	return IP(netip.AddrFrom16(ipb))
}

func (ip IP) AsSlice() []byte {

	if ip.IsZero() {
		panic("uninitialized")
	}
	return netip.Addr(ip).AsSlice()
}

func (ip IP) AsSlice4() []byte {

	if !ip.Is4() {
		panic("expected IPv4 address")
	}
	return ip.AsSlice()
}

func (ip IP) AsSlice6() []byte {

	if !ip.Is6() {
		panic("expected IPv6 address")
	}
	return ip.AsSlice()
}

func (ip IP) AsUint32() uint32 {

	if ip.IsZero() {
		panic("uninitialized")
	}
	ipb := netip.Addr(ip).As4()
	return uint32(be.Uint32(ipb[:]))
}

func (ip IP) AsUint128() Uint128 {

	if ip.IsZero() {
		panic("uninitialized")
	}
	ipb := netip.Addr(ip).As16()
	return Uint128FromBytesBE(ipb[:])
}

func (ip IP) AsUint128Cast() Uint128 {

	if ip.IsZero() {
		panic("uninitialized")
	}
	if ip.Is4() {
		ipb := netip.Addr(ip).As4()
		return Uint128FromUint32(be.Uint32(ipb[:]))
	} else {
		ipb := netip.Addr(ip).As16()
		return Uint128FromBytesBE(ipb[:])
	}
}

func (ip IP) AsUnixSockaddr() unix.Sockaddr {

	if ip.Is4() {
		var addr unix.SockaddrInet4
		copy(addr.Addr[:], ip.AsSlice())
		return &addr
	} else {
		var addr unix.SockaddrInet6
		copy(addr.Addr[:], ip.AsSlice())
		return &addr
	}
}

func IPFromUnixSockaddr(addr unix.Sockaddr) IP {

	switch a := addr.(type) {
	case *unix.SockaddrInet4:
		return IPFromSlice(a.Addr[:])
	case *unix.SockaddrInet6:
		return IPFromSlice(a.Addr[:])
	default:
		return IP{}
	}
}

func (ip IP) Is4() bool {

	if ip.IsZero() {
		panic("uninitialized")
	}
	return netip.Addr(ip).Is4()
}

func (ip IP) Is6() bool {
	return !ip.Is4()
}

func (ip IP) Is4In6() bool {
	return netip.Addr(ip).Is4In6()
}

func (ip IP) IsLinkLocal() bool {
	return !ip.Is4In6() && (netip.Addr(ip).IsLinkLocalUnicast() ||
		netip.Addr(ip).IsLinkLocalMulticast())
}

func (ip IP) IsGlobalUnicast() bool {
	return !ip.Is4In6() && netip.Addr(ip).IsGlobalUnicast()
}

func (ip IP) As4In6() IP {

	if ip.IsZero() || ip.Ver() == 6 {
		return ip
	}
	var bs [16]byte
	bs[10] = 0xff
	bs[11] = 0xff
	copy(bs[12:], ip.AsSlice())
	return IPFromSlice(bs[:])
}

func (ip IP) Un4In6() IP {

	if ip.IsZero() {
		return ip
	}
	return IP(netip.Addr(ip).Unmap())
}

func (ip IP) Len() int {

	if ip.Is4() {
		return 4
	} else {
		return 16
	}
}

func (ip IP) Ver() int {

	if ip.Is4() {
		return 4
	} else {
		return 6
	}
}

func IPLenToVer(l int) int {

	switch l {
	case 4:
		return 4
	case 16:
		return 6
	default:
		return 0
	}
}

func IPVerToLen(ver int) int {

	switch ver {
	case 4:
		return 4
	case 6:
		return 16
	default:
		return 0
	}
}

func (ip IP) ByteFromEnd(i int) byte {

	bs := ip.AsSlice()
	return bs[len(bs) - i - 1]
}

func (a IP) Or(b IP) IP {

	as := a.AsSlice()
	bs := b.AsSlice()
	if len(as) != len(bs) {
		panic("IP addresses are different length")
	}
	var cs [16]byte
	for i := 0; i < len(as); i++ {
		cs[i] = as[i] | bs[i]
	}
	return IPFromSlice(cs[:len(as)])
}

func (a IP) And(b IP) IP {

	as := a.AsSlice()
	bs := b.AsSlice()
	if len(as) != len(bs) {
		panic("IP addresses are different length")
	}
	var cs [16]byte
	for i := 0; i < len(as); i++ {
		cs[i] = as[i] & bs[i]
	}
	return IPFromSlice(cs[:len(as)])
}

func (a IP) XOr(b IP) IP {

	as := a.AsSlice()
	bs := b.AsSlice()
	if len(as) != len(bs) {
		panic("IP addresses are different length")
	}
	var cs [16]byte
	for i := 0; i < len(as); i++ {
		cs[i] = as[i] ^ bs[i]
	}
	return IPFromSlice(cs[:len(as)])
}

func (a IP) Not() IP {

	as := a.AsSlice()
	var bs [16]byte
	for i := 0; i < len(as); i++ {
		bs[i] = ^as[i]
	}
	return IPFromSlice(bs[:len(as)])
}

func (a IP) Add(b IP) IP {

	as := a.AsSlice()
	bs := b.AsSlice()
	if len(as) != len(bs) {
		panic("IP addresses are different length")
	}
	var cs [16]byte
	var carry uint16
	for i := len(as) - 1; i >= 0; i-- {
		carry += uint16(as[i]) + uint16(bs[i])
		cs[i] = uint8(carry)
		carry >>= 8
	}
	return IPFromSlice(cs[:len(as)])
}

func (a IP) Compare(b IP) int {

	switch {

	case a.Is4() && b.Is4():

		ai := a.AsUint32()
		bi := b.AsUint32()
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
		return 0

	case a.Is6() && b.Is6():

		as := a.AsSlice()
		bs := b.AsSlice()
		for i := 0; i < 16; i++ {
			if as[i] < bs[i] {
				return -1
			}
			if as[i] > bs[i] {
				return 1
			}
		}
		return 0

	case a.Is4() && b.Is6():

		return -1

	case a.Is6() && b.Is4():

		return 1

	default:

		panic("unexpected")
	}

}

func IPBits(l, n int) IP {

	if l != 4 && l != 16 {
		panic("invalid IP address length")
	}
	var bs [16]byte
	for i := 0; i < l && n > 0; i++ {
		bs[i] = 0xff
		if n < 8 {
			bs[i] <<= 8 - n
		}
		n -= 8
	}
	return IPFromSlice(bs[:l])
}

func IPNum(l int, n uint32) IP {

	if l != 4 && l != 16 {
		panic("invalid IP address length")
	}
	var bs [16]byte
	be.PutUint32(bs[12:16], n)
	return IPFromSlice(bs[16 - l:])
}

func IPZero(l int) IP {
	return IPNum(l, 0)
}
