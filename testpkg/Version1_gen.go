package testpkg

// Version1 is defined as follow:
//   field    bits
//   -----    ----
//   version  4
type Version1 uint8

// Getters.
func (x Version1) version() uint { return uint(x & 0xF) }

// Setters.
func (x *Version1) versionSet(v uint) *Version1 { *x = *x&^0xF | Version1(v)&0xF; return x }
