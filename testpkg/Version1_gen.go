package testpkg

// Version1 is defined as follow:
//   field     bits
//   -----     ----
//   version   4
//   flag      1
//   (unused)  3
type Version1 uint8

// Getters.
func (x Version1) version() uint { return uint(x & 0xF) }
func (x Version1) flag() bool    { return x>>4&1 != 0 }

// Setters.
func (x *Version1) versionSet(v uint) *Version1 { *x = *x&^0xF | Version1(v)&0xF; return x }
func (x *Version1) flagSet(v bool) *Version1 {
	const b = 1 << 4
	if v {
		*x = *x&^b | b
	} else {
		*x &^= b
	}
	return x
}
