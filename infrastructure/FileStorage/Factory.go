package FileStorage

import "github.com/AntonParaskiv/mockGen/interfaces/FileStorageInterface"

type Factory struct {
}

func NewFactory() (f *Factory) {
	f = new(Factory)
	return
}

func (f *Factory) Create(fileName string) (s FileStorageInterface.Storage) {
	s = New().SetFileName(fileName)
	return
}
