package resource

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func TestNewFromInstance(t *testing.T) {
	marketAppId := "marketAppId"
	reference := "reference"

	instance := publicCloud.Instance{
		Id:                  "id",
		Type:                publicCloud.TYPENAME_C3_2XLARGE,
		Region:              "region",
		Reference:           *publicCloud.NewNullableString(&reference),
		MarketAppId:         *publicCloud.NewNullableString(&marketAppId),
		State:               publicCloud.STATE_CREATING,
		RootDiskSize:        50,
		RootDiskStorageType: publicCloud.STORAGETYPE_CENTRAL,
		Contract: publicCloud.Contract{
			Type: publicCloud.CONTRACTTYPE_MONTHLY,
		},
		Image: publicCloud.Image{
			Id: "UBUNTU_20_04_64BIT",
		},
		Ips: []publicCloud.Ip{
			{
				Ip: "127.0.0.1",
			},
		},
	}

	got, err := NewFromInstance(instance, context.TODO())

	assert.NoError(t, err)

	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, int64(50), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Equal(t, "marketAppId", got.MarketAppId.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())

	image := Image{}
	got.Image.As(context.TODO(), &image, basetypes.ObjectAsOptions{})
	assert.Equal(t, "UBUNTU_20_04_64BIT", image.Id.ValueString())

	contract := Contract{}
	got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(t, "MONTHLY", contract.Type.ValueString())

	var ips []Ip
	got.Ips.ElementsAs(context.TODO(), &ips, false)
	assert.Len(t, ips, 1)
	assert.Equal(t, "127.0.0.1", ips[0].Ip.ValueString())
}

func TestNewFromInstanceDetails(t *testing.T) {
	marketAppId := "marketAppId"
	reference := "reference"

	instance := publicCloud.InstanceDetails{
		Id:                  "id",
		Type:                publicCloud.TYPENAME_C3_2XLARGE,
		Region:              "region",
		Reference:           *publicCloud.NewNullableString(&reference),
		MarketAppId:         *publicCloud.NewNullableString(&marketAppId),
		State:               publicCloud.STATE_CREATING,
		RootDiskSize:        50,
		RootDiskStorageType: publicCloud.STORAGETYPE_CENTRAL,
		Contract: publicCloud.Contract{
			Type: publicCloud.CONTRACTTYPE_MONTHLY,
		},
		Image: publicCloud.Image{
			Id: "UBUNTU_20_04_64BIT",
		},
		Ips: []publicCloud.IpDetails{
			{
				Ip: "127.0.0.1",
			},
		},
	}

	got, err := NewFromInstanceDetails(instance, context.TODO())

	assert.NoError(t, err)

	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, int64(50), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Equal(t, "marketAppId", got.MarketAppId.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())

	image := Image{}
	got.Image.As(context.TODO(), &image, basetypes.ObjectAsOptions{})
	assert.Equal(t, "UBUNTU_20_04_64BIT", image.Id.ValueString())

	contract := Contract{}
	got.Contract.As(context.TODO(), &contract, basetypes.ObjectAsOptions{})
	assert.Equal(t, "MONTHLY", contract.Type.ValueString())

	var ips []Ip
	got.Ips.ElementsAs(context.TODO(), &ips, false)
	assert.Len(t, ips, 1)
	assert.Equal(t, "127.0.0.1", ips[0].Ip.ValueString())
}

func TestInstance_GetLaunchInstanceOpts(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		instance := generateInstanceModel()

		got, err := instance.GetLaunchInstanceOpts(context.TODO())

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

		got, err := instance.GetLaunchInstanceOpts(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, "marketAppId", *got.MarketAppId)
		assert.Equal(t, "reference", *got.Reference)
		assert.Equal(t, int32(55), *got.RootDiskSize)
	})

	t.Run(
		"returns error if invalid rootDiskStorageType is passed",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.RootDiskStorageType = basetypes.NewStringValue("tralala")

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if invalid instanceType is passed",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Type = basetypes.NewStringValue("tralala")

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

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

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

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

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

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

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "555")
		},
	)

	t.Run(
		"returns error if invalid region is passed",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Region = basetypes.NewStringValue("tralala")

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"returns error if Image resource is incorrect",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Image = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, ".Image")
		},
	)

	t.Run(
		"returns error if Contract resource is incorrect",
		func(t *testing.T) {
			instance := generateInstanceModel()
			instance.Contract = basetypes.NewObjectNull(map[string]attr.Type{})

			_, err := instance.GetLaunchInstanceOpts(context.TODO())

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
		Contract{}.AttributeTypes(),
		Contract{
			BillingFrequency: basetypes.NewInt64Value(int64(*billingFrequency)),
			Term:             basetypes.NewInt64Value(int64(*contractTerm)),
			Type:             basetypes.NewStringValue(*contractType),
			State:            basetypes.NewStringUnknown(),
		},
	)

	return contract
}

func generateInstanceModel() Instance {
	image, _ := types.ObjectValueFrom(
		context.TODO(),
		Image{}.AttributeTypes(),
		Image{
			Id: basetypes.NewStringValue("UBUNTU_20_04_64BIT"),
		},
	)

	contract := generateContractObject(nil, nil, nil)

	instance := Instance{
		Id:                  basetypes.NewStringValue("id"),
		Region:              basetypes.NewStringValue("eu-west-3"),
		Type:                basetypes.NewStringValue("lsw.m5a.4xlarge"),
		RootDiskStorageType: basetypes.NewStringValue("CENTRAL"),
		RootDiskSize:        basetypes.NewInt64Value(int64(55)),
		Image:               image,
		Contract:            contract,
		MarketAppId:         basetypes.NewStringValue("marketAppId"),
		Reference:           basetypes.NewStringValue("reference"),
	}

	return instance
}
