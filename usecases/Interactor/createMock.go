package Interactor

import "github.com/AntonParaskiv/mockGen/domain"

func (i *Interactor) createMock(iFace *domain.Interface) (mock *domain.Mock) {
	mock = i.initMock(iFace.Name)
	mock.MethodList = i.createMethodList(iFace.MethodList)

	argList := []*domain.Field{}
	resultList := []*domain.Field{}
	for _, method := range mock.MethodList {
		for _, arg := range method.ArgList {
			i.createFieldExampleValue(arg)
			method.CodeImportList = append(method.CodeImportList, arg.CodeImportList...) // TODO: add unique
			method.TestImportList = append(method.TestImportList, arg.TestImportList...) // TODO: add unique
		}
		for _, result := range method.ResultList {
			i.createFieldExampleValue(result)
			method.CodeImportList = append(method.CodeImportList, result.CodeImportList...) // TODO: add unique
			method.TestImportList = append(method.TestImportList, result.TestImportList...) // TODO: add unique
		}

		argList = append(argList, method.ArgList...)
		resultList = append(resultList, method.ResultList...)
		mock.CodeImportList = append(mock.CodeImportList, method.CodeImportList...) // TODO: add unique
		mock.TestImportList = append(mock.TestImportList, method.TestImportList...) // TODO: add unique
	}

ArgLoop:
	for _, arg := range argList {
		for _, mockField := range mock.Struct.FieldList {
			if mockField.Name == arg.Name {
				continue ArgLoop
			}
		}
		mock.Struct.FieldList = append(mock.Struct.FieldList, arg)
	}

ResultLoop:
	for _, result := range resultList {
		for _, mockField := range mock.Struct.FieldList {
			if mockField.Name == result.Name {
				continue ResultLoop
			}
		}
		mock.Struct.FieldList = append(mock.Struct.FieldList, result)
	}

	for _, field := range mock.Struct.FieldList {
		mock.SetterList = append(mock.SetterList, &domain.Setter{
			Name:  "Set" + field.GetPublicName(),
			Field: field,
		})
	}
	return
}

func (i *Interactor) initMock(iFaceName string) (mock *domain.Mock) {
	mock = &domain.Mock{
		Struct: &domain.Struct{
			Name: iFaceName,
		},
		Constructor: &domain.Constructor{
			Name: "New",
		},
	}

	basePackageName := cutPostfix(i.mockPackage.PackageName, "Mock")
	if mock.Struct.Name != basePackageName {
		mock.Constructor.Name += mock.Struct.GetPublicName()
	}

	return
}
