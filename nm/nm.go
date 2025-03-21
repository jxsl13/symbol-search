package nm

import (
	"strings"
)

var (
	UnknownLibrary  = "unknown"
	UnknownVersion  = "unknown"
	InternalSection = "internal"
	ImportedSection = "imported"
	DynamicSection  = "dynamic"
)

type Symbol struct {
	Name  string
	Value uint64

	Size    uint64
	Section string

	Source string

	// These fields are present only for the dynamic symbol table.
	Version string
	Library string
}

func NewSymbol(name string, value, size uint64, source, section string, version, library string) Symbol {

	// occasionally we get versions with double @@ instead of most likely a single one which
	// is used to split the version from the symbol name
	if strings.Contains(name, "@@") {
		ss := strings.SplitN(name, "@@", 2)
		name = ss[0]
		version = ss[1]
	} else if strings.Contains(name, ":") {
		ss := strings.SplitN(name, ":", 2)
		name = ss[0]
		library = ss[1]
	}

	return Symbol{
		Name:    name,
		Value:   value,
		Size:    size,
		Source:  source,
		Section: section,
		Version: defaultIfEmpty(version, UnknownVersion),
		Library: defaultIfEmpty(library, UnknownLibrary),
	}

}

func (s *Symbol) Header() []any {
	return []any{"Name", "Value", "Size", "Source", "Section", "Version", "Library"}
}

func (s *Symbol) Row() []any {
	return []any{s.Name, s.Value, s.Size, s.Source, s.Section, s.Version, s.Library}
}

func defaultIfEmpty(s, defaultString string) string {
	if s == "" {
		return defaultString
	}
	return s
}
