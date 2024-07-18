package enum

type RootDiskStorageType string

func (r RootDiskStorageType) String() string {
	return string(r)
}

func (r RootDiskStorageType) Values() []string {
	return convertStringEnumToValues(rootDiskStorageTypes)
}

const (
	RootDiskStorageTypeCentral RootDiskStorageType = "CENTRAL"
	RootDiskStorageTypeLocal   RootDiskStorageType = "LOCAL"
)

var rootDiskStorageTypes = []RootDiskStorageType{
	RootDiskStorageTypeCentral,
	RootDiskStorageTypeLocal,
}

func NewRootDiskStorageType(value string) (RootDiskStorageType, error) {
	return findEnumForString(value, rootDiskStorageTypes, RootDiskStorageTypeLocal)
}
