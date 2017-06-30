package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shazow/keyxor/soze"
)

type Encoder func(io.Writer) io.WriteCloser
type Decoder func(io.Reader) io.Reader

func splitKey(path string, num int, enc Encoder) error {
	key, err := os.Open(path)
	if err != nil {
		return err
	}
	defer key.Close()

	outs := make([]io.Writer, 0, num)

	// Open a file for writing for each key component
	for i := 1; i <= num; i++ {
		filename := fmt.Sprintf("%s.%d", filepath.Base(path), i)
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
		defer f.Sync()

		var out io.WriteCloser = f
		if enc != nil {
			// Encoders need to be Close()'d to flush.
			out = enc(out)
			defer out.Close()
		}
		outs = append(outs, out)
	}

	if err := soze.Split(key, outs); err != nil {
		return err
	}

	return nil
}

func mergeKey(paths []string, dec Decoder) error {
	ins := make([]io.Reader, 0, len(paths))

	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		var in io.Reader = f
		if dec != nil {
			in = dec(in)
		}
		ins = append(ins, in)
	}

	return soze.Merge(os.Stdout, ins)
}
