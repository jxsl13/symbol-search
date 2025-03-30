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
		symbol := NewSymbol(s.Name, s.Value, s.Size, DynamicSection, s.Section.String(), s.Version, s.Library)
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
		result = append(result, NewSymbol(s.Name, 0, 0, ImportedSection, "", s.Version, s.Library))
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
		symbol := NewSymbol(s.Name, s.Value, s.Size, InternalSection, s.Section.String(), s.Version, s.Library)
		result = append(result, symbol)
	}

	return result, nil

}

func OpenELF(file io.ReaderAt) (*elf.File, error) {
	return elf.NewFile(file)
}
