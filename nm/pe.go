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
		section := file.Sections[s.SectionNumber].Name
		symbol := NewSymbol(name, uint64(s.Value), 0, InternalSection, section, UnknownVersion, UnknownLibrary)
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
		result = append(result, NewSymbol(name, 0, 0, ImportedSection, "", UnknownVersion, UnknownLibrary))
	}

	return result, nil
}

func OpenPE(file io.ReaderAt) (*pe.File, error) {
	return pe.NewFile(file)
}
