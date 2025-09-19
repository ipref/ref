/* Copyright (c) 2025 Waldemar Augustyn */

package newv1

import . "github.com/ipref/ref"

const ( // v1 constants

	V1_AREC_MIN_LEN = 4 + 4  + 4  + 4  + 16 // header + ea + ip + gw + ref.h + ref.l
	V1_AREC_MAX_LEN = 4 + 16 + 16 + 16 + 16
	// v1 header offsets
	V1_RESERVED = 4
)

func AddrRecEncodedLen(ea_iplen, gw_iplen int) int {
	return 4 + ea_iplen * 2 + gw_iplen + 16 // header + ea + ip + gw + ref.h + ref.l
}

func AddrRecCheck(arec []byte) (ok bool, length, ea_iplen, gw_iplen int) {

	if len(arec) < V1_AREC_MIN_LEN {
		return
	}
	ea_iplen = int(arec[1])
	gw_iplen = int(arec[3])
	if ea_iplen == 0 || IPVerToLen(int(arec[0])) != ea_iplen {
		return
	}
	if gw_iplen == 0 || IPVerToLen(int(arec[2])) != gw_iplen {
		return
	}
	length = AddrRecEncodedLen(ea_iplen, gw_iplen)
	if len(arec) < length {
		return
	}
	ok = true
	return
}

func AddrRecSlices(arec []byte, ea_iplen, gw_iplen int) (
	length int, ea, ip, gw, ref []byte) {

	i := 4
	ea = arec[i : i + ea_iplen]
	i += ea_iplen
	ip = arec[i : i + ea_iplen]
	i += ea_iplen
	gw = arec[i : i + gw_iplen]
	i += gw_iplen
	ref = arec[i : i + 16]
	length = i + 16
	return
}

func AddrRecEncode(arecb []byte, arec AddrRec) int {

	if arec.EA.Len() != arec.IP.Len() {
		panic("unexpected")
	}
	ea_iplen := arec.EA.Len()
	gw_iplen := arec.GW.Len()
	arecb[0] = byte(IPLenToVer(ea_iplen))
	arecb[1] = byte(ea_iplen)
	arecb[2] = byte(IPLenToVer(gw_iplen))
	arecb[3] = byte(gw_iplen)
	length, eab, ipb, gwb, refb := AddrRecSlices(arecb, ea_iplen, gw_iplen)
	copy(eab, arec.EA.AsSlice())
	copy(ipb, arec.IP.AsSlice())
	copy(gwb, arec.GW.AsSlice())
	Uint128(arec.Ref).PutBytesBE(refb)
	return length
}

func AddrRecEncodedLenOf(arec AddrRec) int {

	if arec.EA.Len() != arec.IP.Len() {
		panic("unexpected")
	}
	return AddrRecEncodedLen(arec.EA.Len(), arec.GW.Len())
}

func AddrRecAsSlice(arec AddrRec) []byte {

	arecb := make([]byte, AddrRecEncodedLenOf(arec))
	AddrRecEncode(arecb, arec)
	return arecb
}

func AddrRecDecode(arecb []byte) (bool, int, AddrRec) {

	ok, _, ea_iplen, gw_iplen := AddrRecCheck(arecb)
	if !ok {
		return false, 0, AddrRec{}
	}
	length, eab, ipb, gwb, refb := AddrRecSlices(arecb, ea_iplen, gw_iplen)
	return true, length, AddrRec{
		EA: IPFromSlice(eab),
		IP: IPFromSlice(ipb),
		GW: IPFromSlice(gwb),
		Ref: Ref(Uint128FromBytesBE(refb)),
	}
}
