package model

type TextID string

type Statistics struct {
	TextID               TextID
	AlphabetSymbolsCount int
	AllSymbolsCount      int
	IsDuplicate          bool
}
