package AstRepository

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"go/ast"
)

func createFieldList(astFieldList []*ast.Field) (fieldList []*domain.Field, err error) {
	for _, astField := range astFieldList {
		var field *domain.Field
		field, err = createField(astField)
		if err != nil {
			err = fmt.Errorf("create field from ast field failed: %w", err)
			return
		}
		if field == nil {
			continue
		}
		fieldList = append(fieldList, field)
	}
	return
}

func createField(astField *ast.Field) (field *domain.Field, err error) {
	fieldType, err := getFieldType(astField.Type)
	if err != nil {
		err = fmt.Errorf("get field type failed: %w", err)
		return
	}

	field = &domain.Field{
		Name: getNodeName(astField),
		Type: fieldType,
	}
	return
}
