package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

func Test_newLoadBalancer(t *testing.T) {
	reference := "reference"
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	id := value_object.NewGeneratedUuid()

	entityLoadBalancer := entity.NewLoadBalancer(
		id,
		"type",
		entity.Resources{Cpu: entity.Cpu{Unit: "resources"}},
		"region",
		enum.StateCreating,
		entity.Contract{BillingFrequency: enum.ContractBillingFrequencySix},
		entity.Ips{{Ip: "1.2.3.4"}},
		entity.OptionalLoadBalancerValues{
			Reference:      &reference,
			StartedAt:      &startedAt,
			PrivateNetwork: &entity.PrivateNetwork{Id: "privateNetworkId"},
			Configuration:  &entity.LoadBalancerConfiguration{Balance: enum.BalanceSource},
		},
	)

	got := newLoadBalancer(entityLoadBalancer)

	assert.Equal(t, id.String(), got.Id.ValueString(), "id is set")
	assert.Equal(
		t,
		"type",
		got.Type.ValueString(),
		"type is set",
	)
	assert.Equal(
		t,
		"resources",
		got.Resources.Cpu.Unit.ValueString(),
		"resources is set",
	)
	assert.Equal(
		t,
		"region",
		got.Region.ValueString(),
		"region is set",
	)
	assert.Equal(
		t,
		"reference",
		got.Reference.ValueString(),
		"reference is set",
	)
	assert.Equal(
		t,
		"CREATING",
		got.State.ValueString(),
		"state is set",
	)
	assert.Equal(
		t,
		int64(6),
		got.Contract.BillingFrequency.ValueInt64(),
		"contract is set",
	)
	assert.Equal(
		t,
		"2019-09-08 00:00:00 +0000 UTC",
		got.StartedAt.ValueString(),
		"startedAt is set",
	)
	assert.Equal(
		t,
		"1.2.3.4",
		got.Ips[0].Ip.ValueString(),
		"ips is set",
	)
	assert.Equal(
		t,
		"source",
		got.LoadBalancerConfiguration.Balance.ValueString(),
		"configuration is set",
	)
	assert.Equal(
		t,
		"privateNetworkId",
		got.PrivateNetwork.Id.ValueString(),
		"privateNetwork is set",
	)
}
