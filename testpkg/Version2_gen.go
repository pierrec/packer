package testpkg

// Version2 is defined as follow:
//   field     bits
//   -----     ----
//   version   4
//   flag      1
//   Len       16
//   (unused)  11
type Version2 uint32

// Getters.
func (x Version2) version() uint { return uint(x & 0xF) }
func (x Version2) flag() bool    { return x>>4&1 != 0 }
func (x Version2) Len() int      { return int(x >> 5 & 0xFFFF) }

// Setters.
func (x *Version2) versionSet(v uint) *Version2 { *x = *x&^0xF | Version2(v)&0xF; return x }
func (x *Version2) flagSet(v bool) *Version2 {
	const b = 1 << 4
	if v {
		*x = *x&^b | b
	} else {
		*x &^= b
	}
	return x
}
func (x *Version2) LenSet(v int) *Version2 {
	*x = *x&^(0xFFFF<<5) | (Version2(v) & 0xFFFF << 5)
	return x
}
