package main

import "go/ast"

func createFunctionType(params, results *ast.FieldList) (funcType *ast.FuncType) {
	funcType = &ast.FuncType{
		Params:  params,
		Results: results,
	}
	return
}

func createFuncDecl(recvs *ast.FieldList, name *ast.Ident, args, results *ast.FieldList, bodyLines ...ast.Stmt) (constructorDecl *ast.FuncDecl) {
	constructorDecl = &ast.FuncDecl{
		Recv: recvs,
		Name: name,
		Type: createFunctionType(args, results),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{},
		},
	}

	for _, line := range bodyLines {
		constructorDecl.Body.List = append(constructorDecl.Body.List, line)
	}

	return
}

func createCallExpr(fun ast.Expr, args ...ast.Expr) (callExpr *ast.CallExpr) {
	callExpr = &ast.CallExpr{
		Fun:  fun,
		Args: []ast.Expr{},
	}
	for _, arg := range args {
		callExpr.Args = append(callExpr.Args, arg)
	}
	return
}

func createFuncLit(Type *ast.FuncType, body *ast.BlockStmt) (funcLit *ast.FuncLit) {
	funcLit = &ast.FuncLit{
		Type: Type,
		Body: body,
	}
	return
}
