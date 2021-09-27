package hexdunk

import "testing"

//actually, how to properly test a widget??

type hdtest struct {
	addr int64
	sz   int
}

var hdtests = []hdtest{
	{0, 0},
	{0xf, 1},
	{0xa, 1},
	{0x01, 1},
	{0xff, 2},
	{0xffffff, 6},
	{0xdeadbeef, 8},
	{0x1111000ffffff, 13},
}

func TestNumHexDigits(t *testing.T) {
	for _, tst := range hdtests {
		if numHexDigits(tst.addr) != tst.sz {
			t.Fatalf("numHexDigits(%x) = %d, expected: %d\n", tst.addr, numHexDigits(tst.addr), tst.sz)
		}
	}
}
