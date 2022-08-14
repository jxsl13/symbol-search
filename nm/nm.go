package nm

import (
	"io"
	"os"
	"regexp"
)

var UnknownLibrary = "unknown"

func NewSymbols(file io.ReaderAt) ([]Symbol, error) {
	s, err := SymbolsELF(file)
	if err == nil {
		return s, nil
	}

	s, err = SymbolsPE(file)
	if err == nil {
		return s, nil
	}

	return nil, err
}

func NewFilteredSymbols(file io.ReaderAt, matchers []*regexp.Regexp) ([]Symbol, error) {
	ss, err := NewSymbols(file)
	if err != nil {
		return nil, err
	}

	result := make([]Symbol, 0, len(matchers))
	for _, s := range ss {
		for _, m := range matchers {
			if m.MatchString(s.Name) {
				result = append(result, s)
			}
		}
	}
	return result, nil
}

func GetFilteredSymbols(path string, matchers []*regexp.Regexp) ([]Symbol, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewFilteredSymbols(f, matchers)
}
