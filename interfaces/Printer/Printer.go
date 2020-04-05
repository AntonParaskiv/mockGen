package Printer

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"path/filepath"
	"strings"
)

type Printer struct {
}

func (p *Printer) GenerateCode(mockPackage *domain.GoCodePackage) {
	for _, mockFile := range mockPackage.FileList {

		var mockCode string
		var mockTestCode string

		mockTestFile := &domain.GoCodeFile{
			Name: createTestFilePath(mockFile.Name),
			ImportList: []*domain.Import{
				{Key: "reflect", Path: "reflect"},
				{Key: "testing", Path: "testing"},
			},
		}

		for _, mock := range mockFile.MockList {
			genCodeMock(mock)
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

			mockFile.ImportList = append(mockFile.ImportList, mock.CodeImportList...)         // TODO: add unique
			mockTestFile.ImportList = append(mockTestFile.ImportList, mock.TestImportList...) // TODO: add unique
			mockFile.MockList = append(mockFile.MockList, mock)
		}

		if len(mockCode) == 0 {
			continue
		}

		mockPackageCode := fmt.Sprintf("package %s\n\n", mockPackage.PackageName)

		mockFile.Code = mockPackageCode
		mockFile.Code += createImportList(mockFile.ImportList)
		mockFile.Code += mockCode

		if len(mockTestCode) == 0 {
			continue
		}

		mockTestFile.Code = mockPackageCode
		mockTestFile.Code += createImportList(mockTestFile.ImportList)
		mockTestFile.Code += mockTestCode

		mockPackage.FileList = append(mockPackage.FileList, mockTestFile)
	}

	return
}

func genCodeMock(mock *domain.Mock) {

	genCodeStruct(mock)
	genCodeConstructor(mock)
	genCodeTestConstructor(mock)

	for _, setter := range mock.SetterList {
		genCodeSetter(mock, setter)
		genCodeTestSetter(mock, setter)
	}

	for _, method := range mock.MethodList {
		genCodeMethod(mock, method)
		genCodeTestMethod(mock, method)
	}
	return
}

func genCodeStruct(mock *domain.Mock) {
	var code string
	code += fmt.Sprintf("type %s struct {\n", mock.Struct.Name)
	for _, field := range mock.Struct.FieldList {
		code += fmt.Sprintf("	%s %s\n", field.Name, field.Type)
	}
	code += fmt.Sprintf("}\n\n")
	mock.Struct.Code = code
	return
}

func genCodeConstructor(mock *domain.Mock) {
	var code string
	code += fmt.Sprintf("func %s() (%s *%s) {\n", mock.Constructor.Name, mock.Struct.ReceiverName, mock.Struct.Name)
	code += fmt.Sprintf("	%s = new(%s)\n", mock.Struct.ReceiverName, mock.Struct.Name)
	code += fmt.Sprintf("	return %s\n", mock.Struct.ReceiverName)
	code += fmt.Sprintf("}\n\n")
	mock.Constructor.Code = code
	return
}

func genCodeTestConstructor(mock *domain.Mock) {
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

func genCodeSetter(mock *domain.Mock, setter *domain.Setter) {
	var code string
	code += fmt.Sprintf("func (%s *%s) %s(%s %s) *%s{\n", mock.Struct.ReceiverName, mock.Struct.Name, setter.Name, setter.Field.Name, setter.Field.Type, mock.Struct.Name)
	code += fmt.Sprintf("	%s.%s = %s\n", mock.Struct.ReceiverName, setter.Field.Name, setter.Field.Name)
	code += fmt.Sprintf("	return %s\n", mock.Struct.ReceiverName)
	code += fmt.Sprintf("}\n\n")
	setter.Code = code
	return
}

func genCodeTestSetter(mock *domain.Mock, setter *domain.Setter) {
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

func genCodeMethod(mock *domain.Mock, method *domain.Method) {
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

func genCodeTestMethod(mock *domain.Mock, method *domain.Method) {
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

func createImportList(importList []*domain.Import) (code string) {
	switch len(importList) {
	case 0:
	case 1:
		code = fmt.Sprintf("import %s \"%s\"\n\n", importList[0].Name, importList[0].Path)
	default:
		code = fmt.Sprintf("import (\n")
		for _, Import := range importList {
			code += fmt.Sprintf("	%s \"%s\"\n", Import.Name, Import.Path)
		}
		code += fmt.Sprintf(")\n\n")
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
