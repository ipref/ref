/* Copyright (c) 2018-2019 Waldemar Augustyn */

package ref

import (
	"fmt"
	"math/bits"
	"regexp"
	"strconv"
	"strings"
)

type Ref struct {
	H uint64
	L uint64
}

var re_hexref *regexp.Regexp
var re_decref *regexp.Regexp
var re_dotref *regexp.Regexp

func init() {
	re_hexref = regexp.MustCompile(`^[0-9a-fA-F]+([-][0-9a-fA-F]+)*$`)
	re_decref = regexp.MustCompile(`^[0-9]+([,][0-9]+)+$`)
	re_dotref = regexp.MustCompile(`^([1-9]?[0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])([.]([1-9]?[0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5]))+$`)
}

func (ref *Ref) isZero() bool {
	return ref.H == 0 && ref.L == 0
}

// print ref as dash separated hex quads: 2f-4883-0005-2a1b
func (ref *Ref) String() string {

	var sb strings.Builder

	var writequads = func(word uint64) {
		for ii := 0; ii < 4; ii++ {
			word = bits.RotateLeft64(word, 16)
			if sb.Len() == 0 {
				if quad := word & 0xffff; quad != 0 {
					sb.WriteString(fmt.Sprintf("%x", quad))
				}
			} else {
				sb.WriteString(fmt.Sprintf("-%04x", word&0xffff))
			}
		}
	}

	writequads(ref.H)
	writequads(ref.L)

	return sb.String()
}

// parse reference
func Parse(sss string) (Ref, error) {

	var ref Ref
	var err error
	var val uint64 // go does not allow ref.L, err := something(), need intermediate variable

	// hex

	if re_hexref.MatchString(sss) {

		hex := strings.Replace(sss, "-", "", -1)
		hexlen := len(hex)
		if hexlen < 17 {
			ref.H = 0
			val, err = strconv.ParseUint(hex, 16, 64)
			if err != nil {
				return ref, err
			}
			ref.L = val
			return ref, nil
		} else {
			val, err = strconv.ParseUint(hex[:hexlen-16], 16, 64)
			if err != nil {
				return ref, err
			}
			ref.H = val
			val, err = strconv.ParseUint(hex[hexlen-16:hexlen], 16, 64)
			if err != nil {
				return ref, err
			}
			ref.L = val
			return ref, nil
		}
	}

	// decimal

	if re_decref.MatchString(sss) {

		decstr := strings.Replace(sss, ",", "", -1)
		ref.H = 0
		val, err = strconv.ParseUint(decstr, 10, 64)
		if err != nil {
			return ref, err
		}
		ref.L = val
		return ref, nil
	}

	// dotted decimal

	if re_dotref.MatchString(sss) {
		dot := strings.Split(sss, ".")
		dotlen := len(dot)
		for ii := 0; ii < dotlen; ii++ {
			dec, err := strconv.ParseUint(dot[ii], 10, 8)
			if err != nil {
				return ref, err
			}
			if ii < (dotlen - 8) {
				ref.H <<= 8
				ref.H += uint64(dec)
			} else {
				ref.L <<= 8
				ref.L += uint64(dec)
			}
		}
		return ref, nil
	}

	return ref, fmt.Errorf("unrecognized format")
}
