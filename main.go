package main

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/infrastructure/CodeStorage"
	"github.com/AntonParaskiv/mockGen/interfaces/AstRepository"
	"github.com/AntonParaskiv/mockGen/interfaces/Printer"
	"github.com/AntonParaskiv/mockGen/usecases/Interactor"
	"os"
	"path/filepath"
)

func main() {

	// get arg interface package path
	if len(os.Args[1:]) == 0 {
		usage()
		os.Exit(1)
	}
	interfacePackagePath, err := filepath.Abs(os.Args[1])
	if err != nil {
		err = fmt.Errorf("get absolute interface package path failed: %w", err)
		fmt.Println(err)
		os.Exit(1)
	}

	// init app components
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

	// start generating
	astPackage, err := codeStorage.GetAstPackage(interfacePackagePath)
	if err != nil {
		err = fmt.Errorf("get ast package failed: %w", err)
		fmt.Println(err)
		os.Exit(1)
	}
	if astPackage == nil {
		err = fmt.Errorf("ast package not found")
		fmt.Println(err)
		os.Exit(1)
	}

	interfacePackage, err := astRepository.CreateInterfacePackage(astPackage, interfacePackagePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mockPackage := interactor.CreateMockPackage(interfacePackage)
	printer.GeneratePackageCode(mockPackage)

	err = codeStorage.SaveGoPackage(mockPackage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}

func usage() {
	fmt.Println("Usage: mockGen packagePath")
	return
}
