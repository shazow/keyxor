[![Build Status](https://circleci.com/gh/shazow/keyxor.svg?style=shield&circle-token=5feff4b31427639729e0d42c5f68b61c92d2bbaf)](https://circleci.com/gh/shazow/keyxor)
[![GoDoc](https://godoc.org/github.com/shazow/keyxor?status.svg)](https://godoc.org/github.com/shazow/keyxor)

# Keyxor Söze

`keyxor` is a tool for secret sharing by splitting up a key into multiple pieces which need to be XOR'd together to get the original private key. A key XOR, if you will.

**Status**: Alpha. Implementation is there, needs audit.


## Design

Given an input key, generate N-1 cryptographically secure random values of the key size. XOR the random values against the private key, producing N components (not including the original key) which need to be XOR'd together to produce the original key.

## Usage

```
$ keyxor --help
Usage:
  keyxor [OPTIONS] [merge | split]

Application Options:
  -v, --verbose  Show verbose logging.
      --version  Print version and exit.
      --base32   Encode output and decode input as Base32

Help Options:
  -h, --help     Show this help message

Available commands:
  merge  Merge secure pieces into the original secret
  split  Split a secret into secure pieces

Example:
  $ echo "hunter2" > secret.txt
  $ keyxor split secret.txt --num=3
  $ ls secret.*
  secret.txt   secret.txt.1   secret.txt.2   secret.txt.3
  $ cat secret.txt.{1,2,3}
  [... random symbols]
  $ keyxor merge secret.txt.{1,2,3}
  hunter2
```

```
# Use your fav key generator!
$ ssh-keygen -f key
$ ls
key            key.pub

# Split the private key into 3 pieces
$ keyxor split --num=3 ./key
$ ls
key            key.1            key.2            key.3            key.pub

# Merge the 3 pieces back together
$ keyxor merge ./key.* > key.new
$ shasum key key.new
<same hash>  key
<same hash>  key.new
```

## Is it any good?

It's a very simple tool that doesn't do much heavy lifting but it's very handy if you want to require exactly N secret key components to decrypt something.

This can be applied to a symmetric key or the private half of a public-private keypair, depending on the security model you're trying to achieve.

Keyxor does *not* support M of N secret sharing, like [Shamir's Secret Sharing scheme](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing). Keyxor uses a trivial method which only works for M=N. You need _all_ the pieces to get the original key. Simple is good, right?


## Audit?

Please do. The meat is inside [soze/soze.go](https://github.com/shazow/keyxor/blob/master/soze/soze.go).

* [@mveytsman](https://github.com/mveytsman), after an informal audit on [c2d26906]( https://github.com/shazow/keyxor/tree/c2d26906fbf4120cb2dc92afd9459dca878e8c86):
  > There's not enough code here to be interesting. ¯\\\_(ツ)\_/¯


## License

MIT
