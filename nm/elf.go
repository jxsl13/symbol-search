package nm

import (
	"debug/elf"
	"io"
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
		result = append(result, Symbol{
			Name:    s.Name,
			Version: defaultIfEmpty(s.Version, UnknownLibrary),
			Library: s.Library,
		})
	}

	for _, s := range ds {
		result = append(result, Symbol{
			Name:    s.Name,
			Version: defaultIfEmpty(s.Version, UnknownLibrary),
			Library: s.Library,
		})
	}

	for _, s := range is {
		result = append(result, Symbol{
			Name:    s.Name,
			Version: defaultIfEmpty(s.Version, UnknownLibrary),
			Library: s.Library,
		})
	}

	return result, nil
}
