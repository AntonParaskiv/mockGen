package main

import (
	"fmt"
	"go/ast"
	"path/filepath"
)

func main() {
	interfacePackagePath := "examples/ManagerInterface"

	interfacePackage, err := CreateInterfacePackage(interfacePackagePath)
	if err != nil {
		panic(err)
	}

	mockPackage := CreateMockPackage(interfacePackage)
	err = SaveGoPackage(mockPackage)
	if err != nil {
		panic(err)
	}

	return
}

func CreateInterfacePackage(packagePath string) (interfacePackage *GoCodePackage, err error) {
	astPackage, err := getAstPackage(packagePath)
	if err != nil {
		err = fmt.Errorf("get ast package failed: %w", err)
		return
	}
	if astPackage == nil {
		err = fmt.Errorf("ast package not found")
		return
	}

	interfacePackage = &GoCodePackage{
		Path:        packagePath,
		PackageName: astPackage.Name,
	}

	for fullFileName, astFile := range astPackage.Files {
		if !isFileNameMatchGoCode(fullFileName) {
			continue
		}

		interfaceFile := &GoCodeFile{
			Name:       filepath.Base(fullFileName),
			ImportList: nil, // TODO: fill
		}

		astInterfaceSpecs := getInterfaces(astFile)
		for _, astInterfaceSpec := range astInterfaceSpecs {
			var iFace *Interface
			iFace, err = CreateInterfaceFromAstInterfaceSpec(astInterfaceSpec)
			if err != nil {
				err = fmt.Errorf("create interface from ast interface spec failed: %w", err)
				return
			}
			if iFace == nil {
				continue
			}
			if len(iFace.MethodList) == 0 {
				continue
			}
			interfaceFile.InterfaceList = append(interfaceFile.InterfaceList, iFace)
		}
		if len(interfaceFile.InterfaceList) == 0 {
			continue
		}
		interfacePackage.FileList = append(interfacePackage.FileList, interfaceFile)
	}

	return
}

func CreateInterfaceFromAstInterfaceSpec(astInterfaceSpec *ast.TypeSpec) (iFace *Interface, err error) {
	iFace = &Interface{
		Name: astInterfaceSpec.Name.Name,
	}

	switch astInterfaceType := astInterfaceSpec.Type.(type) {
	case *ast.InterfaceType:
		for _, astMethod := range astInterfaceType.Methods.List {
			method := &Method{
				Name:       getNodeName(astMethod),
				ArgList:    nil,
				ResultList: nil,
			}

			switch astFuncType := astMethod.Type.(type) {
			case *ast.FuncType:
				for _, astArg := range astFuncType.Params.List {
					var arg *Field
					arg, err = createFieldFromAstField(astArg)
					if err != nil {
						err = fmt.Errorf("create field from ast field failed: %w", err)
						return
					}
					method.ArgList = append(method.ArgList, arg)
				}
				for _, astResult := range astFuncType.Results.List {
					var result *Field
					result, err = createFieldFromAstField(astResult)
					if err != nil {
						err = fmt.Errorf("create field from ast field failed: %w", err)
						return
					}
					method.ResultList = append(method.ResultList, result)
				}
			}
			iFace.MethodList = append(iFace.MethodList, method)
		}
	default:
		err = fmt.Errorf("ast spec type is not interface type")
		return
	}

	return
}

func createFieldFromAstField(astField *ast.Field) (field *Field, err error) {
	field = &Field{
		Name: getNodeName(astField),
	}
	fieldType, err := getFieldTypeFromAstFieldType(astField.Type)
	if err != nil {
		err = fmt.Errorf("get field %s type from ast field type failed: %w", field.Name, err)
		return
	}
	field.Type = fieldType

	//switch astType := astField.Type.(type) {
	//case *ast.Ident:
	//	field.Type = astType.Name
	//case *ast.InterfaceType:
	//	if len(astType.Methods.List) > 0 {
	//		err = fmt.Errorf("unsupported type interface{} %s with methods", field.Name)
	//		return
	//	}
	//	field.Type = "interface{}"
	//case *ast.ArrayType:
	//	fmt.Println(astType.Elt)
	//
	//default:
	//	err = fmt.Errorf("unsupported type of %s", field.Name)
	//	return
	//}

	return
}

func getFieldTypeFromAstFieldType(astFieldType ast.Expr) (fieldType string, err error) {
	switch astType := astFieldType.(type) {
	case *ast.Ident:
		fieldType = astType.Name
	case *ast.InterfaceType:
		if len(astType.Methods.List) > 0 {
			err = fmt.Errorf("unsupported type interface{} with methods")
			return
		}
		fieldType = "interface{}"
	case *ast.ArrayType:
		var itemType string
		itemType, err = getFieldTypeFromAstFieldType(astType.Elt)
		if err != nil {
			err = fmt.Errorf("get array item type failed: %w", err)
			return
		}
		fieldType = fmt.Sprintf("[]%s", itemType)
	case *ast.StructType:
		fieldType = fmt.Sprintf("struct {\n")
		for _, item := range astType.Fields.List {
			var itemType string
			itemType, err = getFieldTypeFromAstFieldType(item.Type)
			if err != nil {
				err = fmt.Errorf("get struct field type failed: %w", err)
				return
			}
			fieldType += fmt.Sprintf("	%s %s\n", getNodeName(item), itemType)
		}
		fieldType += fmt.Sprintf("}")
	default:
		err = fmt.Errorf("unsupported type")
		return
	}
	return
}
