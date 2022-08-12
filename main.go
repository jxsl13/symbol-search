package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"symbol-search/nm"
	"sync/atomic"

	"github.com/iafan/cwalk"
)

var (
	Matchers []*regexp.Regexp
	RootPath string
)

func init() {
	if len(os.Args) != 3 {
		log.Fatalln("please provide a symbol name (or , separated list) as the first argument and a path as the second argument")
	}

	strmatchers := strings.Split(os.Args[1], ",")
	for idx, m := range strmatchers {
		strmatchers[idx] = strings.TrimSpace(m)
	}

	for _, matcher := range strmatchers {
		r, err := regexp.Compile(matcher)
		if err != nil {
			log.Fatalln("invalid regular expression: ", matcher)
		}
		Matchers = append(Matchers, r)
	}

	RootPath = os.Args[2]
}

func checkErr(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v", err)
	os.Exit(1)
}

func main() {
	canvas := NewCanvas()
	defer canvas.Close()

	fi, err := os.Stat(RootPath)
	checkErr(err)

	// single file
	if !fi.IsDir() {
		err = WalkFunc(canvas, RootPath, Matchers)("", fi, nil)
		checkErr(err)
		return
	}

	err = WalkDirectory(canvas, RootPath, Matchers)
	checkErr(err)
}

func plural(i int64) string {
	if i != 1 {
		return "s"
	}
	return ""
}

func WalkFunc(canvas *Canvas, rootPath string, matchers []*regexp.Regexp) filepath.WalkFunc {
	permErrCnt := int64(0)
	seenCnt := int64(0)

	t := NewTable()

	return func(path string, info fs.FileInfo, err error) error {
		defer func() {
			visited := atomic.AddInt64(&seenCnt, 1)
			denied := atomic.LoadInt64(&permErrCnt)

			t.SetCaption("%d file%s visited, %d access denied", visited, plural(visited), denied)
			canvas.Paint(t.Render())
		}()
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				atomic.AddInt64(&permErrCnt, 1)
				return nil
			}
			// too many file descriptors or so
			// collect but do not abort
			canvas.Error(err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		// everything below 10kb woul dnot be able to have
		// enough functions in order to execute code
		if info.Size() < 10_000 {
			return nil
		}

		fullPath := filepath.Join(rootPath, path)
		sf, err := nm.GetFilteredSymbols(fullPath, matchers)
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				atomic.AddInt64(&permErrCnt, 1)
			}
			return nil
		}

		if !sf.HasSymbols() {
			return nil
		}

		for _, s := range sf.Symbols {
			t.AppendRow([]interface{}{sf.Path, s.Name, s.Version, s.Library})
		}

		return nil
	}

}

func WalkDirectory(canvas *Canvas, rootPath string, matchers []*regexp.Regexp) error {
	err := cwalk.Walk(rootPath, WalkFunc(canvas, rootPath, matchers))
	return err
}
