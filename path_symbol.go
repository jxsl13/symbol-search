package main

import "github.com/jxsl13/symbol-search/nm"

type PathSymbol struct {
	Type string
	Path string
	nm.Symbol
}

func NewPathSymbolList(filePath string, symbols []nm.Symbol) []PathSymbol {
	if len(symbols) == 0 {
		return []PathSymbol{}
	}

	result := make([]PathSymbol, 0, len(symbols))
	for _, symbol := range symbols {
		result = append(result, PathSymbol{
			Path:   filePath,
			Symbol: symbol,
		})
	}
	return result
}
