package FileStorageInterface

type Factory interface {
	Create(fileName string) (s Storage)
}
