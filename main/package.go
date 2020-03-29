package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
)

func newPackage() (pkg *ast.Package, err error) {
	fSet := token.NewFileSet()
	pkg, err = ast.NewPackage(fSet, nil, nil, nil)
	if err != nil {
		err = fmt.Errorf("create ast package failed: %w", err)
		return
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

func makeWPackages(rPkgs map[string]*ast.Package) (wPkgs map[string]*ast.Package, err error) {
	wPkgs = make(map[string]*ast.Package)
	var wPkg *ast.Package

	for _, rPkg := range rPkgs {
		// create new wPkg
		wPkg, err = newPackage()
		if err != nil {
			return
		}

		// make package name
		wPkgName := createMockPackageName(getNodeName(rPkg))
		wPkg.Name = wPkgName
		wPkgs[wPkgName] = wPkg

		// make files
		wPkg.Files, err = createWFiles(rPkg.Files)
		if err != nil {
			return
		}
	}
	return
}

func savePackages(pkgs map[string]*ast.Package) (err error) {
	var fd *os.File

	for _, pkg := range pkgs {
		fSet := token.NewFileSet()

		for filePath, file := range pkg.Files {
			// create dir
			dirPath := filepath.Dir(filePath)
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				err = fmt.Errorf("create dir %s failed: %w", dirPath, err)
				return
			}

			// save file
			fd, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				err = fmt.Errorf("open file %s failed: %w", filePath, err)
				return
			}

			err = printer.Fprint(fd, fSet, file)
			if err != nil {
				err = fmt.Errorf("print ast to file %s failed: %w", filePath, err)
				return
			}

			err = fd.Close()
			if err != nil {
				err = fmt.Errorf("close file %s failed: %w", filePath, err)
				return
			}
		}

	}
	return
}

func createMockPackageName(packageName string) (mockPackageName string) {
	packageName = cutPostfix(packageName, "Interface")
	mockPackageName = packageName + "Mock"
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

func unMockPackageName(mockPackageName string) (packageName string) {
	packageName = cutPostfix(mockPackageName, "Mock")
	return
}
