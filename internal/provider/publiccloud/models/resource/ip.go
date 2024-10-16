package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Ip struct {
	Ip types.String `tfsdk:"ip"`
}

func (i Ip) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip": types.StringType,
	}
}

func newFromIp(ctx context.Context, sdkIp publicCloud.Ip) (*Ip, error) {
	return &Ip{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}, nil
}

func newFromIpDetails(
	ctx context.Context,
	sdkIpDetails publicCloud.IpDetails,
) (*Ip, error) {
	return &Ip{
		Ip: basetypes.NewStringValue(sdkIpDetails.Ip),
	}, nil
}
