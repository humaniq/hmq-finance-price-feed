package config

var DefaultNativeCurrencyList = []string{"eth", "usd", "eur", "rub", "bnb"}

func DefaultSymbolDiffs() map[string]int {
	mapper := make(map[string]int)
	for key, _ := range SymbolsKnown {
		mapper[key] = 1
	}
	return mapper
}
