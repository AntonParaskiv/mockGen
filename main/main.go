package main

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/infrastructure/CodeStorage"
	"github.com/AntonParaskiv/mockGen/interfaces/AstRepository"
	"github.com/AntonParaskiv/mockGen/interfaces/Printer"
	"github.com/AntonParaskiv/mockGen/usecases/Interactor"
	"os"
)

func main() {

	if len(os.Args[1:]) == 0 {
		usage()
		return
	}

	interfacePackagePath := os.Args[1]
	//interfacePackagePath := "examples/ManagerInterface"

	codeStorage := CodeStorage.Storage{
		FormatEnabled: true,
	}
	astRepository := AstRepository.Repository{
		CodeStorage: codeStorage,
	}
	interactor := Interactor.Interactor{
		AstRepository:       astRepository,
		CreateFieldExamples: true,
	}
	printer := Printer.Printer{
		GenerateTests: true,
	}

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
	printer.GeneratePackageCode(mockPackage)

	err = codeStorage.SaveGoPackage(mockPackage)
	if err != nil {
		panic(err)
	}

	return
}

func usage() {
	fmt.Println("usage")
	return
}
