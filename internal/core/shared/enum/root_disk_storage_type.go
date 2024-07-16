package enum

type RootDiskStorageType string

func (r RootDiskStorageType) String() string {
	return string(r)
}

const (
	RootDiskStorageTypeLocal   RootDiskStorageType = "LOCAL"
	RootDiskStorageTypeCentral RootDiskStorageType = "CENTRAL"
)

var rootDiskStorageTypes = []RootDiskStorageType{
	RootDiskStorageTypeLocal,
	RootDiskStorageTypeCentral,
}

func NewRootDiskStorageType(value string) (RootDiskStorageType, error) {
	return findEnumForString(value, rootDiskStorageTypes, RootDiskStorageTypeLocal)
}
