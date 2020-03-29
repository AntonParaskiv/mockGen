package main

import (
	"fmt"
	"go/ast"
	"path/filepath"
)

func main() {
	interfacePackagePath := "examples/ManagerInterface"

	// TODO: Parse Imports
	interfacePackage := CreateInterfacePackage(interfacePackagePath)
	mockPackage := CreateMockPackage(interfacePackage)
	err := SaveGoPackage(mockPackage)
	if err != nil {
		panic(err)
	}

	return
}

func CreateInterfacePackage(packagePath string) (interfacePackage *GoCodePackage) {
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
			Name:       filepath.Base(fullFileName), // TODO: check
			ImportList: nil,                         // TODO: fill
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
					argName := getNodeName(astArg)
					astIdent, ok := astArg.Type.(*ast.Ident)
					if !ok {
						err = fmt.Errorf("ast type of %s is not ast ident", argName)
					}

					arg := &Field{
						Name: argName,
						Type: astIdent.Name,
					}
					method.ArgList = append(method.ArgList, arg)
				}
				for _, astResult := range astFuncType.Results.List {
					resultName := getNodeName(astResult)
					astIdent, ok := astResult.Type.(*ast.Ident)
					if !ok {
						err = fmt.Errorf("ast type of %s is not ast ident", resultName)
					}

					result := &Field{
						Name: getNodeName(astResult),
						Type: astIdent.Name,
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
