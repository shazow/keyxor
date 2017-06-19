package soze

import (
	"bytes"
	"io"
	"testing"
)

// Black box test
func TestSoze(t *testing.T) {
	input := "hello world"

	inBuf := bytes.NewBufferString(input)
	n := inBuf.Len()
	var out1, out2 bytes.Buffer

	if err := Split(inBuf, []io.Writer{&out1, &out2}); err != nil {
		t.Error(err)
	}

	if want, got := n, out1.Len(); want != got {
		t.Errorf("wrong length: want: %d; got %d", want, got)
	}

	if input == out1.String() {
		t.Errorf("input should not match output: %q", out1.String())
	}

	if out1.String() == out2.String() {
		t.Errorf("outputs should not match: %q", out1.String())
	}
	if out1.Len() != out2.Len() {
		t.Errorf("mismatched split size: %d != %d", out1.Len(), out2.Len())
	}

	var result bytes.Buffer
	if err := Merge(&result, []io.Reader{&out1, &out2}); err != nil {
		t.Error(err)
	}

	if want, got := input, result.String(); want != got {
		t.Errorf("xor'd result does not match: want %q; got %q", want, got)
	}
}

func TestMismatchedSize(t *testing.T) {
	var out bytes.Buffer
	err := Merge(&out, []io.Reader{bytes.NewBuffer([]byte{1, 2}), bytes.NewBuffer([]byte{1, 2, 3})})
	if err == nil {
		t.Fatal("missing mismatch error")
	}
	actualErr, ok := err.(ErrSizeMismatch)
	if !ok {
		t.Errorf("wrong error: %T", actualErr)
	}
	if actualErr.Actual != 3 || actualErr.Expected != 2 {
		t.Errorf("wrong error values: %q", actualErr)
	}
}

func TestMerge(t *testing.T) {
	k1 := []byte{1, 1, 1, 1}
	k2 := []byte{2, 2, 2, 2}

	var out bytes.Buffer
	if err := Merge(&out, []io.Reader{bytes.NewBuffer(k1), bytes.NewBuffer(k2)}); err != nil {
		t.Error(err)
	}

	if got, want := out.Bytes(), []byte{3, 3, 3, 3}; !bytes.Equal(got, want) {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestSplit(t *testing.T) {
	key := []byte{1, 1, 1, 1}
	r := []byte{2, 2, 2, 2, 4, 4, 4, 4}

	var out1, out2, out3 bytes.Buffer
	if err := splitWithRand(bytes.NewBuffer(key), []io.Writer{&out1, &out2, &out3}, bytes.NewBuffer(r)); err != nil {
		t.Error(err)
	}

	// out1 should contain 1^2^4 = 7
	if got, want := out1.Bytes(), []byte{7, 7, 7, 7}; !bytes.Equal(got, want) {
		t.Errorf("out1: got %b; want %b", got, want)
	}

	// out2 should contain 2
	if got, want := out2.Bytes(), []byte{2, 2, 2, 2}; !bytes.Equal(got, want) {
		t.Errorf("out2: got %b; want %b", got, want)
	}

	// out3 should contain 4
	if got, want := out3.Bytes(), []byte{4, 4, 4, 4}; !bytes.Equal(got, want) {
		t.Errorf("out3: got %b; want %b", got, want)
	}

	var result bytes.Buffer
	if err := Merge(&result, []io.Reader{&out1, &out2, &out3}); err != nil {
		t.Errorf("merge: %s", err)
	}

	// result should be 7^2^4 = 1
	if got, want := result.Bytes(), key; !bytes.Equal(got, want) {
		t.Errorf("result: got %b; want %b", got, want)
	}
}
