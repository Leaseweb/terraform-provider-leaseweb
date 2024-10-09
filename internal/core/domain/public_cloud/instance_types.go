package public_cloud

type InstanceTypes []string

func (i InstanceTypes) Contains(name string) bool {
	for _, instanceType := range i {
		if name == instanceType {
			return true
		}
	}

	return false
}
