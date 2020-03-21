package main

import (
	"go/ast"
	"go/token"
)

func createTestSetterName(structName, setterFunctionName string) (setterTestName string) {
	setterTestName = "Test" + structName + "_" + setterFunctionName
	return
}

func createSetterTest(field *ast.Field, pointerStruct *ast.StarExpr, structName, receiverName, functionName, testFunctionName string) (setterTestDecl *ast.FuncDecl) {
	name := createName(testFunctionName)
	args := createFieldList(namedPointerTestingT)

	setterTestDecl = createFuncDecl(nil, name, args, nil,
		createDeclArgStruct(field),
		createSetterStmtTestsDeclare(structName, "want", field, pointerStruct),
		createSetterStmtTestsRun(functionName, receiverName, "want", "got", structName, getNodeName(field)),
		returnStmt,
	)

	return
}

func createDeclArgStruct(field *ast.Field) (declStmt *ast.DeclStmt) {
	structSpec := createStructSpec("args", field)
	declStmt = createDeclStmt(token.TYPE, structSpec)
	return
}

func createSetterStmtTestsDeclare(structName, wantReceiver string, field *ast.Field, pointerStruct *ast.StarExpr) (stmtTestsDeclare *ast.AssignStmt) {
	testTable := createSetterTestTable(structName, wantReceiver, field, pointerStruct)

	// test := []struct{...}{...}
	stmtTestsDeclare = createAssignStmt(
		// tests
		createExprList(createName("tests")),
		// :=
		token.DEFINE,
		// []struct{...}{...}
		createExprList(testTable),
	)

	return
}

func createSetterTestTable(structName, wantReceiver string, field *ast.Field, pointerStruct *ast.StarExpr) (testTable *ast.CompositeLit) {
	testTableFieldList := createFieldList(
		// name string
		createField("name", createName("string")),
		// args string
		createField("args", createName("args")),
		// wantS *Struct
		createField(wantReceiver, pointerStruct),
	)
	testTableRows := createExprList(
		// Struct init
		createTestRowSetting(wantReceiver, structName, field),
	)
	testTable = createTestTable(testTableFieldList, testTableRows)
	return
}

func createSetterStmtTestsRun(functionName, receiverName, wantReceiver, gotReceiver, structName, fieldName string) (runRangeStmt *ast.RangeStmt) {
	ttWantReceiver := createTTSelector(wantReceiver)
	compareResultInit := createAssignStmt(
		createExprList(createName(gotReceiver)),
		token.DEFINE,
		createExprList(createCallExpr(
			createSelectorExpr(createName(receiverName), createName(functionName)),
			createSelectorExpr(
				createSelectorExpr(
					createName("tt"),
					createName("args"),
				),
				createName(fieldName),
			),
		)),
	)
	compareResultCondition := createNotDeepEqualExpr(
		createName(gotReceiver),
		ttWantReceiver,
	)
	compareResultErrorf := createTestCompareResultErrorf(
		functionName+"()",
		wantReceiver,
		createName(gotReceiver),
		ttWantReceiver,
	)

	// s := &Struct{}
	lineSDefineStructLiteral := createAssignStmt(
		// s
		createExprList(createName(receiverName)),
		// :=
		token.DEFINE,
		// &Struct{}
		createExprList(initStructLiteral(structName)),
	)

	compareIfBlock := createIfStmt(
		compareResultInit,
		compareResultCondition,
		compareResultErrorf,
	)
	runTestFunc := createTRunExpr(createSubTestFunction(
		lineSDefineStructLiteral,
		compareIfBlock,
	))
	runRangeStmt = createRunRangeStmt(runTestFunc)
	return
}
