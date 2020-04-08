package AstRepository

import (
	"github.com/AntonParaskiv/mockGen/infrastructure/CodeStorage"
)

type Repository struct {
	CodeStorage CodeStorage.Storage
}
