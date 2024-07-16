package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
)

func Test_newImage(t *testing.T) {
	entityImage := domain.NewImage(
		enum.Ubuntu200464Bit,
		"name",
		"version",
		"family",
		"flavour",
		"architecture",
		[]string{"one"},
		[]string{"storageType"},
	)

	got, diags := newImage(context.TODO(), entityImage)

	assert.Nil(t, diags)

	assert.Equal(
		t,
		"UBUNTU_20_04_64BIT",
		got.Id.ValueString(),
		"id should be set",
	)
	assert.Equal(
		t,
		"name",
		got.Name.ValueString(),
		"name should be set",
	)
	assert.Equal(
		t,
		"version",
		got.Version.ValueString(),
		"version should be set",
	)
	assert.Equal(
		t,
		"family",
		got.Family.ValueString(),
		"family should be set",
	)
	assert.Equal(
		t,
		"flavour",
		got.Flavour.ValueString(),
		"flavour should be set",
	)
	assert.Equal(
		t,
		"architecture",
		got.Architecture.ValueString(),
		"architecture should be set",
	)

	var marketApps []string
	got.MarketApps.ElementsAs(context.TODO(), &marketApps, false)
	assert.Len(t, marketApps, 1)
	assert.Equal(
		t,
		"one",
		marketApps[0],
		"marketApps should be set",
	)

	var storageTypes []string
	got.StorageTypes.ElementsAs(context.TODO(), &storageTypes, false)
	assert.Len(t, storageTypes, 1)
	assert.Equal(
		t,
		"storageType",
		storageTypes[0],
		"storageTypes should be set",
	)
}
