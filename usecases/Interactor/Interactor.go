package Interactor

import (
	"github.com/AntonParaskiv/mockGen/domain"
	"github.com/AntonParaskiv/mockGen/interfaces/AstRepository"
)

type Interactor struct {
	AstRepository       AstRepository.Repository
	mockFile            *domain.GoCodeFile
	interfacePackage    *domain.GoCodePackage
	mockPackage         *domain.GoCodePackage
	CreateFieldExamples bool
}
