// Copyright 2012 <MortalSkulD@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

type test struct {
	site_id, site_salt, encrypt_key, encrypt_salt string
	result                                        string
}

var golden = []test{
	{"id0", "site0", "111", "fuckcsdn", "2jNXfMGoXTSK9pFS"},
	{"id1", "site0", "111", "fuckcsdn", "UndfvMK289BQoU8f"},
	{"id3", "site0", "111", "fuckcsdn", "N9V9FyMJ8tkScBN5"},

	{"id0", "site1", "111", "fuckcsdn", "ykoUohFBBjGtxz7V"},
	{"id1", "site1", "111", "fuckcsdn", "641gNmCY9YFNAQ1p"},

	{"id0", "site2", "111", "fuckcsdn", "4HLACkWRCyDHtqtx"},
	{"id0", "site3", "111", "fuckcsdn", "3eHtu74rMFdeRaVk"},
	{"id0", "site4", "111", "fuckcsdn", "5DSxs623Rciz7bab"},
	{"id0", "site5", "111", "fuckcsdn", "3cfiPrcjdrhwAgM1"},
	{"id0", "site6", "111", "fuckcsdn", "5Las25BPXCjvtywo"},
	{"id0", "site7", "111", "fuckcsdn", "1GK1x3GnRxLSH6DT"},
	{"id0", "site8", "111", "fuckcsdn", "3VqQSgsRRQTeR6vL"},
	{"id0", "site9", "111", "fuckcsdn", "5hHdKchVRPeJkFjU"},

	// ----------------------------------------------------

	{"id0", "site0", "abc", "fuckcsdn", "3d8aTVUZaYhqyxbi"},
	{"id1", "site0", "abc", "fuckcsdn", "5f1qgGWCXj9vh6r3"},
	{"id3", "site0", "abc", "fuckcsdn", "5RkzM4wkGqqqV6dF"},

	{"id0", "site1", "abc", "fuckcsdn", "3LhwqzbNSaMjq33x"},
	{"id1", "site1", "abc", "fuckcsdn", "5NBT7AG9fvj18nRb"},

	{"id0", "site2", "abc", "fuckcsdn", "3L6S4GcMomW92nWH"},
	{"id0", "site3", "abc", "fuckcsdn", "24g39Y5Rc8r2fW1J"},
	{"id0", "site4", "abc", "fuckcsdn", "66X8A3Qwjy9RJnjg"},
	{"id0", "site5", "abc", "fuckcsdn", "2GBeAPF4Ar9JhVDz"},
	{"id0", "site6", "abc", "fuckcsdn", "4FCEsEVqWHFiDFUN"},
	{"id0", "site7", "abc", "fuckcsdn", "3WkbTtfgivycbeGt"},
	{"id0", "site8", "abc", "fuckcsdn", "5v5i7bnwdoGbWUG4"},
	{"id0", "site9", "abc", "fuckcsdn", "4PhkQ28EWm39MSN1"},
}

func TestPwdGen(t *testing.T) {
	for _, g := range golden {
		s := PwdGen(g.site_id, g.site_salt, g.encrypt_key, g.encrypt_salt)
		if s != g.result {
			t.Errorf("Bad result: Need=%v, Got=%v", g.result, s)
		}
	}
}

func ExamplePwdGen() {
	fmt.Println(PwdGen("id0", "site0", "111", "fuckcsdn"))
	// Output: 2jNXfMGoXTSK9pFS
}
