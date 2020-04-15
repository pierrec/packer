package testpkg

// Version3 is defined as follow:
//   field     bits
//   -----     ----
//   version   4
//   _         8
//   Len       16
//   _         4
//   Checksum  32
type Version3 uint64

// Getters.
func (x Version3) version() uint    { return uint(x & 0xF) }
func (x Version3) Len() int         { return int(x >> 12 & 0xFFFF) }
func (x Version3) Checksum() uint32 { return uint32(x >> 32 & 0xFFFFFFFF) }

// Setters.
func (x *Version3) versionSet(v uint) *Version3 { *x = *x&^0xF | Version3(v)&0xF; return x }
func (x *Version3) LenSet(v int) *Version3 {
	*x = *x&^(0xFFFF<<12) | (Version3(v) & 0xFFFF << 12)
	return x
}
func (x *Version3) ChecksumSet(v uint32) *Version3 {
	*x = *x&^(0xFFFFFFFF<<32) | (Version3(v) & 0xFFFFFFFF << 32)
	return x
}
