package nm

import (
	"debug/pe"
	"io"
)

func ReadCOFFSymbolsPE(file *pe.File) ([]Symbol, error) {

	result := make([]Symbol, 0, len(file.Symbols))
	skipAux := uint8(0)
	for _, s := range file.COFFSymbols {
		if skipAux > 0 {
			skipAux--
			continue
		}
		skipAux = s.NumberOfAuxSymbols

		name, err := s.FullName(file.StringTable)
		if err != nil {
			return nil, err
		}

		if s.SectionNumber < 0 || int(s.SectionNumber) >= len(file.Sections) {
			// Skip undefined symbols
			continue
		}

		section := file.Sections[s.SectionNumber].Name
		symbol := NewSymbol(
			TypePE,
			SubTypeInternal,
			name,
			uint64(s.Value),
			0,
			section,
			UnknownVersion,
			UnknownLibrary,
		)

		result = append(result, symbol)
	}

	return result, nil
}

func ReadImportedSymbolsPE(file *pe.File) ([]Symbol, error) {
	imported, err := file.ImportedSymbols()
	if err != nil {
		return nil, err
	}

	result := make([]Symbol, 0, len(imported))
	for _, name := range imported {

		symbol := NewSymbol(
			TypePE,
			SubTypeImported,
			name,
			0,
			0,
			"",
			UnknownVersion,
			UnknownLibrary,
		)
		result = append(result, symbol)
	}

	return result, nil
}

func OpenPE(file io.ReaderAt) (*pe.File, error) {
	return pe.NewFile(file)
}
