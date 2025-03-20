package main

import (
	"testing"

	"github.com/jxsl13/symbol-search/internal/testutils"
)

func TestSymbolSearch(t *testing.T) {

	path := "/DBA/gnu/snk-1.5.7_x64/lib/perl5/site_perl/5.36.1/x86_64-linux-thread-multi/auto/Sybase/"

	out, err := testutils.Execute(
		NewRootCmd(t.Context()),
		"-v",
		"-s",
		"unisem_RegisterCallbacks",
		"-f",
		path,
	)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(out.String())

}
