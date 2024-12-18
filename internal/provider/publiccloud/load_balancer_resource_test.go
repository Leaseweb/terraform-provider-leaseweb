package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		got, err := adaptLoadBalancerDetailsToLoadBalancerResource(
			loadBalancerDetails,
			context.TODO(),
		)

		require.NoError(t, err)
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

		got, err := adaptLoadBalancerDetailsToLoadBalancerResource(
			loadBalancerDetails,
			context.TODO(),
		)

		require.NoError(t, err)
		assert.Equal(t, "reference", got.Reference.ValueString())
	})

}

func Test_loadBalancerResourceModel_GetLaunchLoadBalancerOpts(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		loadBalancer := generateLoadBalancerResourceModel()
		loadBalancer.Reference = basetypes.NewStringPointerValue(nil)

		got, err := loadBalancer.GetLaunchLoadBalancerOpts(context.TODO())

		require.NoError(t, err)
		assert.Equal(t, publiccloud.REGIONNAME_EU_WEST_3, got.Region)
		assert.Equal(t, publiccloud.TYPENAME_C3_2XLARGE, got.Type)
		assert.Equal(t, publiccloud.CONTRACTTYPE_MONTHLY, got.ContractType)
		assert.Equal(t, publiccloud.CONTRACTTERM__3, got.ContractTerm)
		assert.Equal(t, publiccloud.BILLINGFREQUENCY__1, got.BillingFrequency)

		reference, _ := got.GetReferenceOk()
		assert.Nil(t, reference)
	})

	t.Run("optional values are passed", func(t *testing.T) {
		reference := "reference"

		loadBalancer := generateLoadBalancerResourceModel()
		loadBalancer.Reference = basetypes.NewStringPointerValue(&reference)

		got, err := loadBalancer.GetLaunchLoadBalancerOpts(context.TODO())

		require.NoError(t, err)
		assert.Equal(t, "reference", *got.Reference)
	})

	t.Run(
		"returns error if invalid instanceType is passed",
		func(t *testing.T) {
			loadBalancer := generateLoadBalancerResourceModel()
			loadBalancer.Type = basetypes.NewStringValue("tralala")

			_, err := loadBalancer.GetLaunchLoadBalancerOpts(context.TODO())

			require.Error(t, err)
			require.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractType is passed",
		func(t *testing.T) {
			contractType := "tralala"
			loadBalancer := generateLoadBalancerResourceModel()
			contract := generateContractObject(
				nil,
				nil,
				&contractType,
			)
			loadBalancer.Contract = contract

			_, err := loadBalancer.GetLaunchLoadBalancerOpts(context.TODO())

			require.Error(t, err)
			require.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractTerm is passed",
		func(t *testing.T) {
			contractTerm := 555
			loadBalancer := generateLoadBalancerResourceModel()
			contract := generateContractObject(
				nil,
				&contractTerm,
				nil,
			)
			loadBalancer.Contract = contract

			_, err := loadBalancer.GetLaunchLoadBalancerOpts(context.TODO())

			require.Error(t, err)
			require.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if invalid billingFrequency is passed",
		func(t *testing.T) {
			billingFrequency := 555
			loadBalancer := generateLoadBalancerResourceModel()
			contract := generateContractObject(
				&billingFrequency,
				nil,
				nil,
			)
			loadBalancer.Contract = contract

			_, err := loadBalancer.GetLaunchLoadBalancerOpts(context.TODO())

			require.Error(t, err)
			require.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if invalid region is passed",
		func(t *testing.T) {
			loadBalancer := generateLoadBalancerResourceModel()
			loadBalancer.Region = basetypes.NewStringValue("tralala")

			_, err := loadBalancer.GetLaunchLoadBalancerOpts(context.TODO())

			require.Error(t, err)
			require.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if resourceModelContract resource is incorrect",
		func(t *testing.T) {
			loadBalancer := generateLoadBalancerResourceModel()
			loadBalancer.Contract = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := loadBalancer.GetLaunchLoadBalancerOpts(context.TODO())

			require.Error(t, err)
			require.ErrorContains(t, err, ".contractResourceModel")
		},
	)
}

func Test_loadBalancerResourceModel_GetUpdateLoadBalancerOpts(t *testing.T) {
	t.Run("optional values are set", func(t *testing.T) {
		reference := "reference"
		loadBalancerType := string(publiccloud.TYPENAME_C3_2XLARGE)

		loadBalancer := generateLoadBalancerResourceModel()
		loadBalancer.Type = basetypes.NewStringPointerValue(&loadBalancerType)
		loadBalancer.Reference = basetypes.NewStringPointerValue(&reference)

		got, err := loadBalancer.GetUpdateLoadBalancerOpts()

		require.NoError(t, err)
		assert.Equal(t, publiccloud.TYPENAME_C3_2XLARGE, *got.Type)
		assert.Equal(t, "reference", *got.Reference)
	})

	t.Run(
		"returns error if invalid instanceType is passed",
		func(t *testing.T) {
			loadBalancer := generateLoadBalancerResourceModel()
			loadBalancer.Type = basetypes.NewStringValue("tralala")

			_, err := loadBalancer.GetUpdateLoadBalancerOpts()

			require.Error(t, err)
			require.ErrorContains(t, err, "tralala")
		},
	)
}

func generateLoadBalancerResourceModel() loadBalancerResourceModel {
	contract := generateContractObject(
		nil,
		nil,
		nil,
	)

	return loadBalancerResourceModel{
		ID:        basetypes.NewStringValue("305c0bd8-b157-4a9c-885a-e07df86a714f"),
		Region:    basetypes.NewStringValue(string(publiccloud.REGIONNAME_EU_WEST_3)),
		Type:      basetypes.NewStringValue(string(publiccloud.TYPENAME_C3_2XLARGE)),
		Reference: basetypes.NewStringPointerValue(nil),
		Contract:  contract,
	}
}
