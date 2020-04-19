package domain

type Constructor struct {
	Name           string
	Code           string
	CodeTest       string
	CodeImportList []*Import
	TestImportList []*Import
}
