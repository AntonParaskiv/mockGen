package Printer

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
)

func generateSetter(mock *domain.Mock, setter *domain.Setter) (code string) {
	code += fmt.Sprintf("func (%s *%s) %s(%s %s) *%s{\n", mock.Struct.GetReceiverName(), mock.Struct.Name, setter.Name, setter.Field.Name, setter.Field.Type, mock.Struct.Name)
	code += fmt.Sprintf("	%s.%s = %s\n", mock.Struct.GetReceiverName(), setter.Field.Name, setter.Field.Name)
	code += fmt.Sprintf("	return %s\n", mock.Struct.GetReceiverName())
	code += fmt.Sprintf("}\n\n")
	return
}

func generateSetterTest(mock *domain.Mock, setter *domain.Setter) (code string) {

	code += fmt.Sprintf("func Test%s_%s(t *testing.T) {\n", mock.Struct.Name, setter.Name)

	code += fmt.Sprintf("	type args struct {\n")
	code += fmt.Sprintf("		%s %s\n", setter.Field.Name, setter.Field.GetTypeViewStructField())
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("	tests := []struct {\n")
	code += fmt.Sprintf("		name string\n")
	code += fmt.Sprintf("		args args\n")
	code += fmt.Sprintf("		%s	*%s\n", mock.Struct.GetWantName(), mock.Struct.Name)
	code += fmt.Sprintf("	}{\n")

	code += generateSetterTestCase(mock, setter)

	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("	for _, tt := range tests {\n")
	code += fmt.Sprintf("		t.Run(tt.name, func(t *testing.T) {\n")
	code += fmt.Sprintf("			%s := &%s{}\n", mock.Struct.GetReceiverName(), mock.Struct.Name)
	code += fmt.Sprintf("			%s := %s.%s(tt.args.%s)\n", mock.Struct.GetGotName(), mock.Struct.GetReceiverName(), setter.Name, setter.Field.GetNameViewArg())
	code += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", mock.Struct.GetGotName(), mock.Struct.GetWantName())
	code += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", setter.Name, mock.Struct.GetGotName(), mock.Struct.GetWantName())
	code += fmt.Sprintf("			}\n")
	code += fmt.Sprintf("		})\n")
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("}\n\n")
	return
}

func generateSetterTestCase(mock *domain.Mock, setter *domain.Setter) (code string) {
	if len(setter.Field.ExampleValue) == 0 {
		err := fmt.Errorf("no example value for %s", setter.Field.GetNameType())
		err = fmt.Errorf("create setter %s.%s() test case failed: %w", mock.Struct.Name, setter.Name, err)
		fmt.Println(err)
		code += "		// TODO: Add test cases.\n"
		return
	}
	code += fmt.Sprintf("		{\n")
	code += fmt.Sprintf("			name: \"Setting\",\n")
	code += fmt.Sprintf("			args: args{\n")
	code += fmt.Sprintf("			%s: %s,\n", setter.Field.Name, setter.Field.ExampleValue)
	code += fmt.Sprintf("			},\n")
	code += fmt.Sprintf("			%s: &%s{\n", mock.Struct.GetWantName(), mock.Struct.Name)
	code += fmt.Sprintf("			%s: %s,\n", setter.Field.Name, setter.Field.ExampleValue)
	code += fmt.Sprintf("			},\n")
	code += fmt.Sprintf("		},\n")
	return
}
