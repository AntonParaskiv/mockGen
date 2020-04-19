package Interactor

import "github.com/AntonParaskiv/mockGen/domain"

func (i *Interactor) createMockFile(interfaceFile *domain.GoCodeFile) {
	i.mockFile = &domain.GoCodeFile{
		Name:       interfaceFile.Name,
		ImportList: append([]*domain.Import{}, interfaceFile.ImportList...), // TODO: check useless imports
	}

	for _, iFace := range interfaceFile.InterfaceList {
		mock := i.createMock(iFace)

		if i.CreateFieldExamples {
			i.createMockFieldsExampleValues(mock)
		}

		i.mockFile.MockList = append(i.mockFile.MockList, mock)
	}

	i.mockPackage.FileList = append(i.mockPackage.FileList, i.mockFile)
	return
}
