package main

import (
	"go/ast"
	"go/token"
)

// Struct -> s
func getReceiverName(name string) (receiverName string) {
	receiverName = toPrivate(getFirstLetter(name))
	return
}

// Struct -> *Struct
func createPointerStruct(structName string) (pointerStruct *ast.StarExpr) {
	pointerStruct = &ast.StarExpr{
		X: &ast.Ident{
			Name: structName,
		},
	}
	return
}

func createFieldNamedPointerStruct(structName, receiverName string) (field *ast.Field) {
	pointerStruct := createPointerStruct(structName)
	field = createFieldFromExpr(pointerStruct)
	field.Names = createNames(receiverName)
	return
}

func createStructSpec(name string, fields ...*ast.Field) (structSpec *ast.TypeSpec) {
	structSpec = &ast.TypeSpec{
		Name: createName(name),
		Type: &ast.StructType{
			Fields: createFieldList(fields...),
		},
	}
	return
}

func createStructFromInterfaceSpec(interfaceSpec *ast.TypeSpec) (structSpec *ast.TypeSpec, err error) {
	switch specType := interfaceSpec.Type.(type) {
	case *ast.InterfaceType:
		structName := getNodeName(interfaceSpec)
		fieldList := createFieldListFromInterfaceMethods(specType)
		structSpec = createStructSpec(structName, fieldList...)
		return
	}
	return
}

func initStructLiteral(structName string, elts ...ast.Expr) (structLiteral *ast.UnaryExpr) {
	structLiteral = createUnaryExpr(token.AND, createCompositeLit(createName(structName), elts...))
	return
}

func createDeclStruct(structName string, field ...*ast.Field) (declStmt *ast.DeclStmt) {
	structSpec := createStructSpec(structName, field...)
	declStmt = createDeclStmt(token.TYPE, structSpec)
	return
}
