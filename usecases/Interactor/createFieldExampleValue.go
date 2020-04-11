package Interactor

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"strings"
)

func (i *Interactor) createFieldExampleValue(field *domain.Field) {
	switch field.GetTypeType() {
	case domain.FieldTypeString:
		createMyFieldNameExampleValue(field)
	case domain.FieldTypeInterface:
		createMyFieldNameExampleValue(field)
	case domain.FieldTypeInts:
		createIntExampleValue(field)
	case domain.FieldTypeFloat:
		createFloatExampleValue(field)
	case domain.FieldTypeBool:
		createBoolExampleValue(field)
	case domain.FieldTypeRune:
		createRuneExampleValue(field)
	case domain.FieldTypeByte:
		createByteExampleValue(field)
	case domain.FieldTypeError:
		createErrorExampleValue(field)
	case domain.FieldTypeArray:
		i.createArrayExampleValue(field)
	case domain.FieldTypeMap:
		i.createMapExampleValue(field)
	case domain.FieldTypeImportedCustomType:
		i.createImportedCustomExampleValue(field)
	//case domain.FieldTypeLocalCustomType:
	//	i.createImportedCustomExampleValue(field)
	default:
		fmt.Println("unknown type:", field.Type)
		// TODO: struct
		// TODO: custom type
		// TODO: ptr string
		// TODO: ptr int
		// TODO: ptr uint
		// TODO: ptr float
		// TODO: ptr bool
		// TODO: ptr rune
		// TODO: ptr byte
		// TODO: ptr error
		// TODO: ptr array
		// TODO: ptr map
	}
	return
}

func createMyFieldNameExampleValue(field *domain.Field) {
	field.ExampleValue = `"my` + field.GetPublicName() + `"`
	return
}

func createIntExampleValue(field *domain.Field) {
	field.ExampleValue = "100"
	return
}

func createFloatExampleValue(field *domain.Field) {
	field.ExampleValue = "3.14"
	return
}

func createBoolExampleValue(field *domain.Field) {
	field.ExampleValue = "true"
	return
}

func createRuneExampleValue(field *domain.Field) {
	field.ExampleValue = "'X'"
	return
}

func createByteExampleValue(field *domain.Field) {
	field.ExampleValue = "50"
	return
}

func createErrorExampleValue(field *domain.Field) {
	field.ExampleValue = `fmt.Errorf("simulated error")`
	field.TestImportList = append(field.TestImportList, &domain.Import{Path: "fmt"})
	return
}

func (i *Interactor) createArrayExampleValue(field *domain.Field) {
	itemType := field.Type[2:]
	subField := &domain.Field{Type: itemType, Name: itemType + "Example"}
	i.createFieldExampleValue(subField)
	field.ExampleValue = fmt.Sprintf("%s{\n", field.Type)
	field.ExampleValue += fmt.Sprintf("	%s,\n", subField.ExampleValue)
	field.ExampleValue += fmt.Sprintf("}")
	field.CodeImportList = append(field.CodeImportList, subField.CodeImportList...)
}

func (i *Interactor) createMapExampleValue(field *domain.Field) {
	keyType, valueType := getMapKeyValueTypes(field.Type)
	keyField := &domain.Field{Type: keyType, Name: keyType + "Example"}
	valueField := &domain.Field{Type: valueType, Name: valueType + "Example"}
	i.createFieldExampleValue(keyField)
	i.createFieldExampleValue(valueField)
	field.ExampleValue = fmt.Sprintf("%s{\n", field.Type)
	field.ExampleValue += fmt.Sprintf("	%s: %s,\n", keyField.ExampleValue, valueField.ExampleValue)
	field.ExampleValue += fmt.Sprintf("}")
	field.CodeImportList = append(field.CodeImportList, keyField.CodeImportList...)
	field.CodeImportList = append(field.CodeImportList, valueField.CodeImportList...)
}

func (i *Interactor) createImportedCustomExampleValue(field *domain.Field) {
	// check self import
	for _, Import := range field.CodeImportList {
		if Import == i.interfacePackage.SelfImport {
			//mock := i.mockPackage.GetMockByName(field.GetTypeName())
			//_ = mock
			//importList = append(importList, field.CodeImportList...)
			return
		}
	}
	splitted := strings.Split(field.Type, ".")
	if len(splitted) != 2 {
		err := fmt.Errorf("parse imported custom type failed: %s", field.Type)
		fmt.Println(err)
		return
	}
	importKey := splitted[0]
	typeName := splitted[1]
	packagePath := ""
	for _, Import := range i.mockFile.ImportList {
		if Import.GetCallingName() == importKey {
			packagePath = Import.Path
			field.TestImportList = append(field.TestImportList, Import)
		}
	}
	baseType, err := i.AstRepository.GetTypeFieldFromPackagePath(packagePath, typeName)
	if err != nil {
		err = fmt.Errorf("get base type of %s failed: %w", field.Type, err)
		fmt.Printf(err.Error())
		return
	}
	i.createFieldExampleValue(baseType)
	field.BaseType = baseType
	field.ExampleValue = fmt.Sprintf("%s(%s)", field.Type, field.BaseType.ExampleValue)
	field.TestImportList = append(field.TestImportList, baseType.CodeImportList...)
}

func getMapKeyValueTypes(fieldType string) (keyType, valueType string) {
	fieldType = fieldType[3:]
	openedBracketNum := 0
	i := 0
	var char rune

	for i, char = range fieldType {
		if char == '[' {
			openedBracketNum++
		}
		if char == ']' {
			openedBracketNum--
		}
		if openedBracketNum == 0 {
			i++
			break
		}
	}
	keyType = fieldType[1 : i-1]
	valueType = fieldType[i:]
	return
}
