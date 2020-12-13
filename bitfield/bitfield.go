package bitfield

// A BitField represents the pieces that a peer has
type BitField []byte

// HasPiece tells if a bitfield has a particular index set
func (bf BitField) HasPiece (index int) bool {
	byteIndex := index >> 3		// index / 8
	if byteIndex < 0 || byteIndex >= len(bf) { return false }
	offset := index & 0b111		// index % 8
	return bf[byteIndex] >> (7 - offset) & 1 != 0
}

// SetPiece sets a bit in the bitfield
func (bf BitField) SetPiece (index int) {
	byteIndex := index >> 3		// index / 8
	if byteIndex < 0 || byteIndex >= len(bf) { return }
	offset := index & 0b111		// index % 8
	bf[byteIndex] |= 1 << (7 - offset)	// set bit
}
