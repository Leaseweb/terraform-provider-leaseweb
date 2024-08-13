package public_cloud

type InstanceType struct {
	Name         string
	Resources    Resources
	StorageTypes *StorageTypes
	Prices       Prices
}

func (i InstanceType) String() string {
	return i.Name
}

type OptionalInstanceTypeValues struct {
	StorageTypes *StorageTypes
}

func NewInstanceType(
	name string,
	resources Resources,
	prices Prices,
	optional OptionalInstanceTypeValues,
) InstanceType {
	instanceType := InstanceType{
		Name:      name,
		Resources: resources,
		Prices:    prices,
	}

	if optional.StorageTypes != nil {
		instanceType.StorageTypes = optional.StorageTypes
	}

	return instanceType
}
