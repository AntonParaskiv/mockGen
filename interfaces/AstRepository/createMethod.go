package AstRepository

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"go/ast"
)

func (r *Repository) createMethod(astMethod *ast.Field) (method *domain.Method, err error) {
	r.currentMethod = &domain.Method{
		Name: getNodeName(astMethod),
	}

	switch astFuncType := astMethod.Type.(type) {
	case *ast.FuncType:
		if astFuncType.Params != nil {
			r.currentMethod.ArgList, err = r.createFieldList(astFuncType.Params.List)
			if err != nil {
				err = fmt.Errorf("create arg list from ast param list failed: %w", err)
				return
			}
		}
		if astFuncType.Results != nil {
			r.currentMethod.ResultList, err = r.createFieldList(astFuncType.Results.List)
			if err != nil {
				err = fmt.Errorf("create result list from ast result list failed: %w", err)
				return
			}
		}
	default:
		err = fmt.Errorf("invalid ast method type: %v", astMethod.Type)
		return
	}

	method = r.currentMethod
	return
}
