package main

import (
	"go/ast"
	"go/token"
)

func createTestNewName(newName string) (newTestName string) {
	newTestName = "Test" + newName
	return
}

func createConstructorTest(structName, functionName, testFunctionName, wantReceiver, gotReceiver string, pointerStruct *ast.StarExpr) (constructorTestDecl *ast.FuncDecl) {
	constructorTestDecl = &ast.FuncDecl{
		Name: &ast.Ident{
			Name: testFunctionName,
		},
		Type: testFunctionType,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "tests",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.ArrayType{
								Elt: &ast.StructType{
									Fields: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													{
														Name: "name",
													},
												},
												Type: &ast.Ident{
													Name: "string",
												},
											},
											{
												Names: []*ast.Ident{
													{
														Name: wantReceiver,
													},
												},
												Type: pointerStruct,
											},
										},
									},
								},
							},
							Elts: []ast.Expr{
								&ast.CompositeLit{
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "name",
											},
											Colon: 0,
											Value: &ast.BasicLit{
												Kind:  token.STRING,
												Value: `"Success"`,
											},
										},
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: wantReceiver,
											},
											Colon: 0,
											Value: &ast.UnaryExpr{
												Op: token.AND,
												X: &ast.CompositeLit{
													Type: &ast.Ident{
														Name: structName,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.RangeStmt{
					Key: &ast.Ident{
						Name: "_",
					},
					Value: &ast.Ident{
						Name: "tt",
					},
					Tok: token.DEFINE,
					X: &ast.Ident{
						Name: "tests",
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "t",
										},
										Sel: &ast.Ident{
											Name: "Run",
										},
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "tt",
											},
											Sel: &ast.Ident{
												Name: "name",
											},
										},
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{
													List: []*ast.Field{
														{
															Names: []*ast.Ident{
																{
																	Name: "t",
																},
															},
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "testing",
																	},
																	Sel: &ast.Ident{
																		Name: "T",
																	},
																},
															},
														},
													},
												},
											},
											Body: &ast.BlockStmt{
												List: []ast.Stmt{
													&ast.IfStmt{
														Init: &ast.AssignStmt{
															Lhs: []ast.Expr{
																&ast.Ident{
																	Name: gotReceiver,
																},
															},
															Tok: token.DEFINE,
															Rhs: []ast.Expr{
																&ast.CallExpr{
																	Fun: &ast.Ident{
																		Name: functionName,
																	},
																	Lparen:   0,
																	Args:     nil,
																	Ellipsis: 0,
																	Rparen:   0,
																},
															},
														},
														Cond: &ast.UnaryExpr{
															Op: token.NOT,
															X: &ast.CallExpr{
																Fun: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "reflect",
																	},
																	Sel: &ast.Ident{
																		Name: "DeepEqual",
																	},
																},
																Args: []ast.Expr{
																	&ast.Ident{
																		Name: gotReceiver,
																	},
																	&ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "tt",
																		},
																		Sel: &ast.Ident{
																			Name: wantReceiver,
																		},
																	},
																},
															},
														},
														Body: &ast.BlockStmt{
															List: []ast.Stmt{
																&ast.ExprStmt{
																	X: &ast.CallExpr{
																		Fun: &ast.SelectorExpr{
																			X: &ast.Ident{
																				Name: "t",
																			},
																			Sel: &ast.Ident{
																				Name: "Errorf",
																			},
																		},
																		Args: []ast.Expr{
																			&ast.BasicLit{
																				Kind:  token.STRING,
																				Value: `"` + functionName + "() = %v, " + wantReceiver + " %v" + `"`,
																			},
																			&ast.Ident{
																				Name: gotReceiver,
																			},
																			&ast.SelectorExpr{
																				X: &ast.Ident{
																					Name: "tt",
																				},
																				Sel: &ast.Ident{
																					Name: wantReceiver,
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return
}
