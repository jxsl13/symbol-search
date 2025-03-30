# symbol-search

This utility can be used in order to find vulneravilities or rather vulnerable function symbols in binaries and shared libraries.
You can use this utility to search for binaries or shared libraries that might potentially contain vulnerable or exploitable
versions of dynamically linked or statically linked libraries.

Supported file formats:

- Archive files (.tar, .zip, .gz, .bz2, .xz, .7z)
- ELF (Linux) (binary, .so)
- Static libraries (.a)
- PE (Windows) (.exe, .dll)

Not yet suported:

- Mach-O (macOS) (.dylib)

## Installation

```shell
go install github.com/jxsl13/symbol-search@latest
```

## Examples

```sh
symbol-search -s "gnutls_pkcs7_verify" -f /path/to/dir/or/single/file

symbol-search -s "gnutls_pkcs7_verify" -f /path/to/dir/or/single/file -o report.txt

symbol-search -s "gnutls_pkcs7_verify" -f /path/to/dir/or/single/ -o report.txt --concurrency 64

symbol-search -f /path/to/dir/or/single/ -o report.txt --concurrency 64 gnutls_pkcs7_verify callback_blub


# linux only shared library symbols that are that are either loaded at startup or dynamically at runtime
symbol-search --no-pe --no-internal -f "/path/to/dir/or/single/" '.*'

```

## Usage

```shell
$ symbol-search --help
Environment variables:
  SEARCH_DIR           directory to search for files recursively (default: ".")
  FILE_PATH_REGEX      comma separated list regex to match file path in the search dir or in archives (default: "[^.*$]")
  FILE_NAME_REGEX      comma separated list regex to match file name in the search dir or in archives (default: "[^([^\\.]+|.+\\.(so|a|dll|lib|exe))$]")
  INCLUDE_ARCHIVE      search inside archive files (default: "false")
  ARCHIVE_REGEX        regex to match archive files in the search dir (default: "\\.(gz|tgz|xz||zst|bz2|tar|zip|7z)$")
  SYMBOL_NAME_REGEX    comma separated list regex to match symbol name in the search dir or in archives (default: "[.*]")
  CONCURRENCY          number of concurrent workers to use (default: "6")
  OUTPUT_FILE          output file to write the results to
  NO_ELF               do not parse ELF files (Linux binaries) (default: "false")
  NO_PE                do not parse PE files (Windows binaries) (default: "false")
  NO_IMPORTED          do not parse imported symbols (from dll or shared objects) (default: "false")
  NO_DYNAMIC           do not parse dynamic symbols which are loaded at runtime with ldopen (default: "false")
  NO_INTERNAL          do not parse internal symbols from the binary or library itself (default: "false")
  DEBUG                enable debug output (default: "false")

Usage:
  symbol-search [flags]
  symbol-search [command]

Available Commands:
  completion  Generate completion script
  help        Help about any command

Flags:
  -a, --archive-regex string       regex to match archive files in the search dir (default "\\.(gz|tgz|xz||zst|bz2|tar|zip|7z)$")
  -t, --concurrency int            number of concurrent workers to use (default 6)
  -v, --debug                      enable debug output
  -n, --file-name-regex string     comma separated list regex to match file name in the search dir or in archives (default "[^([^\\.]+|.+\\.(so|a|dll|lib|exe))$]")
  -p, --file-path-regex string     comma separated list regex to match file path in the search dir or in archives (default "[^.*$]")
  -h, --help                       help for symbol-search
  -A, --include-archive            search inside archive files
      --no-dynamic                 do not parse dynamic symbols which are loaded at runtime with ldopen
      --no-elf                     do not parse ELF files (Linux binaries)
      --no-imported                do not parse imported symbols (from dll or shared objects)
      --no-internal                do not parse internal symbols from the binary or library itself
      --no-pe                      do not parse PE files (Windows binaries)
  -o, --output-file string         output file to write the results to
  -f, --search-dir string          directory to search for files recursively (default ".")
  -s, --symbol-name-regex string   comma separated list regex to match symbol name in the search dir or in archives (default "[.*]")

Use "symbol-search [command] --help" for more information about a command.
```
