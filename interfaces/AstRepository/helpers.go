package AstRepository

import (
	"fmt"
	"go/ast"
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
	case *ast.ImportSpec:
		if nodeItem.Name != nil {
			name = nodeItem.Name.Name
		}
	default:
		panic(fmt.Sprintf("no getting name case for type %T", node))
	}
	return
}
