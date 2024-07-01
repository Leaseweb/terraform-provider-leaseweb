package model

import (
	"context"
	"testing"
	"time"

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

		sdkLoadBalancer := publicCloud.NewLoadBalancer(
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
		)

		got, gotDiags := newLoadBalancer(context.TODO(), *sdkLoadBalancer)

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
	})
}
