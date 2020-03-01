package main

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"github.com/AntonParaskiv/mockGen/infrastructure/PackageStorage"
	"github.com/AntonParaskiv/mockGen/interfaces/AstRepository"
)

func main() {

	// TODO: init components
	packageStorage := PackageStorage.New()
	astRepository := AstRepository.New().SetPackageStorage(packageStorage)

	// TODO: arg get package path
	packagePath := "examples/ManagerInterface"

	// TODO: get package ast
	pkg := astRepository.ScanPackage(packagePath)
	_ = pkg

	// TODO: generate mock structure

	//result := genStructFromInterface(myInterface)
	//formattedBytes, err := format.Source([]byte(result))
	//if err != nil {
	//	fmt.Println("formatting failed:", err.Error())
	//	return
	//}
	//result = string(formattedBytes)
	//
	//// TODO: save mock
	//fmt.Println(result)
}

func genStructFromInterface(i *domain.Interface) (result string) {
	result += fmt.Sprintf("type %s struct {\n", i.Name)

	fieldList := getFieldList(i.MethodList)
	for _, field := range fieldList {
		result += fmt.Sprintf("	%s %s\n", field.Name, field.Type)
	}
	result += fmt.Sprintf("}\n")
	return
}

func getFieldList(methodList []*domain.Method) (fieldList []*domain.Variable) {
	fieldList = make([]*domain.Variable, 0)
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

func isFieldListContainsVariable(fieldList []*domain.Variable, v *domain.Variable) (isContains bool) {
	for _, field := range fieldList {
		if field.Name == v.Name {
			isContains = true
			break
		}
	}
	return
}
