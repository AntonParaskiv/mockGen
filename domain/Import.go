package domain

import "path/filepath"

type Import struct {
	Name string
	Path string
}

func (i *Import) GetCallingName() (callingName string) {
	if len(i.Name) > 0 {
		callingName = i.Name
		return
	}

	callingName = filepath.Base(i.Path)
	return
}
