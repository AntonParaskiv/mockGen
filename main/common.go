package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

func getNodeName(node ast.Node) (name string) {
	switch nodeItem := node.(type) {
	case *ast.Package:
		name = nodeItem.Name
	case *ast.File:
		name = nodeItem.Name.Name
	case *ast.TypeSpec:
		name = nodeItem.Name.Name
	case *ast.Field:
		name = nodeItem.Names[0].Name
	case *ast.Ident:
		name = nodeItem.Name
	default:
		panic(fmt.Sprintf("no getting name case for type %T", node))
	}
	return
}

func toPublic(name string) (publicName string) {
	firstLetterUpper := strings.ToUpper(getFirstLetter(name))
	publicName = firstLetterUpper + getFollowingLetters(name)
	return
}

func toPrivate(name string) (privateName string) {
	firstLetterLower := strings.ToLower(getFirstLetter(name))
	privateName = firstLetterLower + getFollowingLetters(name)
	return
}

func getFirstLetter(text string) (firstLetter string) {
	firstLetter = text[0:1]
	return
}

func getFollowingLetters(text string) (followingLetters string) {
	followingLetters = text[1:]
	return
}

// Mock -> s
func getReceiverName(name string) (receiverName string) {
	receiverName = toPrivate(getFirstLetter(name))
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

func createTestFilePath(filePath string) (testFilePath string) {
	extension := filepath.Ext(filePath)
	if extension == ".go" {
		filePathLen := len(filePath)
		testFilePath = filePath[:filePathLen-3] + "_test.go"
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

func getAstPackage(packagePath string) (astPackage *ast.Package, err error) {
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

func cutPostfix(text, postfix string) (shortCutText string) {
	lenPostfix := len(postfix)
	if len(text) > lenPostfix {
		startPostfix := len(text) - lenPostfix
		packageNamePostfix := text[startPostfix:]
		if packageNamePostfix == postfix {
			shortCutText = text[0:startPostfix]
		}
	}
	return
}
