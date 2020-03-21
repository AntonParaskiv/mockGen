package main

import (
	"go/ast"
	"go/token"
)

func createTestMethodName(structName, methodName string) (methodTestName string) {
	methodTestName = "Test" + structName + "_" + methodName
	return
}

func createMethodTest(paramList, resultList []*ast.Field, structName, methodName, testMethodName, receiverName string, pointerStruct *ast.StarExpr, namedPointerStruct *ast.Field) (methodTestDecl *ast.FuncDecl) {
	name := createName(testMethodName)
	args := createFieldList(namedPointerTestingT)
	wantReceiver := "want" + toPublic(structName)

	methodTestDecl = createFuncDecl(nil, name, args, nil,
		createDeclStruct("fields", resultList...),
		createDeclStruct("args", paramList...),
		createMethodStmtTestsDeclare(structName, wantReceiver, pointerStruct, resultList),
		createMethodStmtTestsRun(methodName, receiverName, wantReceiver, "got", structName, paramList, resultList),
		returnStmt,
	)

	return
}

func createMethodStmtTestsDeclare(structName, wantReceiver string, pointerStruct *ast.StarExpr, resultList []*ast.Field) (stmtTestsDeclare *ast.AssignStmt) {
	testTable := createMethodTestTable(structName, wantReceiver, pointerStruct, resultList)

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

func createMethodTestTable(structName, wantReceiver string, pointerStruct *ast.StarExpr, methodResultList []*ast.Field) (testTable *ast.CompositeLit) {
	testTableFields := make([]*ast.Field, 0)
	testTableFields = append(testTableFields,
		// name string
		createField("name", createName("string")),
		// fields string
		createField("fields", createName("fields")),
		// args string
		createField("args", createName("args")),
	)

	for _, result := range methodResultList {
		wantResultName := "want" + toPublic(getNodeName(result))
		wantResult := createField(wantResultName, result.Type)
		testTableFields = append(testTableFields, wantResult)
	}

	testTableFields = append(testTableFields,
		// wantS *Struct
		createField(wantReceiver, pointerStruct),
	)

	testTableRows := createExprList(
	// Success
	//createTestRowSetting(wantReceiver, structName, field),
	)
	testTable = createTestTable(
		createFieldList(testTableFields...),
		testTableRows,
	)
	return
}

func createMethodStmtTestsRun(methodName, receiverName, wantReceiver, gotReceiver, structName string, methodParamList, methodResultList []*ast.Field) (runRangeStmt *ast.RangeStmt) {
	runStmts := make([]ast.Stmt, 0)

	settingFieldsExprs := make([]ast.Expr, 0, len(methodResultList))
	for _, result := range methodResultList {
		setting := createKeyValueExpr(
			getNodeName(result),
			createSelectorExpr(
				createSelectorExpr(
					createName("tt"),
					createName("fields"),
				),
				createName(getNodeName(result)),
			),
		)
		settingFieldsExprs = append(settingFieldsExprs, setting)

	}

	runStmts = append(runStmts,
		// s := &Struct{ field: field }
		createAssignStmt(
			// s
			createExprList(createName(receiverName)),
			// :=
			token.DEFINE,
			// &Struct{ field: field }
			createExprList(
				initStructLiteral(structName, settingFieldsExprs...),
			),
		),
	)

	// prepare method execute
	gotResults := make([]ast.Expr, 0)
	for _, result := range methodResultList {
		gotResults = append(gotResults, createName("got"+toPublic(getNodeName(result))))
	}
	args := make([]ast.Expr, 0)
	for _, param := range methodParamList {
		args = append(args, createSelectorExpr(
			createSelectorExpr(
				createName("tt"),
				createName("args"),
			),
			createName(getNodeName(param)),
		),
		)
	}

	runStmts = append(runStmts,
		// gotField1, gotField2 := s.Method(tt.args.arg1, tt.args.arg2)
		createAssignStmt(
			// gotField1, gotField2
			createExprList(gotResults...),
			// :=
			token.DEFINE,
			// s.Method(tt.args.arg1, tt.args.arg2)
			createExprList(createCallExpr(
				createSelectorExpr(createName(receiverName), createName(methodName)),
				args...,
			)),
		),
	)

	// compare wantResults
	for _, param := range methodResultList {
		gotResult := "got" + toPublic(getNodeName(param))
		ttWantResult := createTTSelector(wantReceiver)

		// !reflect.DeepEqual(gotField, tt.wantField)
		compareResultCondition := createNotDeepEqualExpr(
			createName(gotResult),
			ttWantResult,
		)

		// t.Errorf("gotField = %v, want %v", gotField, tt.wantField)
		compareResultErrorf := createTestCompareResultErrorf(
			getNodeName(param),
			wantReceiver,
			createName(gotResult),
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
		createName(receiverName),
		ttWantReceiver,
	)
	compareStructErrorf := createTestCompareResultErrorf(
		structName,
		"want",
		createName(receiverName),
		ttWantReceiver,
	)
	runStmts = append(runStmts,
		// if !reflect.DeepEqual(s, tt.wantStruct) { t.Errorf("Struct = %v, want %v", s, tt.wantStruct) }
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
