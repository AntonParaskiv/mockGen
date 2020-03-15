package main

import (
	"go/ast"
	"go/token"
)

// testing.T
var testingT = &ast.SelectorExpr{
	X: &ast.Ident{
		Name: "testing",
	},
	Sel: &ast.Ident{
		Name: "T",
	},
}

// *testing.T
var pointerTestingT = &ast.StarExpr{
	X: testingT,
}

// t *testing.T
var namedPointerTestingT = &ast.Field{
	Names: []*ast.Ident{
		{
			Name: "t",
		},
	},
	Type: pointerTestingT,
}

// FuncType with args (t *testing.T)
var testFunctionType = &ast.FuncType{
	Params: &ast.FieldList{
		List: []*ast.Field{
			namedPointerTestingT,
		},
	},
}

// reflect.DeepEqual
var reflectDeepEqual = createSelectorExpr(
	createName("reflect"),
	createName("DeepEqual"),
)

// t.Run
var tRun = createSelectorExpr(
	createName("t"),
	createName("Run"),
)

// t.Errorf
var tErrorf = createSelectorExpr(
	createName("t"),
	createName("Errorf"),
)

func createTestName(name string) (testName *ast.KeyValueExpr) {
	key := "name"
	value := createBasicLit(name, token.STRING)
	testName = createKeyValueExpr(key, value)
	return
}

func createTestRowInitStruct(wantReceiver, structName string) (testRow *ast.CompositeLit) {
	testRow = createCompositeLit(nil,
		createTestName(`"Struct init"`),
		createKeyValueExpr(
			wantReceiver,
			initStructLiteral(structName),
		),
	)
	return
}

func createTestTable(fieldList *ast.FieldList, rows []ast.Expr) (testTable *ast.CompositeLit) {
	testTable = createCompositeLit(
		&ast.ArrayType{
			Elt: &ast.StructType{
				Fields: fieldList,
			},
		},
		rows...,
	)
	return
}

func createTTSelector(name string) (ttSelector *ast.SelectorExpr) {
	ttSelector = createSelectorExpr(
		createName("tt"),
		createName(name),
	)
	return
}

// t.Errorf("Result = %v, want %v", gotR, tt.want)
func createTestCompareResultErrorf(resultName, wantName string, result, want ast.Expr) (resultError *ast.ExprStmt) {
	resultError = &ast.ExprStmt{
		X: createCallExpr(
			tErrorf,
			createBasicLit(
				`"`+resultName+" = %v, "+wantName+" %v"+`"`,
				token.STRING,
			),
			result,
			want,
		),
	}
	return
}

func createNotDeepEqualExpr(arg1, arg2 ast.Expr) (expr *ast.UnaryExpr) {
	expr = createUnaryExpr(
		token.NOT,
		createCallExpr(
			reflectDeepEqual,
			arg1,
			arg2,
		),
	)

	return
}

func createSubTestFunction(list ...ast.Stmt) (subTestFunction *ast.FuncLit) {
	subTestFunction = createFuncLit(
		createFunctionType(createFieldList(namedPointerTestingT), nil),
		&ast.BlockStmt{
			List: list,
		},
	)
	return
}

func createTRunExpr(callBackFunc ast.Expr) (tRunExpr *ast.ExprStmt) {
	tRunExpr = &ast.ExprStmt{
		X: createCallExpr(
			tRun,
			createTTSelector("name"),
			callBackFunc,
		),
	}

	return
}

func createRunRangeStmt(stmts ...ast.Stmt) (runRangeStmt *ast.RangeStmt) {
	runRangeStmt = createRangeStmt(
		"_",
		"tt",
		token.DEFINE,
		"tests",
		stmts...,
	)
	return
}
