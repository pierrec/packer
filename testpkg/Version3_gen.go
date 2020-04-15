package testpkg

type Version3 uint64

// Getters.
func (x Version3) version() uint    { return uint(x & 0xF) }
func (x Version3) Len() int         { return int(x >> 4 & 0xFFFF) }
func (x Version3) reserved() int    { return int(x >> 20 & 0x7) }
func (x Version3) Checksum() uint32 { return uint32(x >> 23 & 0xFFFFFFFF) }

// Setters.
func (x *Version3) versionSet(v uint) *Version3 { *x = *x&^0xF | Version3(v)&0xF; return x }
func (x *Version3) LenSet(v int) *Version3 {
	*x = *x&^(0xFFFF<<4) | (Version3(v) & 0xFFFF << 4)
	return x
}
func (x *Version3) reservedSet(v int) *Version3 {
	*x = *x&^(0x7<<20) | (Version3(v) & 0x7 << 20)
	return x
}
func (x *Version3) ChecksumSet(v uint32) *Version3 {
	*x = *x&^(0xFFFFFFFF<<23) | (Version3(v) & 0xFFFFFFFF << 23)
	return x
}
