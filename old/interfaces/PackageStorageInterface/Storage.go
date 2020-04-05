package PackageStorageInterface

type Storage interface {
	GetGoFileList(packagePath string) (fileList []string, err error)
	ReadFile(filePath string) (fileData []byte, err error)
}
