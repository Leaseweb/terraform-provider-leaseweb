package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type loadBalancer struct {
	Id        types.String `tfsdk:"id"`
	Type      types.String `tfsdk:"type"`
	Resources resources    `tfsdk:"resource"`
	Region    types.String `tfsdk:"region"`
	Reference types.String `tfsdk:"reference"`
	State     types.String `tfsdk:"state"`
	Contract  contract     `tfsdk:"contract"`
	StartedAt types.String `tfsdk:"started_at"`
}

func newLoadBalancer(sdkLoadBalancer publicCloud.LoadBalancer) *loadBalancer {
	return &loadBalancer{
		Id:        basetypes.NewStringValue(sdkLoadBalancer.GetId()),
		Type:      basetypes.NewStringValue(sdkLoadBalancer.GetType()),
		Resources: newResources(sdkLoadBalancer.GetResources()),
		Region:    basetypes.NewStringValue(sdkLoadBalancer.GetRegion()),
		Reference: basetypes.NewStringValue(sdkLoadBalancer.GetReference()),
		State:     basetypes.NewStringValue(string(sdkLoadBalancer.GetState())),
		Contract:  newContract(sdkLoadBalancer.GetContract()),
		StartedAt: basetypes.NewStringValue(sdkLoadBalancer.GetStartedAt().String()),
	}
}
