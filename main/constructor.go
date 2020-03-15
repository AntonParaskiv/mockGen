package main

import (
	"go/ast"
	"go/token"
	"strings"
)

func createNewName(structName, packageName string) (newName string) {
	newName = "New"
	if !isStructNameSameAsPackageName(structName, packageName) {
		newName += toPublic(structName)
	}
	return
}

func isStructNameSameAsPackageName(structName, packageName string) (isSame bool) {
	if strings.ToUpper(structName) == strings.ToUpper(packageName) {
		isSame = true
	}
	return
}

func createConstructorAndTest(structSpec *ast.TypeSpec, packageName string) (constructorDecl, constructorTestDecl *ast.FuncDecl) {
	structName := getNodeName(structSpec)
	receiverName := getReceiverName(structName)
	functionName := createNewName(structName, packageName)
	fieldNamedPointerStruct := createFieldNamedPointerStruct(structName, receiverName)
	constructorDecl = createConstructor(structName, receiverName, functionName, fieldNamedPointerStruct)

	testFunctionName := createTestNewName(functionName)
	wantReceiver := "want" + toPublic(receiverName)
	gotReceiver := "got" + toPublic(receiverName)
	pointerToStruct := createPointerStruct(structName)
	constructorTestDecl = createConstructorTest(structName, functionName, testFunctionName, wantReceiver, gotReceiver, pointerToStruct)

	return
}

func createConstructor(structName, receiverName, functionName string, fieldNamedPointerStruct *ast.Field) (constructorDecl *ast.FuncDecl) {
	name := createName(functionName)
	args := createFieldList()
	results := createFieldList(fieldNamedPointerStruct)

	// s = new(Struct)
	lineSAssignNewStruct := createAssignStmt(
		// s
		createExprList(createName(receiverName)),
		// =
		token.ASSIGN,
		// new(Struct)
		createExprList(createCallExpr(createName("new"), createName(structName))),
	)

	constructorDecl = createFuncDecl(name, args, results, lineSAssignNewStruct, returnStmt)
	return
}
