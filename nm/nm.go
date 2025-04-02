package nm

import (
	"strings"
)

const (
	UnknownLibrary  = "unknown"
	UnknownVersion  = "unknown"
	TypeELF         = "elf"
	TypeMachO       = "macho"
	TypePE          = "pe"
	TypeAR          = "ar"
	SubTypeImported = "imported"
	SubTypeDynamic  = "dynamic"
	SubTypeInternal = "internal"
)

type Symbol struct {
	Name  string
	Value uint64

	Size    uint64
	Section string

	// These fields are present only for the dynamic symbol table.
	Version string
	Library string

	ArchiveType string // containing archive type like .a
	Type        string // actual type of the object
	SubType     string // subtype of the symbol (internal, external, dynamic, imported, etc.)
}

func NewSymbol(symbolType, subType string, name string, value, size uint64, section string, version, library string) Symbol {

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
		Section: section,
		Version: defaultIfEmpty(version, UnknownVersion),
		Library: defaultIfEmpty(library, UnknownLibrary),

		ArchiveType: "",
		Type:        symbolType,
		SubType:     subType,
	}
}

func (s *Symbol) SetArchiveType(archiveType string) {
	s.ArchiveType = archiveType
}

func (s *Symbol) Header() []any {
	return []any{"ArchiveType", "Type", "SubType", "Name", "Value", "Size", "Section", "Version", "Library"}
}

func (s *Symbol) Row() []any {
	return []any{s.ArchiveType, s.Type, s.SubType, s.Name, s.Value, s.Size, s.Section, s.Version, s.Library}
}

func defaultIfEmpty(s, defaultString string) string {
	if s == "" {
		return defaultString
	}
	return s
}
