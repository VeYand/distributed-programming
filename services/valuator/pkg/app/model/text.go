package model

type TextID string

type Text struct {
	ID    TextID
	Value string
	Count int
}
