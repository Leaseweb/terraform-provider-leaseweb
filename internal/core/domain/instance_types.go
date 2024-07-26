package domain

type InstanceTypes []InstanceType

func (i InstanceTypes) ToArray() []string {
	var values []string
	for _, instanceType := range i {
		values = append(values, instanceType.String())
	}

	return values
}
