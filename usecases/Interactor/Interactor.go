package Interactor

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"github.com/AntonParaskiv/mockGen/interfaces/AstRepository"
	"strings"
)

type Interactor struct {
	AstRepository    AstRepository.Repository
	mockFile         *domain.GoCodeFile
	interfacePackage *domain.GoCodePackage
	mockPackage      *domain.GoCodePackage
}

func (i *Interactor) createExampleValue(field *domain.Field) (exampleValue string, importList []*domain.Import) {
	switch field.GetTypeType() {
	case domain.FieldTypeString:
		exampleValue = `"my` + toPublic(field.Name) + `"`
	case domain.FieldTypeInterface:
		exampleValue = `"my` + toPublic(field.Name) + `"`
	case domain.FieldTypeInts:
		exampleValue = "100"
	case domain.FieldTypeFloat:
		exampleValue = "3.14"
	case domain.FieldTypeBool:
		exampleValue = "true"
	case domain.FieldTypeRune:
		exampleValue = "'X'"
	case domain.FieldTypeByte:
		exampleValue = "50"
	case domain.FieldTypeError:
		exampleValue = `fmt.Errorf("simulated error")`
		importList = append(importList, &domain.Import{Path: "fmt"})
	case domain.FieldTypeArray:
		itemType := field.Type[2:]
		itemExampleValue, itemImportList := i.createExampleValue(&domain.Field{Type: itemType, Name: itemType + "Example"})
		exampleValue = fmt.Sprintf("%s{\n", field.Type)
		exampleValue += fmt.Sprintf("	%s,\n", itemExampleValue)
		exampleValue += fmt.Sprintf("}")
		importList = append(importList, itemImportList...)
	case domain.FieldTypeMap:
		keyType, valueType := getMapKeyValueTypes(field.Type)
		keyExampleValue, keyImportList := i.createExampleValue(&domain.Field{Type: keyType, Name: keyType + "Example"})
		valueExampleValue, valueImportList := i.createExampleValue(&domain.Field{Type: valueType, Name: valueType + "Example"})

		exampleValue = fmt.Sprintf("%s{\n", field.Type)
		exampleValue += fmt.Sprintf("	%s: %s,\n", keyExampleValue, valueExampleValue)
		exampleValue += fmt.Sprintf("}")

		importList = append(importList, keyImportList...)
		importList = append(importList, valueImportList...)

	case domain.FieldTypeImportedCustomType:

		// check self import
		for _, Import := range field.CodeImportList {
			if Import == i.interfacePackage.SelfImport {
				mock := i.mockPackage.GetMockByName(field.GetTypeName())
				_ = mock
				importList = append(importList, field.CodeImportList...)
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
				importList = append(importList, Import)
			}
		}

		baseType, err := i.AstRepository.GetTypeFieldFromPackagePath(packagePath, typeName)
		if err != nil {
			err = fmt.Errorf("get base type of %s failed: %w", field.Type, err)
			fmt.Printf(err.Error())
			return
		}

		baseExampleValue, baseImportList := i.createExampleValue(baseType)
		baseType.ExampleValue = baseExampleValue
		baseType.TestImportList = baseImportList
		field.BaseType = baseType

		exampleValue = fmt.Sprintf("%s(%s)", field.Type, field.BaseType.ExampleValue)
		importList = append(importList, baseType.TestImportList...)

	case domain.FieldTypeLocalCustomType:
		//baseType, err := i.AstRepository.GetTypeFieldFromPackagePath(i.interfacePackage.FullPath, field.Type)
		//if err != nil {
		//	err = fmt.Errorf("get base type of %s failed: %w", field.Type, err)
		//	fmt.Printf(err.Error())
		//	return
		//}
		//_ = baseType

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
