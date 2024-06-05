package model

import (
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newIso(t *testing.T) {
	sdkIso := publicCloud.NewIso()
	sdkIso.SetId("id")
	sdkIso.SetName("name")

	iso := newIso(*sdkIso)

	assert.Equal(t, "id", iso.Id.ValueString(), "id should be set")
	assert.Equal(t, "name", iso.Name.ValueString(), "name should be set")
}
