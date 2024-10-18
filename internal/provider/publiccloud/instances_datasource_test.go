package publiccloud

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_newDataSourceModelContract(t *testing.T) {
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

	want := dataSourceModelContract{
		BillingFrequency: basetypes.NewInt64Value(1),
		Term:             basetypes.NewInt64Value(3),
		Type:             basetypes.NewStringValue("HOURLY"),
		EndsAt:           basetypes.NewStringValue("2023-12-14 17:09:47 +0000 UTC"),
		State:            basetypes.NewStringValue("ACTIVE"),
	}
	got := newDataSourceModelContract(sdkContract)

	assert.Equal(t, want, got)
}

func Test_newDataSourceModelInstance(t *testing.T) {
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

	got := newDataSourceModelInstance(sdkInstance)

	assert.Equal(t, "id", got.Id.ValueString())
	assert.Equal(t, "region", got.Region.ValueString())
	assert.Equal(t, "reference", got.Reference.ValueString())
	assert.Equal(t, "imageId", got.Image.Id.ValueString())
	assert.Equal(t, "CREATING", got.State.ValueString())
	assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())
	assert.Equal(t, int64(50), got.RootDiskSize.ValueInt64())
	assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
	assert.Len(t, got.Ips, 1)
	assert.Equal(t, "127.0.0.1", got.Ips[0].Ip.ValueString())
	assert.Equal(t, int64(1), got.Contract.Term.ValueInt64())
	assert.Equal(t, "marketAppId", got.MarketAppId.ValueString())
}

func Test_newDataSourceInstances(t *testing.T) {
	sdkInstances := []publicCloud.Instance{
		{Id: "id"},
	}

	got := newDataSourceModelInstances(sdkInstances)

	assert.Len(t, got.Instances, 1)
	assert.Equal(t, "id", got.Instances[0].Id.ValueString())
}

func Test_newDataSourceModelIp(t *testing.T) {
	sdkIp := publicCloud.Ip{
		Ip: "127.0.0.1",
	}

	want := dataSourceModelIp{
		Ip: basetypes.NewStringValue("127.0.0.1"),
	}
	got := newDataSourceModelIp(sdkIp)

	assert.Equal(t, want, got)
}

func Test_newDataSourceModelImage(t *testing.T) {
	sdkImage := publicCloud.Image{
		Id: "imageId",
	}

	want := dataSourceModelImage{
		Id: basetypes.NewStringValue("imageId"),
	}
	got := newDataSourceModelImage(sdkImage)

	assert.Equal(t, want, got)
}

func Test_instancesDataSource_Metadata(t *testing.T) {
	resp := datasource.MetadataResponse{}
	instancesDataSource := NewInstancesDataSource()

	instancesDataSource.Metadata(
		context.TODO(),
		datasource.MetadataRequest{ProviderTypeName: "tralala"},
		&resp,
	)

	assert.Equal(
		t,
		"tralala_public_cloud_instances",
		resp.TypeName,
		"Type name should be tralala_public_cloud_instances",
	)
}
