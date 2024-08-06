package domain

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
)

type StorageTypes []enum.RootDiskStorageType

func (s StorageTypes) ToArray() []string {
	var values []string
	for _, storageType := range s {
		values = append(values, storageType.String())
	}

	return values
}
