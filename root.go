package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jxsl13/archivewalker"
	"github.com/jxsl13/cli-config-boilerplate/cliconfig"
	"github.com/jxsl13/symbol-search/config"
	"github.com/jxsl13/symbol-search/nm"
	"github.com/spf13/cobra"
)

func NewRootCmd(ctx context.Context) *cobra.Command {
	root := RootContext{
		stdout: os.Stdout,
		stderr: os.Stderr,
		ctx:    ctx,
		cfg:    config.NewConfig(),
		table:  NewTable(),
		canvas: nil,
	}

	cmd := cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: "utility for analyzing symbols in binaries, shared objects inside and outside of archives",
		Args:  cobra.MinimumNArgs(0),
	}

	cmd.SetOut(root.stdout)
	cmd.SetErr(root.stderr)

	cmd.PersistentPreRunE = root.PreRunE(&cmd)
	cmd.RunE = root.RunE
	cmd.PersistentPostRunE = root.PostRunE
	cmd.AddCommand(NewCompletionCommand(&cmd))

	return &cmd
}

type RootContext struct {
	stdout io.Writer
	stderr io.Writer

	start     time.Time
	ctx       context.Context
	cfg       config.Config
	table     *SyncTable
	canvas    *Canvas
	lastPaint time.Time

	concurrency chan struct{}

	filesSeen            int64
	permissionError      int64
	failedSymbolReads    int64
	succeededSymbolReads int64
}

func (cli *RootContext) IncFilesSeen() int64 {
	return atomic.AddInt64(&cli.filesSeen, 1)
}

func (cli *RootContext) LoadFilesSeen() int64 {
	return atomic.LoadInt64(&cli.filesSeen)
}

func (cli *RootContext) IncPermissionError() int64 {
	return atomic.AddInt64(&cli.permissionError, 1)
}

func (cli *RootContext) LoadPermissionError() int64 {
	return atomic.LoadInt64(&cli.permissionError)
}

func (cli *RootContext) IncFailedSymbolReads() int64 {
	return atomic.AddInt64(&cli.failedSymbolReads, 1)
}

func (cli *RootContext) LoadFailedSymbolReads() int64 {
	return atomic.LoadInt64(&cli.failedSymbolReads)
}

func (cli *RootContext) IncSucceededSymbolReads() int64 {
	return atomic.AddInt64(&cli.succeededSymbolReads, 1)
}

func (cli *RootContext) LoadSucceededSymbolReads() int64 {
	return atomic.LoadInt64(&cli.succeededSymbolReads)
}

func (cli *RootContext) ElapsedTime() time.Duration {
	return time.Since(cli.start)
}

func (cli *RootContext) Printf(format string, args ...interface{}) {
	fmt.Fprintf(cli.stdout, format, args...)
}

func (cli *RootContext) PreRunE(cmd *cobra.Command) func(*cobra.Command, []string) error {
	cfgParser := cliconfig.RegisterFlags(&cli.cfg, false, cmd, cliconfig.WithoutConfigFile())
	return func(cmd *cobra.Command, args []string) error {

		if len(args) > 0 {
			cli.cfg.SymbolNameRegexList = args
		}

		err := cfgParser()
		if err != nil {
			return err
		}

		cli.start = time.Now()

		// limit the number of concurrent workers
		cli.canvas = NewCanvas(cli.stdout)
		cli.lastPaint = time.Now().Add(-cli.cfg.PaintInterval)
		cli.concurrency = make(chan struct{}, cli.cfg.Concurrency)

		return nil
	}
}

func (cli *RootContext) PostRunE(*cobra.Command, []string) error {
	return nil
}

func (cli *RootContext) RunE(cmd *cobra.Command, args []string) (err error) {
	defer func() {
		err = errors.Join(err, cli.Done())
	}()

	err = filepath.Walk(cli.cfg.SearchDir, func(path string, info fs.FileInfo, err error) error {
		select {
		case cli.concurrency <- struct{}{}:
			go func(path string, info fs.FileInfo, err error) error {
				defer func() {
					<-cli.concurrency
				}()
				return cli.WalkFunc(path, info, err)
			}(path, info, err)

		case <-cli.ctx.Done():
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return err
	}

	// final paint
	cli.Paint()

	if !cli.cfg.WantSave() {
		return nil
	}

	data := cli.canvas.Buffer()
	return os.WriteFile(cli.cfg.OutputFile, []byte(data), 0700)
}

func plural(i int64) string {
	if i != 1 {
		return "s"
	}
	return ""
}

func (cli *RootContext) Done() error {
	select {
	case <-cli.ctx.Done():
		return cli.ctx.Err()
	default:
		return nil
	}
}

func (cli *RootContext) WalkFunc(filePath string, info fs.FileInfo, err error) (rerr error) {
	defer func() {
		if rerr != nil {
			rerr = fmt.Errorf("failed to walk path %s: %w", filePath, rerr)
			rerr = cli.mapFSError(rerr)
		}
	}()

	if err != nil {
		return err
	}

	err = cli.Done()
	if err != nil {
		return filepath.SkipAll
	}

	if info.IsDir() || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return nil
	}

	// make all paths unix paths
	unixPath := toUnixPath(filePath)

	if cli.cfg.IsMatchingArchive(unixPath) {
		err = archivewalker.Walk(
			filePath,
			cli.ArchiveWalker(unixPath),
		)
		if err != nil {
			return err
		}

		return nil
	}

	if !cli.cfg.IsMatchingFileMode(info.Mode()) {
		return nil
	}

	if !cli.cfg.IsMatchingFilePath(unixPath) {
		return nil
	}

	fileName := path.Base(unixPath)
	if !cli.cfg.IsMatchingFileName(fileName) {
		// do not read the file into memory if
		// its path is not expected to be the file that we are looking for
		return nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	cli.SymbolWalker(unixPath, f, info)
	return nil

}

func (cli *RootContext) ArchiveWalker(archivePath string) archivewalker.WalkFunc {
	return func(filePath string, info fs.FileInfo, r io.Reader, err error) error {
		if err != nil {
			return fmt.Errorf("archive walker failed for archive %s: %w", archivePath, err)
		}

		err = cli.Done()
		if err != nil {
			return fmt.Errorf("archive walker failed for archive %s: %w", archivePath, err)
		}

		if info.IsDir() || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		if !cli.cfg.IsMatchingFileMode(info.Mode()) {
			return nil
		}

		// make all paths unix paths
		unixPath := toUnixPath(filePath)
		fileName := path.Base(unixPath)
		if !cli.cfg.IsMatchingFileName(fileName) {
			// do not read the file into memory if
			// its path is not expected to be the file that we are looking for
			return nil
		}

		dir := path.Dir(unixPath)
		if !cli.cfg.IsMatchingFilePath(dir) {
			return nil
		}

		// read file into memory
		f, err := archivewalker.NewFile(r, info.Size())
		if err != nil {
			return fmt.Errorf("archive walker failed fo archive %s: failed to read file %s into memory: %w", archivePath, filePath, err)
		}

		fullUnixPath := strings.Join([]string{archivePath, unixPath}, string(filepath.ListSeparator))
		cli.SymbolWalker(fullUnixPath, f, info)
		return nil
	}
}

func (cli *RootContext) mapFSError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, os.ErrPermission) {
		cli.IncPermissionError()
		return nil
	}

	cli.canvas.Error(err)
	return nil
}

func toUnixPath(filepath string) string {
	return path.Join("/", strings.ReplaceAll(filepath, "\\", "/"))
}

func (cli *RootContext) AppendSymbols(filePath string, symbols []nm.Symbol) {
	for _, s := range symbols {
		cli.table.AppendSymbol(filePath, s)
	}
}

func (cli *RootContext) Paint() {
	var (
		visited       = cli.LoadFilesSeen()
		denied        = cli.LoadPermissionError()
		readFails     = cli.LoadFailedSymbolReads()
		readSucceeded = cli.LoadSucceededSymbolReads()
	)
	caption := fmt.Sprintf("%d file%s visited\n%d access denied\n%d symbol read%s failed\n%d symbol read%s succeeded\ntime elapsed: %s\n",
		visited,
		plural(visited),
		denied,
		readFails,
		plural(readFails),
		readSucceeded,
		plural(readSucceeded),
		cli.ElapsedTime(),
	)

	cli.table.SetCaption(caption)

	clear := !cli.cfg.Debug // keep output history in temrinal if debug is enabled
	cli.canvas.Paint(cli.table, clear)
}

func (cli *RootContext) ThrottledPaint() {
	if time.Since(cli.lastPaint) > cli.cfg.PaintInterval {
		cli.Paint()
		cli.lastPaint = time.Now()
	}
}

// filepath is either a normal file path or a concatenated filespath that points into an archive
func (cli *RootContext) SymbolWalker(filePath string, file archivewalker.File, info fs.FileInfo) {
	numErr := 0
	numTries := 0
	numSymbols := 0
	mode := info.Mode()

	defer func() {
		cli.IncFilesSeen()
		if numErr == numTries {
			cli.IncFailedSymbolReads()
			if cli.cfg.Debug {
				cli.Printf("read failed for file: %s (tries=%d, mode=%#o)\n", filePath, numTries, mode)
			}
		} else {
			cli.IncSucceededSymbolReads()
			if cli.cfg.Debug {
				cli.Printf("read succeeded for file: %s (tries=%d symbols=%d mode=%#o)\n", filePath, numTries, numSymbols, mode)
			}
		}

		cli.ThrottledPaint() // is throttled by the configuration
	}()

	if !cli.cfg.NoELF {
		numTries++
		symbols, eerr := cli.ReadELF(file)
		if eerr == nil {
			numSymbols = len(symbols)
			cli.AppendSymbols(filePath, symbols)
			return
		}
		numErr++
	}

	// must at least be as big as the header size of an ar archive (*.a)
	// https://www.abhirag.com/blog/ar/
	if !cli.cfg.NoStatic && info.Size() > 68 && strings.HasSuffix(filePath, ".a") {
		numTries++
		symbols, aerr := cli.ReadAR(filePath, file)
		if aerr == nil {
			numSymbols = len(symbols)
			cli.AppendSymbols(filePath, symbols)
			return
		}
		numErr++
	}

	if !cli.cfg.NoPE {
		numTries++
		symbols, perr := cli.ReadPE(file)
		if perr == nil {
			numSymbols = len(symbols)
			cli.AppendSymbols(filePath, symbols)
			return
		}
		numErr++
	}
}

func (cli *RootContext) ReadAR(libraryPath string, file archivewalker.File) (_ []nm.Symbol, err error) {
	defer func() {
		_, serr := file.Seek(0, io.SeekStart)
		if serr != nil {
			err = errors.Join(err, serr)
		}
	}()

	result := make([]nm.Symbol, 0, 64)

	err = nm.WalkAR(file, func(filePath string, info fs.FileInfo, r io.Reader, err error) error {
		if err != nil {
			return err
		}

		if filePath == "" {
			return nil
		}

		if info.IsDir() || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		if !strings.HasSuffix(filePath, ".o") {
			return nil
		}

		fullUnixPath := strings.Join([]string{libraryPath, filePath}, string(filepath.ListSeparator))
		buf, err := archivewalker.NewFile(r, info.Size())
		if err != nil {
			return fmt.Errorf("failed to read file %s into memory: %w", fullUnixPath, err)
		}

		symbols, err := cli.ReadMachO(buf)
		if err != nil {
			return nil
		}
		result = append(result, symbols...)
		return nil
	})

	return result, err
}

func (cli *RootContext) ReadMachO(file archivewalker.File) (_ []nm.Symbol, err error) {
	defer func() {
		_, serr := file.Seek(0, io.SeekStart)
		if serr != nil {
			err = errors.Join(err, serr)
		}
	}()

	f, err := nm.OpenMachO(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make([]nm.Symbol, 0, 1)
	if !cli.cfg.NoImported {
		is, err := nm.ReadImportedSymbolsMachO(f)
		if err != nil {
			return nil, err
		}

		for _, s := range is {
			if !cli.cfg.IsMatchingSymbol(s.Name) {
				continue
			}
			result = append(result, s)
		}
	}

	if !cli.cfg.NoDynamic {
		ds, err := nm.ReadDynamicSymbolsMachO(f)
		if err != nil {
			return nil, err
		}
		for _, s := range ds {
			if !cli.cfg.IsMatchingSymbol(s.Name) {
				continue
			}
			result = append(result, s)
		}
	}

	return result, nil
}

func (cli *RootContext) ReadELF(file archivewalker.File) (_ []nm.Symbol, err error) {
	defer func() {
		_, serr := file.Seek(0, io.SeekStart)
		if serr != nil {
			err = errors.Join(err, serr)
		}
	}()

	elf, err := nm.OpenELF(file)
	if err != nil {
		return nil, err
	}

	result := make([]nm.Symbol, 0, 16)
	if !cli.cfg.NoImported {
		is, err := nm.ReadImportedSymbolsELF(elf)
		if err != nil {
			return nil, err
		}

		for _, s := range is {
			if !cli.cfg.IsMatchingSymbol(s.Name) {
				continue
			}
			result = append(result, s)
		}
	}

	if !cli.cfg.NoDynamic {
		ds, err := nm.ReadDynamicSymbolsELF(elf)
		if err != nil {
			return nil, err
		}
		for _, s := range ds {
			if !cli.cfg.IsMatchingSymbol(s.Name) {
				continue
			}
			result = append(result, s)
		}
	}

	if !cli.cfg.NoInternal {
		ss, err := nm.ReadInternalSymbolsELF(elf)
		if err != nil {
			return nil, err
		}
		for _, s := range ss {
			if !cli.cfg.IsMatchingSymbol(s.Name) {
				continue
			}
			result = append(result, s)
		}
	}

	return result, nil
}

func (cli *RootContext) ReadPE(file archivewalker.File) (_ []nm.Symbol, err error) {
	defer func() {
		_, serr := file.Seek(0, io.SeekStart)
		if serr != nil {
			err = errors.Join(err, serr)
		}
	}()

	pe, err := nm.OpenPE(file)
	if err != nil {
		return nil, err
	}

	result := make([]nm.Symbol, 0, 64)
	if !cli.cfg.NoImported {
		is, err := nm.ReadImportedSymbolsPE(pe)
		if err != nil {
			return nil, err
		}

		for _, s := range is {
			if !cli.cfg.IsMatchingSymbol(s.Name) {
				continue
			}
			result = append(result, s)
		}
	}

	if !cli.cfg.NoInternal {
		ss, err := nm.ReadCOFFSymbolsPE(pe)
		if err != nil {
			return nil, err
		}
		for _, s := range ss {
			if !cli.cfg.IsMatchingSymbol(s.Name) {
				continue
			}
			result = append(result, s)
		}
	}

	return result, nil
}
