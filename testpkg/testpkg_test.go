package testpkg

import "testing"

func TestGen(t *testing.T) {
	_i := func(v ...int) []int { return v }
	type tcase struct {
		label string
		x     []int
		run   func(int) int
	}
	for _, tc := range []tcase{
		{"Version1",
			_i(10),
			func(x int) int {
				var v Version1
				return int(v.versionSet(uint(x)).version())
			}},
		{"Version2",
			_i(10, 255, 256, 1000, 1024, 2000, 1<<7),
			func(x int) int {
				var v Version2
				return v.LenSet(x).Len()
			}},
		{"Version3",
			_i(10, 255, 256, 1000, 1024, 2000, 1<<7),
			func(x int) int {
				var v Version3
				return v.LenSet(x).Len()
			}},
		{"Version3",
			_i(0xAAAAAAAA, 0xFFFFFFFF, 0x12345678, 0xDEADBEEF, 0xC001CAFE),
			func(x int) int {
				var v Version3
				v.LenSet(15)
				return int(v.ChecksumSet(uint32(x)).Checksum())
			}},
	} {
		t.Run(tc.label, func(t *testing.T) {
			for _, x := range tc.x {
				if got, want := tc.run(x), x; got != want {
					t.Fatalf("got %d; want %d", got, want)
				}
			}
		})
	}
}
