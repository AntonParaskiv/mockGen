package CodeStorage

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Storage struct {
}

func (s *Storage) GetAstPackage(packagePath string) (astPackage *ast.Package, err error) {
	fSet := token.NewFileSet()
	astPackageList, err := parser.ParseDir(fSet, packagePath, nil, 0)
	if err != nil {
		err = fmt.Errorf("parse ast dir failed: %w", err)
		return
	}

	for _, astPackageItem := range astPackageList {
		astPackage = astPackageItem
		return
	}

	return
}

func (s *Storage) SaveGoPackage(Package *domain.GoCodePackage) (err error) {
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
