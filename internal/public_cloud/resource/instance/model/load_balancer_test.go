package model

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newLoadBalancer(t *testing.T) {
	contractType, _ := publicCloud.NewContractTypeFromValue("HOURLY")
	contract := *publicCloud.NewContract(
		5,
		0,
		*contractType,
		*publicCloud.NewNullableTime(nil),
		time.Now(),
		time.Now(),
		"state",
	)

	t.Run("loadBalancer Conversion works", func(t *testing.T) {
		reference := "reference"
		startedAt, _ := time.Parse(time.RFC3339, "2019-09-08T00:00:00Z")

		sdkLoadBalancerDetails := publicCloud.NewLoadBalancerDetails(
			"id",
			"type",
			*publicCloud.NewResources(
				*publicCloud.NewCpu(0, "cpu"),
				*publicCloud.NewMemory(0, ""),
				*publicCloud.NewNetworkSpeed(0, ""),
				*publicCloud.NewNetworkSpeed(0, ""),
			),
			"region",
			*publicCloud.NewNullableString(&reference),
			"state",
			contract,
			*publicCloud.NewNullableTime(&startedAt),
			[]publicCloud.IpDetails{{Ip: "1.2.3.4"}},
			*publicCloud.NewNullableLoadBalancerConfiguration(&publicCloud.LoadBalancerConfiguration{Balance: "balance"}),
			*publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{Id: "autoScalingGroupId"}),
			*publicCloud.NewNullablePrivateNetwork(&publicCloud.PrivateNetwork{PrivateNetworkId: "privateNetworkId"}),
		)

		got, gotDiags := newLoadBalancer(context.TODO(), *sdkLoadBalancerDetails)

		assert.Nil(t, gotDiags)

		assert.Equal(t, "id", got.Id.ValueString())
		assert.Equal(t, "type", got.Type.ValueString())
		assert.Equal(
			t,
			"{\"unit\":\"cpu\",\"value\":0}",
			got.Resources.Attributes()["cpu"].String(),
		)
		assert.Equal(t, "region", got.Region.ValueString())
		assert.Equal(t, "reference", got.Reference.ValueString())
		assert.Equal(t, "state", got.State.ValueString())

		assert.Equal(
			t,
			"5",
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
			"balance",
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
