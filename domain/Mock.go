package domain

type Mock struct {
	Struct         *Struct
	Constructor    *Constructor
	SetterList     []*Setter
	MethodList     []*Method
	CodeImportList []*Import
	TestImportList []*Import
}
