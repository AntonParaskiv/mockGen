package domain

type File struct {
	Name          string
	InterfaceList []*Interface
	StructList    []*Struct
}

func NewFile() (p *File) {
	p = new(File)
	p.InterfaceList = make([]*Interface, 0)
	return
}
