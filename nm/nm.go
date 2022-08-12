package nm

import (
	"debug/elf"
	"regexp"
)

type Symbol struct {
	Name    string
	Version string
	Library string
}

type SymbolFile struct {
	Path    string
	Symbols []Symbol
}

func (s *SymbolFile) HasSymbols() bool {
	return len(s.Symbols) > 0
}

func GetFilteredSymbols(filename string, matchers []*regexp.Regexp) (*SymbolFile, error) {

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
		for _, matcher := range matchers {
			if matcher.MatchString(s.Name) {
				result = append(result, Symbol{
					Name:    s.Name,
					Version: s.Version,
					Library: s.Library,
				})
			}
		}

	}

	for _, s := range ds {
		for _, matcher := range matchers {
			if matcher.MatchString(s.Name) {
				result = append(result, Symbol{
					Name:    s.Name,
					Version: s.Version,
					Library: s.Library,
				})
			}
		}
	}

	for _, s := range is {
		for _, matcher := range matchers {
			if matcher.MatchString(s.Name) {
				result = append(result, Symbol{
					Name:    s.Name,
					Version: s.Version,
					Library: s.Library,
				})
			}
		}
	}

	m := make(map[Symbol]struct{}, len(result))
	for _, v := range result {
		m[v] = struct{}{}
	}
	result = result[:0]
	for k := range m {
		result = append(result, k)
	}

	return &SymbolFile{
		Path:    filename,
		Symbols: result,
	}, nil
}
