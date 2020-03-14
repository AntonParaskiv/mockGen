package main

import (
	"go/ast"
	"go/token"
)

func createSetterName(fieldName string) (setterName string) {
	setterName = "Set" + toPublic(fieldName)
	return
}

func createSetters(structSpec *ast.TypeSpec) (setterDecls []*ast.FuncDecl) {
	structName := getNodeName(structSpec)
	receiverName := getReceiverName(structName)
	fieldList := structSpec.Type.(*ast.StructType).Fields.List
	pointerToStruct := createFieldFromExpr(createPointerStruct(structName))
	namedPointerToStruct := createFieldNamedPointerStruct(structName, receiverName)

	setterDecls = make([]*ast.FuncDecl, 0)
	for _, field := range fieldList {
		setterName := createSetterName(getNodeName(field))

		setter := &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					namedPointerToStruct,
				},
			},
			Name: &ast.Ident{
				Name: setterName,
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						field,
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						pointerToStruct,
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: receiverName,
								},
								Sel: &ast.Ident{
									Name: getNodeName(field),
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.Ident{
								Name: getNodeName(field),
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: receiverName,
							},
						},
					},
				},
			},
		}
		setterDecls = append(setterDecls, setter)
	}

	return
}
