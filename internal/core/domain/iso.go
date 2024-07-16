package domain

type Iso struct {
	Id   string
	Name string
}

func NewIso(id string, name string) Iso {
	return Iso{Id: id, Name: name}
}
