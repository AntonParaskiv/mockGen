package AstRepository

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"go/ast"
)

func (r *Repository) createInterface(astInterfaceSpec *ast.TypeSpec) (iFace *domain.Interface, err error) {
	methodList := make([]*domain.Method, 0)
	switch astInterfaceType := astInterfaceSpec.Type.(type) {
	case *ast.InterfaceType:
		for _, astMethod := range astInterfaceType.Methods.List {
			var method *domain.Method
			method, err = r.createMethod(astMethod)
			if err != nil {
				err = fmt.Errorf("create method failed: %w", err)
				return
			}
			if method == nil {
				continue
			}
			methodList = append(methodList, method)
		}
	default:
		err = fmt.Errorf("ast spec type is not interface type")
		return
	}

	if len(methodList) == 0 {
		return
	}

	iFace = &domain.Interface{
		Name:       getNodeName(astInterfaceSpec),
		MethodList: methodList,
	}

	return
}
