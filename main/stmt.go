package main

import (
	"go/ast"
	"go/token"
)

var returnStmt = &ast.ReturnStmt{}

func createAssignStmt(leftExpressions []ast.Expr, tok token.Token, rightExpressions []ast.Expr) (stmt *ast.AssignStmt) {
	stmt = &ast.AssignStmt{
		Lhs: leftExpressions,
		Tok: tok,
		Rhs: rightExpressions,
	}
	return
}

func createIfStmt(init ast.Stmt, cond ast.Expr, bodyLines ...ast.Stmt) (ifStmt *ast.IfStmt) {
	ifStmt = &ast.IfStmt{
		Init: init,
		Cond: cond,
		Body: &ast.BlockStmt{
			List: bodyLines,
		},
	}
	return
}
