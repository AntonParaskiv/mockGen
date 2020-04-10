package Printer

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
)

func generateConstructor(mock *domain.Mock) (code string) {
	code += fmt.Sprintf("func %s() (%s *%s) {\n", mock.Constructor.Name, mock.Struct.GetReceiverName(), mock.Struct.Name)
	code += fmt.Sprintf("	%s = new(%s)\n", mock.Struct.GetReceiverName(), mock.Struct.Name)
	code += fmt.Sprintf("	return %s\n", mock.Struct.GetReceiverName())
	code += fmt.Sprintf("}\n\n")
	return
}

func generateConstructorTest(mock *domain.Mock) (code string) {
	code += fmt.Sprintf("func Test%s(t *testing.T) {\n", mock.Constructor.Name)

	code += fmt.Sprintf("	tests := []struct {\n")
	code += fmt.Sprintf("		name string\n")
	code += fmt.Sprintf("		%s *%s\n", mock.Struct.GetWantName(), mock.Struct.Name)
	code += fmt.Sprintf("	}{\n")
	code += fmt.Sprintf("		{\n")
	code += fmt.Sprintf("			name: \"Struct init\",\n")
	code += fmt.Sprintf("			%s: &%s{},\n", mock.Struct.GetWantName(), mock.Struct.Name)
	code += fmt.Sprintf("		},\n")
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("	for _, tt := range tests {\n")
	code += fmt.Sprintf("		t.Run(tt.name, func(t *testing.T) {\n")
	code += fmt.Sprintf("			%s := %s()\n", mock.Struct.GetGotName(), mock.Constructor.Name)
	code += fmt.Sprintf("			if !reflect.DeepEqual(%s, tt.%s) {\n", mock.Struct.GetGotName(), mock.Struct.GetWantName())
	code += fmt.Sprintf("				t.Errorf(\"%s() = %%v, want %%v\", %s, tt.%s)\n", mock.Constructor.Name, mock.Struct.GetGotName(), mock.Struct.GetWantName())
	code += fmt.Sprintf("			}\n")
	code += fmt.Sprintf("		})\n")
	code += fmt.Sprintf("	}\n")

	code += fmt.Sprintf("}\n\n")
	return
}
