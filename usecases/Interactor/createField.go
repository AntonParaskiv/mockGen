package Interactor

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
)

func (i *Interactor) createField(iFaceField *domain.Field) (mockField *domain.Field) {
	mockField = &domain.Field{
		Name:     iFaceField.Name,
		Type:     iFaceField.Type,
		BaseType: iFaceField.BaseType,
	}

	if mockField.GetTypeType() == domain.FieldTypeLocalCustomType {
		mockField.Type = fmt.Sprintf("%s.%s", i.interfacePackage.SelfImport.GetCallingName(), mockField.Type)
		mockField.CodeImportList = append(mockField.CodeImportList, i.interfacePackage.SelfImport)
		//mockField.TestImportList = append(mockField.TestImportList, i.interfacePackage.SelfImport)
	}

	return
}
