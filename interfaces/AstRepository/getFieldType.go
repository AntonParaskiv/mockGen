package AstRepository

import (
	"fmt"
	"go/ast"
)

func getFieldType(astFieldType ast.Node) (fieldType string, err error) {
	switch astType := astFieldType.(type) {

	case *ast.Ident:
		fieldType = getNodeName(astType)

	case *ast.InterfaceType:
		if len(astType.Methods.List) > 0 {
			err = fmt.Errorf("unsupported type interface{} with methods")
			return
		}
		fieldType = "interface{}"

	case *ast.ArrayType:
		var itemType string
		itemType, err = getFieldType(astType.Elt)
		if err != nil {
			err = fmt.Errorf("get array item type failed: %w", err)
			return
		}
		fieldType = fmt.Sprintf("[]%s", itemType)

	case *ast.MapType:
		var keyType, valueType string
		keyType, err = getFieldType(astType.Key)
		if err != nil {
			err = fmt.Errorf("get map key type failed: %w", err)
			return
		}
		valueType, err = getFieldType(astType.Value)
		if err != nil {
			err = fmt.Errorf("get map value type failed: %w", err)
			return
		}
		fieldType = fmt.Sprintf("map[%s]%s", keyType, valueType)

	case *ast.StructType:
		fieldType = fmt.Sprintf("struct {\n")
		for _, item := range astType.Fields.List {
			var itemType string
			itemType, err = getFieldType(item.Type)
			if err != nil {
				err = fmt.Errorf("get struct field type failed: %w", err)
				return
			}
			fieldType += fmt.Sprintf("	%s %s\n", getNodeName(item), itemType)
		}
		fieldType += fmt.Sprintf("}")

	// custom types // TODO: handle imports
	case *ast.SelectorExpr:
		fieldType = fmt.Sprintf("%s.%s", getNodeName(astType.X), getNodeName(astType.Sel))
		// TODO: get baseType

	case *ast.StarExpr:
		var baseFieldType string
		baseFieldType, err = getFieldType(astType.X)
		if err != nil {
			err = fmt.Errorf("get base field type failed: %w", err)
			return
		}
		fieldType = fmt.Sprintf("*%s", baseFieldType)

	default:
		err = fmt.Errorf("unsupported type")
		return
	}

	return
}
