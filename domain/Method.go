package domain

type Method struct {
	Name               string
	ArgList            []*Field
	ResultList         []*Field
	ArgNameTypeList    []string
	ResultNameTypeList []string
	Code               string
	CodeTest           string
	CodeImportList     []*Import
	TestImportList     []*Import
}
