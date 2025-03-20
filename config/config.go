package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func NewConfig() Config {
	return Config{
		SearchDir:           ".",
		FilePathRegexList:   []string{`^.*$`},
		FileNameRegexList:   []string{`^([^\.]+|.+\.(so|a|dll|lib|exe))$`},
		ArchiveRegex:        `\.(gz|tgz|xz||zst|bz2|tar|zip|7z)$`,
		SymbolNameRegexList: []string{".*"},
		PaintInterval:       10 * time.Second,
		Concurrency:         max(1, runtime.NumCPU()),
	}
}

type Config struct {
	SearchDir string `koanf:"search.dir" short:"f" description:"directory to search for files recursively"`

	FilePathRegexList  []string         `koanf:"file.path.regex" short:"p" description:"comma separated list regex to match file path in the search dir or in archives"`
	FilePathRegexpList []*regexp.Regexp `koanf:"-"`

	FileNameRegexList  []string         `koanf:"file.name.regex" short:"n" description:"comma separated list regex to match file name in the search dir or in archives"`
	FileNameRegexpList []*regexp.Regexp `koanf:"-"`

	IncludeArchives bool           `koanf:"include.archive" short:"A" description:"search inside archive files"`
	ArchiveRegex    string         `koanf:"archive.regex" short:"a" description:"regex to match archive files in the search dir"`
	ArchiveRegexp   *regexp.Regexp `koanf:"-"`

	SymbolNameRegexList  []string         `koanf:"symbol.name.regex" short:"s" description:"comma separated list regex to match symbol name in the search dir or in archives"`
	SymbolNameRegexpList []*regexp.Regexp `koanf:"-"`

	Concurrency int `koanf:"concurrency" short:"t" description:"number of concurrent workers to use"`

	OutputFile string `koanf:"output.file" short:"o" description:"output file to write the results to"`

	PaintInterval time.Duration `koanf:"-" short:"i" description:"interval to print progress"`

	NoELF bool `koanf:"no.elf" description:"do not parse ELF files (Linux binaries)"`
	NoPE  bool `koanf:"no.pe" description:"do not parse PE files (Windows binaries)"`

	NoImported bool `koanf:"no.imported" description:"do not parse imported symbols (from dll or shared objects)"`
	NoDynamic  bool `koanf:"no.dynamic" description:"do not parse dynamic symbols which are loaded at runtime with ldopen"`
	NoInternal bool `koanf:"no.internal" description:"do not parse internal symbols from the binary or library itself"`

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

	_, err := os.Stat(cfg.SearchDir)
	if err != nil {
		return fmt.Errorf("invalid search dir: %w", err)
	}

	if len(cfg.FilePathRegexList) == 0 {
		return errors.New("file path regex is required")
	}

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

	for _, fr := range cfg.FileNameRegexList {
		re, err := regexp.Compile(fr)
		if err != nil {
			return fmt.Errorf("invalid file regex: %w", err)
		}
		cfg.FileNameRegexpList = append(cfg.FileNameRegexpList, re)
	}

	if cfg.IncludeArchives {
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

	if cfg.Concurrency < 1 {
		return errors.New("concurrency must be greater than 0")
	}

	if cfg.PaintInterval < time.Second {
		return errors.New("interval must be at least 1 second")
	}

	return nil
}

func (cfg *Config) IsMatchingArchive(path string) bool {
	return cfg.IncludeArchives && cfg.ArchiveRegexp.MatchString(path)
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
	for _, re := range c.FilePathRegexpList {
		if re.MatchString(path) {
			return true
		}
	}
	return false
}

func (c *Config) IsMatchingSymbol(name string) bool {
	for _, re := range c.SymbolNameRegexpList {
		if re.MatchString(name) {
			return true
		}
	}
	return false
}

func (c *Config) WantSave() bool {
	return c.OutputFile != ""
}
