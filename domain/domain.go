package domain

type GoCodePackage struct {
	Path        string
	PackageName string
	FileList    []*GoCodeFile
}

type GoCodeFile struct {
	Name          string
	InterfaceList []*Interface
	MockList      []*Mock
	Code          string
	ImportList    []string
}

type Import struct {
	Name string
	Path string
}

type Interface struct {
	Name       string
	MethodList []*Method
	ImportList []string
}

type Mock struct {
	Struct         *Struct
	Constructor    *Constructor
	SetterList     []*Setter
	MethodList     []*Method
	CodeImportList []string
	TestImportList []string
}

type Struct struct {
	Name         string
	ReceiverName string
	WantName     string
	GotName      string
	FieldList    []*Field
	Code         string
	ImportList   []string
}

type Constructor struct {
	Name           string
	Code           string
	CodeTest       string
	CodeImportList []string
	TestImportList []string
}

type Setter struct {
	Name           string
	Field          *Field
	Code           string
	CodeTest       string
	CodeImportList []string
	TestImportList []string
}

type Method struct {
	Name               string
	ArgList            []*Field
	ResultList         []*Field
	ArgNameTypeList    []string
	ResultNameTypeList []string
	Code               string
	CodeTest           string
	CodeImportList     []string
	TestImportList     []string
}

type Field struct {
	Name           string
	WantName       string
	GotName        string
	Type           string
	BaseType       *Field
	NameType       string
	ExampleValue   string
	CodeImportList []string
	TestImportList []string
}
