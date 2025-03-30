package main

import (
	"testing"

	"github.com/jxsl13/symbol-search/internal/testutils"
)

func TestSymbolSearch(t *testing.T) {
	path := "~/Schreibtisch/github/snk-gnu-linux/3rd/sybase/lib"

	out, err := testutils.Execute(
		NewRootCmd(t.Context()),
		"-t",
		"1",
		"-v",
		"-n",
		".*\\.a",
		"-s",
		"unisem",
		"-f",
		path,
	)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(out.String())
}
