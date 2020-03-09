package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// TODO: arg get package path
	//packagePath := "domain"
	packagePath := "examples/ManagerInterface"

	rPkgs, err := getRPackages(packagePath)
	if err != nil {
		panic(err)
	}

	wPkgs, err := makeWPackages(rPkgs)
	if err != nil {
		panic(err)
	}

	_ = wPkgs

	err = savePackages(wPkgs)
	if err != nil {
		panic(err)
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
		wPkgName := createMockPackageName(getName(rPkg))
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

		mockPackageName := createMockPackageName(getName(rFile))

		wFile := newAstFile()
		wFile.Name = &ast.Ident{
			Name: mockPackageName,
		}

		var structSpec *ast.TypeSpec
		for _, interfaceSpec := range interfaceSpecs {

			// gen struct
			structSpec, err = createStruct(interfaceSpec)
			if err != nil {
				err = fmt.Errorf("create struct from interface %s failed: %w", getName(interfaceSpec), err)
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
			constructorDecl := createConstructor(structSpec, unMockPackageName(mockPackageName))
			wFile.Decls = append(wFile.Decls, constructorDecl)

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
		wFiles[wFilePath] = wFile
	}

	return
}

func createMethods(structSpec *ast.TypeSpec, interfaceSpec *ast.TypeSpec) (structMethodDecls []*ast.FuncDecl) {
	structName := getName(structSpec)
	receiverName := getReceiverName(structName)
	namedPointerToStruct := createNamedPointerToStruct(structName, receiverName)

	structMethodDecls = make([]*ast.FuncDecl, 0)
	for _, interfaceMethod := range interfaceSpec.Type.(*ast.InterfaceType).Methods.List {
		methodName := getName(interfaceMethod)
		methodParams := interfaceMethod.Type.(*ast.FuncType).Params
		methodResults := interfaceMethod.Type.(*ast.FuncType).Results

		bodyList := make([]ast.Stmt, 0)

		for _, param := range interfaceMethod.Type.(*ast.FuncType).Params.List {
			paramName := getName(param)
			setting := &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.Ident{
							Name: receiverName,
						},
						Sel: &ast.Ident{
							Name: paramName,
						},
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.Ident{
						Name: paramName,
					},
				},
			}
			bodyList = append(bodyList, setting)
		}

		for _, result := range interfaceMethod.Type.(*ast.FuncType).Results.List {
			resultName := getName(result)
			returning := &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: resultName,
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.Ident{
							Name: receiverName,
						},
						Sel: &ast.Ident{
							Name: resultName,
						},
					},
				},
			}
			bodyList = append(bodyList, returning)
		}

		returning := &ast.ReturnStmt{}
		bodyList = append(bodyList, returning)

		structMethod := &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					namedPointerToStruct,
				},
			},
			Name: &ast.Ident{
				Name: methodName,
			},
			Type: &ast.FuncType{
				Func:    0,
				Params:  methodParams,
				Results: methodResults,
			},
			Body: &ast.BlockStmt{
				List: bodyList,
			},
		}
		structMethodDecls = append(structMethodDecls, structMethod)
	}
	return
}

func createSetters(structSpec *ast.TypeSpec) (setterDecls []*ast.FuncDecl) {
	structName := getName(structSpec)
	receiverName := getReceiverName(structName)
	fieldList := structSpec.Type.(*ast.StructType).Fields.List
	pointerToStruct := createPointerToStruct(structName)
	namedPointerToStruct := createNamedPointerToStruct(structName, receiverName)

	setterDecls = make([]*ast.FuncDecl, 0)
	for _, field := range fieldList {
		setterName := createSetterName(getName(field))

		setter := &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					namedPointerToStruct,
				},
			},
			Name: &ast.Ident{
				Name: setterName,
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						field,
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						pointerToStruct,
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: receiverName,
								},
								Sel: &ast.Ident{
									Name: getName(field),
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.Ident{
								Name: getName(field),
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: receiverName,
							},
						},
					},
				},
			},
		}
		setterDecls = append(setterDecls, setter)
	}

	return
}

func createNewName(structName, packageName string) (newName string) {
	newName = "New"
	if strings.ToUpper(structName) != strings.ToUpper(packageName) {
		newName += toPublic(structName)
	}
	return
}

func createSetterName(fieldName string) (setterName string) {
	setterName = "Set" + toPublic(fieldName)
	return
}

func getName(node ast.Node) (name string) {
	switch nodeItem := node.(type) {
	case *ast.Package:
		name = nodeItem.Name
	case *ast.File:
		name = nodeItem.Name.Name
	case *ast.TypeSpec:
		name = nodeItem.Name.Name
	case *ast.Field:
		name = nodeItem.Names[0].Name
	default:
		panic(fmt.Sprintf("no getting name case for type %T", node))
	}
	return
}

func toPublic(name string) (publicName string) {
	firstLetterUpper := strings.ToUpper(name[0:1])
	publicName = firstLetterUpper + name[1:]
	return
}

func toPrivate(name string) (privateName string) {
	firstLetterLower := strings.ToLower(name[0:1])
	privateName = firstLetterLower + name[1:]
	return
}

func createPointerToStruct(structName string) (pointerToStruct *ast.Field) {
	pointerToStruct = &ast.Field{
		Type: &ast.StarExpr{
			X: &ast.Ident{
				Name: structName,
			},
		},
	}
	return
}

func createNamedPointerToStruct(structName, receiverName string) (namedPointerToStruct *ast.Field) {
	pointerToStruct := createPointerToStruct(structName)

	namedPointerToStruct = &ast.Field{
		Names: []*ast.Ident{
			{
				Name: receiverName,
			},
		},
		Type: pointerToStruct.Type,
	}
	return
}

func createConstructor(structSpec *ast.TypeSpec, packageName string) (constructorDecl *ast.FuncDecl) {
	structName := getName(structSpec)
	receiverName := getReceiverName(structName)
	functionName := createNewName(structName, packageName)

	//pointerToStruct := createPointerToStruct(structName)
	namedPointerToStruct := createNamedPointerToStruct(structName, receiverName)

	constructorDecl = &ast.FuncDecl{
		Name: &ast.Ident{
			Name: functionName,
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					namedPointerToStruct,
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: receiverName,
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "new",
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: structName,
								},
							},
							Ellipsis: token.NoPos,
						},
					},
				},
				&ast.ReturnStmt{},
			},
		},
	}
	return
}

func createStruct(interfaceSpec *ast.TypeSpec) (structSpec *ast.TypeSpec, err error) {
	switch interfaceSpecType := interfaceSpec.Type.(type) {

	case *ast.InterfaceType:
		// set Name
		structSpec = &ast.TypeSpec{
			Name: interfaceSpec.Name,
		}

		fieldList := createFieldList(interfaceSpecType)
		structSpec.Type = &ast.StructType{
			Fields: &ast.FieldList{
				List: fieldList,
			},
		}

		return

	}
	return
}

func getReceiverName(name string) (receiverName string) {
	receiverName = strings.ToLower(name[0:1])
	return
}

func createFieldList(specType *ast.InterfaceType) (fieldList []*ast.Field) {
	fieldList = make([]*ast.Field, 0)
	for _, method := range specType.Methods.List {

		switch methodType := method.Type.(type) {
		case *ast.FuncType:
			for _, param := range methodType.Params.List {
				if isFieldExist(fieldList, param) {
					continue
				}
				fieldList = append(fieldList, param)
			}

			for _, result := range methodType.Results.List {
				if isFieldExist(fieldList, result) {
					continue
				}
				fieldList = append(fieldList, result)
			}
		}
	}
	return
}

func isFieldExist(fieldList []*ast.Field, wantField *ast.Field) (isExist bool) {
	name := getName(wantField)
	for _, field := range fieldList {
		if getName(field) == name {
			isExist = true
			return
		}
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

func newAstFile() (file *ast.File) {
	file = &ast.File{
		Decls: []ast.Decl{},
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

func createMockFilePath(filePath string) (mockFilePath string) {
	fileName := filepath.Base(filePath)
	dirPath := filepath.Dir(filePath)

	dirPath = cutPostfix(dirPath, "Interface")
	mockDirPath := dirPath + "Mock"

	mockFilePath = filepath.Join(mockDirPath, fileName)
	return
}

func newPackage() (pkg *ast.Package, err error) {
	fSet := token.NewFileSet()
	pkg, err = ast.NewPackage(fSet, nil, nil, nil)
	if err != nil {
		err = fmt.Errorf("create ast package failed: %w", err)
		return
	}
	return
}

func getRPackages(packagePath string) (rPkgs map[string]*ast.Package, err error) {
	fSet := token.NewFileSet()
	rPkgs, err = parser.ParseDir(fSet, packagePath, nil, 0)
	if err != nil {
		err = fmt.Errorf("parse ast dir failed: %w", err)
		return
	}
	return
}

//func main() {
//
//	// TODO: init components
//	packageStorage := PackageStorage.New()
//	astRepository := AstRepository.New().SetPackageStorage(packageStorage)
//
//	// TODO: arg get package path
//	packagePath := "examples/ManagerInterface"
//
//	// TODO: get package ast
//	pkg := astRepository.ScanPackage(packagePath)
//
//	// TODO: generate mock structure
//	result := genStructFromInterface(myInterface)
//	formattedBytes, err := format.Source([]byte(result))
//	if err != nil {
//		fmt.Println("formatting failed:", err.Error())
//		return
//	}
//	result = string(formattedBytes)
//
//	// TODO: save mock
//	fmt.Println(result)
//	return
//}

//func genStructFromInterface(i *domain.Interface) (result string) {
//	result += fmt.Sprintf("type %s struct {\n", i.Name)
//
//	fieldList := getFieldList(i.MethodList)
//	for _, field := range fieldList {
//		result += fmt.Sprintf("	%s %s\n", field.Name, field.Type)
//	}
//	result += fmt.Sprintf("}\n")
//	return
//}

//func getFieldList(methodList []*domain.Method) (fieldList []*domain.Variable) {
//	fieldList = make([]*domain.Variable, 0)
//	for _, method := range methodList {
//		for _, variable := range method.ArgList {
//			if isFieldListContainsVariable(fieldList, variable) {
//				continue
//			}
//			fieldList = append(fieldList, variable)
//		}
//	}
//	for _, method := range methodList {
//		for _, variable := range method.ValueList {
//			if isFieldListContainsVariable(fieldList, variable) {
//				continue
//			}
//			fieldList = append(fieldList, variable)
//		}
//	}
//	return
//}
//
//func isFieldListContainsVariable(fieldList []*domain.Variable, v *domain.Variable) (isContains bool) {
//	for _, field := range fieldList {
//		if field.Name == v.Name {
//			isContains = true
//			break
//		}
//	}
//	return
//}
