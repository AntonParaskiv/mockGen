package old

import (
	"github.com/AntonParaskiv/mockGen/main"
	"go/ast"
	"go/token"
)

func createTestNewName(newName string) (newTestName string) {
	newTestName = "Test" + newName
	return
}

func createConstructorTest(structName, functionName, testFunctionName, wantReceiver, gotReceiver string, pointerStruct *ast.StarExpr) (constructorTestDecl *ast.FuncDecl) {
	name := main.createName(testFunctionName)
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
		createField("name", main.createName("string")),
		// wantS *Mock
		createField(wantReceiver, pointerStruct),
	)
	testTableRows := main.createExprList(
		// Mock init
		createTestRowInitStruct(wantReceiver, structName),
	)
	testTable = createTestTable(testTableFieldList, testTableRows)
	return
}

func createTestRowInitStruct(wantReceiver, structName string) (testRow *ast.CompositeLit) {
	testRow = main.createCompositeLit(nil,
		createTestName(`"Mock init"`),
		main.createKeyValueExpr(
			wantReceiver,
			main.initStructLiteral(structName),
		),
	)
	return
}

func createConstructorStmtTestsDeclare(structName, wantReceiver string, pointerStruct *ast.StarExpr) (stmtTestsDeclare *ast.AssignStmt) {
	testTable := createConstructorTestTable(structName, wantReceiver, pointerStruct)

	// test := []struct{...}{...}
	stmtTestsDeclare = createAssignStmt(
		// tests
		main.createExprList(main.createName("tests")),
		// :=
		token.DEFINE,
		// []struct{...}{...}
		main.createExprList(testTable),
	)

	return
}

func createConstructorStmtTestsRun(functionName, wantReceiver, gotReceiver string) (runRangeStmt *ast.RangeStmt) {
	ttWantReceiver := createTTSelector(wantReceiver)
	compareResultInit := createAssignStmt(
		main.createExprList(main.createName(gotReceiver)),
		token.DEFINE,
		main.createExprList(createCallExpr(main.createName(functionName))),
	)
	compareResultCondition := createNotDeepEqualExpr(
		main.createName(gotReceiver),
		ttWantReceiver,
	)
	compareResultErrorf := createTestCompareResultErrorf(
		functionName+"()",
		wantReceiver,
		main.createName(gotReceiver),
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
