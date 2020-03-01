package domain

type Package struct {
	Name  string
	Files []*File
}

func NewPackage() (p *Package) {
	p = new(Package)
	p.Files = make([]*File, 0)
	return
}
