package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type GoCodePackage struct {
	Path        string
	PackageName string
	FileList    []*GoCodeFile
}

type GoCodeFile struct {
	Name          string
	ImportList    []string
	InterfaceList []*Interface
	MockList      []*Mock
	Code          string
}

type Interface struct {
	Name       string
	MethodList []*Method
}

type Mock struct {
	Struct      *Struct
	Constructor *Constructor
	SetterList  []*Setter
	MethodList  []*Method
}

type Struct struct {
	Name         string
	ReceiverName string
	WantName     string
	GotName      string
	FieldList    []*Field
	Code         string
}

type Constructor struct {
	Name     string
	Code     string
	CodeTest string
}

type Setter struct {
	Name     string
	Field    *Field
	Code     string
	CodeTest string
}

type Method struct {
	Name               string
	ArgList            []*Field
	ResultList         []*Field
	ArgNameTypeList    []string
	ResultNameTypeList []string
	Code               string
	CodeTest           string
}

type Field struct {
	Name         string
	WantName     string
	GotName      string
	Type         string
	NameType     string
	ExampleValue string
}

func CreateMockPackage(interfacePackage *GoCodePackage) (mockPackage *GoCodePackage) {
	mockPackagePath := cutPostfix(interfacePackage.Path, "Interface") + "Mock"
	mockPackageName := cutPostfix(interfacePackage.PackageName, "Interface") + "Mock"

	mockPackage = &GoCodePackage{
		Path:        mockPackagePath,
		PackageName: mockPackageName,
	}

	for _, interfaceFile := range interfacePackage.FileList {
		mockFile, mockTestFile := CreateMockFilesFromInterfaceFile(interfaceFile, mockPackageName)
		mockPackage.FileList = append(mockPackage.FileList, mockFile, mockTestFile)
	}

	return
}

func CreateMockFilesFromInterfaceFile(interfaceFile *GoCodeFile, mockPackageName string) (mockFile *GoCodeFile, mockTestFile *GoCodeFile) {
	var mockCode string
	var mockTestCode string

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
	}

	if len(mockCode) == 0 {
		return
	}

	mockPackageCode := fmt.Sprintf("package %s\n\n", mockPackageName)

	mockFile = &GoCodeFile{
		Name:       interfaceFile.Name,
		ImportList: interfaceFile.ImportList,
	}

	mockFile.Code = mockPackageCode
	mockFile.Code += CreateImportList(mockFile.ImportList)
	mockFile.Code += mockCode

	if len(mockTestCode) == 0 {
		return
	}

	mockTestFileName := createTestFilePath(mockFile.Name)
	mockTestFileImports := interfaceFile.ImportList // TODO: replace with deep copy
	mockTestFileImports = append(mockTestFileImports, "reflect", "testing")

	mockTestFile = &GoCodeFile{
		Name:       mockTestFileName,
		ImportList: mockTestFileImports,
	}

	mockTestFile.Code = mockPackageCode
	mockTestFile.Code += CreateImportList(mockTestFile.ImportList)
	mockTestFile.Code += mockTestCode
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

func CreateMockFromInterface(iFace *Interface, mockPackageName string) (mock *Mock) {
	structName := iFace.Name

	basePackageName := cutPostfix(mockPackageName, "Mock")

	var constructorName string
	if structName == basePackageName {
		constructorName = "New"
	} else {
		constructorName = "New" + toPublic(structName)
	}

	mock = &Mock{
		Struct: &Struct{
			Name:         structName,
			ReceiverName: getReceiverName(structName),
			WantName:     "want" + toPublic(structName),
			GotName:      "got" + toPublic(structName),
		},
		Constructor: &Constructor{
			Name: constructorName,
		},
		MethodList: iFace.MethodList,
	}

	argList := []*Field{}
	resultList := []*Field{}
	for _, method := range mock.MethodList {
		for _, arg := range method.ArgList {
			arg.WantName = "want" + toPublic(arg.Name)
			arg.GotName = "got" + toPublic(arg.Name)
			arg.NameType = arg.Name + " " + arg.Type

			switch {
			case arg.Type == "string":
				arg.ExampleValue = `"my` + toPublic(arg.Name) + `"`
			case arg.Type == "bool":
				arg.ExampleValue = "true"
			case arg.Type == "rune":
				arg.ExampleValue = `"X"`
			case arg.Type == "byte":
				arg.ExampleValue = `50`
			case len(arg.Type) >= 3 && arg.Type[0:3] == "int":
				arg.ExampleValue = "100"
			case len(arg.Type) >= 4 && arg.Type[0:4] == "uint":
				arg.ExampleValue = "200"
			case len(arg.Type) >= 5 && arg.Type[0:5] == "float":
				arg.ExampleValue = "3.14"
			}

			method.ArgNameTypeList = append(method.ArgNameTypeList, arg.NameType)
		}
		for _, result := range method.ResultList {
			result.WantName = "want" + toPublic(result.Name)
			result.GotName = "got" + toPublic(result.Name)
			result.NameType = result.Name + " " + result.Type

			switch {
			case result.Type == "string":
				result.ExampleValue = `"my` + toPublic(result.Name) + `"`
			case result.Type == "bool":
				result.ExampleValue = "true"
			case result.Type == "rune":
				result.ExampleValue = `"X"`
			case result.Type == "byte":
				result.ExampleValue = `50`
			case len(result.Type) >= 3 && result.Type[0:3] == "int":
				result.ExampleValue = "100"
			case len(result.Type) >= 4 && result.Type[0:4] == "uint":
				result.ExampleValue = "200"
			case len(result.Type) >= 5 && result.Type[0:5] == "float":
				result.ExampleValue = "3.14"
			}
			method.ResultNameTypeList = append(method.ResultNameTypeList, result.NameType)
		}

		argList = append(argList, method.ArgList...)
		resultList = append(resultList, method.ResultList...)

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
		mock.SetterList = append(mock.SetterList, &Setter{
			Name:  "Set" + toPublic(field.Name),
			Field: field,
		})
	}
	return
}

func GenCodeMock(mock *Mock) {

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

func GenCodeStruct(mock *Mock) {
	var code string
	code += fmt.Sprintf("type %s struct {\n", mock.Struct.Name)
	for _, field := range mock.Struct.FieldList {
		code += fmt.Sprintf("	%s %s\n", field.Name, field.Type)
	}
	code += fmt.Sprintf("}\n\n")
	mock.Struct.Code = code
	return
}

func GenCodeConstructor(mock *Mock) {
	var code string
	code += fmt.Sprintf("func %s() (%s *%s) {\n", mock.Constructor.Name, mock.Struct.ReceiverName, mock.Struct.Name)
	code += fmt.Sprintf("	%s = new(%s)\n", mock.Struct.ReceiverName, mock.Struct.Name)
	code += fmt.Sprintf("	return %s\n", mock.Struct.ReceiverName)
	code += fmt.Sprintf("}\n\n")
	mock.Constructor.Code = code
	return
}

func GenCodeTestConstructor(mock *Mock) {
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

func GenCodeSetter(mock *Mock, setter *Setter) {
	var code string
	code += fmt.Sprintf("func (%s *%s) %s(%s %s) *%s{\n", mock.Struct.ReceiverName, mock.Struct.Name, setter.Name, setter.Field.Name, setter.Field.Type, mock.Struct.Name)
	code += fmt.Sprintf("	%s.%s = %s\n", mock.Struct.ReceiverName, setter.Field.Name, setter.Field.Name)
	code += fmt.Sprintf("	return %s\n", mock.Struct.ReceiverName)
	code += fmt.Sprintf("}\n\n")
	setter.Code = code
	return
}

func GenCodeTestSetter(mock *Mock, setter *Setter) {
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

func GenCodeMethod(mock *Mock, method *Method) {
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

func GenCodeTestMethod(mock *Mock, method *Method) {
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

func SaveGoPackage(Package *GoCodePackage) (err error) {
	err = os.MkdirAll(Package.Path, 0755)
	if err != nil {
		err = fmt.Errorf("create dir %s failed: %w", Package.Path, err)
		return
	}

	for _, file := range Package.FileList {
		filePath := filepath.Join(Package.Path, file.Name)
		err = ioutil.WriteFile(filePath, []byte(file.Code), 0644)
		if err != nil {
			err = fmt.Errorf("write file %s failed: %w", file.Name, err)
			return
		}
	}
	return
}

func createMockPackagePath(interfacePackagePath string) (mockFilePath string) {

	return
}
