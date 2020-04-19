package Interactor

import (
	"github.com/AntonParaskiv/mockGen/domain"
)

func (i *Interactor) createMethodList(iFaceMethodList []*domain.Method) (mockMethodList []*domain.Method) {
	for _, iFaceMethod := range iFaceMethodList {
		mockMethod := i.createMethod(iFaceMethod)
		if mockMethod == nil {
			continue
		}
		mockMethodList = append(mockMethodList, mockMethod)
	}
	return
}

func (i *Interactor) createMethod(iFaceMethod *domain.Method) (mockMethod *domain.Method) {
	mockMethod = &domain.Method{
		Name: iFaceMethod.Name,
	}

	for _, iFaceArg := range iFaceMethod.ArgList {
		mockArg := i.createField(iFaceArg)
		mockMethod.ArgList = append(mockMethod.ArgList, mockArg)
	}

	for _, iFaceResult := range iFaceMethod.ResultList {
		mockResult := i.createField(iFaceResult)
		mockMethod.ResultList = append(mockMethod.ResultList, mockResult)
	}
	return
}
