package Interface

import "github.com/AntonParaskiv/mockGen/domain/Method"

type Interface struct {
	Name       string
	MethodList []*Method.Method
}
