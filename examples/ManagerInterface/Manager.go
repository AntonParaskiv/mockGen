package ManagerInterface

// TODO: types: err, custom, self

type Manager interface {
	Registration(nickName string, password string) (accountId int64, checkCode string)
	SignIn(accountId int64, password string) (nickName string)
}

//type Factory interface {
//	Create(accountId int64) (manager *Manager)
//}
