package old

import (
	"github.com/AntonParaskiv/mockGen/main"
	"go/ast"
)

func createField(name string, Type ast.Expr) (field *ast.Field) {
	field = &ast.Field{
		Names: main.createNames(name),
		Type:  Type,
	}
	return
}

func createFieldFromExpr(expr ast.Expr) (field *ast.Field) {
	field = &ast.Field{
		Type: expr,
	}
	return
}

func isFieldExistInFieldList(fieldList []*ast.Field, wantField *ast.Field) (isExist bool) {
	name := main.getNodeName(wantField)
	for _, field := range fieldList {
		if main.getNodeName(field) == name {
			isExist = true
			return
		}
	}
	return
}

func createFieldList(fields ...*ast.Field) (fieldList *ast.FieldList) {
	fieldList = &ast.FieldList{
		List: []*ast.Field{},
	}
	for _, field := range fields {
		fieldList.List = append(fieldList.List, field)
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
			methodFieldList := createFieldListFromFunction(methodType)
			mergeUniqueFieldList(methodFieldList, &fieldList)
		}
	}
	return
}
