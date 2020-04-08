package usecases

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

func (i *Interactor) CreateMockPackage(interfacePackage *domain.GoCodePackage) (mockPackage *domain.GoCodePackage) {
	i.interfacePackage = interfacePackage

	mockPackagePath := cutPostfix(interfacePackage.Path, "Interface") + "Mock"
	mockPackageName := cutPostfix(interfacePackage.PackageName, "Interface") + "Mock"

	mockPackage = &domain.GoCodePackage{
		Path:        mockPackagePath,
		PackageName: mockPackageName,
	}

	i.mockPackage = mockPackage

	for _, interfaceFile := range interfacePackage.FileList {
		mockFile := i.createMockFilesFromInterfaceFile(interfaceFile, mockPackageName)
		mockPackage.FileList = append(mockPackage.FileList, mockFile)
	}

	return
}

func (i *Interactor) createMockFilesFromInterfaceFile(interfaceFile *domain.GoCodeFile, mockPackageName string) (mockFile *domain.GoCodeFile) {
	i.mockFile = &domain.GoCodeFile{
		Name:       interfaceFile.Name,
		ImportList: append([]*domain.Import{}, interfaceFile.ImportList...),
	}

	for _, iFace := range interfaceFile.InterfaceList {
		mock := i.createMockFromInterface(iFace, mockPackageName)
		i.mockFile.MockList = append(i.mockFile.MockList, mock)
	}

	return i.mockFile
}

func (i *Interactor) createMockMethodList(iFaceMethodList []*domain.Method) (mockMethodList []*domain.Method) {
	for _, iFaceMethod := range iFaceMethodList {
		mockMethod := &domain.Method{
			Name:               iFaceMethod.Name,
			ResultNameTypeList: nil,
		}

		for _, iFaceArg := range iFaceMethod.ArgList {
			mockArg := &domain.Field{
				Name:     iFaceArg.Name,
				WantName: iFaceArg.WantName,
				GotName:  iFaceArg.GotName,
				Type:     iFaceArg.Type,
				BaseType: iFaceArg.BaseType,
				NameType: iFaceArg.NameType,
			}
			mockMethod.ArgList = append(mockMethod.ArgList, mockArg)
		}

		for _, iFaceResult := range iFaceMethod.ResultList {
			mockResult := &domain.Field{
				Name:     iFaceResult.Name,
				WantName: iFaceResult.WantName,
				GotName:  iFaceResult.GotName,
				Type:     iFaceResult.Type,
				BaseType: iFaceResult.BaseType,
				NameType: iFaceResult.NameType,
			}

			if mockResult.GetTypeType() == domain.FieldTypeLocalCustomType {
				mockResult.Type = fmt.Sprintf("%s.%s", i.interfacePackage.SelfImport.Key, mockResult.Type)
				mockResult.CodeImportList = append(mockResult.CodeImportList, i.interfacePackage.SelfImport)
				//mockResult.TestImportList = append(mockResult.TestImportList, i.interfacePackage.SelfImport)
			}

			mockMethod.ResultList = append(mockMethod.ResultList, mockResult)
		}

		for _, iFaceArgNameType := range iFaceMethod.ArgNameTypeList {
			mockMethod.ArgNameTypeList = append(mockMethod.ArgNameTypeList, iFaceArgNameType)
		}

		for _, iFaceResultNameType := range iFaceMethod.ResultNameTypeList {
			mockMethod.ResultNameTypeList = append(mockMethod.ResultNameTypeList, iFaceResultNameType)
		}

		mockMethodList = append(mockMethodList, mockMethod)
	}

	return
}

func (i *Interactor) createMockFromInterface(iFace *domain.Interface, mockPackageName string) (mock *domain.Mock) {
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
		MethodList: i.createMockMethodList(iFace.MethodList),
	}

	argList := []*domain.Field{}
	resultList := []*domain.Field{}
	for _, method := range mock.MethodList {
		for _, arg := range method.ArgList {
			arg.WantName = "want" + toPublic(arg.Name)
			arg.GotName = "got" + toPublic(arg.Name)
			arg.NameType = arg.Name + " " + arg.Type
			arg.ExampleValue, arg.TestImportList = i.createExampleValue(arg)

			method.ArgNameTypeList = append(method.ArgNameTypeList, arg.NameType)
			method.CodeImportList = append(method.CodeImportList, arg.CodeImportList...) // TODO: add unique
			method.TestImportList = append(method.TestImportList, arg.TestImportList...) // TODO: add unique
		}
		for _, result := range method.ResultList {
			result.WantName = "want" + toPublic(result.Name)
			result.GotName = "got" + toPublic(result.Name)
			result.NameType = result.Name + " " + result.Type
			result.ExampleValue, result.TestImportList = i.createExampleValue(result)

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
		importList = append(importList, &domain.Import{Key: "fmt", Path: "fmt"})
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
			if Import.Key == importKey {
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
