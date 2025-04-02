package nm

import (
	"errors"
	"fmt"
	"io"

	"github.com/blacktop/go-macho"
)

var (
	ErrInvalidMachOFormat = fmt.Errorf("invalid Mach-O format")
)

func OpenMachO(r io.ReaderAt) (*macho.File, error) {
	f, err := macho.NewFile(r)
	if err != nil {
		var machoErr *macho.FormatError
		if errors.As(err, &machoErr) {
			return nil, fmt.Errorf("%w: %w", ErrInvalidMachOFormat, machoErr)
		}
		return nil, err
	}

	if f.Dysymtab == nil || f.Symtab == nil {
		return nil, ErrInvalidMachOFormat
	}

	return f, nil
}

func dynamicSymbolIndexMap(f *macho.File) (map[uint32]struct{}, error) {
	if f.Dysymtab == nil || f.Symtab == nil {
		return nil, ErrInvalidMachOFormat
	}
	dt := f.Dysymtab

	dynamic := make(map[uint32]struct{}, len(dt.IndirectSyms))
	for _, i := range dt.IndirectSyms {
		dynamic[i] = struct{}{}
	}
	return dynamic, nil
}

func ReadImportedSymbolsMachO(f *macho.File) ([]Symbol, error) {

	st := f.Symtab
	dt := f.Dysymtab

	dynamic, err := dynamicSymbolIndexMap(f)
	if err != nil {
		return nil, err
	}

	symbols := st.Syms[dt.Iundefsym : dt.Iundefsym+dt.Nundefsym]

	result := make([]Symbol, 0, len(symbols))

	for i, s := range symbols {

		// skip dynamic symbols
		if _, ok := dynamic[uint32(i)]; ok {
			continue
		}

		lib := s.GetLib(f)
		if lib == "" {
			lib = UnknownLibrary
		}

		if len(f.Sections) >= int(s.Sect) {
			return nil, fmt.Errorf("section %d out of bounds", s.Sect)
		}

		sect := f.Sections[s.Sect]

		symbol := NewSymbol(
			TypeMachO,
			SubTypeImported,
			s.Name,
			s.Value,
			0,
			sect.Name,
			UnknownVersion,
			lib,
		)
		result = append(result, symbol)
	}
	return result, nil
}

func ReadDynamicSymbolsMachO(f *macho.File) ([]Symbol, error) {
	f.ImportedSymbols()

	if f.Dysymtab == nil || f.Symtab == nil {
		return nil, ErrInvalidMachOFormat
	}

	st := f.Symtab
	dt := f.Dysymtab

	result := make([]Symbol, 0, len(dt.IndirectSyms))

	for _, i := range dt.IndirectSyms {
		if i >= uint32(len(st.Syms)) {
			return nil, fmt.Errorf("indirect symbol %d out of bounds", i)
		}
		s := st.Syms[int(i)]

		lib := s.GetLib(f)
		if lib == "" {
			lib = UnknownLibrary
		}

		if len(f.Sections) >= int(s.Sect) {
			return nil, fmt.Errorf("section %d out of bounds", s.Sect)
		}
		sect := f.Sections[s.Sect]
		symbol := NewSymbol(
			TypeMachO,
			SubTypeDynamic,
			s.Name,
			s.Value,
			0,
			sect.Name,
			UnknownVersion,
			lib,
		)
		result = append(result, symbol)
	}

	return result, nil
}
