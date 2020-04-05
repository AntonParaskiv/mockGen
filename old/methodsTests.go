package old

import (
	"github.com/AntonParaskiv/mockGen/main"
	"go/ast"
	"go/token"
)

func createTestMethodName(structName, methodName string) (methodTestName string) {
	methodTestName = "Test" + structName + "_" + methodName
	return
}

func createMethodTest(paramList, resultList []*ast.Field, structName, methodName, testMethodName, receiverName string, pointerStruct *ast.StarExpr, namedPointerStruct *ast.Field) (methodTestDecl *ast.FuncDecl) {
	name := main.createName(testMethodName)
	args := createFieldList(namedPointerTestingT)
	wantReceiver := "want" + main.toPublic(structName)

	methodTestDecl = createFuncDecl(nil, name, args, nil,
		main.createDeclStruct("fields", resultList...),
		main.createDeclStruct("args", paramList...),
		createMethodStmtTestsDeclare(structName, wantReceiver, pointerStruct, paramList, resultList),
		createMethodStmtTestsRun(methodName, receiverName, wantReceiver, "got", structName, paramList, resultList),
		returnStmt,
	)

	return
}

func createMethodStmtTestsDeclare(structName, wantReceiver string, pointerStruct *ast.StarExpr, paramList, resultList []*ast.Field) (stmtTestsDeclare *ast.AssignStmt) {
	testTable := createMethodTestTable(structName, wantReceiver, pointerStruct, paramList, resultList)

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

func createMethodTestTable(structName, wantReceiver string, pointerStruct *ast.StarExpr, methodParamList, methodResultList []*ast.Field) (testTable *ast.CompositeLit) {
	testTableFields := make([]*ast.Field, 0)
	testTableFields = append(testTableFields,
		// name string
		createField("name", main.createName("string")),
		// fields string
		createField("fields", main.createName("fields")),
		// args string
		createField("args", main.createName("args")),
	)

	for _, result := range methodResultList {
		wantResultName := "want" + main.toPublic(main.getNodeName(result))
		wantResult := createField(wantResultName, result.Type)
		testTableFields = append(testTableFields, wantResult)
	}

	testTableFields = append(testTableFields,
		// wantS *Mock
		createField(wantReceiver, pointerStruct),
	)

	testTableRows := main.createExprList(
		// Success
		createTestRowMethod(wantReceiver, structName, methodParamList, methodResultList),
	)
	testTable = createTestTable(
		createFieldList(testTableFields...),
		testTableRows,
	)
	return
}

func createTestRowMethod(wantReceiver, structName string, methodParamList, methodResultList []*ast.Field) (testRow *ast.CompositeLit) {

	// Field: "myField"
	fieldsKeyValue := make([]ast.Expr, 0)
	for _, result := range methodResultList {
		fieldKeyValue := main.createKeyValueExpr(
			main.getNodeName(result),
			generateTestValue(result),
		)
		fieldsKeyValue = append(fieldsKeyValue, fieldKeyValue)
	}

	// arg: "myArg"
	argsKeyValue := make([]ast.Expr, 0)
	for _, param := range methodParamList {
		fieldKeyValue := main.createKeyValueExpr(
			main.getNodeName(param),
			generateTestValue(param),
		)
		argsKeyValue = append(argsKeyValue, fieldKeyValue)
	}

	// result struct fields = fields + args
	resultFields := make([]ast.Expr, 0)
	resultFields = append(resultFields, argsKeyValue...)
	resultFields = append(resultFields, fieldsKeyValue...)

	testRowBody := make([]ast.Expr, 0)
	testRowBody = append(testRowBody,
		createTestName(`"Success"`),
		main.createKeyValueExpr(
			"fields",
			main.createCompositeLit(
				main.createName("fields"),
				fieldsKeyValue...,
			),
		),
		main.createKeyValueExpr(
			"args",
			main.createCompositeLit(
				main.createName("args"),
				argsKeyValue...,
			),
		),
	)

	for _, result := range methodResultList {
		wantResult := "want" + main.toPublic(main.getNodeName(result))
		testRowBody = append(testRowBody,
			main.createKeyValueExpr(
				wantResult,
				generateTestValue(result),
			),
		)
	}

	testRowBody = append(testRowBody,
		main.createKeyValueExpr(
			wantReceiver,
			main.initStructLiteral(structName, resultFields...),
			//initStructLiteral(structName),
		),
	)

	testRow = main.createCompositeLit(nil, testRowBody...)
	return
}

func createMethodStmtTestsRun(methodName, receiverName, wantReceiver, gotReceiver, structName string, methodParamList, methodResultList []*ast.Field) (runRangeStmt *ast.RangeStmt) {
	runStmts := make([]ast.Stmt, 0)

	settingFieldsExprs := make([]ast.Expr, 0, len(methodResultList))
	for _, result := range methodResultList {
		setting := main.createKeyValueExpr(
			main.getNodeName(result),
			main.createSelectorExpr(
				main.createSelectorExpr(
					main.createName("tt"),
					main.createName("fields"),
				),
				main.createName(main.getNodeName(result)),
			),
		)
		settingFieldsExprs = append(settingFieldsExprs, setting)

	}

	runStmts = append(runStmts,
		// s := &Mock{ Field: Field }
		createAssignStmt(
			// s
			main.createExprList(main.createName(receiverName)),
			// :=
			token.DEFINE,
			// &Mock{ Field: Field }
			main.createExprList(
				main.initStructLiteral(structName, settingFieldsExprs...),
			),
		),
	)

	// prepare method execute
	gotResults := make([]ast.Expr, 0)
	for _, result := range methodResultList {
		gotResults = append(gotResults, main.createName("got"+main.toPublic(main.getNodeName(result))))
	}
	args := make([]ast.Expr, 0)
	for _, param := range methodParamList {
		args = append(args, main.createSelectorExpr(
			main.createSelectorExpr(
				main.createName("tt"),
				main.createName("args"),
			),
			main.createName(main.getNodeName(param)),
		),
		)
	}

	runStmts = append(runStmts,
		// gotField1, gotField2 := s.Method(tt.args.arg1, tt.args.arg2)
		createAssignStmt(
			// gotField1, gotField2
			main.createExprList(gotResults...),
			// :=
			token.DEFINE,
			// s.Method(tt.args.arg1, tt.args.arg2)
			main.createExprList(createCallExpr(
				main.createSelectorExpr(main.createName(receiverName), main.createName(methodName)),
				args...,
			)),
		),
	)

	// compare wantResults
	for _, param := range methodResultList {
		gotResult := "got" + main.toPublic(main.getNodeName(param))
		wantResult := "want" + main.toPublic(main.getNodeName(param))
		ttWantResult := createTTSelector(wantResult)

		// !reflect.DeepEqual(gotField, tt.wantField)
		compareResultCondition := createNotDeepEqualExpr(
			main.createName(gotResult),
			ttWantResult,
		)

		// t.Errorf("gotField = %v, want %v", gotField, tt.wantField)
		compareResultErrorf := createTestCompareResultErrorf(
			gotResult,
			"want",
			main.createName(gotResult),
			ttWantResult,
		)

		runStmts = append(runStmts,
			// if !reflect.DeepEqual(gotField, tt.wantField) { t.Errorf("gotField = %v, want %v", gotField, tt.wantField) }
			createIfStmt(
				nil,
				compareResultCondition,
				compareResultErrorf,
			),
		)

	}

	// compare self block
	ttWantReceiver := createTTSelector(wantReceiver)
	compareStructCondition := createNotDeepEqualExpr(
		main.createName(receiverName),
		ttWantReceiver,
	)
	compareStructErrorf := createTestCompareResultErrorf(
		structName,
		"want",
		main.createName(receiverName),
		ttWantReceiver,
	)
	runStmts = append(runStmts,
		// if !reflect.DeepEqual(s, tt.wantStruct) { t.Errorf("Mock = %v, want %v", s, tt.wantStruct) }
		createIfStmt(
			nil,
			compareStructCondition,
			compareStructErrorf,
		),
	)

	// create Run
	runRangeStmt = createRunRangeStmt(
		createTRunExpr(createSubTestFunction(runStmts...)),
	)
	return
}
