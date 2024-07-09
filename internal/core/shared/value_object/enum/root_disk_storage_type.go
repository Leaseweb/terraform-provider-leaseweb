package enum

type RootDiskStorageType string

type RootDiskStorageTypes []RootDiskStorageType

func (r RootDiskStorageType) String() string {
	return string(r)
}

const (
	RootDiskStorageTypeLocal   RootDiskStorageType = "LOCAL"
	RootDiskStorageTypeCentral RootDiskStorageType = "CENTRAL"
)

var RootDiskStorageTypeValues = RootDiskStorageTypes{
	RootDiskStorageTypeLocal,
	RootDiskStorageTypeCentral,
}
