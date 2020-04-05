package old

import (
	"github.com/AntonParaskiv/mockGen/main"
	"go/ast"
	"go/token"
)

func createMethodsAndTests(structSpec *ast.TypeSpec, interfaceSpec *ast.TypeSpec) (methodDecls, methodTestDecls []*ast.FuncDecl) {
	structName := main.getNodeName(structSpec)
	receiverName := main.getReceiverName(structName)
	pointerStruct := main.createPointerStruct(structName)
	namedPointerStruct := main.createFieldNamedPointerStruct(structName, receiverName)

	methodDecls = make([]*ast.FuncDecl, 0)
	methodTestDecls = make([]*ast.FuncDecl, 0)
	for _, interfaceMethod := range interfaceSpec.Type.(*ast.InterfaceType).Methods.List {
		methodName := main.getNodeName(interfaceMethod)
		paramList := interfaceMethod.Type.(*ast.FuncType).Params.List
		resultList := interfaceMethod.Type.(*ast.FuncType).Results.List

		// create method
		method := createMethod(paramList, resultList, methodName, receiverName, namedPointerStruct)
		methodDecls = append(methodDecls, method)

		// create test
		testMethodName := createTestMethodName(structName, methodName)
		methodTestDecl := createMethodTest(paramList, resultList, structName, methodName, testMethodName, receiverName, pointerStruct, namedPointerStruct)
		methodTestDecls = append(methodTestDecls, methodTestDecl)
	}
	return
}

func createMethod(paramList, resultList []*ast.Field, methodName, receiverName string, namedPointerStruct *ast.Field) (method *ast.FuncDecl) {
	name := main.createName(methodName)
	args := createFieldList(paramList...)
	results := createFieldList(resultList...)
	recvs := createFieldList(namedPointerStruct)

	// prepare method body lines
	bodyList := make([]ast.Stmt, 0)
	for _, param := range paramList {
		// s.Field = Field
		lineSFieldAssignField := createAssignStmt(
			// s.Field
			main.createExprList(main.createSelectorExpr(main.createName(receiverName), main.createName(main.getNodeName(param)))),
			// =
			token.ASSIGN,
			// Field
			main.createExprList(main.createName(main.getNodeName(param))),
		)
		bodyList = append(bodyList, lineSFieldAssignField)
	}
	for _, result := range resultList {
		// Field = s.Field
		lineFieldAssignSField := createAssignStmt(
			// Field
			main.createExprList(main.createName(main.getNodeName(result))),
			// =
			token.ASSIGN,
			// s.Field
			main.createExprList(main.createSelectorExpr(main.createName(receiverName), main.createName(main.getNodeName(result)))),
		)
		bodyList = append(bodyList, lineFieldAssignSField)
	}
	bodyList = append(bodyList, returnStmt)

	// create method
	method = createFuncDecl(recvs, name, args, results,
		bodyList...,
	)

	return
}
