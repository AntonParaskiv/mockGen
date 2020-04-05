package old

import (
	"github.com/AntonParaskiv/mockGen/main"
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
var reflectDeepEqual = main.createSelectorExpr(
	main.createName("reflect"),
	main.createName("DeepEqual"),
)

// t.Run
var tRun = main.createSelectorExpr(
	main.createName("t"),
	main.createName("Run"),
)

// t.Errorf
var tErrorf = main.createSelectorExpr(
	main.createName("t"),
	main.createName("Errorf"),
)

func createTestName(name string) (testName *ast.KeyValueExpr) {
	key := "name"
	value := main.createBasicLit(name, token.STRING)
	testName = main.createKeyValueExpr(key, value)
	return
}

func generateTestValue(field *ast.Field) (basicLit *ast.BasicLit) {
	fieldTypeIdent, ok := field.Type.(*ast.Ident)
	if !ok {
		basicLit = main.createBasicLit(`"// TODO: generate value"`, token.STRING)
		return
	}

	fieldTypeName := main.getNodeName(fieldTypeIdent)

	if fieldTypeName == "string" {
		value := "my" + main.toPublic(main.getNodeName(field))
		basicLit = main.createBasicLit(`"`+value+`"`, token.STRING)
		return
	}

	if len(fieldTypeName) >= 3 && fieldTypeName[0:3] == "int" {
		basicLit = main.createBasicLit("100", token.INT)
		return
	}

	if len(fieldTypeName) >= 5 && fieldTypeName[0:5] == "float" {
		basicLit = main.createBasicLit("100", token.INT)
		return
	}

	// TODO: add other types
	basicLit = main.createBasicLit(`"// TODO: generate value"`, token.STRING)
	return
}

func createTestTable(fieldList *ast.FieldList, rows []ast.Expr) (testTable *ast.CompositeLit) {
	testTable = main.createCompositeLit(
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
	ttSelector = main.createSelectorExpr(
		main.createName("tt"),
		main.createName(name),
	)
	return
}

// t.Errorf("Result = %v, want %v", gotR, tt.want)
func createTestCompareResultErrorf(resultName, wantName string, result, want ast.Expr) (resultError *ast.ExprStmt) {
	resultError = &ast.ExprStmt{
		X: createCallExpr(
			tErrorf,
			main.createBasicLit(
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
	expr = main.createUnaryExpr(
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
	runRangeStmt = main.createRangeStmt(
		"_",
		"tt",
		token.DEFINE,
		"tests",
		stmts...,
	)
	return
}
