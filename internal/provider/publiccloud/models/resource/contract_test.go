package resource

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestContract_attributeTypes(t *testing.T) {
	_, diags := types.ObjectValueFrom(
		context.TODO(),
		Contract{}.AttributeTypes(),
		Contract{},
	)

	assert.Nil(t, diags, "attributes should be correct")
}

func Test_newContract(t *testing.T) {
	endsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2023-12-14 17:09:47",
	)
	sdkContract := publicCloud.Contract{
		BillingFrequency: publicCloud.BILLINGFREQUENCY__1,
		Term:             publicCloud.CONTRACTTERM__3,
		Type:             publicCloud.CONTRACTTYPE_HOURLY,
		EndsAt:           *publicCloud.NewNullableTime(&endsAt),
		State:            publicCloud.CONTRACTSTATE_ACTIVE,
	}

	want := Contract{
		BillingFrequency: basetypes.NewInt64Value(1),
		Term:             basetypes.NewInt64Value(3),
		Type:             basetypes.NewStringValue("HOURLY"),
		EndsAt:           basetypes.NewStringValue("2023-12-14 17:09:47 +0000 UTC"),
		State:            basetypes.NewStringValue("ACTIVE"),
	}
	got, err := newContract(context.TODO(), sdkContract)

	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}

func TestIsContractTermValid(t *testing.T) {
	t.Run(
		"false is returned when contract term is monthly and contract term is 0",
		func(t *testing.T) {
			sdkContract := publicCloud.Contract{
				Term: publicCloud.CONTRACTTERM__0,
				Type: publicCloud.CONTRACTTYPE_MONTHLY,
			}

			contract, _ := newContract(context.TODO(), sdkContract)

			got, reason := contract.IsContractTermValid()

			assert.False(t, got)
			assert.Equal(t, ReasonContractTermCannotBeZero, reason)
		},
	)

	t.Run(
		"false is returned when contract term is hourly and contract term is not 0",
		func(t *testing.T) {
			sdkContract := publicCloud.Contract{
				Term: publicCloud.CONTRACTTERM__3,
				Type: publicCloud.CONTRACTTYPE_HOURLY,
			}

			contract, _ := newContract(context.TODO(), sdkContract)

			got, reason := contract.IsContractTermValid()

			assert.False(t, got)
			assert.Equal(t, ReasonContractTermMustBeZero, reason)
		},
	)

	t.Run(
		"true is returned when contract term is valid",
		func(t *testing.T) {
			sdkContract := publicCloud.Contract{
				Term: publicCloud.CONTRACTTERM__0,
				Type: publicCloud.CONTRACTTYPE_HOURLY,
			}

			contract, _ := newContract(context.TODO(), sdkContract)

			got, reason := contract.IsContractTermValid()

			assert.True(t, got)
			assert.Equal(t, ReasonNone, reason)
		},
	)
}
