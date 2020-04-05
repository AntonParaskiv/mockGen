package usecases

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"path/filepath"
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
		mockFile, mockTestFile := CreateMockFilesFromInterfaceFile(interfaceFile, mockPackageName)
		mockPackage.FileList = append(mockPackage.FileList, mockFile, mockTestFile)
	}

	return
}

func CreateMockFilesFromInterfaceFile(interfaceFile *domain.GoCodeFile, mockPackageName string) (mockFile *domain.GoCodeFile, mockTestFile *domain.GoCodeFile) {
	var mockCode string
	var mockTestCode string
	var mockFileImportList []string
	var mockTestFileImportList []string
	mockFileImportList = append(mockFileImportList, interfaceFile.ImportList...)
	mockTestFileImportList = append(mockTestFileImportList, interfaceFile.ImportList...)

	for _, iFace := range interfaceFile.InterfaceList {
		mock := CreateMockFromInterface(iFace, mockPackageName)
		GenCodeMock(mock)
		mockCode += mock.Struct.Code
		mockCode += mock.Constructor.Code
		mockTestCode += mock.Constructor.CodeTest

		for _, setter := range mock.SetterList {
			mockCode += setter.Code
			mockTestCode += setter.CodeTest
		}

		for _, method := range mock.MethodList {
			mockCode += method.Code
			mockTestCode += method.CodeTest
		}

		mockFileImportList = append(mockFileImportList, mock.CodeImportList...)         // TODO: add unique
		mockTestFileImportList = append(mockTestFileImportList, mock.TestImportList...) // TODO: add unique
	}

	if len(mockCode) == 0 {
		return
	}

	mockPackageCode := fmt.Sprintf("package %s\n\n", mockPackageName)

	mockFile = &domain.GoCodeFile{
		Name:       interfaceFile.Name,
		ImportList: mockFileImportList,
	}

	mockFile.Code = mockPackageCode
	mockFile.Code += CreateImportList(mockFile.ImportList)
	mockFile.Code += mockCode

	if len(mockTestCode) == 0 {
		return
	}

	mockTestFileName := createTestFilePath(mockFile.Name)
	mockTestFileImports := mockTestFileImportList // TODO: replace with deep copy
	mockTestFileImports = append(mockTestFileImports, "reflect", "testing")

	mockTestFile = &domain.GoCodeFile{
		Name:       mockTestFileName,
		ImportList: mockTestFileImports,
	}

	mockTestFile.Code = mockPackageCode
	mockTestFile.Code += CreateImportList(mockTestFile.ImportList)
	mockTestFile.Code += mockTestCode
	return
}

func CreateMockFromInterface(iFace *domain.Interface, mockPackageName string) (mock *domain.Mock) {
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

func CreateImportList(importList []string) (code string) {
	switch len(importList) {
	case 0:
	case 1:
		code = fmt.Sprintf("import \"%s\"\n\n", importList[0])
	default:
		code = fmt.Sprintf("import (\n")
		for _, Import := range importList {
			code += fmt.Sprintf("	\"%s\"\n", Import)
		}
		code += fmt.Sprintf(")\n\n")
	}
	return
}

func createExampleValue(fieldType string, fieldName string) (exampleValue string, importList []string) {
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
		importList = append(importList, "fmt")
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

func GenCodeMock(mock *domain.Mock) {

	GenCodeStruct(mock)
	GenCodeConstructor(mock)
	GenCodeTestConstructor(mock)

	for _, setter := range mock.SetterList {
		GenCodeSetter(mock, setter)
		GenCodeTestSetter(mock, setter)
	}

	for _, method := range mock.MethodList {
		GenCodeMethod(mock, method)
		GenCodeTestMethod(mock, method)
	}
	return
}

func GenCodeStruct(mock *domain.Mock) {
	var code string
	code += fmt.Sprintf("type %s struct {\n", mock.Struct.Name)
	for _, field := range mock.Struct.FieldList {
		code += fmt.Sprintf("	%s %s\n", field.Name, field.Type)
	}
	code += fmt.Sprintf("}\n\n")
	mock.Struct.Code = code
	return
}

func GenCodeConstructor(mock *domain.Mock) {
	var code string
	code += fmt.Sprintf("func %s() (%s *%s) {\n", mock.Constructor.Name, mock.Struct.ReceiverName, mock.Struct.Name)
	code += fmt.Sprintf("	%s = new(%s)\n", mock.Struct.ReceiverName, mock.Struct.Name)
	code += fmt.Sprintf("	return %s\n", mock.Struct.ReceiverName)
	code += fmt.Sprintf("}\n\n")
	mock.Constructor.Code = code
	return
}

func GenCodeTestConstructor(mock *domain.Mock) {
	var code string
	code += fmt.Sprintf("func Test%s(t *testing.T) {\n", mock.Constructor.Name)

	code += fmt.Sprintf("	tests := []struct {\n")
	code += fmt.Sprintf("		name string\n")
	code += fmt.Sprintf("		%s *%s\n", mock.Struct.WantName, mock.Struct.Name)
	code += fmt.Sprintf("	}{\n")
	code += fmt.Sprintf("		{\n")
	code += fmt.Sprintf("			name: \"Struct init\",\n")
	code += fmt.Sprintf("			%s: &%s{},\n", mock.Struct.WantName, mock.Struct.Name)
	code += fmt.Sprintf("		},\n")
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("	for _, tt := range tests {\n")
	code += fmt.Sprintf("		t.Run(tt.name, func(t *testing.T) {\n")
	code += fmt.Sprintf("			%s := %s()\n", mock.Struct.GotName, mock.Constructor.Name)
	code += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", mock.Struct.GotName, mock.Struct.WantName)
	code += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", mock.Constructor.Name, mock.Struct.GotName, mock.Struct.WantName)
	code += fmt.Sprintf("			}\n")
	code += fmt.Sprintf("		})\n")
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("}\n\n")
	mock.Constructor.CodeTest = code
	return
}

func GenCodeSetter(mock *domain.Mock, setter *domain.Setter) {
	var code string
	code += fmt.Sprintf("func (%s *%s) %s(%s %s) *%s{\n", mock.Struct.ReceiverName, mock.Struct.Name, setter.Name, setter.Field.Name, setter.Field.Type, mock.Struct.Name)
	code += fmt.Sprintf("	%s.%s = %s\n", mock.Struct.ReceiverName, setter.Field.Name, setter.Field.Name)
	code += fmt.Sprintf("	return %s\n", mock.Struct.ReceiverName)
	code += fmt.Sprintf("}\n\n")
	setter.Code = code
	return
}

func GenCodeTestSetter(mock *domain.Mock, setter *domain.Setter) {
	var code string
	code += fmt.Sprintf("func Test%s_%s(t *testing.T) {\n", mock.Struct.Name, setter.Name)

	code += fmt.Sprintf("	type args struct {\n")
	code += fmt.Sprintf("		%s %s\n", setter.Field.Name, setter.Field.Type)
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("	tests := []struct {\n")
	code += fmt.Sprintf("		name string\n")
	code += fmt.Sprintf("		args args\n")
	code += fmt.Sprintf("		%s	*%s\n", mock.Struct.WantName, mock.Struct.Name)
	code += fmt.Sprintf("	}{\n")
	code += fmt.Sprintf("		{\n")
	code += fmt.Sprintf("			name: \"Setting\",\n")
	code += fmt.Sprintf("			args: args{\n")
	code += fmt.Sprintf("			%s: %s,\n", setter.Field.Name, setter.Field.ExampleValue)
	code += fmt.Sprintf("			},\n")
	code += fmt.Sprintf("			%s: &%s{\n", mock.Struct.WantName, mock.Struct.Name)
	code += fmt.Sprintf("			%s: %s,\n", setter.Field.Name, setter.Field.ExampleValue)
	code += fmt.Sprintf("			},\n")
	code += fmt.Sprintf("		},\n")
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("	for _, tt := range tests {\n")
	code += fmt.Sprintf("		t.Run(tt.name, func(t *testing.T) {\n")
	code += fmt.Sprintf("			%s := &%s{}\n", mock.Struct.ReceiverName, mock.Struct.Name)
	code += fmt.Sprintf("			%s := %s.%s(tt.args.%s)\n", mock.Struct.GotName, mock.Struct.ReceiverName, setter.Name, setter.Field.Name)
	code += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", mock.Struct.GotName, mock.Struct.WantName)
	code += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", setter.Name, mock.Struct.GotName, mock.Struct.WantName)
	code += fmt.Sprintf("			}\n")
	code += fmt.Sprintf("		})\n")
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("}\n\n")
	setter.CodeTest = code
	return
}

func GenCodeMethod(mock *domain.Mock, method *domain.Method) {
	var code string

	argLine := strings.Join(method.ArgNameTypeList, ", ")
	resultLine := strings.Join(method.ResultNameTypeList, ", ")

	code += fmt.Sprintf("func (%s *%s) %s(%s) (%s) {\n", mock.Struct.ReceiverName, mock.Struct.Name, method.Name, argLine, resultLine)
	for _, arg := range method.ArgList {
		code += fmt.Sprintf("	%s.%s = %s\n", mock.Struct.ReceiverName, arg.Name, arg.Name)
	}

	for _, result := range method.ResultList {
		code += fmt.Sprintf("	%s = %s.%s\n", result.Name, mock.Struct.ReceiverName, result.Name)
	}
	code += fmt.Sprintf("	return\n")
	code += fmt.Sprintf("}\n\n")

	method.Code = code
	return
}

func GenCodeTestMethod(mock *domain.Mock, method *domain.Method) {
	var code string

	//tt.args.nickName, tt.args.password
	ttArgList := []string{}
	for _, arg := range method.ArgList {
		ttArgList = append(ttArgList, "tt.args."+arg.Name)
	}
	ttArgsLine := strings.Join(ttArgList, ",")

	// gotAccountId, gotCheckCode
	gotResultList := []string{}
	for _, result := range method.ResultList {
		gotResultList = append(gotResultList, result.GotName)
	}
	gotResultsLine := strings.Join(gotResultList, ",")

	code += fmt.Sprintf("func Test%s_%s(t *testing.T) {\n", mock.Struct.Name, method.Name)

	code += fmt.Sprintf("	type fields struct {\n")
	for _, result := range method.ResultList {
		code += fmt.Sprintf("		%s %s\n", result.Name, result.Type)
	}
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("	type args struct {\n")
	for _, arg := range method.ArgList {
		code += fmt.Sprintf("		%s %s\n", arg.Name, arg.Type)
	}
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("	tests := []struct {\n")
	code += fmt.Sprintf("		name string\n")
	code += fmt.Sprintf("		fields fields\n")
	code += fmt.Sprintf("		args args\n")
	for _, result := range method.ResultList {
		code += fmt.Sprintf("		%s %s\n", result.WantName, result.Type)
	}
	code += fmt.Sprintf("		%s	*%s\n", mock.Struct.WantName, mock.Struct.Name)
	code += fmt.Sprintf("	}{\n")
	code += fmt.Sprintf("		{\n")
	code += fmt.Sprintf("			name: \"Success\",\n")

	code += fmt.Sprintf("			fields: fields{\n")
	for _, result := range method.ResultList {
		code += fmt.Sprintf("		%s: %s,\n", result.Name, result.ExampleValue)
	}
	code += fmt.Sprintf("			},\n")

	code += fmt.Sprintf("			args: args{\n")
	for _, arg := range method.ArgList {
		code += fmt.Sprintf("		%s: %s,\n", arg.Name, arg.ExampleValue)
	}
	code += fmt.Sprintf("			},\n")

	for _, result := range method.ResultList {
		code += fmt.Sprintf("		%s: %s,\n", result.WantName, result.ExampleValue)
	}

	code += fmt.Sprintf("			%s: &%s{\n", mock.Struct.WantName, mock.Struct.Name)
	for _, arg := range method.ArgList {
		code += fmt.Sprintf("			%s: %s,\n", arg.Name, arg.ExampleValue)
	}
	for _, result := range method.ResultList {
		code += fmt.Sprintf("		%s: %s,\n", result.Name, result.ExampleValue)
	}
	code += fmt.Sprintf("			},\n")

	code += fmt.Sprintf("		},\n")
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("	for _, tt := range tests {\n")
	code += fmt.Sprintf("		t.Run(tt.name, func(t *testing.T) {\n")
	code += fmt.Sprintf("			%s := &%s{\n", mock.Struct.ReceiverName, mock.Struct.Name)
	for _, result := range method.ResultList {
		code += fmt.Sprintf("		%s: tt.fields.%s,\n", result.Name, result.Name)
	}
	code += fmt.Sprintf("			}\n")
	code += fmt.Sprintf("			%s := %s.%s(%s)\n", gotResultsLine, mock.Struct.ReceiverName, method.Name, ttArgsLine)

	for _, result := range method.ResultList {
		code += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", result.GotName, result.WantName)
		code += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", result.GotName, result.GotName, result.WantName)
		code += fmt.Sprintf("			}\n")
	}

	code += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", mock.Struct.ReceiverName, mock.Struct.WantName)
	code += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", mock.Struct.Name, mock.Struct.ReceiverName, mock.Struct.WantName)
	code += fmt.Sprintf("			}\n")

	code += fmt.Sprintf("		})\n")
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("}\n\n")

	method.CodeTest = code
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

func createTestFilePath(filePath string) (testFilePath string) {
	extension := filepath.Ext(filePath)
	if extension == ".go" {
		filePathLen := len(filePath)
		testFilePath = filePath[:filePathLen-3] + "_test.go"
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