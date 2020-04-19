package Interactor

import "github.com/AntonParaskiv/mockGen/domain"

func (i *Interactor) CreateMockPackage(interfacePackage *domain.GoCodePackage) (mockPackage *domain.GoCodePackage) {
	i.interfacePackage = interfacePackage
	i.mockPackage = &domain.GoCodePackage{
		Path:        createMockPackagePath(interfacePackage.Path),
		PackageName: createMockPackageName(interfacePackage.PackageName),
	}

	for _, interfaceFile := range interfacePackage.FileList {
		i.createMockFile(interfaceFile)
	}

	return i.mockPackage
}

func createMockPackagePath(interfacePackagePath string) (mockPackagePath string) {
	mockPackagePath = createMockPackageName(interfacePackagePath)
	return
}

func createMockPackageName(interfacePackageName string) (mockPackageName string) {
	mockPackageName = cutPostfix(interfacePackageName, "Interface") + "Mock"
	return
}
