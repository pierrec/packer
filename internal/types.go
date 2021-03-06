// Code generated by `gen.exe`. DO NOT EDIT.

package internal

// UintEntry is defined as follow:
//   field     bits
//   -----     ----
//   Num       3
//   A         3
//   B         3
//   C         3
//   D         3
//   E         3
//   F         3
//   G         3
//   (unused)  8
type UintEntry uint32

// Getters.
func (x UintEntry) Num() uint8 { return uint8(x&0x7) }
func (x UintEntry) A() uint8 { return uint8(x>>3&0x7) }
func (x UintEntry) B() uint8 { return uint8(x>>6&0x7) }
func (x UintEntry) C() uint8 { return uint8(x>>9&0x7) }
func (x UintEntry) D() uint8 { return uint8(x>>12&0x7) }
func (x UintEntry) E() uint8 { return uint8(x>>15&0x7) }
func (x UintEntry) F() uint8 { return uint8(x>>18&0x7) }
func (x UintEntry) G() uint8 { return uint8(x>>21&0x7) }

// Setters.
func (x *UintEntry) NumSet(v uint8) *UintEntry { *x = *x&^0x7 | UintEntry(v)&0x7; return x }
func (x *UintEntry) ASet(v uint8) *UintEntry { *x = *x&^(0x7<<3) | (UintEntry(v)&0x7<<3); return x }
func (x *UintEntry) BSet(v uint8) *UintEntry { *x = *x&^(0x7<<6) | (UintEntry(v)&0x7<<6); return x }
func (x *UintEntry) CSet(v uint8) *UintEntry { *x = *x&^(0x7<<9) | (UintEntry(v)&0x7<<9); return x }
func (x *UintEntry) DSet(v uint8) *UintEntry { *x = *x&^(0x7<<12) | (UintEntry(v)&0x7<<12); return x }
func (x *UintEntry) ESet(v uint8) *UintEntry { *x = *x&^(0x7<<15) | (UintEntry(v)&0x7<<15); return x }
func (x *UintEntry) FSet(v uint8) *UintEntry { *x = *x&^(0x7<<18) | (UintEntry(v)&0x7<<18); return x }
func (x *UintEntry) GSet(v uint8) *UintEntry { *x = *x&^(0x7<<21) | (UintEntry(v)&0x7<<21); return x }
