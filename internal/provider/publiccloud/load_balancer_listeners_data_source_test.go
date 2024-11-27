package publiccloud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptLoadBalancerListenersToLoadBalancerListenersDataSource(t *testing.T) {
	sdkListeners := []publiccloud.LoadBalancerListener{{Id: "id"}}
	got := adaptLoadBalancerListenersToLoadBalancerListenersDataSource(sdkListeners)

	want := loadBalancerListenersDataSourceModel{
		Listeners: []loadBalancerListenerDataSourceModel{
			{
				ID: basetypes.NewStringValue("id"),
			},
		},
	}

	assert.Equal(t, want, got)
}

func Test_loadBalancerListenersDataSourceModel_generateRequest(t *testing.T) {
	listeners := loadBalancerListenersDataSourceModel{
		LoadBalancerID: basetypes.NewStringValue("id"),
	}
	api := publiccloud.PubliccloudAPIService{}

	want := api.GetLoadBalancerListenerList(context.TODO(), "id")

	got := listeners.generateRequest(context.TODO(), &api)

	assert.Equal(t, want, got)
}
