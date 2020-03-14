package main

import (
	"go/ast"
	"go/token"
	"strings"
)

func createNewName(structName, packageName string) (newName string) {
	newName = "New"
	if !isStructNameSameAsPackageName(structName, packageName) {
		newName += toPublic(structName)
	}
	return
}

func isStructNameSameAsPackageName(structName, packageName string) (isSame bool) {
	if strings.ToUpper(structName) == strings.ToUpper(packageName) {
		isSame = true
	}
	return
}

func createConstructorAndTest(structSpec *ast.TypeSpec, packageName string) (constructorDecl, constructorTestDecl *ast.FuncDecl) {
	structName := getNodeName(structSpec)
	receiverName := getReceiverName(structName)
	functionName := createNewName(structName, packageName)
	namedPointerToStruct := createFieldNamedPointerStruct(structName, receiverName)
	constructorDecl = createConstructor(structName, receiverName, functionName, namedPointerToStruct)

	testFunctionName := createTestNewName(functionName)
	wantReceiver := "want" + toPublic(receiverName)
	gotReceiver := "got" + toPublic(receiverName)
	pointerToStruct := createPointerStruct(structName)
	constructorTestDecl = createConstructorTest(structName, functionName, testFunctionName, wantReceiver, gotReceiver, pointerToStruct)

	return
}

func createConstructor(structName, receiverName, functionName string, namedPointerToStruct *ast.Field) (constructorDecl *ast.FuncDecl) {
	constructorDecl = &ast.FuncDecl{
		Name: &ast.Ident{
			Name: functionName,
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					namedPointerToStruct,
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: receiverName,
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "new",
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: structName,
								},
							},
							Ellipsis: token.NoPos,
						},
					},
				},
				&ast.ReturnStmt{},
			},
		},
	}
	return
}
