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
		EndsAt:           *publiccloud.NewNullableTime(&endsAt),
		Term:             publiccloud.CONTRACTTERM__3,
		Type:             publiccloud.CONTRACTTYPE_HOURLY,
		State:            publiccloud.CONTRACTSTATE_ACTIVE,
	}

	want := contractDataSourceModel{
		BillingFrequency: basetypes.NewInt32Value(1),
		EndsAt:           basetypes.NewStringValue("2023-12-14 17:09:47 +0000 UTC"),
		Term:             basetypes.NewInt32Value(3),
		Type:             basetypes.NewStringValue("HOURLY"),
		State:            basetypes.NewStringValue("ACTIVE"),
	}
	got := adaptContractToContractDataSource(sdkContract)

	assert.Equal(t, want, got)
}

func Test_adaptInstanceDetailsToInstanceDataSource(t *testing.T) {
	t.Run("expected model is returned", func(t *testing.T) {
		reference := "reference"
		marketAppId := "marketAppId"

		sdkIso := publiccloud.Iso{
			Id:   "isoId",
			Name: "isoName",
		}

		instanceDetails := publiccloud.InstanceDetails{
			Contract: publiccloud.Contract{
				Term: publiccloud.CONTRACTTERM__1,
			},
			Id: "id",
			Image: publiccloud.Image{
				Id: "imageId",
			},
			Ips: []publiccloud.IpDetails{
				{Ip: "127.0.0.1"},
			},
			Iso:                 *publiccloud.NewNullableIso(&sdkIso),
			MarketAppId:         *publiccloud.NewNullableString(&marketAppId),
			Reference:           *publiccloud.NewNullableString(&reference),
			Region:              "region",
			RootDiskSize:        50,
			RootDiskStorageType: publiccloud.STORAGETYPE_CENTRAL,
			State:               publiccloud.STATE_CREATING,
			Type:                publiccloud.TYPENAME_C3_2XLARGE,
		}

		got := adaptInstanceDetailsToInstanceDataSource(instanceDetails)
		iso := got.ISO

		assert.Equal(t, "id", got.ID.ValueString())
		assert.Len(t, got.IPs, 1)
		assert.Equal(t, "127.0.0.1", got.IPs[0].IP.ValueString())
		assert.Equal(t, "imageId", got.Image.ID.ValueString())
		assert.Equal(t, "isoId", got.ISO.ID.ValueString())
		assert.Equal(t, "isoName", got.ISO.Name.ValueString())
		assert.Equal(t, "marketAppId", got.MarketAppID.ValueString())
		assert.Equal(t, "reference", got.Reference.ValueString())
		assert.Equal(t, "region", got.Region.ValueString())
		assert.Equal(t, int32(50), got.RootDiskSize.ValueInt32())
		assert.Equal(t, "CENTRAL", got.RootDiskStorageType.ValueString())
		assert.Equal(t, "CREATING", got.State.ValueString())
		assert.Equal(t, int32(1), got.Contract.Term.ValueInt32())
		assert.Equal(t, "lsw.c3.2xlarge", got.Type.ValueString())
		assert.NotNil(t, iso)
		assert.Equal(t, "isoId", got.ISO.ID.ValueString())
		assert.Equal(t, "isoName", got.ISO.Name.ValueString())
	})

	t.Run("expected model is returned if iso is nil", func(t *testing.T) {
		instanceDetails := publiccloud.InstanceDetails{
			Iso: *publiccloud.NewNullableIso(nil),
		}

		got := adaptInstanceDetailsToInstanceDataSource(instanceDetails)

		assert.Nil(t, got.ISO)
	})
}

func Test_adaptInstancesToInstancesDatasource(t *testing.T) {
	instanceDetailsList := []publiccloud.InstanceDetails{
		{Id: "id"},
	}

	got := adaptInstancesToInstancesDataSource(instanceDetailsList)

	assert.Len(t, got.Instances, 1)
	assert.Equal(t, "id", got.Instances[0].ID.ValueString())
}

func Test_instanceDetailsList_orderById(t *testing.T) {
	unorderedList := instanceDetailsList{
		{Id: "b"},
		{Id: "c"},
		{Id: "a"},
	}
	got := unorderedList.orderById()

	assert.Equal(
		t,
		instanceDetailsList{
			{Id: "a"},
			{Id: "b"},
			{Id: "c"},
		},
		got,
	)
}
