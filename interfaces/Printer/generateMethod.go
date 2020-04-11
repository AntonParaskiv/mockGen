package Printer

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"strings"
)

func generateMethod(mock *domain.Mock, method *domain.Method) (code string) {
	code += fmt.Sprintf("func (%s *%s) %s(%s) (%s) {\n", mock.Struct.GetReceiverName(), mock.Struct.Name, method.Name, method.GetArgLine(), method.GetResultLine())
	for _, arg := range method.ArgList {
		code += fmt.Sprintf("	%s.%s = %s\n", mock.Struct.GetReceiverName(), arg.Name, arg.Name)
	}
	for _, result := range method.ResultList {
		code += fmt.Sprintf("	%s = %s.%s\n", result.Name, mock.Struct.GetReceiverName(), result.Name)
	}
	code += fmt.Sprintf("	return\n")
	code += fmt.Sprintf("}\n\n")
	return
}

func generateMethodTest(mock *domain.Mock, method *domain.Method) (code string) {

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
		code += fmt.Sprintf("		%s %s\n", result.GetWantName(), result.Type)
	}
	code += fmt.Sprintf("		%s	*%s\n", mock.Struct.GetWantName(), mock.Struct.Name)
	code += fmt.Sprintf("	}{\n")

	code += generateMethodTestCase(mock, method)

	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("	for _, tt := range tests {\n")
	code += fmt.Sprintf("		t.Run(tt.name, func(t *testing.T) {\n")
	code += fmt.Sprintf("			%s := &%s{\n", mock.Struct.GetReceiverName(), mock.Struct.Name)
	for _, result := range method.ResultList {
		code += fmt.Sprintf("		%s: tt.fields.%s,\n", result.Name, result.Name)
	}
	code += fmt.Sprintf("			}\n")
	code += fmt.Sprintf("			%s := %s.%s(%s)\n", createGotResultLine(method.ResultList), mock.Struct.GetReceiverName(), method.Name, createTtArgLine(method.ArgList))

	for _, result := range method.ResultList {
		code += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", result.GetGotName(), result.GetWantName())
		code += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", result.GetGotName(), result.GetGotName(), result.GetWantName())
		code += fmt.Sprintf("			}\n")
	}

	code += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", mock.Struct.GetReceiverName(), mock.Struct.GetWantName())
	code += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", mock.Struct.Name, mock.Struct.GetReceiverName(), mock.Struct.GetWantName())
	code += fmt.Sprintf("			}\n")

	code += fmt.Sprintf("		})\n")
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("}\n\n")
	return
}

func generateMethodTestCase(mock *domain.Mock, method *domain.Method) (code string) {
	var errs []error
	for _, field := range method.ArgList {
		if len(field.ExampleValue) == 0 {
			errs = append(errs, fmt.Errorf("no example value for arg %s", field.GetNameType()))
		}
	}
	for _, field := range method.ResultList {
		if len(field.ExampleValue) == 0 {
			errs = append(errs, fmt.Errorf("no example value for result %s", field.GetNameType()))
		}
	}
	if len(errs) > 0 {
		for _, err := range errs {
			err = fmt.Errorf("create method %s.%s() test case failed: %w", mock.Struct.Name, method.Name, err)
			fmt.Println(err)
		}
		code += "		// TODO: Add test cases.\n"
		return
	}

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
		code += fmt.Sprintf("		%s: %s,\n", result.GetWantName(), result.ExampleValue)
	}
	code += fmt.Sprintf("			%s: &%s{\n", mock.Struct.GetWantName(), mock.Struct.Name)
	for _, arg := range method.ArgList {
		code += fmt.Sprintf("			%s: %s,\n", arg.Name, arg.ExampleValue)
	}
	for _, result := range method.ResultList {
		code += fmt.Sprintf("		%s: %s,\n", result.Name, result.ExampleValue)
	}
	code += fmt.Sprintf("			},\n")
	code += fmt.Sprintf("		},\n")
	return
}

func createTtArgLine(argList []*domain.Field) (ttArgLine string) {
	ttArgList := []string{}
	for _, arg := range argList {
		ttArgList = append(ttArgList, "tt.args."+arg.Name)
	}
	ttArgLine = strings.Join(ttArgList, ", ")
	return
}

func createGotResultLine(resultList []*domain.Field) (gotResultLine string) {
	gotResultList := []string{}
	for _, result := range resultList {
		gotResultList = append(gotResultList, result.GetGotName())
	}
	gotResultLine = strings.Join(gotResultList, ", ")
	return
}
