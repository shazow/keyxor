# Keyxor Söze

`keyxor` is a tool for secret sharing by splitting up the private key into multiple pieces that need to be XOR'd together to get the original private key.

**Status**: Readme-driven development phase.


## Design

- Use [NaCl box](https://nacl.cr.yp.to/box.html) to generate a public key and private key.
- Use cryptographically-secure random data source to generate N-1 values and XOR against the private key, producing N components (not including the original private key) which need to be XOR'd together to produce the original private key.

## Questions

- [ ] Maybe skip the create feature and just do merging/splitting? No need for NaCl business, then. Bring your own PGP or ssh or whatever PKI.

## Usage

Fresh keys:

```
$ keyxor create --num=3
$ ls
key.pub
key.secret.1
key.secret.2
key.secret.3
$ keyxor merge key.secret.* > key.secret
```

Existing key:

```
$ keyxor split key.secret --num=3
key.secret.1
key.secret.2
key.secret.3
```

## License

MIT
