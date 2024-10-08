package public_cloud

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewImage(t *testing.T) {
	state := "state"
	stateReason := "stateReason"
	createdAt := time.Now()
	updatedAt := time.Now()
	storageSize := StorageSize{Unit: "unit"}
	version := "version"
	architecture := "architecture"

	got := NewImage(
		"UBUNTU_24_04_64BIT",
		"name",
		&version,
		"family",
		"flavour",
		&architecture,
		&state,
		&stateReason,
		&Region{Name: "region"},
		&createdAt,
		&updatedAt,
		false,
		&storageSize,
		[]string{"marketApp"},
		[]string{"storageType"},
	)
	want := Image{
		Id:           "UBUNTU_24_04_64BIT",
		Name:         "name",
		Version:      &version,
		Family:       "family",
		Flavour:      "flavour",
		Architecture: &architecture,
		State:        &state,
		StateReason:  &stateReason,
		Region:       &Region{Name: "region"},
		CreatedAt:    &createdAt,
		UpdatedAt:    &updatedAt,
		Custom:       false,
		StorageSize:  &storageSize,
		MarketApps:   []string{"marketApp"},
		StorageTypes: []string{"storageType"},
	}

	assert.Equal(t, want, got)
}
