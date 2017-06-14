package soze

// xorInto xors the src bytes into dst bytes. The destination is assumed to have enough
// space.
func xorInto(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		dst[i] ^= src[i]
	}
}
