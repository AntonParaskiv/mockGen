package ManagerInterface

import "github.com/AntonParaskiv/mockGen/examples"

type Manager interface {
	Registration(nickName string, password examples.Password) (accountId int64, isAdult bool, sign rune, checkCode string, scores map[int][]map[string][]int, err error)
	SignIn(accountId int64, password examples.Password, id interface{}) (nickName string, Byte byte, balance float64, messages []string, messagesId []uint)
}

//type Manager interface {
//	SignIn(accountId int64) (nickName string)
//}

//type Factory interface {
//	Create(accountId int64) (manager Manager)
//}
