package domain

type GoCodeFile struct {
	Name          string
	InterfaceList []*Interface
	MockList      []*Mock
	Code          string
	ImportList    []*Import
}

func (f *GoCodeFile) GetMockByName(mockName string) (mock *Mock) {
	for _, mockListItem := range f.MockList {
		if mockListItem.Struct.Name == mockName {
			mock = mockListItem
			return
		}
	}
	return
}
