package nm

import (
	"io"
	"os"
	"regexp"
)

var (
	UnknownLibrary = "unknown"
	UnknownVersion = "unknown"
)

func NewSymbols(file io.ReaderAt) ([]Symbol, error) {
	s, err := SymbolsELF(file)
	if err == nil {

		return cleanupSymbols(s), nil
	}

	s, err = SymbolsPE(file)
	if err == nil {
		return cleanupSymbols(s), nil
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
			if s.Name == "" {
				continue
			}
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

func cleanupSymbols(ss []Symbol) []Symbol {
	m := make(map[string]Symbol, len(ss))

	for _, s := range ss {
		if s.Name == "" {
			continue
		}
		v, found := m[s.Name]
		if !found {
			m[s.Name] = s
			continue
		}
		m[s.Name] = Symbol{
			Name:    useNotEmpty(s.Name, v.Name),
			Version: useNotEmpty(s.Version, v.Version),
			Library: useNotEmpty(s.Library, v.Library),
		}
	}

	result := make([]Symbol, 0, len(ss))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func useNotEmpty(a, b string) string {
	if a == "" {
		return b
	}
	return b
}
