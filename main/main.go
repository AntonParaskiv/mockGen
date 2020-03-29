package main

import (
	"fmt"
	"go/ast"
	"path/filepath"
)

func main() {
	//interfacePackage := &GoCodePackage{
	//	Path:        "examples/ManagerInterface",
	//	PackageName: "ManagerInterface",
	//	FileList: []*GoCodeFile{
	//		{
	//			Name: "Manager.go",
	//			InterfaceList: []*Interface{
	//				{
	//					Name: "Manager",
	//					MethodList: []*Method{
	//						{
	//							Name: "Registration",
	//							ArgList: []*Field{
	//								{
	//									Name: "nickName",
	//									Type: "string",
	//								},
	//								{
	//									Name: "password",
	//									Type: "string",
	//								},
	//							},
	//							ResultList: []*Field{
	//								{
	//									Name: "accountId",
	//									Type: "int64",
	//								},
	//								{
	//									Name: "checkCode",
	//									Type: "string",
	//								},
	//							},
	//						},
	//						{
	//							Name: "SignIn",
	//							ArgList: []*Field{
	//								{
	//									Name: "accountId",
	//									Type: "int64",
	//								},
	//								{
	//									Name: "password",
	//									Type: "string",
	//								},
	//							},
	//							ResultList: []*Field{
	//								{
	//									Name: "nickName",
	//									Type: "string",
	//								},
	//							},
	//						},
	//					},
	//				},
	//				//{
	//				//	Name: "Factory",
	//				//	MethodList: []*Method{
	//				//		{
	//				//			Name: "Create",
	//				//			ArgList: []*Field{
	//				//				{
	//				//					Name: "accountId",
	//				//					Type: "int64",
	//				//				},
	//				//			},
	//				//			ResultList: []*Field{
	//				//				{
	//				//					Name: "manager",
	//				//					Type: "*Manager",
	//				//				},
	//				//			},
	//				//		},
	//				//	},
	//				//},
	//			},
	//		},
	//	},
	//}

	interfacePackagePath := "examples/ManagerInterface"

	// TODO: Parse Imports
	interfacePackage := CreateInterfacePackage(interfacePackagePath)

	mockPackage := CreateMockPackage(interfacePackage)
	err := SaveGoPackage(mockPackage)
	if err != nil {
		panic(err)
	}

	return

	//// TODO: arg get package path
	////interfacePackagePath := "domain"
	//interfacePackagePath := "examples/ManagerInterface"
	//
	//rPkgs, err := getAstPackage(interfacePackagePath)
	//if err != nil {
	//	panic(err)
	//}
	//
	//wPkgs, err := makeWPackages(rPkgs)
	//if err != nil {
	//	panic(err)
	//}
	//
	//_ = wPkgs
	//
	//err = savePackages(wPkgs)
	//if err != nil {
	//	panic(err)
	//}

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
					arg := &Field{
						Name: getNodeName(astArg),
						//Type: astArg.Type, // TODO: fill
					}
					method.ArgList = append(method.ArgList, arg)
				}
				for _, astResult := range astFuncType.Results.List {
					result := &Field{
						Name: getNodeName(astResult),
						//Type: astResult.Type, // TODO: fill
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
