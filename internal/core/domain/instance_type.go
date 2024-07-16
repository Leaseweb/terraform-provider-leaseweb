package domain

type InstanceType struct {
	Name string
}

func NewInstanceType(name string) InstanceType {
	return InstanceType{Name: name}
}
