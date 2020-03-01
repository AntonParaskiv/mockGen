package PackageStorage

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Storage struct {
}

func New() (s *Storage) {
	s = new(Storage)
	return
}

func (s *Storage) GetGoFileList(packagePath string) (fileList []string, err error) {
	fileInfos, err := ioutil.ReadDir(packagePath)
	if err != nil {
		err = fmt.Errorf("read dir failed: %w", err)
		return
	}

	fileList = make([]string, 0)
	for _, fileInfo := range fileInfos {
		if isFileNameMatchGoCode(fileInfo.Name()) {
			filePath := filepath.Join(packagePath, fileInfo.Name())
			fileList = append(fileList, filePath)
		}
	}

	return
}

func (s *Storage) ReadFile(filePath string) (fileData []byte, err error) {
	fileData, err = ioutil.ReadFile(filePath)
	if err != nil {
		err = fmt.Errorf("read file %s failed: %w", filePath, err)
		return
	}
	return
}

func isFileNameMatchGoCode(fileName string) (isMatch bool) {
	// check non-go files
	fileExtension := filepath.Ext(fileName)
	if fileExtension != ".go" {
		return
	}

	// check test files
	pattern := "_test.go"
	fileNameEndingStartPosition := len(fileName) - len(pattern)
	if fileNameEndingStartPosition < 0 {
		isMatch = true
		return
	}
	fileNameEnding := fileName[fileNameEndingStartPosition:]
	if fileNameEnding != pattern {
		isMatch = true
		return
	}

	return
}
