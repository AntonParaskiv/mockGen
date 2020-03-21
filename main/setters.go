package main

import (
	"go/ast"
	"go/token"
)

func createSetterName(fieldName string) (setterName string) {
	setterName = "Set" + toPublic(fieldName)
	return
}

func createSettersAndTests(structSpec *ast.TypeSpec) (setterDecls, setterTestsDecls []*ast.FuncDecl) {
	structName := getNodeName(structSpec)
	receiverName := getReceiverName(structName)

	fieldList := structSpec.Type.(*ast.StructType).Fields.List
	pointerStruct := createPointerStruct(structName)
	namedPointerStruct := createFieldNamedPointerStruct(structName, receiverName)

	setterDecls = make([]*ast.FuncDecl, 0)
	setterTestsDecls = make([]*ast.FuncDecl, 0)

	for _, field := range fieldList {
		// create setter
		functionName := createSetterName(getNodeName(field))
		setter := createSetter(field, namedPointerStruct, pointerStruct, functionName, receiverName)
		setterDecls = append(setterDecls, setter)

		// create test
		testFunctionName := createTestSetterName(structName, functionName)
		setterTestDecl := createSetterTest(field, pointerStruct, structName, receiverName, functionName, testFunctionName)
		setterTestsDecls = append(setterTestsDecls, setterTestDecl)
	}

	return
}

func createSetter(field, namedPointerStruct *ast.Field, pointerStruct *ast.StarExpr, functionName, receiverName string) (setterDecl *ast.FuncDecl) {
	name := createName(functionName)
	args := createFieldList(field)
	results := createFieldList(createFieldFromExpr(pointerStruct))
	recvs := createFieldList(namedPointerStruct)

	// s.field = field
	lineSFieldAssignField := createAssignStmt(
		// s.field
		createExprList(createSelectorExpr(createName(receiverName), createName(getNodeName(field)))),
		// =
		token.ASSIGN,
		// field
		createExprList(createName(getNodeName(field))),
	)

	// return s
	lineReturnS := createReturnStmt(createName(receiverName))

	setterDecl = createFuncDecl(recvs, name, args, results,
		lineSFieldAssignField,
		lineReturnS,
	)
	return
}
