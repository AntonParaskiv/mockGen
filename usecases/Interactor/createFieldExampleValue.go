package Interactor

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"strings"
)

func (i *Interactor) createFieldExampleValue(field *domain.Field) {
	switch field.GetTypeGroup() {
	case domain.FieldTypeGroupBool:
		createBoolExampleValue(field)
	case domain.FieldTypeGroupString:
		createMyFieldNameExampleValue(field)
	case domain.FieldTypeGroupNumber:
		createIntExampleValue(field)
	case domain.FieldTypeGroupByte:
		createByteExampleValue(field)
	case domain.FieldTypeGroupRune:
		createRuneExampleValue(field)
	default:
		switch field.GetTypeType() {

		case domain.FieldTypeInterface:
			createMyFieldNameExampleValue(field)
		case domain.FieldTypeError:
			createErrorExampleValue(field)
		case domain.FieldTypeArray:
			i.createArrayExampleValue(field)
		case domain.FieldTypeMap:
			i.createMapExampleValue(field)
		case domain.FieldTypeStruct:
			i.createStructExampleValue(field)
		case domain.FieldTypeImportedCustomType:
			i.createImportedCustomExampleValue(field)
		//case domain.FieldTypeLocalCustomType:
		//	i.createImportedCustomExampleValue(field)
		case domain.FieldTypePointer:
			i.createPointerExampleValue(field)
		default:
			fmt.Printf("create field example for %s failed: unknown type\n", field.Type)
		}
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
	if len(subField.ExampleValue) == 0 {
		return
	}
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

func (i *Interactor) createStructExampleValue(Struct *domain.Field) {
	fieldList := Struct.GetStructFieldList()
	field := i.getFirstBasicTypeField(fieldList)
	if field == nil {
		Struct.ExampleValue = fmt.Sprintf("{}")
		return
	}

	i.createFieldExampleValue(field)
	Struct.ExampleValue = fmt.Sprintf("{\n")
	Struct.ExampleValue += fmt.Sprintf("	%s: %s,\n", field.Name, field.ExampleValue)
	Struct.ExampleValue += fmt.Sprintf("}")
}

func (i *Interactor) getFirstBasicTypeField(fieldList []*domain.Field) (field *domain.Field) {
	for _, fieldItem := range fieldList {
		if fieldItem.GetTypeGroup() > 0 {
			field = fieldItem
			return
		}
	}
	return
}

func (i *Interactor) createImportedCustomExampleValue(field *domain.Field) {
	// check self import
	//for _, Import := range field.CodeImportList {
	//	if Import == i.interfacePackage.SelfImport {
	//		//mock := i.mockPackage.GetMockByName(field.GetTypeName())
	//		//_ = mock
	//		//importList = append(importList, field.CodeImportList...)
	//		return
	//	}
	//}

	splitted := strings.Split(field.Type, ".")
	if len(splitted) != 2 {
		err := fmt.Errorf("parse imported custom type failed: %s", field.Type)
		fmt.Println(err)
		return
	}
	importKey := splitted[0]
	typeName := splitted[1]
	packagePath := ""
	for _, Import := range field.CodeImportList {
		if Import.GetCallingName() == importKey {
			packagePath = Import.Path
			field.TestImportList = append(field.TestImportList, Import)
		}
	}
	if len(packagePath) == 0 {
		for _, Import := range i.mockFile.ImportList {
			if Import.GetCallingName() == importKey {
				packagePath = Import.Path
				field.TestImportList = append(field.TestImportList, Import)
			}
		}
	}

	baseType, err := i.AstRepository.GetTypeFieldFromPackagePath(packagePath, typeName)
	if err != nil {
		err = fmt.Errorf("get base type of %s failed: %w", field.Type, err)
		fmt.Println(err.Error())
		return
	}
	i.createFieldExampleValue(baseType)
	field.BaseType = baseType
	if field.BaseType.GetTypeType() == domain.FieldTypeStruct {
		field.ExampleValue = fmt.Sprintf("%s%s", field.Type, field.BaseType.ExampleValue)
	} else {
		field.ExampleValue = fmt.Sprintf("%s(%s)", field.Type, field.BaseType.ExampleValue)
	}

	field.TestImportList = append(field.TestImportList, baseType.CodeImportList...)
}

func (i *Interactor) createPointerExampleValue(field *domain.Field) {
	baseField := &domain.Field{Type: field.Type[1:], Name: field.Name, CodeImportList: field.CodeImportList}
	i.createFieldExampleValue(baseField)

	switch baseField.GetTypeType() {
	case domain.FieldTypeBool:
		field.ExampleValue = `GetPointer.Bool(` + baseField.ExampleValue + `)`
	case domain.FieldTypeString:
		field.ExampleValue = `GetPointer.String(` + baseField.ExampleValue + `)`
	case domain.FieldTypeInt:
		field.ExampleValue = `GetPointer.Int(` + baseField.ExampleValue + `)`
	case domain.FieldTypeInt8:
		field.ExampleValue = `GetPointer.Int8(` + baseField.ExampleValue + `)`
	case domain.FieldTypeInt16:
		field.ExampleValue = `GetPointer.Int16(` + baseField.ExampleValue + `)`
	case domain.FieldTypeInt32:
		field.ExampleValue = `GetPointer.Int32(` + baseField.ExampleValue + `)`
	case domain.FieldTypeInt64:
		field.ExampleValue = `GetPointer.Int64(` + baseField.ExampleValue + `)`
	case domain.FieldTypeUint:
		field.ExampleValue = `GetPointer.Uint(` + baseField.ExampleValue + `)`
	case domain.FieldTypeUint8:
		field.ExampleValue = `GetPointer.Uint8(` + baseField.ExampleValue + `)`
	case domain.FieldTypeUint16:
		field.ExampleValue = `GetPointer.Uint16(` + baseField.ExampleValue + `)`
	case domain.FieldTypeUint32:
		field.ExampleValue = `GetPointer.Uint32(` + baseField.ExampleValue + `)`
	case domain.FieldTypeUint64:
		field.ExampleValue = `GetPointer.Uint64(` + baseField.ExampleValue + `)`
	case domain.FieldTypeUintPtr:
		field.ExampleValue = `GetPointer.UintPtr(` + baseField.ExampleValue + `)`
	case domain.FieldTypeByte:
		field.ExampleValue = `GetPointer.Byte(` + baseField.ExampleValue + `)`
	case domain.FieldTypeRune:
		field.ExampleValue = `GetPointer.Rune(` + baseField.ExampleValue + `)`
	case domain.FieldTypeFloat32:
		field.ExampleValue = `GetPointer.Float32(` + baseField.ExampleValue + `)`
	case domain.FieldTypeFloat64:
		field.ExampleValue = `GetPointer.Float64(` + baseField.ExampleValue + `)`
	case domain.FieldTypeComplex64:
		field.ExampleValue = `GetPointer.Complex64(` + baseField.ExampleValue + `)`
	case domain.FieldTypeComplex128:
		field.ExampleValue = `GetPointer.Complex128(` + baseField.ExampleValue + `)`

	case domain.FieldTypeError:
	case domain.FieldTypeArray:
	case domain.FieldTypeMap:
	case domain.FieldTypeImportedCustomType:
		field.ExampleValue = `&` + baseField.ExampleValue
	case domain.FieldTypeLocalCustomType:
	//case domain.FieldTypeInterface:
	//case domain.FieldTypePointer:
	default:
	}

	return
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
