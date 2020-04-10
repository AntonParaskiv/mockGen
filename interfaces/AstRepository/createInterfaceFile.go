package AstRepository

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

func createInterfaceFile(astFile *ast.File, fullFileName string) (interfaceFile *domain.GoCodeFile, err error) {
	if !isFileNameMatchGoCode(fullFileName) {
		return
	}

	astInterfaceSpecs := getInterfaces(astFile)
	interfaceList := make([]*domain.Interface, 0)

	for _, astInterfaceSpec := range astInterfaceSpecs {
		var iFace *domain.Interface
		iFace, err = createInterface(astInterfaceSpec)
		if err != nil {
			err = fmt.Errorf("create interface failed: %w", err)
			return
		}
		if iFace == nil {
			continue
		}
		interfaceList = append(interfaceList, iFace)
	}
	if len(interfaceList) == 0 {
		return
	}

	interfaceFile = &domain.GoCodeFile{
		Name:          filepath.Base(fullFileName),
		ImportList:    getImportListFromAstFile(astFile),
		InterfaceList: interfaceList,
	}

	return
}

func isFileNameMatchGoCode(fileName string) (isMatch bool) {
	// check non-go files
	fileExtension := filepath.Ext(fileName)
	if fileExtension != ".go" {
		return
	}

	// check test files
	pattern := "_test.go"
	fileNameEndingStartPosition := len(fileName) - len(pattern)
	if fileNameEndingStartPosition < 0 {
		isMatch = true
		return
	}
	fileNameEnding := fileName[fileNameEndingStartPosition:]
	if fileNameEnding != pattern {
		isMatch = true
		return
	}

	return
}

func getInterfaces(f *ast.File) (interfaceSpecs []*ast.TypeSpec) {
	interfaceSpecs = make([]*ast.TypeSpec, 0)

	for _, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			switch decl.Tok {

			// объявления типов
			case token.TYPE:
				spec := decl.Specs[0].(*ast.TypeSpec) // TODO: check array

				switch spec.Type.(type) {

				// тип interface
				case *ast.InterfaceType:
					interfaceSpecs = append(interfaceSpecs, spec)
				}
			}
		}
	}

	return
}

func getImportListFromAstFile(astFile *ast.File) (importList []*domain.Import) {
	for _, spec := range astFile.Imports {
		Import := &domain.Import{
			Name: getNodeName(spec),
			Path: strings.Trim(spec.Path.Value, `"`),
		}
		importList = append(importList, Import)
	}
	return
}
