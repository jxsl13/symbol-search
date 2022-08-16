package nm

import (
	"debug/elf"
	"io"
	"strings"
)

func SymbolsELF(file io.ReaderAt) ([]Symbol, error) {
	f, err := elf.NewFile(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ss, err := f.Symbols()
	if err != nil {
		return nil, err
	}

	ds, err := f.DynamicSymbols()
	if err != nil {
		return nil, err
	}

	is, err := f.ImportedSymbols()
	if err != nil {
		return nil, err
	}

	result := make([]Symbol, 0, len(ss)+len(ds)+len(is))

	for _, s := range ss {
		result = append(result, NewSymbol(s.Name, s.Version, s.Library))
	}

	for _, s := range ds {
		result = append(result, NewSymbol(s.Name, s.Version, s.Library))
	}

	for _, s := range is {
		result = append(result, NewSymbol(s.Name, s.Version, s.Library))
	}

	return result, nil
}

func NewSymbol(name, version, library string) Symbol {

	// occasionally we get versions with double @@ instead of most likely a single one which
	// is used to split the version from the symbol name
	if strings.Contains(name, "@@") {
		ss := strings.SplitN(name, "@@", 2)
		name = ss[0]
		version = ss[1]
	} else if strings.Contains(name, ":") {
		ss := strings.SplitN(name, ":", 2)
		name = ss[0]
		library = ss[1]
	}

	return Symbol{
		Name:    name,
		Version: defaultIfEmpty(version, UnknownVersion),
		Library: defaultIfEmpty(library, UnknownLibrary),
	}

}
