package soze

import (
	"crypto/rand"
	"errors"
	"io"
)

// TODO: use io.ReadFull instead of checking read lengths?

var bufSize = 1024 * 8

// Split takes a input key reader and N output writers, and splits it into N
// key components which need to be XOR'd in order to produce the original
// input key. Split uses crypto/rand.Reader as the randomness source.
func Split(in io.Reader, outs []io.Writer) error {
	return splitWithRand(in, outs, rand.Reader)
}

func splitWithRand(in io.Reader, outs []io.Writer, randReader io.Reader) error {
	if len(outs) < 2 {
		return errors.New("must have at least two outputs")
	}

	inBuf, randBuf, combinedBuf := make([]byte, bufSize), make([]byte, bufSize), make([]byte, bufSize)
	for {
		n, err := in.Read(inBuf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Accumulate the random key components for N-1 keys
		for i, out := range outs {
			if i == 0 {
				// First index contains the combined material, for convenience
				continue
			}
			nRand, err := randReader.Read(randBuf[:n])
			if err != nil {
				return err
			}
			if nRand < n {
				return errors.New("insufficient random data")
			}

			xorInto(combinedBuf[:n], randBuf[:n])

			logger.Print("writing random material", i)
			nOut, err := out.Write(randBuf[:n])
			if err != nil {
				return err
			}
			if nOut != n {
				return errors.New("written output is the wrong length")
			}
		}

		logger.Print("writing combined material")
		xorInto(combinedBuf[:n], inBuf[:n])

		// Write the combined material into the Nth key (first index in this case)
		nOut, err := outs[0].Write(combinedBuf[:n])
		if err != nil {
			return err
		}
		if nOut != n {
			return errors.New("written output is the wrong length")
		}
	}
}

// Merge takes an output key writer and N input key component readers, and
// applies XOR to the inputs to produce the original output key.
func Merge(out io.Writer, ins []io.Reader) error {
	inBuf, outBuf := make([]byte, bufSize), make([]byte, bufSize)
	for {
		var n int
		var err error
		for _, in := range ins {
			n, err = in.Read(inBuf)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			xorInto(outBuf[:n], inBuf[:n])
		}
		if n == 0 {
			return errors.New("no input data")
		}

		nOut, err := out.Write(outBuf[:n])
		if err != nil {
			return err
		}
		if nOut != n {
			return errors.New("written output is the wrong length")
		}
	}
}
