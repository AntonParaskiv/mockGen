package ManagerInterface

type Manager interface {
	Registration(nickName string, password string) (accountId int64, isAdult bool, sign rune, checkCode string, err error)
	SignIn(accountId int64, password string, id interface{}) (nickName string, Byte byte, balance float64)
}

//type Factory interface {
//	Create(accountId int64) (manager *Manager)
//}
