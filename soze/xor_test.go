package soze

import (
	"bytes"
	"testing"
)

func TestXorInto(t *testing.T) {
	tests := []struct {
		dst    []byte
		src    []byte
		result []byte
	}{
		{
			[]byte{1, 1},
			[]byte{0, 1},
			[]byte{1, 0},
		},
		{
			[]byte{1, 2, 3, 4},
			[]byte{2, 3, 4, 5},
			[]byte{3, 1, 7, 1},
		},
	}

	for i, tc := range tests {
		xorInto(tc.dst, tc.src)
		if !bytes.Equal(tc.dst, tc.result) {
			t.Errorf("case %d: got %d; want %d", i, tc.dst, tc.result)
		}
	}
}
