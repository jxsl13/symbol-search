package nm

import (
	"debug/elf"
	"regexp"
	"sort"
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

	sort.Sort(ByName(result))

	return &SymbolFile{
		Path:    filename,
		Symbols: result,
	}, nil
}

type ByPath []*SymbolFile

func (a ByPath) Len() int           { return len(a) }
func (a ByPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPath) Less(i, j int) bool { return a[i].Path < a[j].Path }

type ByName []Symbol

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
