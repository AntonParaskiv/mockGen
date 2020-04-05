package main

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/infrastructure/CodeStorage"
	"github.com/AntonParaskiv/mockGen/interfaces/AstRepository"
	"github.com/AntonParaskiv/mockGen/interfaces/Printer"
	"github.com/AntonParaskiv/mockGen/usecases"
)

func main() {
	interfacePackagePath := "examples/ManagerInterface"

	codeStorage := CodeStorage.Storage{
		FormatEnabled: false,
	}
	astRepository := AstRepository.Repository{
		CodeStorage: codeStorage,
	}
	interactor := usecases.Interactor{
		AstRepository: astRepository,
	}
	printer := Printer.Printer{}

	astPackage, err := codeStorage.GetAstPackage(interfacePackagePath)
	if err != nil {
		err = fmt.Errorf("get ast package failed: %w", err)
		panic(err)
	}
	if astPackage == nil {
		err = fmt.Errorf("ast package not found")
		panic(err)
	}

	interfacePackage, err := astRepository.CreateInterfacePackage(astPackage, interfacePackagePath)
	if err != nil {
		panic(err)
	}

	mockPackage := interactor.CreateMockPackage(interfacePackage)
	printer.GenerateCode(mockPackage)

	err = codeStorage.SaveGoPackage(mockPackage)
	if err != nil {
		panic(err)
	}

	return
}
