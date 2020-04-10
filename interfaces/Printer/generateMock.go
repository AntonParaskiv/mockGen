package Printer

import "github.com/AntonParaskiv/mockGen/domain"

func generateMock(mock *domain.Mock) (code string) {
	code += generateStruct(mock)
	code += generateConstructor(mock)

	for _, setter := range mock.SetterList {
		code += generateSetter(mock, setter)
	}

	for _, method := range mock.MethodList {
		code += generateMethod(mock, method)
	}
	return
}

func generateMockTest(mock *domain.Mock) (code string) {
	code += generateConstructorTest(mock)

	for _, setter := range mock.SetterList {
		code += generateSetterTest(mock, setter)
	}

	for _, method := range mock.MethodList {
		code += generateMethodTest(mock, method)
	}
	return
}
