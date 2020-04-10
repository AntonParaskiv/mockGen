package Printer

import (
	"fmt"
	"github.com/AntonParaskiv/mockGen/domain"
	"path/filepath"
)

func (p *Printer) generateFile(mockFile *domain.GoCodeFile) {
	var code string
	var importList []*domain.Import
	for _, mock := range mockFile.MockList {
		code += generateMock(mock)
		importList = append(importList, mock.CodeImportList...) // TODO: add unique
	}

	if len(code) == 0 {
		return
	}
	mockFile.ImportList = append(mockFile.ImportList, importList...)

	mockFile.Code =
		p.mockPackage.GetPackageLine() +
			createImportList(mockFile.ImportList) +
			code
}

func (p *Printer) generateFileTest(mockFile *domain.GoCodeFile) (mockTestFile *domain.GoCodeFile) {
	var code string
	var importList []*domain.Import
	for _, mock := range mockFile.MockList {
		code += generateMockTest(mock)
		importList = append(importList, mock.TestImportList...) // TODO: add unique
	}

	if len(code) == 0 {
		return
	}

	mockTestFile = &domain.GoCodeFile{
		Name: createTestFilePath(mockFile.Name),
		ImportList: []*domain.Import{
			{Path: "reflect"},
			{Path: "testing"},
		},
	}
	mockTestFile.ImportList = append(mockTestFile.ImportList, importList...)

	mockTestFile.Code =
		p.mockPackage.GetPackageLine() +
			createImportList(mockTestFile.ImportList) +
			code
	return
}

func createTestFilePath(filePath string) (testFilePath string) {
	extension := filepath.Ext(filePath)
	if extension == ".go" {
		filePathLen := len(filePath)
		testFilePath = filePath[:filePathLen-3] + "_test.go"
	}
	return
}

func createImportList(importList []*domain.Import) (code string) {
	switch len(importList) {
	case 0:
	case 1:
		code = fmt.Sprintf("import %s \"%s\"\n\n", importList[0].Name, importList[0].Path)
	default:
		code = fmt.Sprintf("import (\n")
		for _, Import := range importList {
			code += fmt.Sprintf("	%s \"%s\"\n", Import.Name, Import.Path)
		}
		code += fmt.Sprintf(")\n\n")
	}
	return
}
