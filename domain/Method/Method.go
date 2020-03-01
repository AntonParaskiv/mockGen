package Method

import "github.com/AntonParaskiv/mockGen/domain/Variable"

type Method struct {
	Name      string
	ArgList   []*Variable.Variable
	ValueList []*Variable.Variable
}
