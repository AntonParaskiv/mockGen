package old

import (
	"github.com/AntonParaskiv/mockGen/main"
	"go/ast"
	"go/token"
)

func createSetterName(fieldName string) (setterName string) {
	setterName = "Set" + main.toPublic(fieldName)
	return
}

func createSettersAndTests(structSpec *ast.TypeSpec) (setterDecls, setterTestsDecls []*ast.FuncDecl) {
	structName := main.getNodeName(structSpec)
	receiverName := main.getReceiverName(structName)

	fieldList := structSpec.Type.(*ast.StructType).Fields.List
	pointerStruct := main.createPointerStruct(structName)
	namedPointerStruct := main.createFieldNamedPointerStruct(structName, receiverName)

	setterDecls = make([]*ast.FuncDecl, 0)
	setterTestsDecls = make([]*ast.FuncDecl, 0)

	for _, field := range fieldList {
		// create setter
		functionName := createSetterName(main.getNodeName(field))
		setter := createSetter(field, namedPointerStruct, pointerStruct, functionName, receiverName)
		setterDecls = append(setterDecls, setter)

		// create test
		testFunctionName := createTestSetterName(structName, functionName)
		setterTestDecl := createSetterTest(field, pointerStruct, structName, receiverName, functionName, testFunctionName)
		setterTestsDecls = append(setterTestsDecls, setterTestDecl)
	}

	return
}

func createSetter(field, namedPointerStruct *ast.Field, pointerStruct *ast.StarExpr, functionName, receiverName string) (setterDecl *ast.FuncDecl) {
	name := main.createName(functionName)
	args := createFieldList(field)
	results := createFieldList(createFieldFromExpr(pointerStruct))
	recvs := createFieldList(namedPointerStruct)

	// s.Field = Field
	lineSFieldAssignField := createAssignStmt(
		// s.Field
		main.createExprList(main.createSelectorExpr(main.createName(receiverName), main.createName(main.getNodeName(field)))),
		// =
		token.ASSIGN,
		// Field
		main.createExprList(main.createName(main.getNodeName(field))),
	)

	// return s
	lineReturnS := createReturnStmt(main.createName(receiverName))

	setterDecl = createFuncDecl(recvs, name, args, results,
		lineSFieldAssignField,
		lineReturnS,
	)
	return
}
