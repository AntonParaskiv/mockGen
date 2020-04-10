package domain

import "strings"

const (
	FieldTypeUnknown = iota
	FieldTypeString
	FieldTypeInterface
	FieldTypeInts
	FieldTypeFloat
	FieldTypeBool
	FieldTypeRune
	FieldTypeByte
	FieldTypeError
	FieldTypeArray
	FieldTypeMap
	FieldTypeImportedCustomType
	FieldTypeLocalCustomType
)

type Field struct {
	Name           string
	Type           string
	BaseType       *Field
	ExampleValue   string
	CodeImportList []*Import
	TestImportList []*Import
}

func (f *Field) GetTypeType() (typeType int64) {
	switch {
	case f.Type == "string":
		typeType = FieldTypeString
	case f.Type == "interface{}":
		typeType = FieldTypeInterface
	case len(f.Type) >= 3 && f.Type[0:3] == "int": // int must be after interface !
		typeType = FieldTypeInts
	case len(f.Type) >= 4 && f.Type[0:4] == "uint":
		typeType = FieldTypeInts
	case len(f.Type) >= 5 && f.Type[0:5] == "float":
		typeType = FieldTypeFloat
	case f.Type == "bool":
		typeType = FieldTypeBool
	case f.Type == "rune":
		typeType = FieldTypeRune
	case f.Type == "byte":
		typeType = FieldTypeByte
	case f.Type == "error":
		typeType = FieldTypeError
	case len(f.Type) >= 2 && f.Type[0:2] == "[]":
		typeType = FieldTypeArray
	case len(f.Type) >= 4 && f.Type[0:4] == "map[":
		typeType = FieldTypeMap
	default:
		switch len(strings.Split(f.Type, ".")) {
		case 1:
			typeType = FieldTypeLocalCustomType
		case 2:
			typeType = FieldTypeImportedCustomType
		default:
			typeType = FieldTypeUnknown
		}
	}
	return
}

func (f *Field) GetTypePackage() (typePackage string) {
	separatorIndex := strings.Index(f.Type, ".")
	if separatorIndex < 0 {
		return
	}
	typePackage = f.Type[0:separatorIndex]
	return
}

func (f *Field) GetTypeName() (typeName string) {
	separatorIndex := strings.Index(f.Type, ".")
	if separatorIndex < 0 {
		return
	}
	typeName = f.Type[separatorIndex+1:]
	return
}

func (f *Field) GetWantName() (wantName string) {
	wantName = "want" + toPublic(f.Name)
	return
}

func (f *Field) GetGotName() (gotName string) {
	gotName = "got" + toPublic(f.Name)
	return
}

func (f *Field) GetNameType() (nameType string) {
	nameType = f.Name + " " + f.Type
	return
}
