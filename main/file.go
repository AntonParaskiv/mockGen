package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
)

func newAstFile() (file *ast.File) {
	file = &ast.File{
		Decls: []ast.Decl{},
	}
	return
}

func newAstTestFile() (file *ast.File) {
	file = newAstFile()
	file.Imports = []*ast.ImportSpec{
		{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"testing"`,
			},
		},
	}
	return
}

func createWFiles(rFiles map[string]*ast.File) (wFiles map[string]*ast.File, err error) {
	wFiles = make(map[string]*ast.File)

	for rFileRelativePath, rFile := range rFiles {
		if !isFileNameMatchGoCode(rFileRelativePath) {
			continue
		}

		interfaceSpecs := getInterfaces(rFile)
		if len(interfaceSpecs) == 0 {
			continue
		}

		mockPackageName := createMockPackageName(getNodeName(rFile))

		wFile := newAstFile()
		wFile.Name = &ast.Ident{
			Name: mockPackageName,
		}

		wFileTest := newAstTestFile()
		wFileTest.Name = &ast.Ident{
			Name: mockPackageName,
		}

		var structSpec *ast.TypeSpec
		for _, interfaceSpec := range interfaceSpecs {

			// gen struct
			structSpec, err = createStructFromInterfaceSpec(interfaceSpec)
			if err != nil {
				err = fmt.Errorf("create struct from interface %s failed: %w", getNodeName(interfaceSpec), err)
				return
			}
			decl := &ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					structSpec,
				},
			}
			wFile.Decls = append(wFile.Decls, decl)

			// gen constructor
			// TODO: gen test
			constructorDecl, constructorTestDecl := createConstructorAndTest(structSpec, unMockPackageName(mockPackageName))
			wFile.Decls = append(wFile.Decls, constructorDecl)
			wFileTest.Decls = append(wFileTest.Decls, constructorTestDecl)

			// gen setters
			// TODO: gen test
			setterDecls := createSetters(structSpec)
			for _, setterDecl := range setterDecls {
				wFile.Decls = append(wFile.Decls, setterDecl)
			}

			// gen methods
			// TODO: gen test
			methodsDecls := createMethods(structSpec, interfaceSpec)
			for _, methodsDecl := range methodsDecls {
				wFile.Decls = append(wFile.Decls, methodsDecl)
			}
		}

		wFilePath := createMockFilePath(rFileRelativePath)
		wFileTestPath := createTestFilePath(wFilePath)

		wFiles[wFilePath] = wFile
		wFiles[wFileTestPath] = wFileTest

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

func createMockFilePath(filePath string) (mockFilePath string) {
	fileName := filepath.Base(filePath)
	dirPath := filepath.Dir(filePath)

	dirPath = cutPostfix(dirPath, "Interface")
	mockDirPath := dirPath + "Mock"

	mockFilePath = filepath.Join(mockDirPath, fileName)
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
