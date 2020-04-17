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

	i.fillFieldType(mockField)
	return
}

func (i *Interactor) fillFieldType(field *domain.Field) {
	switch field.GetTypeType() {
	case domain.FieldTypePointer:
		baseField := &domain.Field{Type: field.Type[1:], Name: field.Name, CodeImportList: field.CodeImportList}
		i.fillFieldType(baseField)
		field.Type = "*" + baseField.Type
		field.BaseType = baseField
		field.CodeImportList = append(field.CodeImportList, baseField.CodeImportList...)
	case domain.FieldTypeArray:
		baseField := &domain.Field{Type: field.Type[2:], Name: field.Name, CodeImportList: field.CodeImportList}
		i.fillFieldType(baseField)
		field.Type = "[]" + baseField.Type
		field.BaseType = baseField
		field.CodeImportList = append(field.CodeImportList, baseField.CodeImportList...)
	case domain.FieldTypeMap:
		keyType, valueType := getMapKeyValueTypes(field.Type)
		keyField := &domain.Field{Type: keyType, Name: keyType}
		valueField := &domain.Field{Type: valueType, Name: valueType}
		i.createFieldExampleValue(keyField)
		i.createFieldExampleValue(valueField)
		field.Type = fmt.Sprintf("map[%s]%s", keyType, valueType)
		field.CodeImportList = append(field.CodeImportList, keyField.CodeImportList...)
		field.CodeImportList = append(field.CodeImportList, valueField.CodeImportList...)
	case domain.FieldTypeLocalCustomType:
		field.Type = fmt.Sprintf("%s.%s", i.interfacePackage.SelfImport.GetCallingName(), field.Type)
		field.CodeImportList = append(field.CodeImportList, i.interfacePackage.SelfImport)
		//field.TestImportList = append(field.TestImportList, i.interfacePackage.SelfImport)
	}
}
