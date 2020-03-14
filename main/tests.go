package main

import "go/ast"

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
