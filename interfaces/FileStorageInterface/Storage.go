package FileStorageInterface

type Storage interface {
	ReadFile() (fileData []byte, err error)
}
