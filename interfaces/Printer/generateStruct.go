package Printer

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
)

func generateStruct(mock *domain.Mock) (code string) {
	code += fmt.Sprintf("type %s struct {\n", mock.Struct.Name)
	for _, field := range mock.Struct.FieldList {
		code += fmt.Sprintf("	%s %s\n", field.Name, field.Type)
	}
	code += fmt.Sprintf("}\n\n")
	return
}
