package main

import (
	"fmt"
	"io"
	"os"

	"github.com/shazow/keyxor/soze"
)

func splitKey(path string, num int) error {
	key, err := os.Open(path)
	if err != nil {
		return err
	}
	defer key.Close()

	// Keep track of which files we opened and make sure to close them
	outs := make([]io.Writer, 0, num)
	files := make([]*os.File, 0, num)
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	// Open a file for writing for each key component
	for i := 1; i <= num; i++ {
		filename := fmt.Sprintf("%s.%d", path, i)
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		outs = append(outs, f)
		files = append(files, f)
	}

	if err := soze.Split(key, outs); err != nil {
		return err
	}

	// Make sure output files are written to disk
	for _, f := range files {
		if err := f.Sync(); err != nil {
			return err
		}
	}

	return nil
}

func mergeKey(paths []string) error {
	num := len(paths)
	ins := make([]io.Reader, 0, num)
	files := make([]*os.File, 0, num)
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		files = append(files, f)
		ins = append(ins, f)
	}

	return soze.Merge(os.Stdout, ins)
}
