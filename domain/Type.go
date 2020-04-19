package domain

const (
	FieldTypeUnknown = iota
	FieldTypeBool
	FieldTypeString
	FieldTypeInt
	FieldTypeInt8
	FieldTypeInt16
	FieldTypeInt32
	FieldTypeInt64
	FieldTypeUint
	FieldTypeUint8
	FieldTypeUint16
	FieldTypeUint32
	FieldTypeUint64
	FieldTypeUintPtr
	FieldTypeByte
	FieldTypeRune
	FieldTypeFloat32
	FieldTypeFloat64
	FieldTypeComplex64
	FieldTypeComplex128

	FieldTypeInterface
	FieldTypeError
	FieldTypeArray
	FieldTypeMap
	FieldTypeImportedCustomType
	FieldTypeLocalCustomType
	FieldTypePointer
	FieldTypeChan
)

var MapBasicTypeStringToInt = map[string]int64{
	"bool":       FieldTypeBool,
	"string":     FieldTypeString,
	"int":        FieldTypeInt,
	"int8":       FieldTypeInt8,
	"int16":      FieldTypeInt16,
	"int32":      FieldTypeInt32,
	"int64":      FieldTypeInt64,
	"uint":       FieldTypeUint,
	"uint8":      FieldTypeUint8,
	"uint16":     FieldTypeUint16,
	"uint32":     FieldTypeUint32,
	"uint64":     FieldTypeUint64,
	"uintptr":    FieldTypeUintPtr,
	"byte":       FieldTypeByte,
	"rune":       FieldTypeRune,
	"float32":    FieldTypeFloat32,
	"float64":    FieldTypeFloat64,
	"complex64":  FieldTypeComplex64,
	"complex128": FieldTypeComplex128,
}

const (
	FieldTypeGroupUnknown = iota
	FieldTypeGroupBool
	FieldTypeGroupString
	FieldTypeGroupNumber
	FieldTypeGroupByte
	FieldTypeGroupRune
)

var MapBasicTypeToGroup = map[int64]int64{
	FieldTypeUnknown:    FieldTypeGroupUnknown,
	FieldTypeBool:       FieldTypeGroupBool,
	FieldTypeString:     FieldTypeGroupString,
	FieldTypeInt:        FieldTypeGroupNumber,
	FieldTypeInt8:       FieldTypeGroupNumber,
	FieldTypeInt16:      FieldTypeGroupNumber,
	FieldTypeInt32:      FieldTypeGroupNumber,
	FieldTypeInt64:      FieldTypeGroupNumber,
	FieldTypeUint:       FieldTypeGroupNumber,
	FieldTypeUint8:      FieldTypeGroupNumber,
	FieldTypeUint16:     FieldTypeGroupNumber,
	FieldTypeUint32:     FieldTypeGroupNumber,
	FieldTypeUint64:     FieldTypeGroupNumber,
	FieldTypeUintPtr:    FieldTypeGroupNumber,
	FieldTypeByte:       FieldTypeGroupByte,
	FieldTypeRune:       FieldTypeGroupRune,
	FieldTypeFloat32:    FieldTypeGroupNumber,
	FieldTypeFloat64:    FieldTypeGroupNumber,
	FieldTypeComplex64:  FieldTypeGroupNumber,
	FieldTypeComplex128: FieldTypeGroupNumber,
}
