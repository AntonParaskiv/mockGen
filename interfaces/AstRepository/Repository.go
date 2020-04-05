package AstRepository

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"github.com/AntonParaskiv/mockGen/infrastructure/CodeStorage"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Repository struct {
	CodeStorage CodeStorage.Storage
}

func (r *Repository) CreateInterfacePackage(astPackage *ast.Package, packagePath string) (interfacePackage *domain.GoCodePackage, err error) {
	interfacePackage = &domain.GoCodePackage{
		Path:        packagePath,
		PackageName: astPackage.Name,
	}

	for fullFileName, astFile := range astPackage.Files {
		if !isFileNameMatchGoCode(fullFileName) {
			continue
		}

		interfaceFile := &domain.GoCodeFile{
			Name:       filepath.Base(fullFileName),
			ImportList: getImportListFromAstFile(astFile),
		}

		astInterfaceSpecs := getInterfaces(astFile)
		for _, astInterfaceSpec := range astInterfaceSpecs {
			var iFace *domain.Interface
			iFace, err = createInterfaceFromAstInterfaceSpec(astInterfaceSpec)
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

func getImportListFromAstFile(astFile *ast.File) (importList []*domain.Import) {
	for _, spec := range astFile.Imports {
		Import := &domain.Import{
			Name: getNodeName(spec),
			Path: strings.Trim(spec.Path.Value, `"`),
		}
		if Import.Name != "" {
			Import.Key = Import.Name
		} else {
			Import.Key = filepath.Base(Import.Path)
		}
		importList = append(importList, Import)
	}
	return
}

func createInterfaceFromAstInterfaceSpec(astInterfaceSpec *ast.TypeSpec) (iFace *domain.Interface, err error) {
	iFace = &domain.Interface{
		Name: astInterfaceSpec.Name.Name,
	}

	switch astInterfaceType := astInterfaceSpec.Type.(type) {
	case *ast.InterfaceType:
		for _, astMethod := range astInterfaceType.Methods.List {
			method := &domain.Method{
				Name:       getNodeName(astMethod),
				ArgList:    nil,
				ResultList: nil,
			}

			switch astFuncType := astMethod.Type.(type) {
			case *ast.FuncType:
				for _, astArg := range astFuncType.Params.List {
					var arg *domain.Field
					arg, err = createFieldFromAstField(astArg)
					if err != nil {
						err = fmt.Errorf("create field from ast field failed: %w", err)
						return
					}
					method.ArgList = append(method.ArgList, arg)
				}
				for _, astResult := range astFuncType.Results.List {
					var result *domain.Field
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

func createFieldFromAstField(astField *ast.Field) (field *domain.Field, err error) {
	field = &domain.Field{
		Name: getNodeName(astField),
	}
	fieldType, err := getFieldTypeFromAstFieldType(astField.Type)
	if err != nil {
		err = fmt.Errorf("get field %s type from ast field type failed: %w", field.Name, err)
		return
	}
	field.Type = fieldType

	return
}

func getFieldTypeFromAstFieldType(astFieldType ast.Node) (fieldType string, err error) {
	switch astType := astFieldType.(type) {
	case *ast.Ident:
		fieldType = getNodeName(astType)
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
	case *ast.MapType:
		var keyType, valueType string
		keyType, err = getFieldTypeFromAstFieldType(astType.Key)
		if err != nil {
			err = fmt.Errorf("get map key type failed: %w", err)
			return
		}
		valueType, err = getFieldTypeFromAstFieldType(astType.Value)
		if err != nil {
			err = fmt.Errorf("get map value type failed: %w", err)
			return
		}
		fieldType = fmt.Sprintf("map[%s]%s", keyType, valueType)
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
	// custom types // TODO: handle imports
	case *ast.SelectorExpr:
		fieldType = fmt.Sprintf("%s.%s", getNodeName(astType.X), getNodeName(astType.Sel))
		// TODO: get baseType

	case *ast.StarExpr:
		var baseFieldType string
		baseFieldType, err = getFieldTypeFromAstFieldType(astType.X)
		if err != nil {
			err = fmt.Errorf("get base field type failed: %w", err)
			return
		}
		fieldType = fmt.Sprintf("*%s", baseFieldType)
	default:
		err = fmt.Errorf("unsupported type")
		return
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
	case *ast.ImportSpec:
		if nodeItem.Name != nil {
			name = nodeItem.Name.Name
		}
	default:
		panic(fmt.Sprintf("no getting name case for type %T", node))
	}
	return
}

func (r *Repository) getTypeDeclarationFromPackage(packagePath, typeName string) (typeSpec *ast.TypeSpec) {
	fullPackagePath := filepath.Join(os.Getenv("GOPATH"), "src", packagePath)
	astPackage, err := r.CodeStorage.GetAstPackage(fullPackagePath)
	if err != nil {
		err = fmt.Errorf("get ast package %s failed: %w", packagePath, err)
		return
	}

	for _, astFile := range astPackage.Files {
		for _, decl := range astFile.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				switch decl.Tok {
				case token.TYPE:
					spec := decl.Specs[0].(*ast.TypeSpec) // TODO: check array
					if getNodeName(spec) == typeName {
						typeSpec = spec
						return
					}
				}
			}
		}
	}

	return
}

func (r *Repository) GetTypeFieldFromPackage(packagePath, typeName string) (field *domain.Field, err error) {
	typeSpec := r.getTypeDeclarationFromPackage(packagePath, typeName)

	field = &domain.Field{
		Name: getNodeName(typeSpec),
	}

	fieldType, err := getFieldTypeFromAstFieldType(typeSpec.Type)
	if err != nil {
		err = fmt.Errorf("get field %s type from ast field type failed: %w", field.Name, err)
		return
	}
	field.Type = fieldType

	return
}
