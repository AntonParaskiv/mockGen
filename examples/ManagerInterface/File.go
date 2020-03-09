package ManagerInterface

type File struct {
	Name   string
	Path   string
	isOpen bool
}

func NewFile() (f *File) {
	f = new(File)
	return
}

func (f *File) SetName(name string) *File {
	f.Name = name
	return f
}
