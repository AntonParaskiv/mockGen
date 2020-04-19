package AstRepository

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"go/ast"
)

func createFieldList(astFieldList []*ast.Field) (fieldList []*domain.Field, err error) {
	for _, astField := range astFieldList {
		var fields []*domain.Field
		fields, err = createFields(astField)
		if err != nil {
			err = fmt.Errorf("create fields from ast field failed: %w", err)
			return
		}
		if fields == nil {
			continue
		}
		fieldList = append(fieldList, fields...)
	}
	return
}

func createFields(astField *ast.Field) (fields []*domain.Field, err error) {
	fieldType, err := getFieldType(astField.Type)
	if err != nil {
		err = fmt.Errorf("get field type failed: %w", err)
		return
	}

	for _, ident := range astField.Names {
		field := &domain.Field{
			Name: ident.Name,
			Type: fieldType,
		}
		fields = append(fields, field)
	}
	return
}
