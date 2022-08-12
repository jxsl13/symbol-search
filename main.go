package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"symbol-search/nm"
	"sync/atomic"

	"github.com/jxsl13/cwalk"
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

func main() {
	canvas := NewCanvas()
	defer canvas.Close()

	err := cwalk.Walk(RootPath, WalkFunc(canvas, Matchers))
	checkErr(err)
}

func plural(i int64) string {
	if i != 1 {
		return "s"
	}
	return ""
}

func WalkFunc(canvas *Canvas, matchers []*regexp.Regexp) filepath.WalkFunc {
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

		sf, err := nm.GetFilteredSymbols(path, matchers)
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
			t.AppendRow([]interface{}{path, s.Name, s.Version, s.Library})
		}
		return nil
	}

}

func checkErr(err error) {
	if err != nil {
		if er, ok := err.(cwalk.WalkerErrorList); ok {
			fmt.Println("We encountered a few errors along the way which ")
			m := make(map[string]struct{}, len(er.ErrorList))

			for _, e := range er.ErrorList {
				m[e.Error()] = struct{}{}
			}
			er.ErrorList = er.ErrorList[:0]

			sl := make([]string, 0, len(m))
			for k := range m {
				sl = append(sl, k)
			}

			sort.Strings(sl)

			for _, s := range sl {
				fmt.Println(s)
			}

		} else {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}
	}
}
