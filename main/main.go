package main

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain/Interface"
	"github.com/AntonParaskiv/mockGen/domain/Method"
	"github.com/AntonParaskiv/mockGen/domain/Variable"
	"github.com/AntonParaskiv/mockGen/infrastructure/FileStorage"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

func main() {
	// TODO: init components
	fileStorageFactory := FileStorage.NewFactory()

	// TODO: arg get package path
	fileName := "examples/ManagerInterface/Manager.go"

	// TODO: get package ast

	fileStorage := fileStorageFactory.Create(fileName)
	fileData, err := fileStorage.ReadFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", fileData, 0)
	if err != nil {
		fmt.Println("parse failed:", err.Error())
		return
	}
	myInterface := getInterface(f)

	// TODO: generate mock structure
	result := genStructFromInterface(myInterface)
	formattedBytes, err := format.Source([]byte(result))
	if err != nil {
		fmt.Println("formatting failed:", err.Error())
		return
	}
	result = string(formattedBytes)

	// TODO: save mock
	fmt.Println(result)
}

func getInterface(f *ast.File) (myInterface *Interface.Interface) {

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

					myInterface = &Interface.Interface{}
					myInterface.Name = spec.Name.Name

					// методы интерфейса
					for _, method := range specType.Methods.List {

						switch methodType := method.Type.(type) {
						case *ast.FuncType:
							myMethod := &Method.Method{}
							myMethod.Name = method.Names[0].Name // TODO: check

							for _, param := range methodType.Params.List {
								myVariable := &Variable.Variable{}
								myVariable.Name = param.Names[0].Name          // TODO: check
								myVariable.Type = param.Type.(*ast.Ident).Name // TODO: check
								myMethod.ArgList = append(myMethod.ArgList, myVariable)
							}

							for _, result := range methodType.Results.List {
								myVariable := &Variable.Variable{}
								myVariable.Name = result.Names[0].Name          // TODO: check
								myVariable.Type = result.Type.(*ast.Ident).Name // TODO: check
								myMethod.ValueList = append(myMethod.ValueList, myVariable)
							}

							myInterface.MethodList = append(myInterface.MethodList, myMethod)
						}

					}

				}
			}
		}
	}
	return
}

func genStructFromInterface(i *Interface.Interface) (result string) {
	result += fmt.Sprintf("type %s struct {\n", i.Name)

	fieldList := getFieldList(i.MethodList)
	for _, field := range fieldList {
		result += fmt.Sprintf("	%s %s\n", field.Name, field.Type)
	}
	result += fmt.Sprintf("}\n")
	return
}

func getFieldList(methodList []*Method.Method) (fieldList []*Variable.Variable) {
	fieldList = make([]*Variable.Variable, 0)
	for _, method := range methodList {
		for _, variable := range method.ArgList {
			if isFieldListContainsVariable(fieldList, variable) {
				continue
			}
			fieldList = append(fieldList, variable)
		}
	}
	for _, method := range methodList {
		for _, variable := range method.ValueList {
			if isFieldListContainsVariable(fieldList, variable) {
				continue
			}
			fieldList = append(fieldList, variable)
		}
	}
	return
}

func isFieldListContainsVariable(fieldList []*Variable.Variable, v *Variable.Variable) (isContains bool) {
	for _, field := range fieldList {
		if field.Name == v.Name {
			isContains = true
			break
		}
	}
	return
}
