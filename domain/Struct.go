package domain

type Struct struct {
	Name         string
	ReceiverName string
	WantName     string
	GotName      string
	FieldList    []*Field
	Code         string
	ImportList   []*Import
}
