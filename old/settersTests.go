package old

import (
	"github.com/AntonParaskiv/mockGen/main"
	"go/ast"
	"go/token"
)

func createTestSetterName(structName, setterFunctionName string) (setterTestName string) {
	setterTestName = "Test" + structName + "_" + setterFunctionName
	return
}

func createSetterTest(field *ast.Field, pointerStruct *ast.StarExpr, structName, receiverName, functionName, testFunctionName string) (setterTestDecl *ast.FuncDecl) {
	name := main.createName(testFunctionName)
	args := createFieldList(namedPointerTestingT)

	setterTestDecl = createFuncDecl(nil, name, args, nil,
		main.createDeclStruct("args", field),
		createSetterStmtTestsDeclare(structName, "want", field, pointerStruct),
		createSetterStmtTestsRun(functionName, receiverName, "want", "got", structName, main.getNodeName(field)),
		returnStmt,
	)

	return
}

func createSetterStmtTestsDeclare(structName, wantReceiver string, field *ast.Field, pointerStruct *ast.StarExpr) (stmtTestsDeclare *ast.AssignStmt) {
	testTable := createSetterTestTable(structName, wantReceiver, field, pointerStruct)

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

func createSetterTestTable(structName, wantReceiver string, field *ast.Field, pointerStruct *ast.StarExpr) (testTable *ast.CompositeLit) {
	testTableFieldList := createFieldList(
		// name string
		createField("name", main.createName("string")),
		// args string
		createField("args", main.createName("args")),
		// wantS *Mock
		createField(wantReceiver, pointerStruct),
	)
	testTableRows := main.createExprList(
		// Mock init
		createTestRowSetting(wantReceiver, structName, field),
	)
	testTable = createTestTable(testTableFieldList, testTableRows)
	return
}

func createTestRowSetting(wantReceiver, structName string, field *ast.Field) (testRow *ast.CompositeLit) {

	// Field: "myField"
	fieldKeyValue := main.createKeyValueExpr(
		main.getNodeName(field),
		generateTestValue(field),
	)

	testRow = main.createCompositeLit(nil,
		createTestName(`"Setting"`),
		main.createKeyValueExpr(
			"args",
			main.createCompositeLit(
				main.createName("args"),
				fieldKeyValue,
			),
		),
		main.createKeyValueExpr(
			wantReceiver,
			main.initStructLiteral(structName, fieldKeyValue),
		),
	)
	return
}

func createSetterStmtTestsRun(functionName, receiverName, wantReceiver, gotReceiver, structName, fieldName string) (runRangeStmt *ast.RangeStmt) {
	ttWantReceiver := createTTSelector(wantReceiver)
	compareResultInit := createAssignStmt(
		main.createExprList(main.createName(gotReceiver)),
		token.DEFINE,
		main.createExprList(createCallExpr(
			main.createSelectorExpr(main.createName(receiverName), main.createName(functionName)),
			main.createSelectorExpr(
				main.createSelectorExpr(
					main.createName("tt"),
					main.createName("args"),
				),
				main.createName(fieldName),
			),
		)),
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

	// s := &Mock{}
	lineSDefineStructLiteral := createAssignStmt(
		// s
		main.createExprList(main.createName(receiverName)),
		// :=
		token.DEFINE,
		// &Mock{}
		main.createExprList(main.initStructLiteral(structName)),
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
