package Printer

import (
	"github.com/AntonParaskiv/mockGen/domain"
)

func (p *Printer) GeneratePackageCode(mockPackage *domain.GoCodePackage) {
	p.mockPackage = mockPackage
	for _, mockFile := range p.mockPackage.FileList {
		// TODO: check if go code file (not test)
		p.generateFile(mockFile)

		if p.GenerateTests {
			mockFileTest := p.generateFileTest(mockFile)
			if mockFileTest != nil {
				mockPackage.FileList = append(mockPackage.FileList, mockFileTest)
			}
		}
	}
	return
}
