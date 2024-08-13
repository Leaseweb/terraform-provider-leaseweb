package public_cloud

import (
	"testing"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/stretchr/testify/assert"
)

func TestStorageTypes_ToArray(t *testing.T) {
	storageTypes := StorageTypes{
		enum.RootDiskStorageTypeLocal,
		enum.RootDiskStorageTypeCentral,
	}

	got := storageTypes.ToArray()
	want := []string{"LOCAL", "CENTRAL"}

	assert.Equal(t, want, got)
}
