package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func getNodeName(node ast.Node) (name string) {
	switch nodeItem := node.(type) {
	case *ast.Package:
		name = nodeItem.Name
	case *ast.File:
		name = nodeItem.Name.Name
	case *ast.TypeSpec:
		name = nodeItem.Name.Name
	case *ast.Field:
		name = nodeItem.Names[0].Name
	case *ast.Ident:
		name = nodeItem.Name
	default:
		panic(fmt.Sprintf("no getting name case for type %T", node))
	}
	return
}

func toPublic(name string) (publicName string) {
	firstLetterUpper := strings.ToUpper(getFirstLetter(name))
	publicName = firstLetterUpper + getFollowingLetters(name)
	return
}

func toPrivate(name string) (privateName string) {
	firstLetterLower := strings.ToLower(getFirstLetter(name))
	privateName = firstLetterLower + getFollowingLetters(name)
	return
}

func getFirstLetter(text string) (firstLetter string) {
	firstLetter = text[0:1]
	return
}

func getFollowingLetters(text string) (followingLetters string) {
	followingLetters = text[1:]
	return
}

func createName(name string) (names *ast.Ident) {
	names = &ast.Ident{
		Name: name,
	}
	return
}

func createNames(name string) (names []*ast.Ident) {
	names = []*ast.Ident{
		{
			Name: name,
		},
	}
	return
}

func createExprList(exprs ...ast.Expr) (exprList []ast.Expr) {
	for _, expr := range exprs {
		exprList = append(exprList, expr)
	}
	return
}

func createKeyValueExpr(key string, value ast.Expr) (keyValueExpr *ast.KeyValueExpr) {
	keyValueExpr = &ast.KeyValueExpr{
		Key:   createName(key),
		Value: value,
	}
	return
}

func createSelectorExpr(expression ast.Expr, fieldSelector *ast.Ident) (selectorExpr *ast.SelectorExpr) {
	selectorExpr = &ast.SelectorExpr{
		X:   expression,
		Sel: fieldSelector,
	}
	return
}

func createBasicLit(value string, kind token.Token) (basicLit *ast.BasicLit) {
	basicLit = &ast.BasicLit{
		Kind:  kind,
		Value: value,
	}
	return
}

func createUnaryExpr(operator token.Token, operand ast.Expr) (unaryExpr *ast.UnaryExpr) {
	unaryExpr = &ast.UnaryExpr{
		Op: operator,
		X:  operand,
	}
	return
}

func createCompositeLit(Type ast.Expr, elts ...ast.Expr) (compositeLit *ast.CompositeLit) {
	compositeLit = &ast.CompositeLit{
		Type: Type,
		Elts: elts,
	}
	return
}

func createRangeStmt(key, value string, tok token.Token, rangeName string, stmts ...ast.Stmt) (rangeStmt *ast.RangeStmt) {
	rangeStmt = &ast.RangeStmt{
		Key:   createName(key),
		Value: createName(value),
		Tok:   tok,
		X:     createName(rangeName),
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}
	return
}
