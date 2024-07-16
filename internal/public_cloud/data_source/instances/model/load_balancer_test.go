package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

func Test_newLoadBalancer(t *testing.T) {
	reference := "reference"
	startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
	id := value_object.NewGeneratedUuid()

	entityLoadBalancer := domain.NewLoadBalancer(
		id,
		"type",
		domain.Resources{Cpu: domain.Cpu{Unit: "resources"}},
		"region",
		enum.StateCreating,
		domain.Contract{BillingFrequency: enum.ContractBillingFrequencySix},
		domain.Ips{{Ip: "1.2.3.4"}},
		domain.OptionalLoadBalancerValues{
			Reference:      &reference,
			StartedAt:      &startedAt,
			PrivateNetwork: &domain.PrivateNetwork{Id: "privateNetworkId"},
			Configuration:  &domain.LoadBalancerConfiguration{Balance: enum.BalanceSource},
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
