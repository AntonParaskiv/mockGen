package ManagerInterface

import (
	helpers "github.com/AntonParaskiv/mockGen/examples"
)

type Manager interface {
	Registration(nickName string, password helpers.Password) (accountId int64, isAdult bool, sign rune, checkCode string, scores map[int][]map[string][]int, err error)
	SignIn(accountId int64, password helpers.Password, id interface{}) (nickName string, Byte byte, balance float64, messages []string, messagesId []uint)
}

//type Factory interface {
//	Create(accountId int64) (manager *Manager)
//}
