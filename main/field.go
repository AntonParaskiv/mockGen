package main

import "go/ast"

func createFieldFromExpr(expr ast.Expr) (field *ast.Field) {
	field = &ast.Field{
		Type: expr,
	}
	return
}

func isFieldExistInFieldList(fieldList []*ast.Field, wantField *ast.Field) (isExist bool) {
	name := getNodeName(wantField)
	for _, field := range fieldList {
		if getNodeName(field) == name {
			isExist = true
			return
		}
	}
	return
}

func mergeUniqueFieldList(from []*ast.Field, to *[]*ast.Field) {
	for _, field := range from {
		if isFieldExistInFieldList(*to, field) {
			continue
		}
		*to = append(*to, field)
	}
}

func createFieldListFromFunction(funcType *ast.FuncType) (fieldList []*ast.Field) {
	mergeUniqueFieldList(funcType.Params.List, &fieldList)
	mergeUniqueFieldList(funcType.Results.List, &fieldList)
	return
}

func createFieldListFromInterfaceMethods(interfaceType *ast.InterfaceType) (fieldList []*ast.Field) {
	for _, field := range interfaceType.Methods.List {
		switch methodType := field.Type.(type) {
		case *ast.FuncType:
			fieldList = createFieldListFromFunction(methodType)
		}
	}
	return
}
