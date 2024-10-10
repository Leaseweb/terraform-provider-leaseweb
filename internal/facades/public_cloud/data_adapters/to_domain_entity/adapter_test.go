package to_domain_entity

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/resources/public_cloud/model"
	"github.com/stretchr/testify/assert"
)

// TODO Enable SSH key support
//var defaultSshKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDWvBbugarDWMkELKmnzzYaxPkDpS9qDokehBM+OhgrgyTWssaREYPDHsRjq7Ldv/8kTdK9i+f9HMi/BTskZrd5npFtO2gfSgFxeUALcqNDcjpXvQJxLUShNFmtxPtQLKlreyWB1r8mcAQBC/jrWD5I+mTZ7uCs4CNV4L0eLv8J1w=="

func TestAdaptToCreateInstanceOpts(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		instance := generateInstanceModel(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		got, err := AdaptToCreateInstanceOpts(
			instance,
			[]string{string(publicCloud.TYPENAME_M5A_4XLARGE)},
			context.TODO(),
		)

		assert.NoError(t, err)
		assert.Equal(t, "region", got.Region)
		assert.Equal(t, "lsw.m5a.4xlarge", got.Type)
		assert.Equal(t, enum.StorageTypeCentral, got.RootDiskStorageType)
		assert.Equal(t, "UBUNTU_20_04_64BIT", got.Image.Id)
		assert.Equal(t, enum.ContractTypeMonthly, got.Contract.Type)
		assert.Equal(t, enum.ContractTermThree, got.Contract.Term)
		assert.Equal(t,
			enum.ContractBillingFrequencyOne,
			got.Contract.BillingFrequency,
		)
	})

	t.Run("optional values are passed", func(t *testing.T) {
		instance := generateInstanceModel(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		got, err := AdaptToCreateInstanceOpts(
			instance,
			[]string{string(publicCloud.TYPENAME_M5A_4XLARGE)},
			context.TODO(),
		)

		assert.NoError(t, err)
		assert.Equal(t, "marketAppId", *got.MarketAppId)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, 55, got.RootDiskSize.Value)
		// TODO Enable SSH key support
		//assert.Equal(t, defaultSshKey, got.SshKey.String())
	})

	t.Run(
		"returns error if invalid rootDiskStorageType is passed",
		func(t *testing.T) {
			rootDiskStorageType := "tralala"
			instance := generateInstanceModel(
				&rootDiskStorageType,
				nil,
				nil,
				nil,
				nil,
				nil,
			)

			_, err := AdaptToCreateInstanceOpts(
				instance,
				[]string{string(publicCloud.TYPENAME_M5A_4XLARGE)},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid instanceType is passed",
		func(t *testing.T) {
			instanceType := "tralala"
			instance := generateInstanceModel(
				nil,
				nil,
				nil,
				nil,
				nil,
				&instanceType,
			)

			_, err := AdaptToCreateInstanceOpts(
				instance,
				[]string{string(publicCloud.TYPENAME_M5A_4XLARGE)},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractType is passed",
		func(t *testing.T) {
			contractType := "tralala"
			instance := generateInstanceModel(
				nil,
				&contractType,
				nil,
				nil,
				nil,
				nil,
			)

			_, err := AdaptToCreateInstanceOpts(
				instance,
				[]string{string(publicCloud.TYPENAME_M5A_4XLARGE)},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid contractTerm is passed",
		func(t *testing.T) {
			contractTerm := 555
			instance := generateInstanceModel(
				nil,
				nil,
				&contractTerm,
				nil,
				nil,
				nil,
			)

			_, err := AdaptToCreateInstanceOpts(
				instance,
				[]string{string(publicCloud.TYPENAME_M5A_4XLARGE)},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if invalid billingFrequency is passed",
		func(t *testing.T) {
			billingFrequency := 555
			instance := generateInstanceModel(
				nil,
				nil,
				nil,
				&billingFrequency,
				nil,
				nil,
			)

			_, err := AdaptToCreateInstanceOpts(
				instance,
				[]string{string(publicCloud.TYPENAME_M5A_4XLARGE)},
				context.TODO(),
			)

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

	  	_, err := AdaptToCreateInstanceOpts(
	  		instance,
	  		[]string{string(publicCloud.TYPENAME_M5A_4XLARGE)},
	  		context.TODO(),
	  	)

	  	assert.Error(t, err)
	  	assert.ErrorContains(t, err, "ssh key is invalid")
	  })
	*/

	t.Run(
		"returns error if invalid rootDiskSize is passed",
		func(t *testing.T) {
			rootDiskSize := 1
			instance := generateInstanceModel(
				nil,
				nil,
				nil,
				nil,
				&rootDiskSize,
				nil,
			)

			_, err := AdaptToCreateInstanceOpts(
				instance,
				[]string{string(publicCloud.TYPENAME_M5A_4XLARGE)},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "1")
		},
	)

	t.Run(
		"returns error if Instance cannot be created",
		func(t *testing.T) {
			instanceType := "instanceType"
			instance := generateInstanceModel(
				nil,
				nil,
				nil,
				nil,
				nil,
				&instanceType,
			)

			_, err := AdaptToCreateInstanceOpts(
				instance,
				[]string{},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "instanceType")
		},
	)
}

func TestAdaptToUpdateInstanceOpts(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		id := "id"
		instance := generateInstanceModel(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		instance.Id = basetypes.NewStringValue(id)

		got, diags := AdaptToUpdateInstanceOpts(
			instance,
			[]string{},
			"lsw.m5a.4xlarge",
			context.TODO(),
		)

		assert.Nil(t, diags)
		assert.Equal(t, id, got.Id)
	})

	t.Run("optional values are set", func(t *testing.T) {
		instance := generateInstanceModel(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		got, diags := AdaptToUpdateInstanceOpts(
			instance,
			[]string{},
			"lsw.m5a.4xlarge",
			context.TODO(),
		)

		assert.Nil(t, diags)
		assert.Equal(t, string(publicCloud.TYPENAME_M5A_4XLARGE), got.Type)
		assert.Equal(t, enum.ContractTypeMonthly, got.Contract.Type)
		assert.Equal(t, enum.ContractTermThree, got.Contract.Term)
		assert.Equal(
			t,
			enum.ContractBillingFrequencyOne,
			got.Contract.BillingFrequency,
		)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, 55, got.RootDiskSize.Value)
	})

	t.Run(
		"returns error if Instance cannot be created",
		func(t *testing.T) {
			instance := generateInstanceModel(
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
			)

			_, err := AdaptToUpdateInstanceOpts(
				instance,
				[]string{},
				"",
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "lsw.m5a.4xlarge")
		},
	)
}

func generateInstanceModel(
	rootDiskStorageType *string,
	contractType *string,
	contractTerm *int,
	billingFrequency *int,
	// TODO Enable SSH key support
	//sshKey *string,
	rootDiskSize *int,
	instanceTypeName *string,
) model.Instance {
	defaultRootDiskStorageType := "CENTRAL"
	defaultContractType := "MONTHLY"
	defaultContractTerm := 3
	defaultBillingFrequency := 1
	defaultRootDiskSize := 55
	defaultInstanceTypeName := "lsw.m5a.4xlarge"

	if rootDiskStorageType == nil {
		rootDiskStorageType = &defaultRootDiskStorageType
	}
	if contractType == nil {
		contractType = &defaultContractType
	}
	if contractTerm == nil {
		contractTerm = &defaultContractTerm
	}
	if billingFrequency == nil {
		billingFrequency = &defaultBillingFrequency
	}
	if rootDiskSize == nil {
		rootDiskSize = &defaultRootDiskSize
	}
	// TODO Enable SSH key support
	/**
	  if sshKey == nil {
	  	sshKey = &defaultSshKey
	  }
	*/
	if instanceTypeName == nil {
		instanceTypeName = &defaultInstanceTypeName
	}

	image, _ := types.ObjectValueFrom(
		context.TODO(),
		model.Image{}.AttributeTypes(),
		model.Image{
			Id: basetypes.NewStringValue("UBUNTU_20_04_64BIT"),
		},
	)

	contract, _ := types.ObjectValueFrom(
		context.TODO(),
		model.Contract{}.AttributeTypes(),
		model.Contract{
			BillingFrequency: basetypes.NewInt64Value(int64(*billingFrequency)),
			Term:             basetypes.NewInt64Value(int64(*contractTerm)),
			Type:             basetypes.NewStringValue(*contractType),
			State:            basetypes.NewStringUnknown(),
		},
	)

	instance := model.Instance{
		Id:                  basetypes.NewStringValue("id"),
		Region:              basetypes.NewStringValue("region"),
		Type:                basetypes.NewStringValue(*instanceTypeName),
		RootDiskStorageType: basetypes.NewStringValue(*rootDiskStorageType),
		RootDiskSize:        basetypes.NewInt64Value(int64(*rootDiskSize)),
		Image:               image,
		Contract:            contract,
		MarketAppId:         basetypes.NewStringValue("marketAppId"),
		Reference:           basetypes.NewStringValue("reference"),
		//SshKey:              basetypes.NewStringValue(*sshKey),
	}

	return instance
}
