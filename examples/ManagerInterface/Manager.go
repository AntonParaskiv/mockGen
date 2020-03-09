package ManagerInterface

type Manager interface {
	Registration(nickName string, password string) (accountId int64)
	SignIn(accountId int64, password string) (nickName string, err error)
}

type factory interface {
	Create(accountId int64) (manager *Manager)
}
