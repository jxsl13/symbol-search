package nm

import "debug/elf"

func GetSymbols(filename string) (*SymbolFile, error) {
	f, err := elf.Open(filename)
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
			Version: s.Version,
			Library: s.Library,
		})
	}

	for _, s := range ds {
		result = append(result, Symbol{
			Name:    s.Name,
			Version: s.Version,
			Library: s.Library,
		})
	}

	for _, s := range is {
		result = append(result, Symbol{
			Name:    s.Name,
			Version: s.Version,
			Library: s.Library,
		})
	}

	return &SymbolFile{
		Path:    filename,
		Symbols: result,
	}, nil
}
