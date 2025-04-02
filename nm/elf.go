package nm

import (
	"debug/elf"
	"errors"
	"io"
)

func ReadDynamicSymbolsELF(file *elf.File) ([]Symbol, error) {
	ds, err := file.DynamicSymbols()
	if err != nil {
		if errors.Is(err, elf.ErrNoSymbols) {
			return []Symbol{}, nil
		}
		return nil, err
	}

	result := make([]Symbol, 0, len(ds))
	for _, s := range ds {
		symbol := NewSymbol(
			TypeELF,
			SubTypeDynamic,
			s.Name,
			s.Value,
			s.Size,
			s.Section.String(),
			s.Version,
			s.Library,
		)
		result = append(result, symbol)
	}

	return result, nil
}

func ReadImportedSymbolsELF(file *elf.File) ([]Symbol, error) {
	is, err := file.ImportedSymbols()
	if err != nil {
		if errors.Is(err, elf.ErrNoSymbols) {
			return []Symbol{}, nil
		}
		return nil, err
	}

	result := make([]Symbol, 0, len(is))
	for _, s := range is {
		symbol := NewSymbol(
			TypeELF,
			SubTypeImported,
			s.Name,
			0,
			0,
			"",
			s.Version,
			s.Library,
		)
		result = append(result, symbol)
	}

	return result, nil
}

func ReadInternalSymbolsELF(file *elf.File) ([]Symbol, error) {
	ss, err := file.Symbols()
	if err != nil {
		if errors.Is(err, elf.ErrNoSymbols) {
			return []Symbol{}, nil
		}
		return nil, err
	}

	result := make([]Symbol, 0, len(ss))
	for _, s := range ss {
		symbol := NewSymbol(
			TypeELF,
			SubTypeInternal,
			s.Name,
			s.Value,
			s.Size,
			s.Section.String(),
			s.Version,
			s.Library,
		)
		result = append(result, symbol)
	}

	return result, nil

}

func OpenELF(file io.ReaderAt) (*elf.File, error) {
	return elf.NewFile(file)
}
