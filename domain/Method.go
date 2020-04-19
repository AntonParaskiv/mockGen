package domain

import "strings"

type Method struct {
	Name           string
	ArgList        []*Field
	ResultList     []*Field
	Code           string
	CodeTest       string
	CodeImportList []*Import
	TestImportList []*Import
}

func (m *Method) GetPrivateName() (privateName string) {
	privateName = toPrivate(m.Name)
	return
}

func (m *Method) GetArgLine() (argLine string) {
	argLine = getLinedFieldList(m.ArgList)
	return
}

func (m *Method) GetResultLine() (resultLine string) {
	resultLine = getLinedFieldList(m.ResultList)
	return
}

func getLinedFieldList(fieldList []*Field) (linedFieldNameTypeList string) {
	fieldNameTypeList := getFieldNameTypeList(fieldList)
	linedFieldNameTypeList = strings.Join(fieldNameTypeList, ", ")
	return
}

func getFieldNameTypeList(fieldList []*Field) (fieldNameTypeList []string) {
	for _, field := range fieldList {
		fieldNameTypeList = append(fieldNameTypeList, field.GetNameType())
	}
	return
}
