package model

import (
	"github.com/gofrs/uuid"
)

type TextID uuid.UUID

type Text struct {
	ID    TextID
	Value string
}
