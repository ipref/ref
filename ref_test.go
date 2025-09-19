/* Copyright (c) 2025 Waldemar Augustyn */

package ref

import "testing"

func TestRefParsing(t *testing.T) {

	test_cases := []struct {
		str       string
		formatted string
		valid     bool
		ref       Ref
	}{
		{"0-0", "0", true, Ref{}},
		{"1", "1", true, Ref(UINT128_1)},
		{"12", "12", true, Ref(Uint128FromUint64(12))},
		{"0-12", "18", true, Ref(Uint128FromUint64(0x12))},
		{"12-0", "12-0", true, Ref(Uint128FromUint64(0x12 << 16))},
		{"a0--12", "a0-0-0-0-0-0-0-12", true,
			Ref(Uint128FromUint64(0xa0).Lsh(112).Or(Uint128FromUint64(0x12)))},
		{"0--", "0", true, Ref{}},
		{"22-33--", "22-33-0-0-0-0-0-0", true,
			Ref(Uint128FromUint64(0x220033).Lsh(96))},
		{"123456789012345678901234567890123456789", "5ce0-e9a5-6015-fec5-aadf-a328-ae39-8115", true,
			Ref(MustParseUint128("123456789012345678901234567890123456789", 10))},
		{"33-0012--1", "33-12-0-0-0-0-0-1", true,
			Ref(Uint128FromUint64(0x330012).Lsh(96).Or(UINT128_1))},
		{"33-0012--a-12", "33-12-0-0-0-0-a-12", true,
			Ref(Uint128FromUint64(0x330012).Lsh(96).Or(Uint128FromUint64(0xa0012)))},
		{"--0", "", false, Ref{}},
		{"12ab", "", false, Ref{}},
	}

	for i, c := range test_cases {

		ref, err := ParseRef(c.str)
		if c.valid {
			if err != nil {
				t.Errorf("case %v: unexpected error parsing ref %q: %v", i, c.str, err)
			}
			if ref != c.ref {
				t.Errorf("case %v: parsing %q: expected %v, got %v", i, c.str, c.ref, ref)
			}
			if s := ref.String(); s != c.formatted {
				t.Errorf("case %v: expected %q, got %q", i, c.formatted, s)
			}
		} else {
			if err == nil {
				t.Errorf("case %v: expected error when parsing %q", i, c.str)
			}
		}
	}
}
