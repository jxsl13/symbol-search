# symbol-search

This utility can be used in order to find vulneravilities or rather vulnerable function symbols in binaries and shared libraries.
You can use this utility to search for binaries or shared libraries that might potentially contain vulnerable or exploitable
versions of dynamically linked or statically linked libraries.

The first argument expects a comma separated list of regular expressions.
The second argument is the root path that is used to start the search for binaries or shared libraries.
```sh
./symbol-search "gnutls_pkcs7_verify" /path/to/dir/or/single/file

./symbol-search "gnutls_pkcs7_verify" /path/to/dir/or/single/file -o report.txt
```

## Installation

```shell
go install github.com/jxsl13/symbol-search@latest
```

## Supported file formats
 - ELF
 - PE/COFF
 - Inside of archives:
    - tgz/gz, tar, xz, zip, 7z
