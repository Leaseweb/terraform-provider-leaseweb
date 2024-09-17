package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageType_String(t *testing.T) {
	got := StorageTypeLocal.String()

	assert.Equal(t, "LOCAL", got)

}

func TestNewStorageType(t *testing.T) {
	want := StorageTypeCentral
	got, err := NewStorageType("CENTRAL")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestStorageType_Values(t *testing.T) {
	want := []string{"CENTRAL", "LOCAL"}
	got := StorageTypeCentral.Values()

	assert.EqualValues(t, want, got)
}
