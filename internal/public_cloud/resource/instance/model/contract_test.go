package model

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newContract(t *testing.T) {
	endsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2023-12-14 17:09:47",
	)
	renewalsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2022-12-14 17:09:47",
	)
	createdAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2021-12-14 17:09:47",
	)

	sdkContract := publicCloud.NewContract(
		1,
		4,
		"HOURLY",
		*publicCloud.NewNullableTime(&endsAt),
		renewalsAt,
		createdAt,
		"RUNNING",
	)
	got, diags := newContract(context.TODO(), *sdkContract)

	assert.Nil(t, diags)

	assert.Equal(
		t,
		int64(1),
		got.BillingFrequency.ValueInt64(),
		"billingFrequency should be set",
	)
	assert.Equal(
		t,
		int64(4),
		got.Term.ValueInt64(),
		"term should be set",
	)
	assert.Equal(
		t,
		"HOURLY",
		got.Type.ValueString(),
		"type should be set",
	)
	assert.Equal(
		t,
		"2023-12-14 17:09:47 +0000 UTC",
		got.EndsAt.ValueString(),
		"endsAt should be set",
	)
	assert.Equal(
		t,
		"2022-12-14 17:09:47 +0000 UTC",
		got.RenewalsAt.ValueString(),
		"renewalsAt should be set",
	)
	assert.Equal(
		t,
		"2021-12-14 17:09:47 +0000 UTC",
		got.CreatedAt.ValueString(),
		"createdAt should be set",
	)
	assert.Equal(
		t,
		"RUNNING",
		got.State.ValueString(),
		"state should be set",
	)
}

func TestContract_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		Contract{}.AttributeTypes(),
		Contract{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}
