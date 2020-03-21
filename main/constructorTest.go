package main

import (
	"go/ast"
	"go/token"
)

func createTestNewName(newName string) (newTestName string) {
	newTestName = "Test" + newName
	return
}

func createConstructorTest(structName, functionName, testFunctionName, wantReceiver, gotReceiver string, pointerStruct *ast.StarExpr) (constructorTestDecl *ast.FuncDecl) {
	name := createName(testFunctionName)
	args := createFieldList(namedPointerTestingT)
	results := createFieldList()

	constructorTestDecl = createFuncDecl(nil, name, args, results,
		createConstructorStmtTestsDeclare(structName, wantReceiver, pointerStruct),
		createConstructorStmtTestsRun(functionName, wantReceiver, gotReceiver),
	)

	return
}

func createConstructorTestTable(structName, wantReceiver string, pointerStruct *ast.StarExpr) (testTable *ast.CompositeLit) {
	testTableFieldList := createFieldList(
		// name string
		createField("name", createName("string")),
		// wantS *Struct
		createField(wantReceiver, pointerStruct),
	)
	testTableRows := createExprList(
		// Struct init
		createTestRowInitStruct(wantReceiver, structName),
	)
	testTable = createTestTable(testTableFieldList, testTableRows)
	return
}

func createConstructorStmtTestsDeclare(structName, wantReceiver string, pointerStruct *ast.StarExpr) (stmtTestsDeclare *ast.AssignStmt) {
	testTable := createConstructorTestTable(structName, wantReceiver, pointerStruct)

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

func createConstructorStmtTestsRun(functionName, wantReceiver, gotReceiver string) (runRangeStmt *ast.RangeStmt) {
	ttWantReceiver := createTTSelector(wantReceiver)
	compareResultInit := createAssignStmt(
		createExprList(createName(gotReceiver)),
		token.DEFINE,
		createExprList(createCallExpr(createName(functionName))),
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
	compareIfBlock := createIfStmt(
		compareResultInit,
		compareResultCondition,
		compareResultErrorf,
	)
	runTestFunc := createTRunExpr(createSubTestFunction(compareIfBlock))
	runRangeStmt = createRunRangeStmt(runTestFunc)
	return
}
