package publiccloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_adaptLoadBalancerListenersToLoadBalancerListenersDataSource(t *testing.T) {
	sdkListeners := []publicCloud.LoadBalancerListener{{Id: "id"}}
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
