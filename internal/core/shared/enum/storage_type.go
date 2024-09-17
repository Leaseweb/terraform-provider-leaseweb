package enum

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type StorageType string

func (r StorageType) String() string {
	return string(r)
}

func (r StorageType) Values() []string {
	return enum_utils.ConvertStringEnumToValues(storageTypes)
}

const (
	StorageTypeCentral StorageType = "CENTRAL"
	StorageTypeLocal   StorageType = "LOCAL"
)

var storageTypes = []StorageType{
	StorageTypeCentral,
	StorageTypeLocal,
}

func NewStorageType(value string) (StorageType, error) {
	return enum_utils.FindEnumForString(value, storageTypes, StorageTypeLocal)
}
