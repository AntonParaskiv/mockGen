package AstRepository

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"go/ast"
)

func createMethod(astMethod *ast.Field) (method *domain.Method, err error) {
	argList := make([]*domain.Field, 0)
	resultList := make([]*domain.Field, 0)

	switch astFuncType := astMethod.Type.(type) {
	case *ast.FuncType:
		if astFuncType.Params != nil {
			argList, err = createFieldList(astFuncType.Params.List)
			if err != nil {
				err = fmt.Errorf("create arg list from ast param list failed: %w", err)
				return
			}
		}
		if astFuncType.Results != nil {
			resultList, err = createFieldList(astFuncType.Results.List)
			if err != nil {
				err = fmt.Errorf("create result list from ast result list failed: %w", err)
				return
			}
		}
	default:
		err = fmt.Errorf("invalid ast method type: %v", astMethod.Type)
		return
	}

	method = &domain.Method{
		Name:       getNodeName(astMethod),
		ArgList:    argList,
		ResultList: resultList,
	}
	return
}
