package model

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

	sdkContract := publicCloud.NewContract()
	sdkContract.SetBillingFrequency(1)
	sdkContract.SetTerm(4)
	sdkContract.SetType("HOURLY")
	sdkContract.SetEndsAt(endsAt)
	sdkContract.SetRenewalsAt(renewalsAt)
	sdkContract.SetCreatedAt(createdAt)
	sdkContract.SetState("RUNNING")

	contract := newContract(sdkContract)

	assert.Equal(
		t,
		int64(1),
		contract.BillingFrequency.ValueInt64(),
		"billingFrequency should be set",
	)
	assert.Equal(
		t,
		int64(4),
		contract.Term.ValueInt64(),
		"term should be set",
	)
	assert.Equal(
		t,
		"HOURLY",
		contract.Type.ValueString(),
		"type should be set",
	)
	assert.Equal(
		t,
		"2023-12-14 17:09:47 +0000 UTC",
		contract.EndsAt.ValueString(),
		"endsAt should be set",
	)
	assert.Equal(
		t,
		"2022-12-14 17:09:47 +0000 UTC",
		contract.RenewalsAt.ValueString(),
		"renewalsAt should be set",
	)
	assert.Equal(
		t,
		"2021-12-14 17:09:47 +0000 UTC",
		contract.CreatedAt.ValueString(),
		"createdAt should be set",
	)
	assert.Equal(
		t,
		"RUNNING",
		contract.State.ValueString(),
		"state should be set",
	)
}

func TestContract_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(context.TODO(), Contract{}.attributeTypes(), Contract{})

	assert.Nil(t, diags, "attributes should be correct")
}
