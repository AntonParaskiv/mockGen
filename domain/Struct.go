package domain

type Struct struct {
	Name       string
	FieldList  []*Field
	Code       string
	ImportList []*Import
}

func (s *Struct) GetPublicName() (publicName string) {
	publicName = toPublic(s.Name)
	return
}

func (s *Struct) GetReceiverName() (receiverName string) {
	receiverName = getReceiverName(s.Name)
	return
}

func (s *Struct) GetWantName() (wantName string) {
	wantName = "want" + toPublic(s.Name)
	return
}

func (s *Struct) GetGotName() (gotName string) {
	gotName = "got" + toPublic(s.Name)
	return
}
