# Keyxor Söze

`keyxor` is a tool for secret sharing by splitting up the private key into multiple pieces that need to be XOR'd together to get the original private key.

**Status**: Readme-driven development phase.


## Design

Given an input secret key, generate N-1 cryptographically secure random values of the key size. XOR the random values against the private key, producing N components (not including the original private key) which need to be XOR'd together to produce the original private key.

## Usage

```
$ ssh-keygen -f key
$ ls
key            key.pub
$ keyxor split ./key
$ ls
key            key.1            key.2            key.3            key.pub
$ keyxor merge ./key.* > key.new
$ shasum key key.new
<same hash>  key
<same hash>  key.new
```

## Is it any good?

Yes.

It's a very simple tool that doesn't do much heavy lifting but it's very handy if you want to require exactly N secret key components to decrypt something.


## License

MIT
