package old

import (
	"github.com/AntonParaskiv/mockGen/main"
	"go/ast"
	"go/token"
	"strings"
)

func createNewName(structName, packageName string) (newName string) {
	newName = "New"
	if !isStructNameSameAsPackageName(structName, packageName) {
		newName += main.toPublic(structName)
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
	structName := main.getNodeName(structSpec)
	receiverName := main.getReceiverName(structName)
	functionName := createNewName(structName, packageName)
	fieldNamedPointerStruct := main.createFieldNamedPointerStruct(structName, receiverName)
	constructorDecl = createConstructor(structName, receiverName, functionName, fieldNamedPointerStruct)

	testFunctionName := createTestNewName(functionName)
	wantReceiver := "want" + main.toPublic(receiverName)
	gotReceiver := "got" + main.toPublic(receiverName)
	pointerToStruct := main.createPointerStruct(structName)
	constructorTestDecl = createConstructorTest(structName, functionName, testFunctionName, wantReceiver, gotReceiver, pointerToStruct)

	return
}

func createConstructor(structName, receiverName, functionName string, fieldNamedPointerStruct *ast.Field) (constructorDecl *ast.FuncDecl) {
	name := main.createName(functionName)
	args := createFieldList()
	results := createFieldList(fieldNamedPointerStruct)

	// s = new(Mock)
	lineSAssignNewStruct := createAssignStmt(
		// s
		main.createExprList(main.createName(receiverName)),
		// =
		token.ASSIGN,
		// new(Mock)
		main.createExprList(createCallExpr(main.createName("new"), main.createName(structName))),
	)

	constructorDecl = createFuncDecl(nil, name, args, results,
		lineSAssignNewStruct,
		returnStmt,
	)
	return
}
