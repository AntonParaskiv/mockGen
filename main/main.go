package main

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"github.com/AntonParaskiv/mockGen/interfaces/AstRepository"
	"github.com/AntonParaskiv/mockGen/usecases"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	interfacePackagePath := "examples/ManagerInterface"

	astRepository := AstRepository.Repository{}
	interactor := usecases.Interactor{}

	interfacePackage, err := astRepository.CreateInterfacePackage(interfacePackagePath)
	if err != nil {
		panic(err)
	}

	mockPackage := interactor.CreateMockPackage(interfacePackage)

	err = SaveGoPackage(mockPackage)
	if err != nil {
		panic(err)
	}

	return
}

func SaveGoPackage(Package *domain.GoCodePackage) (err error) {
	err = os.MkdirAll(Package.Path, 0755)
	if err != nil {
		err = fmt.Errorf("create dir %s failed: %w", Package.Path, err)
		return
	}

	for _, file := range Package.FileList {
		filePath := filepath.Join(Package.Path, file.Name)

		//var formattedCode []byte
		//formattedCode, err = imports.Process("", []byte(file.Code), &imports.Options{
		//	Fragment:   false,
		//	AllErrors:  true,
		//	Comments:   true,
		//	TabIndent:  true,
		//	TabWidth:   8,
		//	FormatOnly: false,
		//})
		//if err != nil {
		//	err = fmt.Errorf("format file %s code failed: %s", file.Name, err)
		//	return
		//}

		//err = ioutil.WriteFile(filePath, formattedCode, 0644)
		err = ioutil.WriteFile(filePath, []byte(file.Code), 0644)
		if err != nil {
			err = fmt.Errorf("write file %s failed: %w", file.Name, err)
			return
		}
	}
	return
}
