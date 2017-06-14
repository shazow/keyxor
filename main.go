package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/shazow/keyxor/soze"
)

// Version of keyxor, overwritten during build.
var Version = "dev"

var numSplit = flag.Int("num", 2, "number of components to split the key into")
var versionFlag = flag.Bool("version", false, "print the version and exit")

func exit(code int, msg string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", a...)
	os.Exit(code)
}

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

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s\n", Version)
		os.Exit(0)
	}

	switch cmd := flag.Arg(0); cmd {
	case "":
		exit(1, "must supply a command: split, merge")
	case "split":
		path := flag.Arg(1)
		if path == "" {
			exit(1, "split: specify a path to a key to split")
		}
		if err := splitKey(flag.Arg(1), *numSplit); err != nil {
			exit(2, "split failed: %s", err)
		}
	case "merge":
		args := flag.Args()
		if len(args) < 3 {
			exit(1, "merge: need at least 2 key component paths to merge")
		}
		if err := mergeKey(args[1:len(args)]); err != nil {
			exit(3, "merge failed: %s", err)
		}
	}
}
