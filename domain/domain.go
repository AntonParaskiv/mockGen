package domain

type Import struct {
	Key  string
	Name string
	Path string
}

type Interface struct {
	Name       string
	MethodList []*Method
	ImportList []*Import
}

type Mock struct {
	Struct         *Struct
	Constructor    *Constructor
	SetterList     []*Setter
	MethodList     []*Method
	CodeImportList []*Import
	TestImportList []*Import
}

type Struct struct {
	Name         string
	ReceiverName string
	WantName     string
	GotName      string
	FieldList    []*Field
	Code         string
	ImportList   []*Import
}

type Constructor struct {
	Name           string
	Code           string
	CodeTest       string
	CodeImportList []*Import
	TestImportList []*Import
}

type Setter struct {
	Name           string
	Field          *Field
	Code           string
	CodeTest       string
	CodeImportList []*Import
	TestImportList []*Import
}

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
