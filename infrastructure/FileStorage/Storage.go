package FileStorage

import (
	"fmt"
	"io/ioutil"
)

type Storage struct {
	fileName string
}

func New() (s *Storage) {
	s = new(Storage)
	return
}

func (s *Storage) SetFileName(fileName string) *Storage {
	s.fileName = fileName
	return s
}

func (s *Storage) ReadFile() (fileData []byte, err error) {
	fileData, err = ioutil.ReadFile(s.fileName)
	if err != nil {
		err = fmt.Errorf("read file %s failed: %w", s.fileName, err)
		return
	}
	return
}
