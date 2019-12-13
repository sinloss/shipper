package util

const chars = "0123456789abcdef"

// Hexchar produces the high and low bytes of the hex char representing
// the given byte
func Hexchar(b byte) (high byte, low byte) {
	return chars[b>>4], chars[b&0xf]
}
