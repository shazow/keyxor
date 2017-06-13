# Keyxor Söze

`keyxor` is a tool for secret sharing by splitting up a key into multiple pieces which need to be XOR'd together to get the original private key. A key XOR, if you will.

**Status**: Readme-driven development phase.


## Design

Given an input key, generate N-1 cryptographically secure random values of the key size. XOR the random values against the private key, producing N components (not including the original key) which need to be XOR'd together to produce the original key.

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

This can be applied to a symmetric key or the private half of a public-private keypair, depending on the security model you're trying to achieve.


## Someday?

- Maybe use NaCl box to have built-in key generating and encrypting/decrypting functionality? (Worth it?)


## License

MIT
