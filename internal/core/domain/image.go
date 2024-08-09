package domain

import (
	"time"
)

type Image struct {
	Id           string
	Name         string
	Version      *string
	Family       string
	Flavour      string
	Architecture *string
	State        *string
	StateReason  *string
	Region       *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	Custom       bool
	StorageSize  *StorageSize
	MarketApps   []string
	StorageTypes []string
}

func NewImage(
	id string,
	name string,
	version *string,
	family string,
	flavour string,
	architecture *string,
	State *string,
	stateReason *string,
	region *string,
	createdAt *time.Time,
	updatedAt *time.Time,
	custom bool,
	storageSize *StorageSize,
	marketApps []string,
	storageTypes []string,
) Image {
	return Image{
		Id:           id,
		Name:         name,
		Version:      version,
		Family:       family,
		Flavour:      flavour,
		Architecture: architecture,
		State:        State,
		StateReason:  stateReason,
		Region:       region,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Custom:       custom,
		StorageSize:  storageSize,
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
	}
}
