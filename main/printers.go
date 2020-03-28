package main

import (
	"fmt"
	"strings"
)

type Mock struct {
	Name         string
	ReceiverName string
	WantName     string
	GotName      string
	FieldList    []*Field
	Constructor  *Constructor
	SetterList   []*Setter
	MethodList   []*Method
}

type Field struct {
	Name         string
	WantName     string
	GotName      string
	Type         string
	NameType     string
	ExampleValue string
}

type Constructor struct {
	Name string
}

type Setter struct {
	Name  string
	Field *Field
}

type Method struct {
	Name               string
	ArgList            []*Field
	ResultList         []*Field
	ArgNameTypeList    []string
	ResultNameTypeList []string
}

func PrintMock(mock *Mock) (result string) {
	result += PrintStruct(mock)
	result += PrintConstructor(mock)
	resultTest := PrintConstructorTest(mock)

	for _, setter := range mock.SetterList {
		result += PrintSetters(mock, setter)
		resultTest += PrintSetterTest(mock, setter)
	}

	for _, method := range mock.MethodList {
		result += PrintMethod(mock, method)
		resultTest += PrintMethodTest(mock, method)
	}

	result += resultTest
	return
}

func PrintStruct(mock *Mock) (result string) {
	result += fmt.Sprintf("type %s struct {\n", mock.Name)
	for _, field := range mock.FieldList {
		result += fmt.Sprintf("	%s %s\n", field.Name, field.Type)
	}
	result += fmt.Sprintf("}\n\n")
	return
}

func PrintConstructor(mock *Mock) (result string) {
	result += fmt.Sprintf("func %s() (%s *%s) {\n", mock.Constructor.Name, mock.ReceiverName, mock.Name)
	result += fmt.Sprintf("	%s = new(%s)\n", mock.ReceiverName, mock.Name)
	result += fmt.Sprintf("	return %s\n", mock.ReceiverName)
	result += fmt.Sprintf("}\n\n")
	return
}

func PrintConstructorTest(mock *Mock) (result string) {
	result += fmt.Sprintf("func Test%s(t *testing.T) {\n", mock.Constructor.Name)

	result += fmt.Sprintf("	tests := []struct {\n")
	result += fmt.Sprintf("		name string\n")
	result += fmt.Sprintf("		%s *%s\n", mock.WantName, mock.Name)
	result += fmt.Sprintf("	}{\n")
	result += fmt.Sprintf("		{\n")
	result += fmt.Sprintf("			name: \"Struct init\",\n")
	result += fmt.Sprintf("			%s: &%s{},\n", mock.WantName, mock.Name)
	result += fmt.Sprintf("		},\n")
	result += fmt.Sprintf("	}\n")

	result += fmt.Sprintf("	for _, tt := range tests {\n")
	result += fmt.Sprintf("		t.Run(tt.name, func(t *testing.T) {\n")
	result += fmt.Sprintf("			%s := %s()\n", mock.GotName, mock.Constructor.Name)
	result += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", mock.GotName, mock.WantName)
	result += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", mock.Constructor.Name, mock.GotName, mock.WantName)
	result += fmt.Sprintf("			}\n")
	result += fmt.Sprintf("		})\n")
	result += fmt.Sprintf("	}\n")

	result += fmt.Sprintf("}\n\n")

	return
}

func PrintSetters(mock *Mock, setter *Setter) (result string) {
	result += fmt.Sprintf("func (%s *%s) %s(%s %s) *%s{\n", mock.ReceiverName, mock.Name, setter.Name, setter.Field.Name, setter.Field.Type, mock.Name)
	result += fmt.Sprintf("	%s.%s = %s\n", mock.ReceiverName, setter.Field.Name, setter.Field.Name)
	result += fmt.Sprintf("	return %s\n", mock.ReceiverName)
	result += fmt.Sprintf("}\n\n")
	return
}

func PrintSetterTest(mock *Mock, setter *Setter) (result string) {
	result += fmt.Sprintf("func Test%s_%s(t *testing.T) {\n", mock.Name, setter.Name)

	result += fmt.Sprintf("	type args struct {\n")
	result += fmt.Sprintf("		%s %s\n", setter.Field.Name, setter.Field.Type)
	result += fmt.Sprintf("	}\n")

	result += fmt.Sprintf("	tests := []struct {\n")
	result += fmt.Sprintf("		name string\n")
	result += fmt.Sprintf("		args args\n")
	result += fmt.Sprintf("		%s	*%s\n", mock.WantName, mock.Name)
	result += fmt.Sprintf("	}{\n")
	result += fmt.Sprintf("		{\n")
	result += fmt.Sprintf("			name: \"Setting\",\n")
	result += fmt.Sprintf("			args: args{\n")
	result += fmt.Sprintf("			%s: %s,\n", setter.Field.Name, setter.Field.ExampleValue)
	result += fmt.Sprintf("			},\n")
	result += fmt.Sprintf("			%s: &%s{\n", mock.WantName, mock.Name)
	result += fmt.Sprintf("			%s: %s,\n", setter.Field.Name, setter.Field.ExampleValue)
	result += fmt.Sprintf("			},\n")
	result += fmt.Sprintf("		},\n")
	result += fmt.Sprintf("	}\n")

	result += fmt.Sprintf("	for _, tt := range tests {\n")
	result += fmt.Sprintf("		t.Run(tt.name, func(t *testing.T) {\n")
	result += fmt.Sprintf("			%s := &%s{}\n", mock.ReceiverName, mock.Name)
	result += fmt.Sprintf("			%s := %s.%s(tt.args.%s)\n", mock.GotName, mock.ReceiverName, setter.Name, setter.Field.Name)
	result += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", mock.GotName, mock.WantName)
	result += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", setter.Name, mock.GotName, mock.WantName)
	result += fmt.Sprintf("			}\n")
	result += fmt.Sprintf("		})\n")
	result += fmt.Sprintf("	}\n")

	result += fmt.Sprintf("}\n\n")

	return
}

func PrintMethod(mock *Mock, method *Method) (output string) {
	argLine := strings.Join(method.ArgNameTypeList, ", ")
	resultLine := strings.Join(method.ResultNameTypeList, ", ")

	output += fmt.Sprintf("func (%s *%s) %s(%s) (%s) {\n", mock.ReceiverName, mock.Name, method.Name, argLine, resultLine)
	for _, arg := range method.ArgList {
		output += fmt.Sprintf("	%s.%s = %s\n", mock.ReceiverName, arg.Name, arg.Name)
	}

	for _, result := range method.ResultList {
		output += fmt.Sprintf("	%s = %s.%s\n", result.Name, mock.ReceiverName, result.Name)
	}
	output += fmt.Sprintf("	return\n")
	output += fmt.Sprintf("}\n\n")
	return
}

func PrintMethodTest(mock *Mock, method *Method) (output string) {
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

	output += fmt.Sprintf("func Test%s_%s(t *testing.T) {\n", mock.Name, method.Name)

	output += fmt.Sprintf("	type fields struct {\n")
	for _, result := range method.ResultList {
		output += fmt.Sprintf("		%s %s\n", result.Name, result.Type)
	}
	output += fmt.Sprintf("	}\n")

	output += fmt.Sprintf("	type args struct {\n")
	for _, arg := range method.ArgList {
		output += fmt.Sprintf("		%s %s\n", arg.Name, arg.Type)
	}
	output += fmt.Sprintf("	}\n")

	output += fmt.Sprintf("	tests := []struct {\n")
	output += fmt.Sprintf("		name string\n")
	output += fmt.Sprintf("		fields fields\n")
	output += fmt.Sprintf("		args args\n")
	for _, result := range method.ResultList {
		output += fmt.Sprintf("		%s %s\n", result.WantName, result.Type)
	}
	output += fmt.Sprintf("		%s	*%s\n", mock.WantName, mock.Name)
	output += fmt.Sprintf("	}{\n")
	output += fmt.Sprintf("		{\n")
	output += fmt.Sprintf("			name: \"Success\",\n")

	output += fmt.Sprintf("			fields: fields{\n")
	for _, result := range method.ResultList {
		output += fmt.Sprintf("		%s: %s,\n", result.Name, result.ExampleValue)
	}
	output += fmt.Sprintf("			},\n")

	output += fmt.Sprintf("			args: args{\n")
	for _, arg := range method.ArgList {
		output += fmt.Sprintf("		%s: %s,\n", arg.Name, arg.ExampleValue)
	}
	output += fmt.Sprintf("			},\n")

	for _, result := range method.ResultList {
		output += fmt.Sprintf("		%s: %s,\n", result.WantName, result.ExampleValue)
	}

	output += fmt.Sprintf("			%s: &%s{\n", mock.WantName, mock.Name)
	for _, arg := range method.ArgList {
		output += fmt.Sprintf("			%s: %s,\n", arg.Name, arg.ExampleValue)
	}
	for _, result := range method.ResultList {
		output += fmt.Sprintf("		%s: %s,\n", result.Name, result.ExampleValue)
	}
	output += fmt.Sprintf("			},\n")

	output += fmt.Sprintf("		},\n")
	output += fmt.Sprintf("	}\n")

	output += fmt.Sprintf("	for _, tt := range tests {\n")
	output += fmt.Sprintf("		t.Run(tt.name, func(t *testing.T) {\n")
	output += fmt.Sprintf("			%s := &%s{\n", mock.ReceiverName, mock.Name)
	for _, result := range method.ResultList {
		output += fmt.Sprintf("		%s: tt.fields.%s,\n", result.Name, result.Name)
	}
	output += fmt.Sprintf("			}\n")
	output += fmt.Sprintf("			%s := %s.%s(%s)\n", gotResultsLine, mock.ReceiverName, method.Name, ttArgsLine)

	for _, result := range method.ResultList {
		output += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", result.GotName, result.WantName)
		output += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", result.GotName, result.GotName, result.WantName)
		output += fmt.Sprintf("			}\n")
	}

	output += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", mock.ReceiverName, mock.WantName)
	output += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", mock.Name, mock.ReceiverName, mock.WantName)
	output += fmt.Sprintf("			}\n")

	output += fmt.Sprintf("		})\n")
	output += fmt.Sprintf("	}\n")

	output += fmt.Sprintf("}\n\n")

	return
}
