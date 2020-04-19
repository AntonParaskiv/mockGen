package AstRepository

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"go/ast"
)

func (r *Repository) createFieldList(astFieldList []*ast.Field) (fieldList []*domain.Field, err error) {
	for fieldIndex, astField := range astFieldList {
		var fields []*domain.Field
		fields, err = r.createFields(astField, fieldIndex+1)
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

func (r *Repository) createFields(astField *ast.Field, fieldNumber int) (fields []*domain.Field, err error) {
	fieldType, err := getFieldType(astField.Type)
	if err != nil {
		err = fmt.Errorf("get field type failed: %w", err)
		return
	}

	if len(astField.Names) == 0 {
		field := &domain.Field{
			Name: fmt.Sprintf("%sResult%d", r.currentMethod.GetPrivateName(), fieldNumber),
			Type: fieldType,
		}
		fields = append(fields, field)
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
