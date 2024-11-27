package publiccloud

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptContractToContractDataSource(t *testing.T) {
	endsAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2023-12-14 17:09:47",
	)
	sdkContract := publiccloud.Contract{
		BillingFrequency: publiccloud.BILLINGFREQUENCY__1,
		Term:             publiccloud.CONTRACTTERM__3,
		Type:             publiccloud.CONTRACTTYPE_HOURLY,
		EndsAt:           *publiccloud.NewNullableTime(&endsAt),
		State:            publiccloud.CONTRACTSTATE_ACTIVE,
	}

	want := contractDataSourceModel{
		BillingFrequency: basetypes.NewInt32Value(1),
		Term:             basetypes.NewInt32Value(3),
		Type:             basetypes.NewStringValue("HOURLY"),
		EndsAt:           basetypes.NewStringValue("2023-12-14 17:09:47 +0000 UTC"),
		State:            basetypes.NewStringValue("ACTIVE"),
	}
	got := adaptContractToContractDataSource(sdkContract)

	assert.Equal(t, want, got)
}

func Test_adaptInstanceToInstanceDataSource(t *testing.T) {
	reference := "reference"
	marketAppId := "marketAppId"

	sdkInstance := publiccloud.Instance{
		Id:        "id",
		Region:    "region",
		Reference: *publiccloud.NewNullableString(&reference),
		Image: publiccloud.Image{
			Id: "imageId",
		},
		State:               publiccloud.STATE_CREATING,
		Type:                publiccloud.TYPENAME_C3_2XLARGE,
		RootDiskSize:        50,
		RootDiskStorageType: publiccloud.STORAGETYPE_CENTRAL,
		Ips: []publiccloud.Ip{
			{Ip: "127.0.0.1"},
		},
		Contract: publiccloud.Contract{
			Term: publiccloud.CONTRACTTERM__1,
		},
		MarketAppId: *publiccloud.NewNullableString(&marketAppId),
	}

	got := adaptInstanceToInstanceDataSource(sdkInstance)

	assert.Equal(t, "id", got.ID.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "imageId", got.Image.ID.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())
	assert.Equal(t, int32(50), got.RootDiskSize.ValueInt32())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Len(t, got.IPs, 1)
	assert.Equal(t, "127.0.0.1", got.IPs[0].IP.ValueString())
	assert.Equal(t, int32(1), got.Contract.Term.ValueInt32())
	assert.Equal(t, "marketAppId", got.MarketAppID.ValueString())
}

func Test_adaptInstancesToInstancesDatasource(t *testing.T) {
	sdkInstances := []publiccloud.Instance{
		{Id: "id"},
	}

	got := adaptInstancesToInstancesDataSource(sdkInstances)

	assert.Len(t, got.Instances, 1)
	assert.Equal(t, "id", got.Instances[0].ID.ValueString())
}
