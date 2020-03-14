package main

import (
	"go/ast"
	"go/token"
)

func createMethods(structSpec *ast.TypeSpec, interfaceSpec *ast.TypeSpec) (structMethodDecls []*ast.FuncDecl) {
	structName := getNodeName(structSpec)
	receiverName := getReceiverName(structName)
	namedPointerToStruct := createFieldNamedPointerStruct(structName, receiverName)

	structMethodDecls = make([]*ast.FuncDecl, 0)
	for _, interfaceMethod := range interfaceSpec.Type.(*ast.InterfaceType).Methods.List {
		methodName := getNodeName(interfaceMethod)
		methodParams := interfaceMethod.Type.(*ast.FuncType).Params
		methodResults := interfaceMethod.Type.(*ast.FuncType).Results

		bodyList := make([]ast.Stmt, 0)

		for _, param := range interfaceMethod.Type.(*ast.FuncType).Params.List {
			paramName := getNodeName(param)
			setting := &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.Ident{
							Name: receiverName,
						},
						Sel: &ast.Ident{
							Name: paramName,
						},
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.Ident{
						Name: paramName,
					},
				},
			}
			bodyList = append(bodyList, setting)
		}

		for _, result := range interfaceMethod.Type.(*ast.FuncType).Results.List {
			resultName := getNodeName(result)
			returning := &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: resultName,
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.Ident{
							Name: receiverName,
						},
						Sel: &ast.Ident{
							Name: resultName,
						},
					},
				},
			}
			bodyList = append(bodyList, returning)
		}

		returning := &ast.ReturnStmt{}
		bodyList = append(bodyList, returning)

		structMethod := &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					namedPointerToStruct,
				},
			},
			Name: &ast.Ident{
				Name: methodName,
			},
			Type: &ast.FuncType{
				Func:    0,
				Params:  methodParams,
				Results: methodResults,
			},
			Body: &ast.BlockStmt{
				List: bodyList,
			},
		}
		structMethodDecls = append(structMethodDecls, structMethod)
	}
	return
}
