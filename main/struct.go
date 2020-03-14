package main

import (
	"go/ast"
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

func createStructSpec(name string, fieldList []*ast.Field) (structSpec *ast.TypeSpec) {
	structSpec = &ast.TypeSpec{
		Name: createName(name),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: fieldList,
			},
		},
	}
	return
}

func createStructFromInterfaceSpec(interfaceSpec *ast.TypeSpec) (structSpec *ast.TypeSpec, err error) {
	switch specType := interfaceSpec.Type.(type) {
	case *ast.InterfaceType:
		structName := getNodeName(interfaceSpec)
		fieldList := createFieldListFromInterfaceMethods(specType)
		structSpec = createStructSpec(structName, fieldList)
		return
	}
	return
}
