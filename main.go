package main

import (
	"encoding/base32"
	"fmt"
	"io"
	"os"

	flags "github.com/jessevdk/go-flags"
)

// Version of keyxor, overwritten during build.
var Version = "dev"

// Options contains the flag options
type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose logging."`
	Version bool   `long:"version" description:"Print version and exit."`
	Base32  bool   `long:"base32" description:"Encode output and decode input as Base32"`

	Split struct {
		Num  int `short:"n" long:"num" description:"Number of pieces." default:"3"`
		Args struct {
			Key string `positional-arg-name:"SECRET"`
		} `positional-args:"yes" required:"yes"`
	} `command:"split" description:"Split a secret into secure pieces"`

	Merge struct {
		Args struct {
			Pieces []string `required:"2" positional-arg-name:"PIECES"`
		} `positional-args:"yes" required:"yes"`
	} `command:"merge" description:"Merge secure pieces into the original secret"`
}

const usageExample = `Example:
  $ echo "hunter2" > secret.txt
  $ keyxor split secret.txt --num=3
  $ ls secret.*
  secret.txt   secret.txt.1   secret.txt.2   secret.txt.3
  $ cat secret.txt.{1,2,3}
  [... random symbols]
  $ keyxor merge secret.txt.{1,2,3}
  hunter2`

func main() {
	options := Options{}
	parser := flags.NewParser(&options, flags.Default)
	parser.SubcommandsOptional = true

	args, err := parser.Parse()
	if err != nil {
		if flagErr, ok := err.(*flags.Error); ok && flagErr.Type == flags.ErrHelp && parser.Active == nil {
			// Print usage when run with just --help
			exit(0, usageExample)
		}
		if args == nil {
			exit(1, err.Error())
		}
		os.Exit(1)
		return
	}

	if options.Version {
		fmt.Printf("%s\n", Version)
		os.Exit(0)
	}

	if parser.Active == nil {
		parser.WriteHelp(os.Stderr)
		exit(1, "\nMust specify a command.")
	}

	var enc Encoder
	var dec Decoder
	if options.Base32 {
		encoding := base32.StdEncoding

		enc = func(w io.Writer) io.WriteCloser {
			return base32.NewEncoder(encoding, w)
		}
		dec = func(r io.Reader) io.Reader {
			return base32.NewDecoder(encoding, r)
		}
	}

	switch parser.Active.Name {
	case "split":
		if err := splitKey(options.Split.Args.Key, options.Split.Num, enc); err != nil {
			exit(2, "split failed: %s", err)
		}
	case "merge":
		if err := mergeKey(options.Merge.Args.Pieces, dec); err != nil {
			exit(3, "merge failed: %s", err)
		}
	default:
		exit(1, "invalid command: %q", parser.Active)
	}
}

// exit prints an error and exits with the given code
func exit(code int, msg string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", a...)
	os.Exit(code)
}
