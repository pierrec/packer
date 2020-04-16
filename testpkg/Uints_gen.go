package testpkg

// Uints is defined as follow:
//   field     bits
//   -----     ----
//   Uint8     8
//   Uint16    16
//   Uint32    32
//   (unused)  8
type Uints uint64

// Getters.
func (x Uints) Uint8() uint8   { return uint8(x & 0xFF) }
func (x Uints) Uint16() uint16 { return uint16(x >> 8 & 0xFFFF) }
func (x Uints) Uint32() uint32 { return uint32(x >> 24 & 0xFFFFFFFF) }

// Setters.
func (x *Uints) Uint8Set(v uint8) *Uints   { *x = *x&^0xFF | Uints(v)&0xFF; return x }
func (x *Uints) Uint16Set(v uint16) *Uints { *x = *x&^(0xFFFF<<8) | (Uints(v) & 0xFFFF << 8); return x }
func (x *Uints) Uint32Set(v uint32) *Uints {
	*x = *x&^(0xFFFFFFFF<<24) | (Uints(v) & 0xFFFFFFFF << 24)
	return x
}
