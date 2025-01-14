package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptLoadBalancerDetailsToLoadBalancerResource(t *testing.T) {
	t.Run("required fields are set", func(t *testing.T) {
		loadBalancerDetails := publiccloud.LoadBalancerDetails{
			Id:        "id",
			Region:    "region",
			Type:      publiccloud.TYPENAME_C3_2XLARGE,
			Reference: *publiccloud.NewNullableString(nil),
			Contract: publiccloud.Contract{
				Type: publiccloud.CONTRACTTYPE_MONTHLY,
			},
		}

		diags := diag.Diagnostics{}

		got := adaptLoadBalancerDetailsToLoadBalancerResource(
			loadBalancerDetails,
			context.TODO(),
			&diags,
		)

		assert.False(t, diags.HasError())
		assert.Equal(t, "id", got.ID.ValueString())
		assert.Equal(t, "region", got.Region.ValueString())
		assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())
		assert.Nil(t, got.Reference.ValueStringPointer())

		contract := contractResourceModel{}
		got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
		assert.Equal(t, "MONTHLY", contract.Type.ValueString())
	})

	t.Run("optional fields are set", func(t *testing.T) {
		reference := "reference"

		loadBalancerDetails := publiccloud.LoadBalancerDetails{
			Id:        "id",
			Region:    "region",
			Type:      publiccloud.TYPENAME_C3_2XLARGE,
			Reference: *publiccloud.NewNullableString(&reference),
			Contract: publiccloud.Contract{
				Type: publiccloud.CONTRACTTYPE_MONTHLY,
			},
		}

		diags := diag.Diagnostics{}

		got := adaptLoadBalancerDetailsToLoadBalancerResource(
			loadBalancerDetails,
			context.TODO(),
			&diags,
		)

		assert.False(t, diags.HasError())
		assert.Equal(t, "reference", got.Reference.ValueString())
	})
}
