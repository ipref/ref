/* Copyright (c) 2025 Waldemar Augustyn */

package oldv1

import (
	. "github.com/ipref/ref"
	"encoding/binary"
)

const ( // v1 constants

	V1_AREC_MAX_LEN = 16 + 16 + 16 + 16 // ea + ip + gw + ref.h + ref.l
	// v1 header offsets
	V1_IPVER    = 4 // high nibble is the ea IP ver, low nibble is gw IP ver
	V1_RESERVED = 5
)

var be = binary.BigEndian

func AddrRecEncodedLen(ea_iplen, gw_iplen int) int {
	return ea_iplen * 2 + gw_iplen + 16 // ea + ip + gw + ref.h + ref.l
}

func AddrRecSlices(ea_iplen, gw_iplen int, arec []byte) (ea, ip, gw, refh, refl []byte) {
	i := 0
	ea = arec[i : i + ea_iplen]
	i += ea_iplen
	ip = arec[i : i + ea_iplen]
	i += ea_iplen
	gw = arec[i : i + gw_iplen]
	i += gw_iplen
	refh = arec[i : i + 8]
	i += 8
	refl = arec[i : i + 8]
	return
}

func AddrRecEncode(arecb []byte, arec AddrRec) {
	if arec.EA.Len() != arec.IP.Len() {
		panic("unexpected")
	}
	eab, ipb, gwb, refhb, reflb := AddrRecSlices(arec.EA.Len(), arec.GW.Len(), arecb)
	copy(eab, arec.EA.AsSlice())
	copy(ipb, arec.IP.AsSlice())
	copy(gwb, arec.GW.AsSlice())
	be.PutUint64(refhb, arec.Ref.H)
	be.PutUint64(reflb, arec.Ref.L)
}

func AddrRecEncodedLenOf(arec AddrRec) int {
	if arec.EA.Len() != arec.IP.Len() {
		panic("unexpected")
	}
	return AddrRecEncodedLen(arec.EA.Len(), arec.GW.Len())
}

func AddrRecDecode(ea_iplen, gw_iplen int, arecb []byte) (arec AddrRec) {
	eab, ipb, gwb, refhb, reflb := AddrRecSlices(ea_iplen, gw_iplen, arecb)
	arec.EA = IPFromSlice(eab)
	arec.IP = IPFromSlice(ipb)
	arec.GW = IPFromSlice(gwb)
	arec.Ref.H = be.Uint64(refhb)
	arec.Ref.L = be.Uint64(reflb)
	return
}
