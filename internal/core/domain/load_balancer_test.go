package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/shared/enum"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

func TestNewLoadBalancer(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		id := value_object.NewGeneratedUuid()
		instanceType, _ := value_object.NewInstanceType(
			"instanceType",
			[]string{"instanceType"},
		)

		got := NewLoadBalancer(
			id,
			*instanceType,
			Resources{Cpu: Cpu{Unit: "cpu"}},
			"region",
			enum.StateRunning,
			Contract{Type: enum.ContractTypeMonthly},
			Ips{{Ip: "1.2.3.4"}},
			OptionalLoadBalancerValues{},
		)

		assert.Equal(t, id, got.Id)
		assert.Equal(t, *instanceType, got.Type)
		assert.Equal(t, "cpu", got.Resources.Cpu.Unit)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, enum.StateRunning, got.State)
		assert.Equal(t, enum.ContractTypeMonthly, got.Contract.Type)
		assert.Equal(t, "1.2.3.4", got.Ips[0].Ip)

		assert.Nil(t, got.Reference)
		assert.Nil(t, got.StartedAt)
		assert.Nil(t, got.PrivateNetwork)
		assert.Nil(t, got.Configuration)
	})

	t.Run("optional values are set", func(t *testing.T) {
		reference := "reference"
		startedAt := time.Now()
		instanceType, _ := value_object.NewInstanceType(
			"instanceType",
			[]string{"instanceType"},
		)

		got := NewLoadBalancer(
			value_object.NewGeneratedUuid(),
			*instanceType,
			Resources{},
			"",
			enum.StateRunning,
			Contract{Type: enum.ContractTypeMonthly},
			Ips{},
			OptionalLoadBalancerValues{
				Reference:      &reference,
				StartedAt:      &startedAt,
				PrivateNetwork: &PrivateNetwork{Id: "privateNetworkId"},
				Configuration:  &LoadBalancerConfiguration{TargetPort: 54},
			},
		)

		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, startedAt, *got.StartedAt)
		assert.Equal(t, "privateNetworkId", got.PrivateNetwork.Id)
		assert.Equal(t, 54, got.Configuration.TargetPort)
	})

}
