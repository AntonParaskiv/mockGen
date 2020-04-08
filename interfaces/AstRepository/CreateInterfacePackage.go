package AstRepository

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"go/ast"
	"os"
	"path/filepath"
)

func (r *Repository) CreateInterfacePackage(astPackage *ast.Package, packagePath string) (interfacePackage *domain.GoCodePackage, err error) {
	interfacePackage, err = initInterfacePackage(packagePath, astPackage.Name)
	if err != nil {
		err = fmt.Errorf("init interface package failed: %w", err)
		return
	}

	for fullFileName, astFile := range astPackage.Files {
		var interfaceFile *domain.GoCodeFile
		interfaceFile, err = createInterfaceFile(astFile, fullFileName)
		if err != nil {
			err = fmt.Errorf("create interface file failed: %w", err)
			return
		}
		if interfaceFile == nil {
			continue
		}

		interfacePackage.FileList = append(interfacePackage.FileList, interfaceFile)
	}

	return
}

func initInterfacePackage(packagePath, packageName string) (interfacePackage *domain.GoCodePackage, err error) {
	fullPackagePath, err := filepath.Abs(packagePath)
	if err != nil {
		err = fmt.Errorf("create full package path failed: %w", err)
		return
	}

	goPathSrc := filepath.Join(os.Getenv("GOPATH"), "src")

	goPathSrcPackagePath, err := filepath.Rel(goPathSrc, fullPackagePath)
	if err != nil {
		err = fmt.Errorf("create GOPATH/src package path failed: %w", err)
		return
	}

	selfImport := &domain.Import{
		Key:  filepath.Base(goPathSrcPackagePath),
		Name: "",
		Path: filepath.ToSlash(goPathSrcPackagePath),
	}

	interfacePackage = &domain.GoCodePackage{
		Path:        packagePath,
		FullPath:    fullPackagePath,
		PackageName: packageName,
		SelfImport:  selfImport,
	}
	return
}
