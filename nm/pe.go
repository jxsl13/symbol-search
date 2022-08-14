package nm

import (
	"debug/pe"
	"io"
	"strconv"
)

func SymbolsPE(file io.ReaderAt) ([]Symbol, error) {
	f, err := pe.NewFile(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make([]Symbol, 0, len(f.Symbols))
	for _, s := range f.Symbols {
		result = append(result, Symbol{
			Name:    s.Name,
			Version: strconv.Itoa(int(s.Value)),
			Library: UnknownLibrary,
		})
	}

	return result, nil
}
