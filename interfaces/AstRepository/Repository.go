package AstRepository

import (
	"bytes"
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"github.com/AntonParaskiv/mockGen/interfaces/PackageStorageInterface"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"path/filepath"
)

type Repository struct {
	packageStorage PackageStorageInterface.Storage
}

func New() (r *Repository) {
	r = new(Repository)
	return
}

func (r *Repository) SetPackageStorage(packageStorage PackageStorageInterface.Storage) *Repository {
	r.packageStorage = packageStorage
	return r
}

func (r *Repository) ScanPackage(path string) (resultPackage *domain.Package) {

	fileList, err := r.packageStorage.GetGoFileList(path)
	if err != nil {
		err = fmt.Errorf("get go file list failed: %w", err)
		return
	}

	resultPackage = domain.NewPackage()
	resultPackage.Name = filepath.Base(path)

	for _, filePath := range fileList {
		file, err := r.ScanGoCodeFile(filePath)
		if err != nil {
			err = fmt.Errorf("scan file %s failed: %w", filePath, err)
			return
		}
		resultPackage.Files = append(resultPackage.Files, file)
	}

	return
}

func (r *Repository) ScanGoCodeFile(path string) (file *domain.File, err error) {
	data, err := r.packageStorage.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("read file %s failed: %w", path, err)
		return
	}

	fileSet := token.NewFileSet()
	fileAst, err := parser.ParseFile(fileSet, "", data, 0)
	if err != nil {
		err = fmt.Errorf("parse ast failed: %w", err)
		return
	}

	fileAst.Scope.

	fileSetWrite := token.NewFileSet()
	var buf bytes.Buffer
	err = printer.Fprint(&buf, fileSetWrite, fileAst)
	if err != nil {
		err = fmt.Errorf("print ast failed: %w", err)
		return
	}

	return
	file = domain.NewFile()
	file.Name = filepath.Base(path)
	file.InterfaceList = getInterfaceList(fileAst)

	return
}

func getInterfaceList(f *ast.File) (interfaceList []*domain.Interface) {
	interfaceList = make([]*domain.Interface, 0)

	for _, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			switch decl.Tok {

			// объявления типов
			case token.TYPE:
				spec := decl.Specs[0].(*ast.TypeSpec) // TODO: check

				switch specType := spec.Type.(type) {

				// тип interface
				case *ast.InterfaceType:
					iFace := &domain.Interface{}
					iFace.Name = spec.Name.Name

					// методы интерфейса
					for _, method := range specType.Methods.List {

						switch methodType := method.Type.(type) {
						case *ast.FuncType:
							myMethod := &domain.Method{}
							myMethod.Name = method.Names[0].Name // TODO: check

							for _, param := range methodType.Params.List {
								myVariable := &domain.Variable{}
								myVariable.Name = param.Names[0].Name          // TODO: check
								myVariable.Type = param.Type.(*ast.Ident).Name // TODO: check
								myMethod.ArgList = append(myMethod.ArgList, myVariable)
							}

							for _, result := range methodType.Results.List {
								myVariable := &domain.Variable{}
								myVariable.Name = result.Names[0].Name          // TODO: check
								myVariable.Type = result.Type.(*ast.Ident).Name // TODO: check
								myMethod.ValueList = append(myMethod.ValueList, myVariable)
							}

							iFace.MethodList = append(iFace.MethodList, myMethod)
						}

					}

					interfaceList = append(interfaceList, iFace)
				}
			}
		}
	}

	return
}
