package main

import (
	"testing"

	"github.com/jxsl13/symbol-search/internal/testutils"
)

func TestSymbolSearch(t *testing.T) {
	path := "/home/behm015"

	out, err := testutils.Execute(
		NewRootCmd(t.Context()),
		"-v",
		"-s",
		".*",
		"-f",
		path,
	)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(out.String())

}
