package testpkg

// Version2 is defined as follow:
//   field     bits
//   -----     ----
//   version   4
//   Len       16
//   (unused)  12
type Version2 uint32

// Getters.
func (x Version2) version() uint { return uint(x & 0xF) }
func (x Version2) Len() int      { return int(x >> 4 & 0xFFFF) }

// Setters.
func (x *Version2) versionSet(v uint) *Version2 { *x = *x&^0xF | Version2(v)&0xF; return x }
func (x *Version2) LenSet(v int) *Version2 {
	*x = *x&^(0xFFFF<<4) | (Version2(v) & 0xFFFF << 4)
	return x
}
