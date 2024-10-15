package to_opts

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/models/resource"
	"github.com/stretchr/testify/assert"
)

func TestAdaptToLaunchInstanceOpts(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		instance := generateInstanceModel()

		got, err := AdaptToLaunchInstanceOpts(instance, context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, publicCloud.REGIONNAME_EU_WEST_3, got.Region)
		assert.Equal(t, publicCloud.TYPENAME_M5A_4XLARGE, got.Type)
		assert.Equal(t, publicCloud.STORAGETYPE_CENTRAL, got.RootDiskStorageType)
		assert.Equal(t, "UBUNTU_20_04_64BIT", got.ImageId)
		assert.Equal(t, publicCloud.CONTRACTTYPE_MONTHLY, got.ContractType)
		assert.Equal(t, publicCloud.CONTRACTTERM__3, got.ContractTerm)
		assert.Equal(t, publicCloud.BILLINGFREQUENCY__1, got.BillingFrequency)
	})

	t.Run("optional values are passed", func(t *testing.T) {
		instance := generateInstanceModel()

		got, err := AdaptToLaunchInstanceOpts(instance, context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, "marketAppId", *got.MarketAppId)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, int32(55), *got.RootDiskSize)
		// TODO Enable SSH key support
		//assert.Equal(t, defaultSshKey, got.SshKey.String())
	})

	t.Run(
		"returns error if invalid rootDiskStorageType is passed",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.RootDiskStorageType = basetypes.NewStringValue("tralala")

			_, err := AdaptToLaunchInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid instanceType is passed",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Type = basetypes.NewStringValue("tralala")

			_, err := AdaptToLaunchInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractType is passed",
		func(t *testing.T) {
			contractType := "tralala"
			instance := generateInstanceModel()
			contract := generateContractObject(
				nil,
				nil,
				&contractType,
			)
			instance.Contract = contract

			_, err := AdaptToLaunchInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractTerm is passed",
		func(t *testing.T) {
			contractTerm := 555
			instance := generateInstanceModel()
			contract := generateContractObject(
				nil,
				&contractTerm,
				nil,
			)
			instance.Contract = contract

			_, err := AdaptToLaunchInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if invalid billingFrequency is passed",
		func(t *testing.T) {
			billingFrequency := 555
			instance := generateInstanceModel()
			contract := generateContractObject(
				&billingFrequency,
				nil,
				nil,
			)
			instance.Contract = contract

			_, err := AdaptToLaunchInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	// TODO Enable SSH key support
	/**
	  t.Run("returns error if invalid sshKey is passed", func(t *testing.T) {
	  	sshKey := "tralala"
	  	instance := generateInstanceModel(
	  		nil,
	  		nil,
	  		nil,
	  		nil,
	  		&sshKey,
	  		nil,
	  		nil,
	  	)

	  	_, err := AdaptToLaunchInstanceOpts(instance, context.TODO())

	  	assert.Error(t, err)
	  	assert.ErrorContains(t, err, "ssh key is invalid")
	  })
	*/

	t.Run(
		"returns error if invalid region is passed",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Region = basetypes.NewStringValue("tralala")

			_, err := AdaptToLaunchInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if Image resource is incorrect",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Image = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := AdaptToLaunchInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, ".Image")
		},
	)

	t.Run(
		"returns error if Contract resource is incorrect",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Contract = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := AdaptToLaunchInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, ".Contract")
		},
	)
}

func TestAdaptToUpdateInstanceOpts(t *testing.T) {
	t.Run("optional values are set", func(t *testing.T) {
		instance := generateInstanceModel()

		got, diags := AdaptToUpdateInstanceOpts(instance, context.TODO())

		assert.Nil(t, diags)
		assert.Equal(t, publicCloud.TYPENAME_M5A_4XLARGE, *got.Type)
		assert.Equal(t, publicCloud.CONTRACTTYPE_MONTHLY, *got.ContractType)
		assert.Equal(t, publicCloud.CONTRACTTERM__3, *got.ContractTerm)
		assert.Equal(t, publicCloud.BILLINGFREQUENCY__1, *got.BillingFrequency)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, int32(55), *got.RootDiskSize)
	})

	t.Run(
		"returns error if invalid instanceType is passed",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Type = basetypes.NewStringValue("tralala")

			_, err := AdaptToUpdateInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractType is passed",
		func(t *testing.T) {
			contractType := "tralala"
			instance := generateInstanceModel()
			contract := generateContractObject(
				nil,
				nil,
				&contractType,
			)
			instance.Contract = contract

			_, err := AdaptToUpdateInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractTerm is passed",
		func(t *testing.T) {
			contractTerm := 555
			instance := generateInstanceModel()
			contract := generateContractObject(
				nil,
				&contractTerm,
				nil,
			)
			instance.Contract = contract

			_, err := AdaptToUpdateInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if invalid billingFrequency is passed",
		func(t *testing.T) {
			billingFrequency := 555
			instance := generateInstanceModel()
			contract := generateContractObject(
				&billingFrequency,
				nil,
				nil,
			)
			instance.Contract = contract

			_, err := AdaptToUpdateInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if Contract resource is incorrect",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Contract = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := AdaptToUpdateInstanceOpts(instance, context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, ".Contract")
		},
	)
}

func generateContractObject(
	billingFrequency *int,
	contractTerm *int,
	contractType *string,
) types.Object {
	defaultContractType := "MONTHLY"
	defaultContractTerm := 3
	defaultBillingFrequency := 1

	if contractType == nil {
		contractType = &defaultContractType
	}
	if contractTerm == nil {
		contractTerm = &defaultContractTerm
	}
	if billingFrequency == nil {
		billingFrequency = &defaultBillingFrequency
	}

	contract, _ := types.ObjectValueFrom(
		context.TODO(),
		resource.Contract{}.AttributeTypes(),
		resource.Contract{
			BillingFrequency: basetypes.NewInt64Value(int64(*billingFrequency)),
			Term:             basetypes.NewInt64Value(int64(*contractTerm)),
			Type:             basetypes.NewStringValue(*contractType),
			State:            basetypes.NewStringUnknown(),
		},
	)

	return contract
}

func generateInstanceModel(
// TODO Enable SSH key support
// sshKey *string,
) resource.Instance {
	// TODO Enable SSH key support
	/**
	  if sshKey == nil {
	  	sshKey = &defaultSshKey
	  }
	*/

	image, _ := types.ObjectValueFrom(
		context.TODO(),
		resource.Image{}.AttributeTypes(),
		resource.Image{
			Id: basetypes.NewStringValue("UBUNTU_20_04_64BIT"),
		},
	)

	contract := generateContractObject(nil, nil, nil)

	instance := resource.Instance{
		Id:                  basetypes.NewStringValue("id"),
		Region:              basetypes.NewStringValue("eu-west-3"),
		Type:                basetypes.NewStringValue("lsw.m5a.4xlarge"),
		RootDiskStorageType: basetypes.NewStringValue("CENTRAL"),
		RootDiskSize:        basetypes.NewInt64Value(int64(55)),
		Image:               image,
		Contract:            contract,
		MarketAppId:         basetypes.NewStringValue("marketAppId"),
		Reference:           basetypes.NewStringValue("reference"),
		//SshKey:              basetypes.NewStringValue(*sshKey),
	}

	return instance
}
