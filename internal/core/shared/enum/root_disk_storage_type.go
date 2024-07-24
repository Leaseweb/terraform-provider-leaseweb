package enum

import (
	"terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type RootDiskStorageType string

func (r RootDiskStorageType) String() string {
	return string(r)
}

func (r RootDiskStorageType) Values() []string {
	return enum_utils.ConvertStringEnumToValues(rootDiskStorageTypes)
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
	return enum_utils.FindEnumForString(value, rootDiskStorageTypes, RootDiskStorageTypeLocal)
}
