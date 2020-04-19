package domain

import "fmt"

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

func (p *GoCodePackage) GetPackageLine() (packageLine string) {
	packageLine = fmt.Sprintf("package %s\n\n", p.PackageName)
	return
}
