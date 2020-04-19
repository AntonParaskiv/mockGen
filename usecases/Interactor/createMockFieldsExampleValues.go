package Interactor

import "github.com/AntonParaskiv/mockGen/domain"

func (i *Interactor) createMockFieldsExampleValues(mock *domain.Mock) {
	for _, field := range mock.Struct.FieldList {
		i.createFieldExampleValue(field)
		field.CodeImportList = append(field.CodeImportList, field.CodeImportList...) // TODO: add unique
		field.TestImportList = append(field.TestImportList, field.TestImportList...) // TODO: add unique
	}
	return
}
