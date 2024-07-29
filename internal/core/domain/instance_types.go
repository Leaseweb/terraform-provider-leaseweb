package domain

type InstanceTypes []InstanceType

func (i InstanceTypes) ContainsName(name string) bool {
	for _, instanceType := range i {
		if name == instanceType.Name {
			return true
		}
	}

	return false
}

func (i InstanceTypes) ToArray() []string {
	var values []string
	for _, instanceType := range i {
		values = append(values, instanceType.String())
	}

	return values
}
