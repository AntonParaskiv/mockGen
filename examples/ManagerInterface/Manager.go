package ManagerInterface

type Bool bool
type String string
type Int64 int64
type Uint64 uint64
type UintPtr uintptr
type Complex128 complex128
type Rune rune
type Byte byte

type Manager interface {
	Registration(nickName chan []int)
	//RegistrationBasic(isAdult bool, nickName string) (accountId int64, userId uint64, xVar uintptr, complex complex128, sign rune, flags byte)
	//RegistrationBasicPointer(isAdult *bool, nickName *string) (accountId *int64, userId *uint64, xVar *uintptr, complex *complex128, sign *rune, flags *byte)
	//RegistrationBasicArray(isAdult []bool, nickName []string) (accountId []int64, userId []uint64, xVar []uintptr, complex []complex128, sign []rune, flags []byte)
	//RegistrationBasicPointerArray(isAdult []*bool, nickName []*string) (accountId []*int64, userId []*uint64, xVar []*uintptr, complex []*complex128, sign []*rune, flags []*byte)
	//RegistrationBasicMap(isAdult map[bool]bool, nickName map[string]string) (accountId map[int64]int64, userId map[uint64]uint64, xVar map[uintptr]uintptr, complex map[complex128]complex128, sign map[rune]rune, flags map[byte]byte)
	//RegistrationBasicMapPointer(isAdult map[bool]*bool, nickName map[string]*string) (accountId map[int64]*int64, userId map[uint64]*uint64, xVar map[uintptr]*uintptr, complex map[complex128]*complex128, sign map[rune]*rune, flags map[byte]*byte)
	//RegistrationBasicCustom(isAdult Bool, nickName String) (accountId Int64, userId Uint64, xVar UintPtr, complex Complex128, sign Rune, flags Byte)

	// TODO:
	//RegistrationBasicCustomPointer(isAdult *Bool, nickName *String) (accountId *Int64, userId *Uint64, xVar *UintPtr, complex *Complex128, sign *Rune, flags *Byte)
	//RegistrationBasicCustomArray(isAdult []Bool, nickName []String) (accountId []Int64, userId []Uint64, xVar []UintPtr, complex []Complex128, sign []Rune, flags []Byte)
}

//type Manager interface {
//	Registration(nickName *string, password examples.Password) (accountId int64, isAdult bool, sign rune, checkCode string, scores map[int][]map[string][]int, err error)
//	SignIn(accountId int64, password examples.Password, id interface{}) (nickName *string, Byte byte, balance float64, messages []string, messagesId []uint)
//}

//type Factory interface {
//	Create(accountId int64) (manager *Manager)
//}
