# Keyxor Söze

`keyxor` is a tool for secret sharing by splitting up the private key into multiple pieces that need to be XOR'd together to get the original private key.

**Status**: Readme-driven development phase.


## Design

- Use [NaCl box](https://nacl.cr.yp.to/box.html) to generate a public key and private key.
- Use cryptographically-secure random data source to generate N-1 values and XOR against the private key, producing N components (not including the original private key) which need to be XOR'd together to produce the original private key.


## Usage

```
$ keyxor create --num=3
$ ls
key.pub
key.1.secret
key.2.secret
key.3.secret
$ keyxor xor key.*.secret > key.secret
```


## License

MIT
