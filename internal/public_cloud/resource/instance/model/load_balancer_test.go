package model

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

func Test_newLoadBalancer(t *testing.T) {
	t.Run("loadBalancer Conversion works", func(t *testing.T) {
		reference := "reference"
		startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")
		id := value_object.NewGeneratedUuid()

		loadBalancer := domain.NewLoadBalancer(
			id,
			"type",
			domain.Resources{Cpu: domain.Cpu{Unit: "cpu"}},
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

		got, gotDiags := newLoadBalancer(context.TODO(), loadBalancer)

		assert.Nil(t, gotDiags)

		assert.Equal(t, id.String(), got.Id.ValueString())
		assert.Equal(t, "type", got.Type.ValueString())
		assert.Equal(
			t,
			"{\"unit\":\"cpu\",\"value\":0}",
			got.Resources.Attributes()["cpu"].String(),
		)
		assert.Equal(t, "region", got.Region.ValueString())
		assert.Equal(t, "reference", got.Reference.ValueString())
		assert.Equal(t, "CREATING", got.State.ValueString())

		assert.Equal(
			t,
			"6",
			got.Contract.Attributes()["billing_frequency"].String(),
		)

		assert.Equal(
			t,
			"2019-09-08 00:00:00 +0000 UTC",
			got.StartedAt.ValueString(),
		)

		var ips []Ip
		got.Ips.ElementsAs(
			context.TODO(),
			&ips,
			false,
		)
		assert.Equal(t, "1.2.3.4", ips[0].Ip.ValueString())

		loadBalancerConfiguration := LoadBalancerConfiguration{}
		got.LoadBalancerConfiguration.As(
			context.TODO(),
			&loadBalancerConfiguration,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			"source",
			loadBalancerConfiguration.Balance.ValueString(),
		)

		privateNetwork := PrivateNetwork{}
		got.PrivateNetwork.As(
			context.TODO(),
			&privateNetwork,
			basetypes.ObjectAsOptions{},
		)
		assert.Equal(
			t,
			"privateNetworkId",
			privateNetwork.Id.ValueString(),
		)
	})
}
