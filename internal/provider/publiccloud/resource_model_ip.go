package publiccloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type ResourceModelIp struct {
	Ip types.String `tfsdk:"ip"`
}

func (i ResourceModelIp) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip": types.StringType,
	}
}

func newResourceModelIpFromIp(ctx context.Context, sdkIp publicCloud.Ip) (*ResourceModelIp, error) {
	return &ResourceModelIp{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}, nil
}

func newResourceModelIpFromIpDetails(
	ctx context.Context,
	sdkIpDetails publicCloud.IpDetails,
) (*ResourceModelIp, error) {
	return &ResourceModelIp{
		Ip: basetypes.NewStringValue(sdkIpDetails.Ip),
	}, nil
}
