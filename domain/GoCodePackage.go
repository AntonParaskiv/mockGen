package domain

type GoCodePackage struct {
	Path        string
	FullPath    string
	PackageName string
	FileList    []*GoCodeFile
	SelfImport  *Import
}

func (p *GoCodePackage) GetMockByName(mockName string) (mock *Mock) {
	for _, file := range p.FileList {
		mock = file.GetMockByName(mockName)
		if mock != nil {
			return
		}
	}
	return
}
