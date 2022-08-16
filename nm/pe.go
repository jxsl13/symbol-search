package nm

import (
	"debug/pe"
	"io"
)

func SymbolsPE(file io.ReaderAt) ([]Symbol, error) {
	f, err := pe.NewFile(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make([]Symbol, 0, len(f.COFFSymbols))
	for _, s := range f.COFFSymbols {
		name, err := s.FullName(f.StringTable)
		if err == nil {
			result = append(result, NewSymbol(name, "", ""))
		} else {
			result = append(result, NewSymbol(string(s.Name[:]), "", ""))
		}
	}

	imported, err := f.ImportedSymbols()
	if err == nil {
		for _, s := range imported {
			result = append(result, NewSymbol(s, "", ""))
		}
	}

	return result, nil
}
