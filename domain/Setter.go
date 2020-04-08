package domain

type Setter struct {
	Name           string
	Field          *Field
	Code           string
	CodeTest       string
	CodeImportList []*Import
	TestImportList []*Import
}
