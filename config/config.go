package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jxsl13/symbol-search/nm"
)

func NewConfig() Config {
	return Config{
		SearchDir:            ".",
		FilePathRegexList:    []string{},
		FileNameRegexList:    []string{`^([^\.]+|.+\.(so|a|dll|lib|exe|dylib))$`},
		FileModeMaskList:     []string{"0500", "0444"},
		SectionNameRegexList: []string{".*"},
		ArchiveRegex:         `\.(gz|tgz|xz||zst|bz2|tar|zip|7z)$`,
		SymbolNameRegexList:  []string{},
		PaintInterval:        10 * time.Second,
	}
}

type Config struct {
	SearchDir string `koanf:"search.dir" short:"f" description:"directory, file or archive to search for symbols recursively"`

	FilePathRegexList  []string         `koanf:"file.path.regex" short:"p" description:"optional comma separated list regex to match the file's parent path in the search dir or in archives"`
	FilePathRegexpList []*regexp.Regexp `koanf:"-"`

	FileNameRegexList  []string         `koanf:"file.name.regex" short:"n" description:"mandatory comma separated list regex to match file name in the search dir or in archives"`
	FileNameRegexpList []*regexp.Regexp `koanf:"-"`

	FileModeMaskList   []string      `koanf:"file.mode" short:"m" description:"optional comma separated list of file mode masks to match against (e.g. 0555, 0755, 0640, mode&mask == mask)"`
	FileFSModeMaskList []fs.FileMode `koanf:"-"`

	NoArchives    bool           `koanf:"no.archive" short:"A" description:"disables searching inside of archives"`
	ArchiveRegex  string         `koanf:"archive.regex" short:"a" description:"regex to match archive files in the search dir"`
	ArchiveRegexp *regexp.Regexp `koanf:"-"`

	SymbolNameRegexList  []string         `koanf:"symbol.name.regex" short:"s" description:"mandatory comma separated list of regex to match symbol name in binaries or libraries"`
	SymbolNameRegexpList []*regexp.Regexp `koanf:"-"`

	SectionNameRegexList  []string         `koanf:"section.name.regex" short:"S" description:"optional comma separated list of regex to match section name in binaries or libraries"`
	SectionNameRegexpList []*regexp.Regexp `koanf:"-"`

	OutputFile string `koanf:"output.file" short:"o" description:"output file to write the results to"`

	PaintInterval time.Duration `koanf:"-" short:"i" description:"interval to print progress"`

	Debug bool `koanf:"debug" short:"v" description:"enable debug output"`
}

func (cfg *Config) Validate() error {
	if cfg.SearchDir == "" {
		return errors.New("search dir is required")
	}

	if cfg.SearchDir == "." || strings.HasPrefix(cfg.SearchDir, "./") || strings.HasPrefix(cfg.SearchDir, ".\\") {
		pwd, err := os.Getwd()
		if err == nil {
			cfg.SearchDir = pwd
		}
	}

	// resolve home directory
	if cfg.SearchDir == "~" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		cfg.SearchDir = dir
	} else if strings.HasPrefix(cfg.SearchDir, "~/") || strings.HasPrefix(cfg.SearchDir, "~\\") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		cfg.SearchDir = filepath.Join(dir, cfg.SearchDir[2:])
	}

	_, err := os.Stat(cfg.SearchDir)
	if err != nil {
		return fmt.Errorf("invalid search dir: %w", err)
	}

	// can be empty
	for _, fr := range cfg.FilePathRegexList {
		re, err := regexp.Compile(fr)
		if err != nil {
			return fmt.Errorf("invalid file path regex: %w", err)
		}
		cfg.FilePathRegexpList = append(cfg.FilePathRegexpList, re)
	}

	if len(cfg.FileNameRegexList) == 0 {
		return errors.New("file name regex is required")
	}

	// must not be empty
	for _, fr := range cfg.FileNameRegexList {
		re, err := regexp.Compile(fr)
		if err != nil {
			return fmt.Errorf("invalid file regex: %w", err)
		}
		cfg.FileNameRegexpList = append(cfg.FileNameRegexpList, re)
	}

	// can be empty
	for _, mode := range cfg.FileModeMaskList {
		m, err := Mode(mode)
		if err != nil {
			return fmt.Errorf("invalid file mode: %w", err)
		}
		cfg.FileFSModeMaskList = append(cfg.FileFSModeMaskList, m)
	}

	if !cfg.NoArchives {
		if cfg.ArchiveRegex == "" {
			return errors.New("archive regex is required when including archives")
		}
		re, err := regexp.Compile(cfg.ArchiveRegex)
		if err != nil {
			return fmt.Errorf("invalid archive regex: %w", err)
		}
		cfg.ArchiveRegexp = re
	}

	if len(cfg.SymbolNameRegexList) == 0 {
		return errors.New("symbol name regex is required")
	}

	for _, sr := range cfg.SymbolNameRegexList {
		re, err := regexp.Compile(sr)
		if err != nil {
			return fmt.Errorf("invalid symbol name regex: %w", err)
		}
		cfg.SymbolNameRegexpList = append(cfg.SymbolNameRegexpList, re)
	}

	if len(cfg.SectionNameRegexList) == 0 {
		return errors.New("section name regex is required")
	}
	for _, sr := range cfg.SectionNameRegexList {
		re, err := regexp.Compile(sr)
		if err != nil {
			return fmt.Errorf("invalid section name regex: %w", err)
		}
		cfg.SectionNameRegexpList = append(cfg.SectionNameRegexpList, re)
	}

	if cfg.PaintInterval < time.Second {
		return errors.New("interval must be at least 1 second")
	}

	return nil
}

func (cfg *Config) IsMatchingArchive(path string) bool {
	return !cfg.NoArchives && cfg.ArchiveRegexp.MatchString(path)
}

func (cfg *Config) IsMatchingFileName(path string) bool {
	for _, re := range cfg.FileNameRegexpList {
		if re.MatchString(path) {
			return true
		}
	}
	return false
}

func (c *Config) IsMatchingFilePath(path string) bool {
	if len(c.FilePathRegexpList) == 0 {
		return true
	}

	for _, re := range c.FilePathRegexpList {
		if re.MatchString(path) {
			return true
		}
	}
	return false
}

func (c *Config) IsMatchingFileMode(mode fs.FileMode) bool {
	if len(c.FileFSModeMaskList) == 0 {
		return true
	}

	for _, mask := range c.FileFSModeMaskList {
		if mode&mask == mask {
			return true
		}
	}
	return false
}

func (c *Config) IsMatchingSymbol(symbol nm.Symbol) bool {
	return c.IsMatchingSymbolName(symbol.Name) &&
		c.IsMatchingSymbolSection(symbol.Section)
}

func (c *Config) IsMatchingSymbolName(name string) bool {
	for _, re := range c.SymbolNameRegexpList {
		if re.MatchString(name) {
			return true
		}
	}
	return false
}

func (c *Config) IsMatchingSymbolSection(name string) bool {
	for _, re := range c.SectionNameRegexpList {
		if re.MatchString(name) {
			return true
		}
	}
	return false
}

func (c *Config) WantSave() bool {
	return c.OutputFile != ""
}

// Mode converts a string representation of a file mode to a fs.FileMode
// The passed string is the same as the one used by unix chmod command.
// Example: "0755" or "01777"
func Mode(strMode string) (fs.FileMode, error) {
	if !strings.HasPrefix(strMode, "0") {
		strMode = "0" + strMode
	}

	u, err := strconv.ParseUint(strMode, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid mode: %s", strMode)
	}

	uintMode := fs.FileMode(u)
	resultMode := fs.FileMode(0)

	if uintMode&01000 != 0 {
		resultMode |= os.ModeSticky
	}

	resultMode |= os.ModePerm & uintMode

	return resultMode, nil
}
