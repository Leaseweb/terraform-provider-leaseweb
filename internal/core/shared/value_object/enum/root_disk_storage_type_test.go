package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootDiskStorageType_String(t *testing.T) {
	got := RootDiskStorageTypeLocal.String()

	assert.Equal(t, "LOCAL", got)

}

func TestNewRootDiskStorageType(t *testing.T) {
	want := RootDiskStorageTypeCentral
	got, err := NewRootDiskStorageType("CENTRAL")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
