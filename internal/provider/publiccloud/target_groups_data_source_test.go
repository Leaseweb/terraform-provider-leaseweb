package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v3/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptTargetGroupToTargetGroupDataSource(t *testing.T) {
	sdkTargetGroup := generateTargetGroup()

	got := adaptTargetGroupToTargetGroupDataSource(sdkTargetGroup)
	want := targetGroupDataSourceModel{
		ID:       basetypes.NewStringValue("id"),
		Name:     basetypes.NewStringValue("name"),
		Protocol: basetypes.NewStringValue("HTTP"),
		Port:     basetypes.NewInt32Value(80),
		Region:   basetypes.NewStringValue("eu-west-3"),
	}

	assert.Equal(t, want, got)
}

func Test_adaptTargetGroupsToTargetGroupsDataSource(t *testing.T) {
	targetGroups := []publiccloud.TargetGroup{
		generateTargetGroup(),
	}

	got := adaptTargetGroupsToTargetGroupsDataSource(targetGroups)
	want := targetGroupsDataSourceModel{
		TargetGroups: []targetGroupDataSourceModel{
			{
				ID:       basetypes.NewStringValue("id"),
				Name:     basetypes.NewStringValue("name"),
				Protocol: basetypes.NewStringValue("HTTP"),
				Port:     basetypes.NewInt32Value(80),
				Region:   basetypes.NewStringValue("eu-west-3"),
			},
		},
	}

	assert.Equal(t, want, got)
}

func generateTargetGroup() publiccloud.TargetGroup {
	return publiccloud.TargetGroup{
		Id:       "id",
		Name:     "name",
		Protocol: publiccloud.PROTOCOL_HTTP,
		Port:     80,
		Region:   publiccloud.REGIONNAME_EU_WEST_3,
	}
}

func Test_targetGroupsDataSourceModel_generateRequest(t *testing.T) {
	t.Run("request is generated as expected", func(t *testing.T) {
		targetGroups := targetGroupsDataSourceModel{
			ID:       basetypes.NewStringValue("id"),
			Name:     basetypes.NewStringValue("name"),
			Protocol: basetypes.NewStringValue("HTTP"),
			Port:     basetypes.NewInt32Value(80),
			Region:   basetypes.NewStringValue("eu-west-3"),
		}
		api := publiccloud.PubliccloudAPIService{}

		want := api.GetTargetGroupList(context.TODO()).
			Id("id").
			Name("name").
			Protocol("HTTP").
			Port(80).
			Region("eu-west-3")

		got, err := targetGroups.generateRequest(context.TODO(), &api)

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("invalid protocol returns an error", func(t *testing.T) {
		targetGroups := targetGroupsDataSourceModel{
			Protocol: basetypes.NewStringValue("tralala"),
		}

		_, err := targetGroups.generateRequest(context.TODO(), &publiccloud.PubliccloudAPIService{})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})

	t.Run("invalid region returns an error", func(t *testing.T) {
		targetGroups := targetGroupsDataSourceModel{
			Region: basetypes.NewStringValue("tralala"),
		}

		_, err := targetGroups.generateRequest(context.TODO(), &publiccloud.PubliccloudAPIService{})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "tralala")
	})
}
