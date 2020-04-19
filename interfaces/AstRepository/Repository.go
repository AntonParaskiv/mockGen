package AstRepository

import (
	"github.com/AntonParaskiv/mockGen/domain"
	"github.com/AntonParaskiv/mockGen/infrastructure/CodeStorage"
)

type Repository struct {
	CodeStorage   CodeStorage.Storage
	currentMethod *domain.Method
}
