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
  SEARCH_DIR            directory, file or archive to search for symbols recursively (default: ".")
  FILE_PATH_REGEX       optional comma separated list regex to match the file's parent path in the search dir or in archives (default: "[]")
  FILE_NAME_REGEX       mandatory comma separated list regex to match file name in the search dir or in archives (default: "[^([^\\.]+|.+\\.(so|a|dll|lib|exe|dylib))$]")
  FILE_MODE             optional comma separated list of file mode masks to match against (e.g. 0555, 0755, 0640, mode&mask == mask) (default: "[0500 0444]")
  NO_ARCHIVE            disables searching inside of archives (default: "false")
  ARCHIVE_REGEX         regex to match archive files in the search dir (default: "\\.(gz|tgz|xz||zst|bz2|tar|zip|7z)$")
  SYMBOL_NAME_REGEX     mandatory comma separated list of regex to match symbol name in binaries or libraries (default: "[]")
  SECTION_NAME_REGEX    optional comma separated list of regex to match section name in binaries or libraries (default: "[.*]")
  OUTPUT_FILE           output file to write the results to
  DEBUG                 enable debug output (default: "false")

Usage:
  symbol-search [flags]
  symbol-search [command]

Available Commands:
  completion  Generate completion script
  help        Help about any command

Flags:
  -a, --archive-regex string        regex to match archive files in the search dir (default "\\.(gz|tgz|xz||zst|bz2|tar|zip|7z)$")
  -v, --debug                       enable debug output
  -m, --file-mode string            optional comma separated list of file mode masks to match against (e.g. 0555, 0755, 0640, mode&mask == mask) (default "[0500 0444]")
  -n, --file-name-regex string      mandatory comma separated list regex to match file name in the search dir or in archives (default "[^([^\\.]+|.+\\.(so|a|dll|lib|exe|dylib))$]")
  -p, --file-path-regex string      optional comma separated list regex to match the file's parent path in the search dir or in archives (default "[]")
  -h, --help                        help for symbol-search
  -A, --no-archive                  disables searching inside of archives
  -o, --output-file string          output file to write the results to
  -f, --search-dir string           directory, file or archive to search for symbols recursively (default ".")
  -S, --section-name-regex string   optional comma separated list of regex to match section name in binaries or libraries (default "[.*]")
  -s, --symbol-name-regex string    mandatory comma separated list of regex to match symbol name in binaries or libraries (default "[]")

Use "symbol-search [command] --help" for more information about a command.
```
