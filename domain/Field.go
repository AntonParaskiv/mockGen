package domain

import "strings"

type Field struct {
	Name           string
	Type           string
	BaseType       *Field
	ExampleValue   string
	CodeImportList []*Import
	TestImportList []*Import
}

func (f *Field) GetPublicName() (publicName string) {
	publicName = toPublic(f.Name)
	return
}

//////////////////// basic types
// +bool
// +string
// +int  int8  int16  int32  int64
// +uint uint8 uint16 uint32 uint64 uintptr
// +byte // alias for uint8
// +rune // alias for int32 // represents a Unicode code point
// +float32 float64
// +complex64 complex128

//////////////////// advanced types
// +*
// +[]
// +map[]
// custom type

//////////////////// special types
// struct
// interface {}
// error

func (f *Field) GetTypeType() (typeType int64) {
	// check basic type
	basicType, ok := MapBasicTypeStringToInt[f.Type]
	if ok {
		typeType = basicType
		return
	}

	switch f.Type {
	case "interface{}":
		typeType = FieldTypeInterface
	case "error":
		typeType = FieldTypeError
	default:
		switch {
		case len(f.Type) >= 2 && f.Type[0:2] == "[]":
			typeType = FieldTypeArray
		case len(f.Type) >= 4 && f.Type[0:4] == "map[":
			typeType = FieldTypeMap
		case len(f.Type) > 1 && f.Type[0] == '*':
			typeType = FieldTypePointer
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
	}
	return
}

func (f *Field) GetTypeGroup() (typeGroup int64) {
	typeGroup = MapBasicTypeToGroup[f.GetTypeType()]
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
