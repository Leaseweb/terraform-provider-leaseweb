package public_cloud

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
)

type StorageTypes []enum.StorageType

func (s StorageTypes) ToArray() []string {
	var values []string
	for _, storageType := range s {
		values = append(values, storageType.String())
	}

	return values
}
