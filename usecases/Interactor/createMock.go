package Interactor

import "github.com/AntonParaskiv/mockGen/domain"

func (i *Interactor) createMock(iFace *domain.Interface) (mock *domain.Mock) {
	// init mock
	mock = i.initMock(iFace.Name)

	// create methods
	mock.MethodList = i.createMethodList(iFace.MethodList)

	// create struct fields
	for _, method := range mock.MethodList {
		syncFieldLists(&mock.Struct.FieldList, &method.ArgList)
		syncFieldLists(&mock.Struct.FieldList, &method.ResultList)
	}

	// create setters
	for _, field := range mock.Struct.FieldList {
		mock.SetterList = append(mock.SetterList, &domain.Setter{
			Name:  "Set" + field.GetPublicName(),
			Field: field,
		})
	}
	return
}

func syncFieldLists(structFieldList, methodFieldList *[]*domain.Field) {
methodFieldListLoop:
	for m, methodField := range *methodFieldList {
		for _, structField := range *structFieldList {
			if structField.Name == methodField.Name {
				// TODO: check type
				(*methodFieldList)[m] = structField
				continue methodFieldListLoop
			}
		}
		*structFieldList = append(*structFieldList, methodField)
		//method.CodeImportList = append(method.CodeImportList, methodField.CodeImportList...) // TODO: add unique
		//method.TestImportList = append(method.TestImportList, methodField.TestImportList...) // TODO: add unique
	}
}

func (i *Interactor) initMock(iFaceName string) (mock *domain.Mock) {
	mock = &domain.Mock{
		Struct: &domain.Struct{
			Name:      iFaceName,
			FieldList: []*domain.Field{},
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
