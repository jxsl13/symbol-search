# symbol-search


This utility can be used in order to find vulneravilities or rather vulnerable function symbols in binaries and shared libraries.
You can use this utility to search for binaries or shared libraries that might potentially contain vulnerable or exploitable
versions of dynamically linked or statically linked libraries.

The first argument expects a comma separated list of regular expressions.
The second argument is the root path that is used to start the search for binaries or shared libraries.
```sh
./symbol-search "gnutls,gnu" /path/to/dir/or/single/file
```

## Supported file formats
 - ELF

## TODO:

- Add an automatic zip,tar,tar.gz,tgz,tar.xz file reader that extracts files in memory and checks their symbols.
- Support for the Windows file format: COFF (https://pkg.go.dev/debug/pe)