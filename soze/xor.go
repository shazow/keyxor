package soze

// xorInto xors the src bytes into dst bytes. The destination is assumed to have enough
// space.
func xorInto(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		dst[i] ^= src[i]
	}
}

// zero will zero out every byte in slice buf.
func zero(buf []byte) {
	for i := 0; i < len(buf); i++ {
		buf[i] = 0
	}
}
