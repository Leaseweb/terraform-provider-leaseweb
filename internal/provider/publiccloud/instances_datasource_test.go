package publiccloud

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptContractToContractDataSource(t *testing.T) {
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

	want := contractDataSourceModel{
		BillingFrequency: basetypes.NewInt64Value(1),
		Term:             basetypes.NewInt64Value(3),
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

	sdkInstance := publicCloud.Instance{
		Id:        "id",
		Region:    "region",
		Reference: *publicCloud.NewNullableString(&reference),
		Image: publicCloud.Image{
			Id: "imageId",
		},
		State:               publicCloud.STATE_CREATING,
		Type:                publicCloud.TYPENAME_C3_2XLARGE,
		RootDiskSize:        50,
		RootDiskStorageType: publicCloud.STORAGETYPE_CENTRAL,
		Ips: []publicCloud.Ip{
			{Ip: "127.0.0.1"},
		},
		Contract: publicCloud.Contract{
			Term: publicCloud.CONTRACTTERM__1,
		},
		MarketAppId: *publicCloud.NewNullableString(&marketAppId),
	}

	got := adaptInstanceToInstanceDataSource(sdkInstance)

	assert.Equal(t, "id", got.ID.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "imageId", got.Image.ID.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())
	assert.Equal(t, int64(50), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Len(t, got.IPs, 1)
	assert.Equal(t, "127.0.0.1", got.IPs[0].IP.ValueString())
	assert.Equal(t, int64(1), got.Contract.Term.ValueInt64())
	assert.Equal(t, "marketAppId", got.MarketAppID.ValueString())
}

func Test_adaptInstancesToInstancesDatasource(t *testing.T) {
	sdkInstances := []publicCloud.Instance{
		{Id: "id"},
	}

	got := adaptInstancesToInstancesDataSource(sdkInstances)

	assert.Len(t, got.Instances, 1)
	assert.Equal(t, "id", got.Instances[0].ID.ValueString())
}

func Test_adaptImageToImageDatasource(t *testing.T) {
	sdkImage := publicCloud.Image{
		Id: "imageId",
	}

	want := imageDataSourceModel{
		ID: basetypes.NewStringValue("imageId"),
	}
	got := adaptImageToImageDataSource(sdkImage)

	assert.Equal(t, want, got)
}
