package ManagerInterface

type Manager interface {
	Registration(nickName string, password string) (accountId int64)
	SignIn(accountId int64, password string) (err error)
}
