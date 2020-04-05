package usecases

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"strings"
)

type Interactor struct {
}

func (i *Interactor) CreateMockPackage(interfacePackage *domain.GoCodePackage) (mockPackage *domain.GoCodePackage) {
	mockPackagePath := cutPostfix(interfacePackage.Path, "Interface") + "Mock"
	mockPackageName := cutPostfix(interfacePackage.PackageName, "Interface") + "Mock"

	mockPackage = &domain.GoCodePackage{
		Path:        mockPackagePath,
		PackageName: mockPackageName,
	}

	for _, interfaceFile := range interfacePackage.FileList {
		mockFile := createMockFilesFromInterfaceFile(interfaceFile, mockPackageName)
		mockPackage.FileList = append(mockPackage.FileList, mockFile)
	}

	return
}

func createMockFilesFromInterfaceFile(interfaceFile *domain.GoCodeFile, mockPackageName string) (mockFile *domain.GoCodeFile) {
	mockFile = &domain.GoCodeFile{
		Name:       interfaceFile.Name,
		ImportList: append([]*domain.Import{}, interfaceFile.ImportList...),
	}

	for _, iFace := range interfaceFile.InterfaceList {
		mock := createMockFromInterface(iFace, mockPackageName)
		mockFile.MockList = append(mockFile.MockList, mock)
	}

	return
}

func createMockFromInterface(iFace *domain.Interface, mockPackageName string) (mock *domain.Mock) {
	structName := iFace.Name

	basePackageName := cutPostfix(mockPackageName, "Mock")

	var constructorName string
	if structName == basePackageName {
		constructorName = "New"
	} else {
		constructorName = "New" + toPublic(structName)
	}

	mock = &domain.Mock{
		Struct: &domain.Struct{
			Name:         structName,
			ReceiverName: getReceiverName(structName),
			WantName:     "want" + toPublic(structName),
			GotName:      "got" + toPublic(structName),
		},
		Constructor: &domain.Constructor{
			Name: constructorName,
		},
		MethodList: iFace.MethodList,
	}

	argList := []*domain.Field{}
	resultList := []*domain.Field{}
	for _, method := range mock.MethodList {
		for _, arg := range method.ArgList {
			arg.WantName = "want" + toPublic(arg.Name)
			arg.GotName = "got" + toPublic(arg.Name)
			arg.NameType = arg.Name + " " + arg.Type
			arg.ExampleValue, arg.TestImportList = createExampleValue(arg.Type, arg.Name)

			method.ArgNameTypeList = append(method.ArgNameTypeList, arg.NameType)
			method.CodeImportList = append(method.CodeImportList, arg.CodeImportList...) // TODO: add unique
			method.TestImportList = append(method.TestImportList, arg.TestImportList...) // TODO: add unique
		}
		for _, result := range method.ResultList {
			result.WantName = "want" + toPublic(result.Name)
			result.GotName = "got" + toPublic(result.Name)
			result.NameType = result.Name + " " + result.Type
			result.ExampleValue, result.TestImportList = createExampleValue(result.Type, result.Name)

			method.ResultNameTypeList = append(method.ResultNameTypeList, result.NameType)
			method.CodeImportList = append(method.CodeImportList, result.CodeImportList...) // TODO: add unique
			method.TestImportList = append(method.TestImportList, result.TestImportList...) // TODO: add unique
		}

		argList = append(argList, method.ArgList...)
		resultList = append(resultList, method.ResultList...)
		mock.CodeImportList = append(mock.CodeImportList, method.CodeImportList...) // TODO: add unique
		mock.TestImportList = append(mock.TestImportList, method.TestImportList...) // TODO: add unique
	}

ArgLoop:
	for _, arg := range argList {
		for _, mockField := range mock.Struct.FieldList {
			if mockField.Name == arg.Name {
				continue ArgLoop
			}
		}
		mock.Struct.FieldList = append(mock.Struct.FieldList, arg)
	}

ResultLoop:
	for _, result := range resultList {
		for _, mockField := range mock.Struct.FieldList {
			if mockField.Name == result.Name {
				continue ResultLoop
			}
		}
		mock.Struct.FieldList = append(mock.Struct.FieldList, result)
	}

	for _, field := range mock.Struct.FieldList {
		mock.SetterList = append(mock.SetterList, &domain.Setter{
			Name:  "Set" + toPublic(field.Name),
			Field: field,
		})
	}
	return
}

func createExampleValue(fieldType string, fieldName string) (exampleValue string, importList []*domain.Import) {
	switch {
	case fieldType == "string":
		exampleValue = `"my` + toPublic(fieldName) + `"`
	case fieldType == "interface{}":
		exampleValue = `"my` + toPublic(fieldName) + `"`
	case len(fieldType) >= 3 && fieldType[0:3] == "int": // int must be after interface !
		exampleValue = "100"
	case len(fieldType) >= 4 && fieldType[0:4] == "uint":
		exampleValue = "200"
	case len(fieldType) >= 5 && fieldType[0:5] == "float":
		exampleValue = "3.14"
	case fieldType == "bool":
		exampleValue = "true"
	case fieldType == "rune":
		exampleValue = "'X'"
	case fieldType == "byte":
		exampleValue = "50"
	case fieldType == "error":
		exampleValue = `fmt.Errorf("simulated error")`
		importList = append(importList, &domain.Import{Key: "fmt", Path: "fmt"})
	case len(fieldType) >= 2 && fieldType[0:2] == "[]":
		itemType := fieldType[2:]
		itemExampleValue, itemImportList := createExampleValue(itemType, itemType+"Example")
		exampleValue = fmt.Sprintf("%s{\n", fieldType)
		exampleValue += fmt.Sprintf("	%s,\n", itemExampleValue)
		exampleValue += fmt.Sprintf("}")
		importList = append(importList, itemImportList...)
	case len(fieldType) >= 4 && fieldType[0:4] == "map[":
		keyType, valueType := getMapKeyValueTypes(fieldType)
		keyExampleValue, keyImportList := createExampleValue(keyType, keyType+"Example")
		valueExampleValue, valueImportList := createExampleValue(valueType, valueType+"Example")

		exampleValue = fmt.Sprintf("%s{\n", fieldType)
		exampleValue += fmt.Sprintf("	%s: %s,\n", keyExampleValue, valueExampleValue)
		exampleValue += fmt.Sprintf("}")

		importList = append(importList, keyImportList...)
		importList = append(importList, valueImportList...)

	default:
		fmt.Println("unknown type:", fieldType)

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

func cutPostfix(text, postfix string) (shortCutText string) {
	lenPostfix := len(postfix)
	if len(text) > lenPostfix {
		startPostfix := len(text) - lenPostfix
		packageNamePostfix := text[startPostfix:]
		if packageNamePostfix == postfix {
			shortCutText = text[0:startPostfix]
		}
	}
	return
}

func toPublic(name string) (publicName string) {
	firstLetterUpper := strings.ToUpper(getFirstLetter(name))
	publicName = firstLetterUpper + getFollowingLetters(name)
	return
}

func toPrivate(name string) (privateName string) {
	firstLetterLower := strings.ToLower(getFirstLetter(name))
	privateName = firstLetterLower + getFollowingLetters(name)
	return
}

func getFirstLetter(text string) (firstLetter string) {
	firstLetter = text[0:1]
	return
}

func getFollowingLetters(text string) (followingLetters string) {
	followingLetters = text[1:]
	return
}

// Mock -> s
func getReceiverName(name string) (receiverName string) {
	receiverName = toPrivate(getFirstLetter(name))
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
