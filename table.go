package main

import (
	"github.com/jedib0t/go-pretty/v6/table"
)

func NewTable() *SyncTable {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Path", "Name", "Version", "Library"})
	t.AppendSeparator()
	t.SortBy([]table.SortBy{
		{
			Number: 1,
			Mode:   table.Asc,
		},
		{
			Number: 2,
			Mode:   table.Asc,
		},
		{
			Number: 3,
			Mode:   table.Asc,
		},
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 3, AutoMerge: false},
		{Number: 4, AutoMerge: false},
	})
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	return &SyncTable{
		t: t,
	}
}
